#### 怎么使用?
```
go run Application.go http://localhost:8848/eureka/apps /data/prometheus/target.json
```
第一个参数是eureka地址(目前不支持集群,选取集群中一个节点即可)
第二个参数是供prometheus识别的json文件生成路径

#### prometheus 中如何配置?
```
global:
  scrape_interval: 10s  
  scrape_timeout: 10s
  evaluation_interval: 10s

scrape_configs:
- job_name: prometheus
  static_configs:
  - targets: ['localhost:9090']
  
- job_name: 'eureka_microservers'
  scrape_interval: 15s
  metrics_path: '/actuator/prometheus'
  file_sd_configs:
  - files:
    - 'target.json'
```
注意此处target.json的路径应该和上面生成的路径一致

#### 如何在docker中运行?
执行`sh run_docker.sh`进行构建和执行(需要按需修改脚本中的eureka地址和生产模板文件地址)