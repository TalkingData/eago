package dao

import (
	"eago/auth/model"
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewProduct 新建产品线
func NewProduct(name, alias, description string, disabled *bool) (*model.Product, error) {
	prod := model.Product{
		Name:        name,
		Alias:       alias,
		Disabled:    disabled,
		Description: &description,
	}

	if res := db.Create(&prod); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":        name,
			"alias":       alias,
			"disabled":    disabled,
			"description": description,
			"error":       res.Error,
		}, "An error occurred while db.Create.")
		return nil, res.Error
	}

	return &prod, nil
}

// RemoveProduct 删除产品线
func RemoveProduct(productId int) error {
	res := db.Delete(model.Product{}, "id=?", productId)
	if res.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}

	return res.Error
}

// SetProduct 更新产品线
func SetProduct(id int, name, alias, description string, disabled bool) (*model.Product, error) {
	p := model.Product{}

	res := db.Model(&model.Product{}).
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
			"error":       res.Error,
		}, "An error occurred while db.Model.Where.Updates.First.")
		return nil, res.Error
	}

	return &p, nil
}

// GetProduct 查询单个产品线
func GetProduct(query Query) (*model.Product, bool) {
	var (
		p = model.Product{}
		d = db
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&p); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.First.")
		return nil, false
	}

	return &p, true
}

//  GetProductCount 查询产品线数量
func GetProductCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Product{})

	for k, v := range query {
		d = d.Where(k, v)
	}

	if res := d.Count(&count); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return count, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Count.")
		return count, false
	}
	return count, true
}

// ListProducts 查询产品线
func ListProducts(query Query) (*[]model.Product, bool) {
	var d = db
	ps := make([]model.Product, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ps); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return &ps, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return &ps, true
}

// PagedListProducts 查询产品线-分页
func PagedListProducts(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&model.Product{})
	ps := make([]model.Product, pageSize)

	for k, v := range query {
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
			"error": err,
		}, "An error occurred while pagination.GormPaging.")
		return nil, false
	}

	return pg, true
}

// AddProductUser 关联表操作::添加用户至产品线
func AddProductUser(productId, userId int, isOwner bool) bool {
	up := model.UserProduct{
		ProductId: productId,
		UserId:    userId,
		IsOwner:   &isOwner,
	}

	if res := db.Create(&up); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":    userId,
			"product_id": productId,
			"is_owner":   isOwner,
			"joined_at":  up.JoinedAt,
			"error":      res.Error,
		}, "An error occurred while db.Create.")
		return false
	}

	return true
}

// RemoveProductUser 关联表操作::移除产品线中用户
func RemoveProductUser(productId, userId int) bool {
	res := db.Delete(model.UserProduct{}, "user_id=? AND product_id=?", userId, productId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"product_id": productId,
			"user_id":    userId,
			"error":      res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// RemoveUserProducts 关联表操作::移除用户所有产品线
func RemoveUserProducts(userId int) bool {
	res := db.Delete(model.UserProduct{}, "user_id=?", userId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id": userId,
			"error":   res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// SetProductUserIsOwner 关联表操作::设置用户是否是产品线Owner
func SetProductUserIsOwner(productId, userId int, isOwner bool) bool {
	res := db.Model(&model.UserProduct{}).
		Where("user_id=? AND product_id=?", userId, productId).
		Update("is_owner", isOwner)
	if res.Error != nil {
		log.WarnWithFields(log.Fields{
			"product_id": productId,
			"user_id":    userId,
			"is_owner":   isOwner,
			"error":      res.Error,
		}, "An error occurred while db.Model.Where.Update.")
		return false
	}

	return true
}

// GetProductUserCount 关联表操作::列出产品线中用户数量
func GetProductUserCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.User{}).
		Select("users.id AS id, " +
			"users.username AS username, " +
			"up.is_owner AS is_owner, " +
			"up.joined_at AS joined_at").
		Joins("LEFT JOIN user_products AS up ON users.id = up.user_id")

	for k, v := range query {
		d = d.Where(k, v)
	}

	if res := d.Count(&count); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return count, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Select.Joins.Where.Count.")
		return count, false
	}
	return count, true
}

// ListProductUsers 关联表操作::列出产品线中所有用户
func ListProductUsers(productId int, query Query) (*[]model.MemberUser, bool) {
	var d = db.Model(&model.User{})
	mus := make([]model.MemberUser, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	res := d.Select("users.id AS id, "+
		"users.username AS username, "+
		"up.is_owner AS is_owner, "+
		"up.joined_at AS joined_at").
		Joins("LEFT JOIN user_products AS up ON users.id = up.user_id").
		Where("product_id=?", productId).
		Find(&mus)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error,
			}, "Record not found.")
			return &mus, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Select.Joins.Where.Find.")
		return nil, false
	}

	return &mus, true
}
