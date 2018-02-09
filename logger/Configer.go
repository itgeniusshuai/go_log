package logger

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"../common"
	"time"
)

// 文件总配置文件
type ConfigInfo struct{
	ConsoleOuters []*ConsoleLogOuter `yaml:"consoleOuters"`
	TimeCutOuters []*TimeCutFileLogOuter `yaml:"timeCutOuters"`
	CapacityCutOuters []*CapacityCutFileLogOuter `yaml:"capacityCutOuters"`
	FixedFileOuters []*FixedFileLogOuter `yaml:"fixedFileOuter"`
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

	// 遍历所有的outer
	for _,logOuter := range configInfo.TimeCutOuters{
		logOuter.LogFileOuterInterface = logOuter
	}
	// 遍历所有的outer
	for _,logOuter := range configInfo.CapacityCutOuters{
		logOuter.LogFileOuterInterface = logOuter
	}

	for _,logOuter := range configInfo.FixedFileOuters{
		logOuter.LogFileOuterInterface = logOuter
	}
	// 启动所有的日志器监听log
	go func (){
		for {
			select {
			case v := <- common.LogChannel:
				// 遍历所有的outer
				for _,logOuter := range configInfo.ConsoleOuters{
					logOuter.Println(v)
				}
				// 遍历所有的outer
				for _,logOuter := range configInfo.TimeCutOuters{
					logOuter.Println(v)
				}
				// 遍历所有的outer
				for _,logOuter := range configInfo.CapacityCutOuters{
					logOuter.Println(v)
				}

				for _,logOuter := range configInfo.FixedFileOuters{
					logOuter.Println(v)
				}

			}
		}
	}()
	// 启动定时写文件程序，缓存区不满，每30s清空缓存
	go func (){
		ticker := time.NewTicker(time.Second*30)
		ticks := ticker.C
		for _ = range ticks{
			for _,logOuter := range configInfo.CapacityCutOuters{
				logOuter.WriteFile()
			}
			for _,logOuter := range configInfo.TimeCutOuters{
				logOuter.WriteFile()
			}
		}
	}()
	return nil
}



