package dto

import "eago/common/utils"

type FormWithoutBody struct {
	Id          uint32            `json:"id"`
	Name        string            `json:"name"`
	Disabled    *bool             `json:"disabled"`
	Description *string           `json:"description"`
	CreatedAt   *utils.CustomTime `json:"created_at"`
	CreatedBy   string            `json:"created_by"`
	UpdatedAt   *utils.CustomTime `json:"updated_at"`
	UpdatedBy   *string           `json:"updated_by"`
}
