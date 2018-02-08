package main

import (
	"../logger"
	"strconv"
)

func main(){
	logger.InitConfig("conf/conf.yml")
	for i := 0; i < 1000000; i++{
		logger.Info("mmmmmmmmmmmmmmmmmmm"+strconv.Itoa(i))
	}
}






