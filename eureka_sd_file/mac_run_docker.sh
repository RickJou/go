#!/usr/bin/env bash

docker build -t eureka_sd_file .

#在本地启动时,该程序向系统申请的最小内存是68m,所以容器要有适当的内存可用
docker run -itd \
--rm  \
--memory="100m" \
--memory-swap="105m" \
--cpus="0.1" \
-v /Users/alanzhou/Desktop/dev_tools/docker/prometheus:/data/prometheus \
-p 6060:6060 \
--name eureka_sd_file \
eureka_sd_file main http://192.168.1.34:8848/eureka/apps /tmp/prometheus/target.json