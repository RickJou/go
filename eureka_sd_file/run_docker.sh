#!/bin/bash

docker build -t eureka_sd_file .

docker run -itd \
--memory="7m" \
--memory-swap="7m" \
--cpus="0.1" \
-v /data/prometheus:/data/prometheus \
--name eureka_sd_file \
eureka_sd_file main http://192.168.1.34:8848/eureka/apps /data/prometheus/target.json