package main

import "../logger"

func main(){
	logger.InitConfig("conf/conf.yml")
	for i := 0; i < 10000; i++{
		logger.Info("mmmmmmmmmmmmmmmmmmm"+string(i))
	}
}






