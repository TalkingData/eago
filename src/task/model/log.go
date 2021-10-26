package model

import (
	"eago/common/utils"
)

type Log struct {
	Id        int              `json:"id" gorm:"type:int(11) NOT NULL AUTO_INCREMENT;primaryKey"`
	ResultId  int              `json:"result_id" gorm:"type:int(11) NOT NULL;index"`
	Content   string           `json:"content" gorm:"type:text NOT NULL"`
	CreatedAt *utils.LocalTime `json:"created_at" gorm:"type:datetime NOT NULL"`
}

// TableName 获取数据库表名
func (*Log) TableName() string {
	return "_tmp_logs"
}
