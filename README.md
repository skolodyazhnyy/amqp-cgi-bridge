# AMQP FastCGI bridge

AMQP FastCGI bridge is an AMQP consumer which consumes messages from AMQP server and process them using FastCGI server.
It was built to address problems of consuming AMQP messages with PHP. [More information](https://medium.com/@sergey.kolodyazhnyy/consuming-amqp-messages-in-php-6650c06936fa).

## Install

You can download binary from [GitHub release page](https://github.com/skolodyazhnyy/amqp-cgi-bridge/releases), or using go get:

```
go get github.com/skolodyazhnyy/amqp-cgi-bridge
```

## Usage

Application requires simple YAML configuration file which contains AMQP URL and list of queues to consume messages.

```yaml
# AMQP URI (see https://www.rabbitmq.com/uri-spec.html)
amqp_url: "amqp://localhost"

# an array of consumers
consumers:
  - # a queue to consume messages
    queue: "messages"
    # address of FastCGI server and name of the script to run to process messages
    fastcgi:
      net: "tcp"
      addr: "127.0.0.1:9000"
      script_name: "/path/to/script.php"
    # number of messages to be processed in parallel
    parallelism: 10
    # additional environment variables
    env:
      QUEUE_NAME: "messages"
```

Then, you need to configure and start PHP-FPM server (or any other FastCGI server) to process messages.

Your PHP script to process messages will work more or less same way as with Web Server, message body will be delivered
in request body, and AMQP headers will be available through `$_SERVER` variable.

