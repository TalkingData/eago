package dao

import (
	"gorm.io/gorm"
)

var db *gorm.DB

type Query map[string]interface{}

// Init
func Init(d *gorm.DB) {
	if d == nil {
		panic("Got a nil gorm db object.")
	}
	db = d
}
