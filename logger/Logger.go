package logger

import (
	"time"
	"runtime"
	"../common"
)

// 日志类
func Warning(msg string){
	log(msg,common.Warning)
}

func Info(msg string){
	log(msg,common.Info)
}

func Debug(msg string){
	log(msg,common.Debug)
}

func Error(msg string){
	log(msg,common.Error)
}

func Fatal(msg string){
	log(msg,common.Fatal)
}

func log(msg string, level common.LogLevel){
	pc,file,line,_ := runtime.Caller(1)
	methodName := runtime.FuncForPC(pc).Name()
	logInfo := common.LogInfo{
		Level:level,
		FileName: file,
		LineNum:line,
		Time:time.Now(),
		Message:msg,
		MethodName:methodName,
	}
	sendToLogChan(&logInfo)
}

func sendToLogChan(info *common.LogInfo){
	common.LogChannel <- info
}

