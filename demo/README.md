# AMQP FastCGI bridge Demo

This folder contains a demo environment for [amqp-cgi-bridge](https://github.com/skolodyazhnyy/amqp-cgi-bridge), it meant to help you quickly setup amqp-cgi-bridge and try it out.

## Setup

Demo environment uses [Docker](http://www.docker.com) to run AMQP and PHP-FPM servers, so you would need [Docker Engine](https://www.docker.com/community-edition#/download) and [docker-compose](https://docs.docker.com/compose/) before you continue.

> Alternatively, you can setup both servers on your own using your OS package manager. If you choose this option, make sure to modify `config.yml` to point to the addresses for AMQP and PHP-FPM servers. Also make sure `script_name` in each section has proper path to a script in `scripts` folder.

In case you don't have [docker](https://www.docker.com/community-edition#/download) and [docker-compose](https://docs.docker.com/compose/), follow installation guide in [docker documentation](https://docs.docker.com/).
Then, simply run command below in this folder:

```
docker-composer up -d
```

Docker compose will create two containers for AMQP and PHP-FPM servers.

Open [RabbitMQ management console](http://localhost:15672), login using default credentials (username: guest, password: guest) and go to Queue section.

Create queues: `logger`, `failure` and `reject` with default parameters. These queues are used for different demonstrations, you can find them in `config.yml`.

Now, download latest amqp-cgi-bridge from [GitHub](https://github.com/skolodyazhnyy/amqp-cgi-bridge/releases) into the demo folder. And run it using,

```
./amqp-cgi-bridge -config config.yml
```

You will see

