package log

import (
	"io"
	"log"
	"os"
)

var Logger logger

type logger struct {
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
}

func init() {
	//日志输出文件
	file, err := os.OpenFile("./log/log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Faild to open error logger file:", err)
	}
	//自定义日志格式
	Logger.Info = log.New(io.MultiWriter(file, os.Stderr), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Error = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Warn = log.New(io.MultiWriter(file, os.Stderr), "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)

}
