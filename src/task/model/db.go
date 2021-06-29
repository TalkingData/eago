package model

import (
	"eago/task/conf"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type Query map[string]interface{}

// InitDb 初始化数据库
func InitDb() error {
	var err error

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Config.MysqlUser,
		conf.Config.MysqlPassword,
		conf.Config.MysqlAddress,
		conf.Config.MysqlDbName,
	)

	m := mysql.New(mysql.Config{
		DSN:                      dsn,
		DefaultStringSize:        200,
		DisableDatetimePrecision: true,
	})

	db, err = gorm.Open(m, &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}
