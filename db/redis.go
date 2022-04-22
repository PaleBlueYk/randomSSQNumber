package db

import (
	"github.com/PaleBlueYk/randomSSQNumber/config"
	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func ConnectRDB() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.AppConf.Redis.Addr,
		Password: config.AppConf.Redis.Password,
		DB:       config.AppConf.Redis.DB,
	})
}
