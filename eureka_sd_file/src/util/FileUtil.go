package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//create new file
func CreateNewFile(path string, fileContent []byte) {
	var parentFold = filepath.Dir(path)
	os.MkdirAll(parentFold,0777)
	os.Create(path)
	var err = ioutil.WriteFile(path, fileContent, 0777)
	if err != nil {
		fmt.Printf("文件打开失败=%v\n", err)
	} else {
		fmt.Printf("写入配置文件成功,chmod(%s,0777)", path)
		//授权目标文件状态
		os.Chmod(path, 0777)
	}
}

//remove file if exists
func RemoveFile(path string) {
	if Exists(path) {
		fmt.Printf("配置文件存在删除文件:%s", path)
		os.RemoveAll(path)
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
