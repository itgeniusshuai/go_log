package common

import "time"

type LogLevel int
// 日志输出级别
const(
	Debug = iota
	Warning
	Info
	Error
	Fatal
)

var LogChannel = make(chan *LogInfo,1024)

type LogInfo struct{
	// 日志级别
	Level LogLevel `json:"level"`
	// 文件名
	FileName string `json:"fileName"`
	// 日志对应方法名
	MethodName string `json:"methodName"`
	// 文件的行号
	LineNum int `json:"lineNum"`
	// 时间戳
	Time time.Time `json:"time"`
	// 日志内容
	Message string `json:"message"`
}

// 根据日志级别获取对应的级别字符串
func GetLogLevelStr(level LogLevel) string{
	var msg string
	switch level {
	case Debug:
		msg = "Debug"
	case Warning:
		msg = "Warning"
	case Info:
		msg = "Info"
	case Error:
		msg = "Error"
	case Fatal:
		msg = "Fatal"
	}
	return msg
}
