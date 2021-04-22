package database

import (
	"eago-common/api-suite/pagination"
	"eago-common/log"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var ProductModel productModel

type productModel struct{}

type UserProduct struct {
	Id        int    `json:"id" swaggerignore:"true"`
	UserId    int    `json:"user_id" binding:"required"`
	ProductId int    `json:"product_id" swaggerignore:"true"`
	IsOwner   *bool  `json:"is_owner" binding:"required"`
	JoinedAt  MyTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime" swaggerignore:"true"`
}

type Product struct {
	Id          int     `json:"id" swaggerignore:"true"`
	Name        string  `json:"name" binding:"required"`
	Alias       string  `json:"alias" binding:"required"`
	Disabled    *bool   `json:"disabled" binding:"required"`
	Description string  `json:"description" binding:"required"`
	CreatedAt   MyTime  `json:"created_at" swaggerignore:"true"`
	UpdatedAt   *MyTime `json:"updated_at" swaggerignore:"true"`
}

// New 新建产品线
func (pm *productModel) New(name, alias, description string, disabled bool) *Product {
	var prod = Product{
		Name:        name,
		Alias:       alias,
		Disabled:    &disabled,
		Description: description,
	}

	if res := db.Create(&prod); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":        name,
			"alias":       alias,
			"disabled":    disabled,
			"description": description,
			"error":       res.Error.Error(),
		}, "Error in productModel.New.")
		return nil
	}

	return &prod
}

// Remove 删除产品线
func (pm *productModel) Remove(productId int) bool {
	res := db.Delete(Product{}, "id=?", productId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    productId,
			"error": res.Error.Error(),
		}, "Error in productModel.Remove.")
		return false
	}

	return true
}

// Set 更新产品线
func (pm *productModel) Set(id int, name, alias, description string, disabled bool) (*Product, bool) {
	var p = Product{}

	res := db.Model(&Product{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"alias":       alias,
			"disabled":    disabled,
			"description": description,
		}).
		First(&p)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":          id,
			"name":        name,
			"alias":       alias,
			"disabled":    disabled,
			"description": description,
			"error":       res.Error.Error(),
		}, "Error in productModel.Set.")
		return nil, false
	}

	return &p, true
}

// List 查询产品线
func (pm *productModel) List(query *Query) (*[]Product, bool) {
	var d = db
	ps := make([]Product, 0)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ps); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in productModel.List.")
			return &ps, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in productModel.List.")
		return nil, false
	}

	return &ps, true
}

// PagedList 查询产品线-分页
func (pm *productModel) PagedList(query *Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&Product{})
	ps := make([]Product, 0)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &ps)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": err.Error(),
		}, "Error in productModel.PagedList.")
		return nil, false
	}

	return pg, true
}

// AddUser 关联表操作::添加用户至产品线
func (pm *productModel) AddUser(userId, productId int, isOwner bool) bool {
	var up = UserProduct{
		UserId:    userId,
		ProductId: productId,
		IsOwner:   &isOwner,
	}

	if res := db.Create(&up); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":    userId,
			"product_id": productId,
			"is_owner":   isOwner,
			"joined_at":  up.JoinedAt,
			"error":      res.Error.Error(),
		}, "Error in productModel.AddUser.")
		return false
	}

	return true
}

// RemoveUser 关联表操作::移除产品线中用户
func (pm *productModel) RemoveUser(userId, productId int) bool {
	res := db.Delete(UserProduct{}, "user_id=? AND product_id=?", userId, productId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":    userId,
			"product_id": productId,
			"error":      res.Error.Error(),
		}, "Error in productModel.RemoveUser.")
		return false
	}

	return true
}

// 关联表操作::设置用户是否是产品线Owner
func (pm *productModel) SetUserIsOwner(userId, productId int, isOwner bool) bool {
	res := db.Model(&UserProduct{}).
		Where("user_id=? AND product_id=?", userId, productId).
		Update("IsOwner", isOwner)
	if res.Error != nil {
		log.WarnWithFields(log.Fields{
			"user_id":    userId,
			"product_id": productId,
			"is_owner":   isOwner,
			"error":      res.Error.Error(),
		}, "Error in productModel.SetUserIsOwner.")
		return false
	}

	return true
}

// ListUsers 关联表操作::列出产品线中所有用户
func (pm *productModel) ListUsers(productId int, query *Query) (*[]memberUser, bool) {
	var d = db.Model(&User{})
	mus := make([]memberUser, 0)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	res := d.Select("users.id AS id, users.username AS username, up.is_owner AS is_owner, up.joined_at AS joined_at").
		Joins("LEFT JOIN user_products AS up ON users.id = up.user_id").
		Where("product_id=?", productId).
		Find(&mus)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in productModel.ListUsers.")
			return &mus, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in productModel.ListUsers.")
		return nil, false
	}

	return &mus, true
}
