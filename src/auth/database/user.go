package database

import (
	"eago-common/api-suite/pagination"
	"eago-common/log"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

var UserModel userModel

type userModel struct{}

type userMember struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	IsOwner   bool   `json:"is_owner"`
	CreatedAt MyTime `json:"joined_at"`
}

type userProductMember struct {
	userMember

	Alias    string `json:"alias"`
	Disabled bool   `json:"disabled"`
}

type memberUser struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	IsOwner   bool   `json:"is_owner"`
	CreatedAt MyTime `json:"joined_at"`
}

type User struct {
	Id          int     `json:"id" swaggerignore:"true"`
	Username    string  `json:"username" swaggerignore:"true"`
	Password    string  `json:"-" swaggerignore:"true"`
	Email       string  `json:"email" binding:"required" `
	Phone       string  `json:"phone" binding:"required"`
	IsSuperuser bool    `json:"is_superuser" swaggerignore:"true"`
	Disabled    bool    `json:"disabled" swaggerignore:"true"`
	LastLogin   *MyTime `json:"last_login" swaggerignore:"true"`
	CreatedAt   MyTime  `json:"created_at" swaggerignore:"true"`
	UpdatedAt   *MyTime `json:"updated_at" swaggerignore:"true"`
}

// 新建用户
func (um *userModel) New(username string, email string, login bool) *User {
	var u = User{}
	u.Username = username
	u.Email = email

	// 判断时候设置最近登录时间
	if login {
		u.LastLogin = &MyTime{time.Now()}
	}

	if res := db.Create(&u); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"username":     u.Username,
			"email":        u.Email,
			"is_superuser": false,
			"disabled":     false,
			"last_login":   u.LastLogin,
			"error":        res.Error.Error(),
		}, "Error in userModel.New.")
		return nil
	}

	return &u
}

// 删除用户
func (um *userModel) Delete(userId int) bool {
	res := db.Delete(User{}, "id=?", userId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    userId,
			"error": res.Error.Error(),
		}, "Error in userModel.Delete.")
		return false
	}

	return true
}

// 更新最后登录时间
func (um *userModel) SetLastLogin(query *Query) bool {
	var d = db.Model(&User{})

	for k, v := range *query {
		d = d.Where(k, v)
	}
	res := d.Update("LastLogin", time.Now())
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in userModel.SetLastLogin.")
		return false
	}

	return true
}

// 更新用户为禁用状态
func (um *userModel) SetDisabled(id int) bool {
	res := db.Model(&User{}).
		Where("id=?", id).
		Update("Disabled", true)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error.Error(),
		}, "Error in userModel.SetDisabled.")
		return false
	}

	return true
}

// 更新用户
func (um *userModel) Set(id int, email string, phone string) (*User, bool) {
	var u = User{}

	res := db.Model(&User{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"email": email,
			"phone": phone,
		}).
		First(&u)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error.Error(),
		}, "Error in userModel.Set.")
		return nil, false
	}

	return &u, true
}

// 查询单个用户
func (um *userModel) Get(query *Query) (*User, bool) {
	var (
		u = User{}
		d = db
	)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	if res := d.First(&u); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in userModel.Get.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in userModel.Get.")
		return nil, false
	}

	return &u, true
}

// 查询用户
func (um *userModel) List(query *Query) (*[]User, bool) {
	var d = db
	us := make([]User, 0)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	if res := d.Find(&us); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in userModel.List.")
			return &us, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in userModel.List.")
		return nil, false
	}

	return &us, true
}

// 查询用户-分页
func (um *userModel) PagedList(query *Query, page int, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&User{})
	us := make([]User, 0)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &us)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": err.Error(),
		}, "Error in userModel.PagedList.")
		return nil, false
	}

	return pg, true
}

// 查询用户是否是Admin
func (um *userModel) IsSuperuser(userId int) bool {
	var u = User{}

	res := db.Where("id=?", userId).First(&u)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"user_id": userId,
				"error":   res.Error.Error(),
			}, "Record not found in userModel.IsSuperuser.")
			return false
		}
		log.ErrorWithFields(log.Fields{
			"user_id": userId,
			"error":   res.Error.Error(),
		}, "Error in userModel.IsSuperuser.")
		return false
	}

	if u.IsSuperuser {
		return true
	}

	return false
}

// 关联表操作::列出用户所有角色
func (um *userModel) ListRoles(userId int) (*[]Role, bool) {
	rs := make([]Role, 0)

	res := db.Model(&Role{}).
		Joins("LEFT JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_id=?", userId).
		Find(&rs)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in userModel.ListRoles.")
			return &rs, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in userModel.ListRoles.")
		return nil, false
	}

	return &rs, true
}

// 关联表操作::列出用户所有产品线
func (um *userModel) ListProducts(userId int) (*[]userProductMember, bool) {
	ups := make([]userProductMember, 0)

	res := db.Model(&Product{}).
		Select("products.id AS id, products.name AS name, products.alias, products.disabled, up.is_owner, up.created_at").
		Joins("LEFT JOIN user_products AS up ON products.id = up.product_id").
		Where("user_id=?", userId).
		Find(&ups)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in userModel.ListProducts.")
			return &ups, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in userModel.ListProducts.")
		return nil, false
	}

	return &ups, true
}

// 关联表操作::列出用户所有组
func (um *userModel) ListGroups(userId int) (*[]userMember, bool) {
	ugs := make([]userMember, 0)

	res := db.Model(&Group{}).
		Select("groups.id AS id, groups.name AS name, ug.is_owner AS is_owner, ug.created_at AS created_at").
		Joins("LEFT JOIN user_groups AS ug ON groups.id = ug.group_id").
		Where("user_id=?", userId).
		Find(&ugs)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in userModel.ListGroups.")
			return &ugs, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in userModel.ListGroups.")
		return nil, false
	}

	return &ugs, true
}

// 关联表操作::获得用户所在部门
func (um *userModel) GetDepartment(userId int) (*userMember, bool) {
	var uMember = userMember{}

	res := db.Model(&Group{}).
		Select("groups.id AS id, groups.name AS name, ug.is_owner AS is_owner, ug.created_at AS created_at").
		Joins("LEFT JOIN user_groups AS ug ON groups.id = ug.group_id").
		Where("user_id=?", userId).
		First(&uMember)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in userModel.ListGroups.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in userModel.ListGroups.")
		return nil, false
	}

	return &uMember, true
}
