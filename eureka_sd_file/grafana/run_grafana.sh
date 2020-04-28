#!/usr/bin/env bash

mkdir -p /data/grafana && chmod -R +777 /data/grafana && chmod g+s /data/grafana

docker run -d \
--name=grafana \
--memory="1024m" \
--memory-swap="1124m" \
-p 3000:3000 \
--volume "/data/grafana/var/lib/grafana:/var/lib/grafana" \
--user $(id -u):$(id -g) \
grafana/grafana