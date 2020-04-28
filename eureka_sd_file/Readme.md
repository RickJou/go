### docker启动脚本

#### 启动不同环境脚本
```
sh deploy/run_{evn}_eureka_sd.sh 
```

#### 启动prometheus
```
sh prometheus/run_prometeus.sh
```


#### 启动grafana
```
sh prometheus/run_grafana.sh
```

### 代码调试

#### 调试代码
```
go run Application.go http://localhost:8848/eureka/apps /data/prometheus/target.json
```
第一个参数是eureka地址(目前不支持集群,选取集群中一个节点即可)
第二个参数是供prometheus识别的json文件生成路径
