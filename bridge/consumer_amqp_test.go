package bridge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

type delivery struct {
	headers map[string]string
	body    []byte
}

// TestAMQPConsumer implements few integration tests for AMQP Consumer
func TestAMQPConsumer(t *testing.T) {
	url := os.Getenv("TEST_AMQP_URL")
	if url == "" {
		t.Skip("This test requires AMQP server, use environment variable TEST_AMQP_URL to set AMQP URL")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		t.Fatalf("Unable to connect to AMQP server: %v", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Unable to open AMQP channel: %v", err)
	}

	defer ch.Close()

	// create a queue for test
	queue, err := ch.QueueDeclare("", false, false, false, false, amqp.Table{})
	if err != nil {
		t.Fatalf("Unable to create test queue: %v", err)
	}

	// we can not use auto-delete because we run multiple tests
	defer ch.QueueDelete(queue.Name, false, false, false)

	// initialize AMQP consumer
	ctx := context.Background()
	dvs := make(chan delivery)

	queues := []Queue{
		{
			Name:        queue.Name,
			Parallelism: 1,
			Processor: func(c context.Context, h map[string]string, b []byte) error {
				dvs <- delivery{h, b}
				return nil
			},
		},
	}

	// make sure messages published to the queue are delivered to AMQP consumer
	t.Run("consume", func(t *testing.T) {
		cons := NewAMQPConsumer(ctx, url, queues, &nilLogger{})
		defer cons.Stop()

		payload := []byte("test message")

		// publish test message to the queue
		if err := ch.Publish("", queue.Name, false, false, amqp.Publishing{Body: []byte(payload)}); err != nil {
			t.Fatalf("Unable to publish message: %v", err)
		}

		// check if message has arrived
		select {
		case d := <-dvs:
			if !reflect.DeepEqual(d.body, payload) {
				t.Errorf("Message is received, but it's payload does not match expected value: want %s, got %s", payload, d.body)
			}
		case <-time.After(2 * time.Second):
			t.Errorf("Message is not received from AMQP server")
		}
	})

	// make sure AMQP consumer will re-connect after connection has been lost
	// this test uses RabbitMQ management plugin to simulate connection closing,
	// it's quite slow, but automates testing of an important feature
	t.Run("re-connect", func(t *testing.T) {
		mgmt := os.Getenv("TEST_MANAGEMENT_URL")
		if mgmt == "" {
			t.Skip("This test requires RabbitMQ Management API, use environment variable TEST_MANAGEMENT_URL to set API URL, it should include authorization credentials (eq. http://guest:guest@localhost:15672).")
		}

		cons := NewAMQPConsumer(ctx, url, queues, &nilLogger{})
		defer cons.Stop()

		// wait a bit so AMQP consumer has time to connect and appear in Management API
		time.Sleep(5 * time.Second)

		// interrupt connection
		if err := closeLatestConnection(mgmt); err != nil {
			t.Fatalf("Unable to interrupt connection: %v", err)
		}

		// publish test message to the queue
		if err := ch.Publish("", queue.Name, false, false, amqp.Publishing{Body: []byte{}}); err != nil {
			t.Fatalf("Unable to publish message: %v", err)
		}

		// check if message has arrived
		select {
		case <-dvs:
		case <-time.After(10 * time.Second):
			t.Errorf("Message is not received from AMQP server")
		}
	})

	// make sure AMQP consumer finishes message processing before stopping
	t.Run("finish processing", func(t *testing.T) {
		cons := NewAMQPConsumer(ctx, url, queues, &nilLogger{})

		if err := ch.Publish("", queue.Name, false, false, amqp.Publishing{Body: []byte{}}); err != nil {
			t.Fatalf("Unable to publish message: %v", err)
		}

		// give consumer a second to receive message, it won't finish processing it until we start reading dvs channel
		time.Sleep(time.Second)

		// signal consumer to stop
		stopped := make(chan struct{})
		go func() {
			cons.Stop()
			close(stopped)
		}()

		// make sure consumer hasn't stop
		select {
		case <-stopped:
			t.Fatal("consumer has been stopped before message has been processed")
		case <-time.After(time.Second):
		}

		// check if message has arrived
		select {
		case <-dvs:
		case <-time.After(time.Second):
			t.Errorf("Message is not received from AMQP server")
		}

		// make sure consumer has stopped after message got processed
		select {
		case <-stopped:
		case <-time.After(time.Second):
			t.Fatal("consumer has not been stopped after message has been processed")
		}
	})
}

func closeLatestConnection(url string) error {
	latest, err := getLatestConnection(url)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, url+"/api/connections/"+latest, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("unable to close connection, status code: %v", resp.StatusCode)
	}

	return nil
}

func getLatestConnection(url string) (string, error) {
	var conns []struct {
		Name        string `json:"name"`
		ConnectedAt int64  `json:"connected_at"`
	}

	resp, err := http.Get(url + "/api/connections")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("unable to fetch list of connections, status code %v", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&conns); err != nil {
		return "", err
	}

	var latest string
	var connected int64

	for _, c := range conns {
		if c.ConnectedAt > connected {
			latest = c.Name
			connected = c.ConnectedAt
		}
	}

	if latest == "" {
		return "", errors.New("no active connections")
	}

	return latest, nil
}
