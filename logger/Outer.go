package logger

import (
	"../common"
	"fmt"
	"strings"
	"regexp"
)

// 日志输出接口
type LogOuter interface{
	Println(logInfo *common.LogInfo)
}

// 控制台输出器
type ConsoleLogOuter struct{
	msgFormat string
}

// 文件输出器
type FileLogOuter struct{
	msgFormat string
	filePath string
	buff string
	buffSize int64

}

// 时间切分文件输出器
type TimeCutFileLogOuter struct{
	FileLogOuter
	TimeFormat string
}

// 容量切分文件输出器
type CapacityCutFileLogOuter struct{
	FileLogOuter
	capacity string
}

func (this *ConsoleLogOuter) Println(logInfo *common.LogInfo){
	msgFormat := this.msgFormat
	msg := parseMsgFormat(msgFormat,logInfo)
	fmt.Println(msg)
}

func (this *TimeCutFileLogOuter) Println(logInfo *common.LogInfo){
	msgFormat := this.msgFormat
	msg := parseMsgFormat(msgFormat,logInfo)
	fmt.Println(msg)
}

func (this *CapacityCutFileLogOuter) Println(logInfo *common.LogInfo){
	msgFormat := this.msgFormat
	msg := parseMsgFormat(msgFormat,logInfo)
	fmt.Println(msg)
}

/**
	%m method
	%t(yyyy-MM-dd HH:mm:ss) time
	%l level
	%n lineNum
	%msg msg
	%fn filename
 */
func parseMsgFormat(msgFormat string,logInfo *common.LogInfo) string{
	msg := msgFormat
	msg = strings.Replace(msgFormat,"%msg",logInfo.Message,-1)
	msg = strings.Replace(msgFormat,"%m",logInfo.MethodName,-1)
	msg = strings.Replace(msgFormat,"%l",common.GetLogLevelStr(logInfo.Level),-1)
	msg = strings.Replace(msgFormat,"%n",string(logInfo.LineNum),-1)
	msg = strings.Replace(msgFormat,"%fn",logInfo.FileName,-1)
	r,_ := regexp.Compile("%t\\(([^)]+)\\)")
	msg = r.ReplaceAllStringFunc(msg,func(str string)string{
		rr := r.FindAllStringSubmatch(str,-1)
		return logInfo.Time.Format(rr[0][1])
	})
	return msg
}