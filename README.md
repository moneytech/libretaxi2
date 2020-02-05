## Running services

```
docker-compose up -d
```

Will run PostgreSQL and RabbitMQ with default credentials (see connection strings below).

## Setting up RabbitMQ (for development and production)

`rabbitmq:3-management` contains UI plugin for queue management. Plugin port is 8080 (15672 in container).
Login **guest/guest**.

Login to RabbitUI here: http://localhost:8080

There is only one queue at the moment:

* `messages` queue, http://localhost:8080/#/queues/%2F/messages - picked up by message handler, enqueued by libretaxi

Note that there is one message producer, and one message consumer threads (goroutines) in application.

Port 5672 is RabbitMQ itself.

## LibreTaxi settings

Init settings for `./libretaxi.yml`:

```
telegram_token: YOUR_TOKEN
db_conn_str: postgres://libretaxi:libretaxi@localhost/libretaxi
rabbit_url: amqp://127.0.0.1/
admin_channel_chat_id: -1001324105405
```

Admin channel is the place where you shadow ban spamers. 
See https://stackoverflow.com/a/41779623/337085 for how to get id for you private channel.
You'll need to invite `@get_id_bot` and type `/my_id@get_id_bot`. You'll see chat id.

## Running

When all services are running, run libretaxi:

```
dep ensure
go build
./libretaxi
```
