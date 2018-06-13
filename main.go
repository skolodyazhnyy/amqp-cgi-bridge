package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/skolodyazhnyy/amqp-cgi-bridge/bridge"
	"github.com/skolodyazhnyy/amqp-cgi-bridge/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
)

var version = "unknown"
var commit = "unknown"

var config struct {
	AMQPURL   string `yaml:"amqp_url"`
	Consumers []struct {
		Queue       string
		Parallelism int
		Env         map[string]string
		FastCGI     struct {
			Net        string
			Addr       string
			ScriptName string `yaml:"script_name"`
		}
	}
}

func load(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, v)
}

func main() {
	// parse flags
	filename := flag.String("config", "config.yml", "Configuration")
	logfmt := flag.String("log", "text", "Log format: json or text")
	printVersion := flag.Bool("v", false, "Print version")
	flag.Parse()

	if *printVersion {
		fmt.Println("Version", version)
		fmt.Println("Commit", commit)
		os.Exit(0)
	}

	logger := log.NewX(*logfmt, os.Stdout, log.DefaultTextFormat).With(log.R{
		"app":     "amqp-cgi-bridge",
		"version": version,
	})

	if err := load(*filename, &config); err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()
	var queues []bridge.Queue

	for _, c := range config.Consumers {
		p := bridge.NewFastCGIProcessor(
			c.FastCGI.Net,
			c.FastCGI.Addr,
			c.FastCGI.ScriptName,
			logger.Channel("fastcgi").With(log.R{
				"script_name": c.FastCGI.ScriptName,
			}),
		)

		if c.Env != nil {
			p = bridge.ProcessorWithEnv(p, c.Env)
		}

		queues = append(queues, bridge.Queue{
			Name:        c.Queue,
			Parallelism: c.Parallelism,
			Processor:   p,
		})
	}

	cons := bridge.NewAMQPConsumer(ctx, config.AMQPURL, queues, logger.Channel("amqp"))

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)

	s := <-signals
	logger.Infof("Signal %v received, stopping...", s)

	cons.Stop()
}
