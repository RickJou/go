FROM golang:1.13

WORKDIR /go/src/app
# 复制项目文件到/go/src/app中
COPY . .

#本项目代码中的所有外部依赖都下载安装
RUN go get -d -v ./...
#编译生产叫做main的执行文件(/go/bin/main)
RUN go install -v ./...

#运行本项目(默认比那以后会在当前目录生产main可执行文件)
CMD ["main"]