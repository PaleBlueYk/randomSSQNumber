package main

import (
	"fmt"
	"github.com/PaleBlueYk/randomSSQNumber/config"
	"github.com/PaleBlueYk/randomSSQNumber/db"
	"github.com/PaleBlueYk/randomSSQNumber/pkg/enum"
	v1 "github.com/PaleBlueYk/randomSSQNumber/routers/v1"
	"github.com/gin-gonic/gin"
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
	if err := db.ConnectMysql(); err != nil {
		logger.Error("mysql连接失败: ", err)
		return
	}
	db.ConnectRDB()

	engine := gin.Default()
	v1.Routers(engine)
	_ = engine.Run(fmt.Sprintf(":%d", config.AppConf.Server.Port))
}
