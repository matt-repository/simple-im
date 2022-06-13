package main

import (
	"flag"
	"fmt"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	"goland-IM-System/server/comm/frame"
	"goland-IM-System/server/log"
	"net/http"
	"time"

	"github.com/panjf2000/gnet"
)

//阿里云 服务端填写的是内网ip
//客户端 填写的是外网ip
var serverIp = "127.0.0.1"
var serverPort = 6999

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "port", 6999, "设置服务器端口是（默认是6999）")
}

func main() {
	//main.exe -ip=127.0.0.1 -port=8888
	flag.Parse()
	log.Logger.Info.Println("GO-IM-System  启动")
	fmt.Println(serverIp, serverPort)

	server.Start()
}
