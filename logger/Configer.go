package logger

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

// 文件总配置文件
type ConfigInfo struct{
	logOuters []LogOuter
}
var LogConfigPath = "../conf/conf.yml"

func initConfig(logConfigPath string) error{
	b,err := ioutil.ReadFile(logConfigPath)
	if err != nil{
		return err
	}
	configInfo := ConfigInfo{}
	yaml.Unmarshal(b,&configInfo)

}



