package model

import (
	"gorm.io/gorm"
)

var db *gorm.DB

type Query map[string]interface{}

// SetDb 设置数据库
func SetDb(d *gorm.DB) {
	if d == nil {
		panic("Got a nil orm db object.")
	}
	db = d
}
