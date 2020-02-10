#!/bin/bash

echo "Wiping the old version..." && \
ssh ro@libretaxi.bot "cd ~/go/src/libretaxi && sudo docker-compose down;  rm -rf ~/go/src/libretaxi/* ; mkdir -p ~/go/src/libretaxi" && \
git reset ./db/init.sql &&
git checkout -- ./db/init.sql &&\
cat ./migration/seed.sql >> ./db/init.sql &&\
rm tmp.zip ; \
zip -r tmp.zip * -x ./libretaxi -x ./migration/* -x ./.git/* && \
echo "Uploading ..." && \
scp tmp.zip ro@libretaxi.bot:/home/ro/go/src/libretaxi/ && \
ssh ro@libretaxi.bot "cd ~/go/src/libretaxi && unzip tmp.zip && ~/go/bin/dep ensure && go build && sudo docker-compose up -d" \
rm tmp.zip \
echo "Done." \
echo "Setting remote port redirects:" \
echo "DB: 127.0.0.1, port 25432" \
echo "RabbitMQ UI: http://127.0.0.1:28080/ (guest/guest)" \
ssh -L 25432:127.0.0.1:15432 -L 28080:127.0.0.1:8080 ro@libretaxi.bot \
