package main

import (
	"../logger"
	"strconv"
	"time"
)

func main(){
	logger.InitConfig("conf/conf.yml")
	for i := 0; i < 1000000; i++{
		logger.Info("mmmmmmmmmmmmmmmmmmm"+strconv.Itoa(i))
		time.Sleep(time.Millisecond*300)
	}
}






