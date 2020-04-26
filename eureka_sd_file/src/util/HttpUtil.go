package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//发送http请求
func SendHttpGet(url string, httpHeaders map[string]string) string {
	var transport http.RoundTripper = &http.Transport{
		DisableKeepAlives: true,
	}
	//set url and headers
	var client = http.Client{Transport: transport}
	var req, err = http.NewRequest("GET", url, nil)
	for k, v := range httpHeaders {
		req.Header.Set(k, v)
	}

	//send http request
	var resp, err2 = client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	//地址连接不上时,返回为nil
	if resp == nil {
		return ""
	}
	//response handler
	if resp.StatusCode == http.StatusOK {
		var bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Printf("返回成功,读取流失败:%s", err)
		}
		return string(bodyBytes)
	} else if err2 != nil {
		fmt.Printf("返回码[%s],错误:[%s]", resp.Status, err)
	}
	return ""
}
