package model

import (
	"eago/common/utils"
)

// Product 产品线-数据库模型
type Product struct {
	Id uint32 `json:"id"`

	Name        string  `json:"name"`
	Alias       string  `json:"alias"`
	Disabled    *bool   `json:"disabled"`
	Description *string `json:"description"`

	CreatedAt *utils.CustomTime `json:"created_at"`
	UpdatedAt *utils.CustomTime `json:"updated_at"`
}

// UserProduct 用户产品线关联-数据库模型
type UserProduct struct {
	Id uint32 `json:"id"`

	UserId    uint32            `json:"user_id"`
	ProductId uint32            `json:"product_id"`
	IsOwner   *bool             `json:"is_owner"`
	JoinedAt  *utils.CustomTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}
