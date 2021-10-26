package orm

import (
	"gorm.io/gorm"
)

var db *gorm.DB

type DbOption func(d *gorm.DB)

// Close 关闭数据库
func Close() {
	if db == nil {
		return
	}

	d, _ := db.DB()
	_ = d.Close()
}
