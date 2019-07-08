package client_system

import (
	"os"
	"log"
)

const (
	LOG_INFO    string = "[info] "
	LOG_WARNING string = "[warning] "
	LOG_ERROR   string = "[error] "
	LOG_FATAL   string = "[fatal] "
)

func LogWrite(logType string, str string) {
	logFile, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("日志文件 run.log 初始化失败")
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags|log.Llongfile)
	logger.SetFlags(log.LstdFlags)
	logger.SetPrefix(logType)
	logger.Println(str)
}

func LogWriteData(str string) {
	logFile, err := os.OpenFile("./runtime/log/msg.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("消息日志文件 run.log 初始化失败")
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags|log.Llongfile)
	logger.SetFlags(log.LstdFlags)
	logger.SetPrefix("")
	logger.Println(str)
}

