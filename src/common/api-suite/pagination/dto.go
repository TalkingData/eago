package pagination

import (
	"gorm.io/gorm"
)

// GormParams 分页请求结构
type GormParams struct {
	Db       *gorm.DB
	Page     int
	PageSize int
	OrderBy  []string
}

// Paginator 分页器结构
type Paginator struct {
	Page     int         `json:"page"`
	Pages    int         `json:"pages"`
	PageSize int         `json:"page_size"`
	Total    int64       `json:"total"`
	Data     interface{} `json:"data"`
	Offset   int
}
