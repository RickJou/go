package eureka

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"util"
)

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

//将json结果解析为对象
func ParseJsonToJobs(resp string) Jobs {
	/*
	 eureka rest api 获取实例信息后 转换成 prometheus 能够识别的 file_sd_config文件
	 eureka 返回的格式如下: localhost:8848/eureka/apps
	 eureka rest api 文档只有netflix[https://github.com/Netflix/eureka/wiki/Eureka-REST-operations]
	 spring cloud netflix则没有找到文档,只能参照netflix的文档,和spring cloud eureka源码进行推测
	 gjson语法文档[https://github.com/tidwall/gjson/blob/master/SYNTAX.md]
	*/

	//解析实例地址
	var allInstances = gjson.Get(resp, "applications.application.#.instance")
	//构造file_sd_config需要的文件结构内容
	var lable = Lable{Job: "eureka_microservers"}
	var job = Job{}
	job.Labels = lable

	//读取地址
	for _, addressArr := range allInstances.Array() {
		for _, address := range addressArr.Array() {
			var ip = gjson.Get(address.Raw, "ipAddr").Str
			var port = gjson.Get(address.Raw, "port.$").Raw
			job.Targets = append(job.Targets, ip+":"+port)
		}
	}
	//绑定最终结构
	var jobs = Jobs{}
	jobs.Jobs = append(jobs.Jobs, job)
	return jobs
}

func InstanceToPrometheusFileSDConfig(url string, targetFile string) {
	var httpHeader = map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	var resp = util.SendHttpGet(url, httpHeader)
	if resp == "" {
		//连接失败
		fmt.Println("接口超时或无数据返回,请检查eureka地址是否正确.")
		return
	}

	//解析json
	var jobs = ParseJsonToJobs(resp)

	//最终格式化好后的配置文件内容
	var jsonFileStr, _ = json.Marshal(jobs)
	var jsonFilePrettyContent = gjson.Get(string(jsonFileStr), "jobs.@pretty").String()
	jsonFilePrettyContent = jsonFilePrettyContent + ""
	//删除原文件,写入到新文件
	util.RemoveFile(targetFile)
	util.CreateNewFile(targetFile, []byte(jsonFilePrettyContent))
}
