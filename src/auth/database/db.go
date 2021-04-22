package database

import (
	"database/sql/driver"
	"eago-auth/conf"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
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

type MyTime struct {
	time.Time
}

func (t *MyTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(conf.TIMESTAMP_FORMAT))
	return []byte(formatted), nil
}

func (t *MyTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *MyTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = MyTime{Time: value}
		return nil
	}
	return fmt.Errorf("Can not convert %v to timestamp", v)
}
