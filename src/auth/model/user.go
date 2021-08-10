package model

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/common/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type userMember struct {
	Id       int              `json:"id"`
	Name     string           `json:"name"`
	IsOwner  bool             `json:"is_owner"`
	JoinedAt *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

type userProductMember struct {
	Id       int              `json:"id"`
	Name     string           `json:"name"`
	Alias    string           `json:"alias"`
	Disabled bool             `json:"disabled"`
	IsOwner  bool             `json:"is_owner"`
	JoinedAt *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

type userDepartmentMember struct {
	Id       int              `json:"id"`
	Name     string           `json:"name"`
	ParentId *int             `json:"parent_id"`
	IsOwner  bool             `json:"is_owner"`
	JoinedAt *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

type memberUser struct {
	Id       int              `json:"id"`
	Username string           `json:"username"`
	IsOwner  bool             `json:"is_owner"`
	JoinedAt *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

type User struct {
	Id          int              `json:"id" swaggerignore:"true"`
	Username    string           `json:"username" swaggerignore:"true"`
	Password    string           `json:"-" swaggerignore:"true"`
	Email       string           `json:"email" binding:"required,email" `
	Phone       string           `json:"phone" binding:"required"`
	IsSuperuser bool             `json:"is_superuser" swaggerignore:"true"`
	Disabled    bool             `json:"disabled" swaggerignore:"true"`
	LastLogin   *utils.LocalTime `json:"last_login" swaggerignore:"true"`
	CreatedAt   *utils.LocalTime `json:"created_at" swaggerignore:"true"`
	UpdatedAt   *utils.LocalTime `json:"updated_at" swaggerignore:"true"`
}

// NewUser 新建用户
func NewUser(username, email string, login bool) *User {
	var u = User{}
	u.Username = username
	u.Email = email

	// 判断时候设置最近登录时间
	if login {
		u.LastLogin = &utils.LocalTime{time.Now()}
	}

	if res := db.Create(&u); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"username":     u.Username,
			"email":        u.Email,
			"is_superuser": false,
			"disabled":     false,
			"last_login":   u.LastLogin,
			"error":        res.Error.Error(),
		}, "Error in NewUser.")
		return nil
	}

	return &u
}

// RemoveUser 删除用户
func RemoveUser(userId int) bool {
	res := db.Delete(User{}, "id=?", userId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    userId,
			"error": res.Error.Error(),
		}, "Error in RemoveUser.")
		return false
	}

	return true
}

// SetUserLastLogin 更新最后登录时间
func SetUserLastLogin(id int) bool {
	var d = db.Model(&User{}).Where("id=?", id)

	res := d.Update("last_login", &utils.LocalTime{time.Now()})
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error.Error(),
		}, "Error in SetUserLastLogin.")
		return false
	}

	return true
}

// SetUserDisabled 更新用户为禁用状态
func SetUserDisabled(id int) bool {
	res := db.Model(&User{}).
		Where("id=?", id).
		Update("disabled", true)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error.Error(),
		}, "Error in SetUserDisabled.")
		return false
	}

	return true
}

// SetUser 更新用户
func SetUser(id int, email, phone string) (*User, bool) {
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
		}, "Error in SetUser.")
		return nil, false
	}

	return &u, true
}

// GetUser 查询单个用户
func GetUser(query Query) (*User, bool) {
	var (
		u = User{}
		d = db
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&u); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in GetUser.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in GetUser.")
		return nil, false
	}

	return &u, true
}

// ListUsers 查询用户
func ListUsers(query Query) (*[]User, bool) {
	var d = db
	us := make([]User, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&us); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in ListUsers.")
			return &us, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in ListUsers.")
		return nil, false
	}

	return &us, true
}

// PagedListUsers 查询用户-分页
func PagedListUsers(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&User{})
	us := make([]User, 0)

	for k, v := range query {
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
		}, "Error in PagedListUsers.")
		return nil, false
	}

	return pg, true
}

// UserIsSuperuser 查询用户是否是Admin
func UserIsSuperuser(userId int) bool {
	var u = User{}

	res := db.Where("id=?", userId).First(&u)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"user_id": userId,
				"error":   res.Error.Error(),
			}, "Record not found in UserIsSuperuser.")
			return false
		}
		log.ErrorWithFields(log.Fields{
			"user_id": userId,
			"error":   res.Error.Error(),
		}, "Error in UserIsSuperuser.")
		return false
	}

	return u.IsSuperuser
}

// ListUserRoles 关联表操作::列出用户所有角色
func ListUserRoles(userId int) (*[]Role, bool) {
	rs := make([]Role, 0)

	res := db.Model(&Role{}).
		Joins("LEFT JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_id=?", userId).
		Find(&rs)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in ListUserRoles.")
			return &rs, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in ListUserRoles.")
		return nil, false
	}

	return &rs, true
}

// ListUserProducts 关联表操作::列出用户所有产品线
func ListUserProducts(userId int) (*[]userProductMember, bool) {
	ups := make([]userProductMember, 0)

	res := db.Model(&Product{}).
		Select("products.id AS id, products.name AS name, products.alias, products.disabled, up.is_owner AS is_owner, up.joined_at").
		Joins("LEFT JOIN user_products AS up ON products.id = up.product_id").
		Where("user_id=?", userId).
		Find(&ups)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in ListUserProducts.")
			return &ups, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in ListUserProducts.")
		return nil, false
	}

	return &ups, true
}

// ListUserGroups 关联表操作::列出用户所有组
func ListUserGroups(userId int) (*[]userMember, bool) {
	ugs := make([]userMember, 0)

	res := db.Model(&Group{}).
		Select("groups.id AS id, groups.name AS name, ug.is_owner AS is_owner, ug.joined_at AS joined_at").
		Joins("LEFT JOIN user_groups AS ug ON groups.id = ug.group_id").
		Where("user_id=?", userId).
		Find(&ugs)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in ListGroups.")
			return &ugs, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in ListGroups.")
		return nil, false
	}

	return &ugs, true
}

// GetUserDepartment 关联表操作::获得用户所在部门
func GetUserDepartment(userId int) (*userDepartmentMember, bool) {
	var uDeptMember = userDepartmentMember{}

	res := db.Model(&Department{}).
		Select("departments.id AS id, departments.name AS name, departments.parent_id AS parent_id, ud.is_owner AS is_owner, ud.joined_at AS joined_at").
		Joins("LEFT JOIN user_departments AS ud ON departments.id = ug.department_id").
		Where("user_id=?", userId).
		First(&uDeptMember)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in GetUserDepartment.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in GetUserDepartment.")
		return nil, false
	}

	return &uDeptMember, true
}
