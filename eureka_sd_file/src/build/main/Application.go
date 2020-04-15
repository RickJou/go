package main

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)
import _ "net/http/pprof"

//prometheus file_sd_config 结构体
/*
[
  {
    "targets": [
      "localhost:9100"
    ],
    "labels": {
      "job": "node"
    }
  },
  {
    "targets": [
      "localhost:9200"
    ],
    "labels": {
      "job": "node"
    }
  }
]
*/
type Jobs struct {
	Jobs []Job `json:"jobs"`
}
type Job struct {
	Targets []string `json:"targets"`
	Labels  Lable    `json:"labels"`
}
type Lable struct {
	Job string `json:"job"`
}

/*
 * 外部依赖:
 * go get -v -u github.com/tidwall/gjson
 * 测试运行启动参数:
 * go run Application.go http://localhost:8848/eureka/apps /tmp/target.json
 * 编译:
 * CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build application.go
 * CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build application.go
 */
func main() {
	// 开启pprof
	/*go func() {
		ip := "0.0.0.0:6060"
		if err := http.ListenAndServe(ip, nil); err != nil {
			log.Printf("start pprof failed on %s\n", ip)
			os.Exit(1)
		}
	}()*/

	var url = ""
	var targetFile = ""
	log.Printf("args:[%s],\n args length:[%s]", os.Args, len(os.Args))
	if len(os.Args) == 3 {
		url = os.Args[1]
		targetFile = os.Args[2]
		log.Printf("使用自定义配置eureka地址[%s]和配置文件生成地址[%s]", url, targetFile)
	} else {
		url = "http://localhost:8848/eureka/apps"
		targetFile = "/tmp/target.json"
		log.Printf("使用默认配置eureka地址[%s]和配置文件生成地址[%s]", url, targetFile)
	}

	loopLoadConfig(url, targetFile)
	log.Println("哈?结束了!")
}

func loopLoadConfig(url string, targetFile string) {
	for now := range time.Tick(10 * time.Second) {
		log.Printf("定时加载eureka最新实例配置开始..." + now.String())
		eurekaInstanceToPrometheusFileSDConfig(url, targetFile)
		runtime.GC()
	}
}

/*
eureka rest api 获取实例信息后 转换成 prometheus 能够识别的 file_sd_config文件
eureka 返回的格式如下: localhost:8848/eureka/apps
eureka rest api 文档只有netflix[https://github.com/Netflix/eureka/wiki/Eureka-REST-operations]
spring cloud netflix则没有找到文档,只能参照netflix的文档,和spring cloud eureka源码进行推测
gjson语法文档[https://github.com/tidwall/gjson/blob/master/SYNTAX.md]
*/
func eurekaInstanceToPrometheusFileSDConfig(url string, targetFile string) {
	var httpHeader = map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	var resp = SendHttpGet(url, httpHeader)
	//var res = gjson.Get(resp,"applications.application")
	log.Println(resp)

	//解析实例地址
	var allInstance = gjson.Get(resp, "applications.application.#.instance.#.instanceId")
	//log.Println(allInstance)

	//构造file_sd_config需要的文件结构内容
	var lable = Lable{Job: "eureka_microservers"}
	var job = Job{}
	job.Labels = lable

	//读取地址
	for _, addressArr := range allInstance.Array() {
		for _, address := range addressArr.Array() {
			//log.Printf("append instance address %s \n", address)
			job.Targets = append(job.Targets, address.String())
		}
	}
	//绑定最终结构
	var jobs = Jobs{}
	jobs.Jobs = append(jobs.Jobs, job)

	var jsonFileStr, _ = json.Marshal(jobs)
	//最终格式化好后的配置文件内容
	var jsonFilePrettyContent = gjson.Get(string(jsonFileStr), "jobs.@pretty").String()
	log.Printf("\n", jsonFilePrettyContent)

	//写入到指定文件
	if Exists(targetFile) {
		log.Printf("配置文件存在删除文件:%s", targetFile)
		os.RemoveAll(targetFile)
	}

	var err = ioutil.WriteFile(targetFile, []byte(jsonFilePrettyContent), 0666)
	if err != nil {
		log.Printf("文件打开失败=%v\n", err)
	} else {
		log.Printf("写入配置文件成功,chmod(%s,0666)", targetFile)
		//授权目标文件状态
		os.Chmod(targetFile, 0666)
	}

}

//file or dir exists
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//发送http请求
func SendHttpGet(url string, httpHeaders map[string]string) string {
	//set url and headers
	var client = &http.Client{}
	var req, err = http.NewRequest("GET", url, nil)
	for k, v := range httpHeaders {
		req.Header.Set(k, v)
	}
	//send http request
	var resp, err2 = client.Do(req)

	//error handler
	if err2 != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//response handler
	if resp.StatusCode == http.StatusOK {
		var bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panicf("返回成功,读取流失败:%s", err)
		}
		return string(bodyBytes)
	} else if err != nil {
		log.Panicf("返回码[%s],错误:[%s]", resp.Status, err)
	}
	return ""
}
