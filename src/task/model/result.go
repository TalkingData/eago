package model

import (
	"eago/common/utils"
	"errors"
)

type Result struct {
	Id           uint32            `json:"id" gorm:"type:int(11) unsigned NOT NULL AUTO_INCREMENT;primaryKey"`
	TaskCodename string            `json:"task_codename" gorm:"type:varchar(100);not null;index"`
	Status       int32             `json:"status" gorm:"type:int(11) NOT NULL;index"`
	Caller       string            `json:"caller" gorm:"type:varchar(100) NOT NULL;index"`
	Worker       string            `json:"worker" gorm:"type:varchar(100) NOT NULL;default:'';index"`
	Timeout      *int64            `json:"timeout" gorm:"type:bigint(20) NOT NULL"`
	Arguments    string            `json:"arguments" gorm:"type:json NOT NULL"`
	StartAt      *utils.CustomTime `json:"start_at" gorm:"type:datetime NOT NULL;index"`
	EndAt        *utils.CustomTime `json:"end_at" gorm:"type:datetime"`
}

// TableName 获取数据库表名
func (*Result) TableName() string {
	return "_tmp_results"
}

func (r *Result) GetPartition(partTsFormat string) (string, error) {
	if r.StartAt == nil {
		return "", errors.New("result's start_at field is nil")
	}

	return r.StartAt.Format(partTsFormat), nil
}
