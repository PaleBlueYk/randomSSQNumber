package db

import (
	"fmt"
	"github.com/PaleBlueYk/randomSSQNumber/config"
	"github.com/PaleBlueYk/randomSSQNumber/model"
	"github.com/paleblueyk/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Mysql *gorm.DB

// ConnectMysql 连接mysql
func ConnectMysql() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConf.Mysql.User,
		config.AppConf.Mysql.Password,
		config.AppConf.Mysql.Host,
		config.AppConf.Mysql.Port,
		config.AppConf.Mysql.DB)
	Mysql, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		logger.Error(err)
		return
	}

	Migrate()
	return nil
}

func Migrate() {
	if err := Mysql.AutoMigrate(
		&model.NumSaveData{},
	); err != nil {
		logger.Error(err)
	}
}
