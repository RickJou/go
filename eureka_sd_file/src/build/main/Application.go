package main

import (
	"eureka"
	"fmt"
	"os"
	"runtime"
	"time"
)
import _ "net/http/pprof"

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
			fmt.Printf("start pprof failed on %s\n", ip)
			os.Exit(1)
		}
	}()*/

	var url = ""
	var targetFile = ""
	fmt.Printf("args:[%s],\n args length:[%s]", os.Args, len(os.Args))
	if len(os.Args) == 3 {
		url = os.Args[1]
		targetFile = os.Args[2]
		fmt.Printf("使用自定义配置eureka地址[%s]和配置文件生成地址[%s]", url, targetFile)
	} else {
		url = "http://localhost:8761/eureka/apps"
		targetFile = "/tmp/target.json"
		fmt.Printf("使用默认配置eureka地址[%s]和配置文件生成地址[%s]", url, targetFile)
	}

	loopLoadConfig(url, targetFile)
	fmt.Println("哈?结束了!")
}

func loopLoadConfig(url string, targetFile string) {
	/*for range time.Tick(10 * time.Second) {
		fmt.Printf("定时加载eureka最新实例配置开始...")
		eureka.InstanceToPrometheusFileSDConfig(url, targetFile)
	}*/
	for {
		time.Sleep(time.Duration(10) * time.Second)
		fmt.Printf("定时加载eureka最新实例配置开始...")
		eureka.InstanceToPrometheusFileSDConfig(url, targetFile)
		PrintMemUsage()
	}

}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("\nAlloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tMallocs = %v MiB", bToMb(m.Mallocs))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
