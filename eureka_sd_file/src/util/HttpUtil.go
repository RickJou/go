package util

import (
	"io/ioutil"
	"log"
	"net/http"
)

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
