package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/PaleBlueYk/randomSSQNumber/pkg/enum"
	"github.com/paleblueyk/logger"
	"path/filepath"
)

// App 配置文件结构
type App struct {
	Server   Server
	WxPusher WxPusher
}

type Server struct {
	Port int `toml:"Port"`
}

type WxPusher struct {
	Url      string `toml:"Url"`
	AppToken string `toml:"AppToken"`
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
	logger.Info(AppConf)
	return true
}
