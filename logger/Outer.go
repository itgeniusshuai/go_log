package logger

import (
	"../common"
	"fmt"
	"strings"
	"regexp"
	"sync"
	"time"
	"os"
	"sync/atomic"
	"../helpers"
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
	fileNamePrefix string
	buff string
	buffSize int
	rwLock sync.RWMutex
	buffLock sync.RWMutex
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
	lastFileId int64
}

func (this *FileLogOuter) wirteFile(){
	this.rwLock.Lock()
	defer this.rwLock.Unlock()
	_,err := os.Open(this.filePath)
	if err != nil{
		os.MkdirAll(this.filePath,os.ModePerm)
	}
	fileName := this.filePath + "/" + this.getFileName() + ".log"
	file, err := os.Open(fileName)
	if err != nil{
		file,_ = os.Create(fileName)
	}
	file.WriteString(this.buff)
}

func (this *FileLogOuter) getFileName() string{
	switch v := this.(type) {
	case TimeCutFileLogOuter:
		return time.Now().Format(v.TimeFormat)
	case CapacityCutFileLogOuter:
		atomic.AddInt64(&v.lastFileId,1)
		return v.fileNamePrefix + "_" + helpers.GetString(v.lastFileId)
	}
	return nil
}

func (this *ConsoleLogOuter) Println(logInfo *common.LogInfo){
	msgFormat := this.msgFormat
	msg := parseMsgFormat(msgFormat,logInfo)
	fmt.Println(msg)
}

func (this *FileLogOuter) Println(logInfo *common.LogInfo){
	msgFormat := this.msgFormat
	msg := parseMsgFormat(msgFormat,logInfo)
	this.buffLock.Lock()
	this.buff = this.buff + msg
	defer this.buffLock.Unlock()
	if this.buffSize <= 0 || len(this.buff) >= this.buffSize {
		this.wirteFile(logInfo)
	}
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