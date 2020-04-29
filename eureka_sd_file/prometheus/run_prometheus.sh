#!/usr/bin/env bash

mkdir -p /data/prometheus
mkdir -p /data/prometheus/tsdb
chmod -R 777 /data/prometheus

cp -f ./prometheus.yml /data/prometheus/

docker run \
--name="prometheus" \
--memory="1024m" \
--memory-swap="1124m" \
--cpus="1" \
-p 9090:9090 \
-v /data/prometheus:/etc/prometheus \
-d \
prom/prometheus \
--web.enable-lifecycle \
--config.file=/etc/prometheus/prometheus.yml \
--storage.tsdb.path=data\ \
--storage.tsdb.retention.time=366d \
--web.console.libraries=/usr/share/prometheus/console_libraries \
--web.console.templates=/usr/share/prometheus/consoles