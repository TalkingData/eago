package orm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitMysql 初始化Mysql
func InitMysql(address, user, password, dbName string, opts ...DbOption) (d *gorm.DB) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		address,
		dbName,
	)

	m := mysql.New(mysql.Config{
		DSN:                      dsn,
		DefaultStringSize:        200,
		DisableDatetimePrecision: true,
	})

	d, err := gorm.Open(m, &gorm.Config{})
	if err != nil {
		panic(err.Error())
		return nil
	}

	for _, o := range opts {
		o(d)
	}

	db = d
	return db
}

// MysqlDefaultLogMode 设置Mysql默认LogMOde
func MysqlDefaultLogMode(level logger.LogLevel) DbOption {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		d.Logger = logger.Default.LogMode(level)
	}
}

// MysqlMaxIdleConns 设置Mysql最大空闲连接数量
func MysqlMaxIdleConns(count int) DbOption {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		sqlDB, _ := d.DB()
		sqlDB.SetMaxIdleConns(count)
	}
}

// MysqlMaxOpenConns 设置Mysql最大打开连接数量
func MysqlMaxOpenConns(count int) DbOption {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		sqlDB, _ := d.DB()
		sqlDB.SetMaxOpenConns(count)
	}
}
