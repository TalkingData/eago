package model

import (
	"eago/common/utils"
)

type Product struct {
	Id          int              `json:"id"`
	Name        string           `json:"name"`
	Alias       string           `json:"alias"`
	Disabled    *bool            `json:"disabled"`
	Description *string          `json:"description"`
	CreatedAt   *utils.LocalTime `json:"created_at"`
	UpdatedAt   *utils.LocalTime `json:"updated_at"`
}

type UserProduct struct {
	Id        int              `json:"id"`
	UserId    int              `json:"user_id"`
	ProductId int              `json:"product_id"`
	IsOwner   *bool            `json:"is_owner"`
	JoinedAt  *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}
