package model

import (
	"eago/common/utils"
)

type Log struct {
	Id        uint64            `json:"id" gorm:"type:bigint(20) unsigned NOT NULL AUTO_INCREMENT;primaryKey"`
	ResultId  uint32            `json:"result_id" gorm:"type:int(11) unsigned NOT NULL;index"`
	Content   string            `json:"content" gorm:"type:text NOT NULL"`
	CreatedAt *utils.CustomTime `json:"created_at" gorm:"type:datetime NOT NULL"`
}

// TableName 获取数据库表名
func (*Log) TableName() string {
	return "_tmp_logs"
}
