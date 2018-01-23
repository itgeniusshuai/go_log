package logger

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"../common"
	"time"
)

// 文件总配置文件
type ConfigInfo struct{
	consoleOuters []ConsoleLogOuter
	timeCutOuters []TimeCutFileLogOuter
	capacityCutOuters []CapacityCutFileLogOuter
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

	// 启动所有的日志器监听log
	go func (){
		for {
			select {
			case v := <- common.LogChannel:
				// 遍历所有的outer
				for _,logOuter := range configInfo.consoleOuters{
					logOuter.Println(v)
				}
				// 遍历所有的outer
				for _,logOuter := range configInfo.timeCutOuters{
					logOuter.Println(v,logOuter)
				}
				// 遍历所有的outer
				for _,logOuter := range configInfo.capacityCutOuters{
					logOuter.Println(v,logOuter)
				}
			}
		}
	}()
	// 启动定时写文件程序
	go func (){
		ticker := time.NewTicker(time.Second*30)
		ticks := ticker.C
		for _ = range ticks{
			for _,logOuter := range configInfo.capacityCutOuters{
				logOuter.writeFile()
			}
			for _,logOuter := range configInfo.timeCutOuters{
				logOuter.writeFile()
			}
		}
	}()
	return nil
}



