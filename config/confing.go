package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/paleblueyk/logger"
	"path/filepath"
	"github.com/PaleBlueYk/randomSSQNumber/pkg/enum"
)

// App 配置文件结构
type App struct {
	WxPusher
}

type WxPusher struct {
	AppToken string
}

var AppConf App

// InitConfig 初始化配置文件
func InitConfig(env enum.ENVType) bool {
	var configFile string
	switch env {
	case enum.Prod:
		configFile = "config.toml"
	case enum.Dev:
		configFile = "config_dev.toml"
	case enum.Test:
		configFile = "config_test.toml"
	}
	if _, err := toml.DecodeFile(filepath.FromSlash(fmt.Sprintf("config/%s", configFile)), &AppConf); err != nil {
		logger.Error(err)
		return false
	}
	return true
}
