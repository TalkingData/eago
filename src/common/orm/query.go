package orm

import (
	"gorm.io/gorm"
)

type Query map[string]interface{}

// NewQueryByMapStrStr 从map[string]string生成Query
func NewQueryByMapStrStr(in map[string]string) (q Query) {
	for k, v := range in {
		q[k] = v
	}
	return
}

func (q Query) Where(db *gorm.DB) *gorm.DB {
	for k, v := range q {
		db = db.Where(k, v)
	}

	return db
}
