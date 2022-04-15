package main

import (
	"github.com/PaleBlueYk/randomSSQNumber/config"
	"github.com/PaleBlueYk/randomSSQNumber/pkg/enum"
	"github.com/PaleBlueYk/randomSSQNumber/service"
	"github.com/paleblueyk/logger"
	"os"
)

func main() {
	env := enum.Prod
	logger.Info(os.Args)
	if len(os.Args) > 1 {
		env = enum.ENVType(os.Args[1])
	}
	if !config.InitConfig(env) {
		logger.Info("读取配置文件失败")
		return
	}
	for i := 0; i < 5; i++ {
		logger.Info(service.GenNumber())
	}
}
