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
// 获取文件名称，没有创建
// 写入缓冲区
// 缓冲区溢出写入文件
// 清空缓冲区
// 如果是容量切分，更新当前容量
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
	Capacity int64
	lastFileId int64
	currentCapacity int64
}

func (this *TimeCutFileLogOuter) writeFile(){
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
	this.buff = ""
}

func (this *CapacityCutFileLogOuter) writeFile(){
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
	this.buff = ""
	this.currentCapacity = int64(len(this.buff)) + this.currentCapacity
}

func (this *TimeCutFileLogOuter) getFileName() string{
	return time.Now().Format(this.TimeFormat)
}

func (this *CapacityCutFileLogOuter) getFileName() string{
		if this.currentCapacity >= this.Capacity{
			atomic.AddInt64(&this.lastFileId,1)
			this.currentCapacity = int64(0)
		}
		return this.fileNamePrefix + "_" + helpers.GetString(this.lastFileId)
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
		this.writeFile()
	}
}

func (this *FileLogOuter) writeFile(){
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