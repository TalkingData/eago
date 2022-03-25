package model

import (
	"eago/common/utils"
	"eago/task/conf"
	"errors"
)

type Result struct {
	Id           int              `json:"id" gorm:"type:int(11) NOT NULL AUTO_INCREMENT;primaryKey"`
	TaskCodename string           `json:"task_codename" gorm:"type:varchar(100);not null;index"`
	Status       int              `json:"status" gorm:"type:int(11) NOT NULL;index"`
	Caller       string           `json:"caller" gorm:"type:varchar(100) NOT NULL;index"`
	Worker       string           `json:"worker" gorm:"type:varchar(100) NOT NULL;default:'';index"`
	Timeout      *int             `json:"timeout" gorm:"type:int(11) NOT NULL"`
	Arguments    string           `json:"arguments" gorm:"type:json NOT NULL"`
	StartAt      *utils.LocalTime `json:"start_at" gorm:"type:datetime NOT NULL;index"`
	EndAt        *utils.LocalTime `json:"end_at" gorm:"type:datetime"`
}

// TableName 获取数据库表名
func (*Result) TableName() string {
	return "_tmp_results"
}

func (r *Result) GetPartition() (string, error) {
	if r.StartAt == nil {
		return "", errors.New("result's start_at field is nil")
	}

	return r.StartAt.Format(conf.TASK_PARTITION_TIMESTAMP_FORMAT), nil
}
