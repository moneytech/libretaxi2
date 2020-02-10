#!/bin/bash

git reset ./db/init.sql &&
git checkout -- ./db/init.sql &&\
cat ./migration/seed.sql >> ./db/init.sql &&\
rm tmp.zip ; \
zip -r tmp.zip * -x ./libretaxi -x ./migration/* -x ./.git/* && \
echo "Uploading ..." && \
scp tmp.zip ro@libretaxi.bot:/home/ro/go/src/libretaxi/ && \
echo "Done"

