package bridge

import (
	"context"
	"github.com/streadway/amqp"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

type Processor func(ctx context.Context, headers map[string]interface{}, body []byte) error

type Queue struct {
	Name        string
	Parallelism int
	Requeue     bool
	Processor   Processor
}

type AMQPConsumer struct {
	url    string
	queues []Queue
	log    logger
	wg     sync.WaitGroup
	ctx    context.Context
	cancel func()
}

// NewAMQPConsumer constructs AMQP consumer and starts message processing routine
func NewAMQPConsumer(ctx context.Context, url string, queues []Queue, log logger) *AMQPConsumer {
	ctx, cancel := context.WithCancel(ctx)

	c := &AMQPConsumer{
		url:    url,
		queues: queues,
		log:    log,
		ctx:    ctx,
		cancel: cancel,
	}

	c.run()

	return c
}

// Run message processing routine
func (c *AMQPConsumer) run() {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()

		// re-connect loop: re-initialize connection to AMQP server in case an error occurs
		for {
			if err := c.serve(); err != nil {
				c.log.Errorf("An error occurred while serving AMQP connection: %v", err)
			}

			if isStopping(c.ctx) {
				return
			}

			c.log.Infof("Waiting 5 seconds before re-connect")

			if isStoppingWithTimeout(c.ctx, 5*time.Second) {
				return
			}
		}
	}()
}

// Stop AMQP consumer and wait for all routines to gracefully finish
func (c *AMQPConsumer) Stop() {
	c.cancel()
	c.wg.Wait()
}

// Serve handles AMQP connection attempt. When this method returns all resources used for serving connection should be
// released: go routines stopped, connections closed etc.
func (c *AMQPConsumer) serve() error {
	c.log.Infof("Connecting to AMQP server")
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}

	defer c.log.Infof("AMQP connection was closed")
	defer conn.Close()

	// create wait group for all individual queue consumers
	wg := sync.WaitGroup{}

	// create context for current connection attempt
	ctx, cancel := context.WithCancel(c.ctx)

	for _, queue := range c.queues {
		wg.Add(1)

		go func(queue Queue) {
			defer wg.Done()

			// consumer re-start loop: restarts consumer in case an error occurs
			for {
				if err := c.consume(ctx, queue, conn); err != nil {
					c.log.Errorf("An error occurred while consuming messages from %v: %v", queue.Name, err)
				}

				if isStopping(ctx) {
					return
				}

				c.log.Infof("Waiting 5 seconds before re-starting consumer for %v", queue.Name)

				if isStoppingWithTimeout(ctx, 5*time.Second) {
					return
				}
			}
		}(queue)
	}

	// handle connection closing notification
	closing := make(chan *amqp.Error)
	conn.NotifyClose(closing)

	select {
	case <-closing:
		cancel()
	case <-ctx.Done():
	}

	// wait for all individual queue consumers to stop
	wg.Wait()

	return nil
}

// Consume messages from individual queue. When this method returns all resources used by individual queue consumer
// should be released: go routines stopped, connections closed etc.
func (c *AMQPConsumer) consume(ctx context.Context, queue Queue, conn *amqp.Connection) error {
	c.log.Infof("Starting consumer for queue %v", queue.Name)
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	defer c.log.Infof("Consumer for queue %v has stopped", queue.Name)
	defer ch.Close()

	if err := ch.Qos(queue.Parallelism, 0, false); err != nil {
		return err
	}

	dv, err := ch.Consume(queue.Name, "", false, false, false, false, amqp.Table{})
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)

	sem := make(chan struct{}, queue.Parallelism)

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case d, ok := <-dv:
			if !ok {
				break loop
			}

			sem <- struct{}{}

			eg.Go(func() error {
				defer func() {
					<-sem
				}()

				c.log.Debugf("Processing message with ID #%v (%v)", d.MessageId, d.DeliveryTag)

				err := queue.Processor(ctx, d.Headers, d.Body)
				switch err {
				case nil:
					c.log.Debugf("Message with ID #%v (%v) successfully processed", d.MessageId, d.DeliveryTag)

					if err := d.Ack(false); err != nil {
						return err
					}
				case ErrProcessorInternal:
					// we couldn't deliver message to processor, so it make sense to put it back to the queue
					c.log.Debugf("Message with ID #%v (%v) not processed due an internal error", d.MessageId, d.DeliveryTag)

					// hold message for a little and put back to the queue, otherwise same message will be consumed right away
					// putting a lot of pressure on CPU
					select {
					case <-time.After(time.Second):
					case <-ctx.Done():
					}

					if err := d.Reject(true); err != nil {
						return err
					}
				default:
					c.log.Debugf("Message with ID #%v (%v) processed with error: %v", d.MessageId, d.DeliveryTag, err)

					if err := d.Reject(queue.Requeue); err != nil {
						return err
					}
				}

				return nil
			})
		}
	}

	return eg.Wait()
}

func isStopping(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func isStoppingWithTimeout(ctx context.Context, duration time.Duration) bool {
	select {
	case <-ctx.Done():
		return true
	case <-time.After(duration):
		return false
	}
}
