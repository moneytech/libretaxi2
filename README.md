## Install

Init settings for `./libretaxi.yml`:

```
telegram_token: YOUR_TOKEN
db_conn_str: postgres://libretaxi:libretaxi@localhost/libretaxi
```

Then run:

```
dep ensure
go build
./libretaxi
```

