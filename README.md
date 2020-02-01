## Setting up RabbitMQ (for development and production)

RabbitMQ is required for local development and prod. It's pretty easy to install though:

```
./bin/run_rabbit.sh
```

RabbitMQ works based on hostname:

> One of the important things to note about RabbitMQ is that it stores data based on what it calls the "Node Name", which defaults to the hostname

So `my-rabbit` in `run_rabbit.sh` is node name, and name (`libretaxi-rabbit`) is docker name.

`rabbitmq:3-management` contains UI plugin for queue management. Plugin port is 8080 (15672 in container).
Login **guest/guest**.

Login to RabbitUI here: http://localhost:8080

There is one queue:

* `messages` queue, http://localhost:8080/#/queues/%2F/messages - picked up by message handler, enqueued by libretaxi

Note that there is one libretaxi (message producer), and one message handler (message consumer).

Port 5672 is RabbitMQ itself.

**Important:** most likely you will run host/master application outside Docker, and workers will work
from docker. You need to make sure that you cast ports from docker's RabbitMQ to your 127.0.0.1, so host
application could reach RabbitMQ (workers should reach too, of course, but they should by default).

## LibreTaxi settings

Init settings for `./libretaxi.yml`:

```
telegram_token: YOUR_TOKEN
db_conn_str: postgres://libretaxi:libretaxi@localhost/libretaxi
rabbit_url: amqp://127.0.0.1/
```

## Running

Run:

```
dep ensure
go build
./libretaxi
```
