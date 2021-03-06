package logger

import (
	"../common"
	"fmt"
	"strings"
	"regexp"
	"sync"
	"time"
	"os"
	"../helpers"
	"encoding/json"
	"pkg/errors"
)

// 日志输出接口
type LogOuterInterface interface{
	Println(logInfo *common.LogInfo)
}

// 文件输出器扩展方法接口
type LogFileOuterInterface interface {
	GetFileName() (string)
}

// 消息格式，当type为json的时候Format无效
type MsgFormat struct{
	Type string `yaml:"type"`
	Format string `yaml:"format"`
}

// 输出器基类
type LogOuter struct {
	MsgFormat MsgFormat `yaml:"msgFormat"`
	LessLevel common.LogLevel `yaml:"lessLevel"`
}
// 控制台输出器
type ConsoleLogOuter struct{
	LogOuter `yaml:"logOuter"`
}

// 文件输出器
// 获取文件名称，没有创建
// 写入缓冲区
// 缓冲区溢出写入文件
// 清空缓冲区
// 如果是容量切分，更新当前容量
type FileLogOuter struct{
	LogOuter `yaml:"logOuter"`
	FilePath string `yaml:"filePath"`
	FileNamePrefix string `yaml:"fileNamePrefix"`
	Buff string `yaml:"buff"`
	BuffSize int `yaml:"buffSize"`
	RwLock sync.RWMutex `yaml:"rwLock"`
	BuffLock sync.RWMutex `yaml:"buffLock"`
	LessLevel common.LogLevel `yaml:"lessLevel"`

	LogFileOuterInterface
}

// 固定文件不切分文件输出器
type FixedFileLogOuter  struct {
	FileLogOuter `yaml:"fileLogOuter"`
}

// 时间切分文件输出器
type TimeCutFileLogOuter struct{
	FileLogOuter `yaml:"fileLogOuter"`
	TimeFormat string `yaml:"timeFormat"`
}

// 容量切分文件输出器
type CapacityCutFileLogOuter struct{
	FileLogOuter `yaml:"fileLogOuter"`
	Capacity int64 `yaml:"capacity"`
	lastFileId int64
}

// 写入文件
func (this *FileLogOuter) WriteFile(){
	this.RwLock.Lock()
	defer this.RwLock.Unlock()
	_,err := os.Open(this.FilePath)
	if err != nil{
		os.MkdirAll(this.FilePath,os.ModePerm)
	}
	fileName := this.FilePath + string(os.PathSeparator) + this.GetFileName() + ".log"
	file, err := os.OpenFile(fileName,os.O_APPEND,os.ModePerm)
	defer file.Close()
	if err != nil{
		file,_ = os.Create(fileName)
	}
	file.WriteString(this.Buff)
	this.Buff = ""
}

// 获取要写入文件名
func (this *FixedFileLogOuter) GetFileName() string{
	return this.FileNamePrefix
}

// 获取要写入文件名
func (this *TimeCutFileLogOuter) GetFileName() string{
	return this.FileNamePrefix+"_"+time.Now().Format(this.TimeFormat)
}

// 获取要写入文件名
func (this *CapacityCutFileLogOuter) GetFileName() string{
		fileName := this.FilePath + string(os.PathSeparator) + this.FileNamePrefix + "_" + helpers.GetString(this.lastFileId) + ".log"
		fileSize,_ := getFileSize(fileName)
		for fileSize > this.Capacity{
			this.lastFileId ++
			fileName = this.FilePath + string(os.PathSeparator) + this.FileNamePrefix + "_" + helpers.GetString(this.lastFileId) + ".log"
			fileSize,_ = getFileSize(fileName)
		}
		return this.FileNamePrefix + "_" + helpers.GetString(this.lastFileId)
}

// 控制台输出器顶级接口实现
func (this *ConsoleLogOuter) Println(logInfo *common.LogInfo){
	if this.LessLevel > logInfo.Level{
		return
	}
	msgFormat := this.MsgFormat
	msg := parseMsgFormat(msgFormat,logInfo)
	fmt.Println(msg)
}

// 输出器顶级实现
func (this *FileLogOuter) Println(logInfo *common.LogInfo){
	if this.LessLevel > logInfo.Level{
		return
	}
	msgFormat := this.MsgFormat
	msg := parseMsgFormat(msgFormat,logInfo)
	this.BuffLock.Lock()
	this.Buff = this.Buff + msg
	defer this.BuffLock.Unlock()
	if this.BuffSize <= 0 || len(this.Buff) >= this.BuffSize {
		this.WriteFile()
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
func parseMsgFormat(msgFormat MsgFormat,logInfo *common.LogInfo) string{
	msgType := msgFormat.Type
	format := msgFormat.Format
	var msg string
	switch msgType {
	case "json":
		b,_ := json.Marshal(logInfo)
		msg = string(b)
	case "string":
		msg = strings.Replace(format,"%msg",logInfo.Message,-1)
		msg = strings.Replace(msg,"%m",logInfo.MethodName,-1)
		msg = strings.Replace(msg,"%l",common.GetLogLevelStr(logInfo.Level),-1)
		msg = strings.Replace(msg,"%n",string(logInfo.LineNum),-1)
		msg = strings.Replace(msg,"%fn",logInfo.FileName,-1)
		r,_ := regexp.Compile("%t\\(([^)]+)\\)")
		msg = r.ReplaceAllStringFunc(msg,func(str string)string{
			rr := r.FindAllStringSubmatch(str,-1)
			return logInfo.Time.Format(rr[0][1])
		})
	default:
		b,_ := json.Marshal(logInfo)
		msg = string(b)
	}
	return msg
}

// 获取文件大小
func getFileSize(fileName string) (int64,error){
	file,err := os.Open(fileName)
	if err != nil{
		fmt.Println("open file faild")
		return -1,errors.New("can't open file: "+fileName)
	}
	fileInfo,_ := file.Stat()
	fileSize := fileInfo.Size()
	return fileSize,nil
}