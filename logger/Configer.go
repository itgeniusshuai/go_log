package logger

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"../common"
	"fmt"
)

// 文件总配置文件
type ConfigInfo struct{
	logOuters []LogOuter
}
var LogConfigPath = "../conf/conf.yml"

func InitConfig(logConfigPath string) error{
	if logConfigPath ==  ""{
		logConfigPath = LogConfigPath
	}
	b,err := ioutil.ReadFile(logConfigPath)
	if err != nil{
		return err
	}
	configInfo := ConfigInfo{}
	yaml.Unmarshal(b,&configInfo)

	// 启动所有的日志器监听binlog
	for {
		select {
		case v := <- common.LogChannel:
			// 遍历所有的outer
			for _,logOuter := range configInfo.logOuters{
				logOuter.Println(v)
			}
		}
	}
	return nil
}



