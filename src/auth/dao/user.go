package dao

import (
	"eago/auth/model"
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/common/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// NewUser 新建用户
func NewUser(username, email string, login bool) *model.User {
	u := model.User{
		Username: username,
		Email:    email,
	}

	// 判断时候设置最近登录时间
	if login {
		u.LastLogin = &utils.LocalTime{Time: time.Now()}
	}

	if res := db.Create(&u); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"username":     u.Username,
			"email":        u.Email,
			"is_superuser": false,
			"disabled":     false,
			"last_login":   u.LastLogin,
			"error":        res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &u
}

// RemoveUser 删除用户
func RemoveUser(userId int) bool {
	res := db.Delete(model.User{}, "id=?", userId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    userId,
			"error": res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// SetUserLastLogin 更新最后登录时间
func SetUserLastLogin(id int) bool {
	var d = db.Model(&model.User{}).Where("id=?", id)

	res := d.Update("last_login", &utils.LocalTime{Time: time.Now()})
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error,
		}, "An error occurred while db.Model.Update.")
		return false
	}

	return true
}

// SetUserDisabled 更新用户为禁用状态
func SetUserDisabled(id int) bool {
	res := db.Model(&model.User{}).
		Where("id=?", id).
		Update("disabled", true)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error,
		}, "An error occurred while b.Model.Where.Update.")
		return false
	}

	return true
}

// SetUser 更新用户
func SetUser(id int, email, phone string) (*model.User, bool) {
	u := model.User{}

	res := db.Model(&model.User{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"email": email,
			"phone": phone,
		}).
		First(&u)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error,
		}, "An error occurred while db.Model.Where.Updates.First.")
		return nil, false
	}

	return &u, true
}

// GetUser 查询单个用户
func GetUser(query Query) (*model.User, bool) {
	var (
		u = model.User{}
		d = db
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&u); res.Error != nil {
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

	return &u, true
}

// GetUserCount 查询用户数量
func GetUserCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.User{})

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

// ListUsers 查询用户
func ListUsers(query Query) (*[]model.User, bool) {
	var d = db
	us := make([]model.User, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&us); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return &us, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return &us, true
}

// PagedListUsers 查询用户-分页
func PagedListUsers(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&model.User{})
	us := make([]model.User, 0)

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
			"error": err,
		}, "An error occurred while pagination.GormPaging.")
		return nil, false
	}

	return pg, true
}

// UserIsSuperuser 查询用户是否是Admin
func UserIsSuperuser(userId int) bool {
	u := model.User{}

	res := db.Where("id=?", userId).First(&u)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"user_id": userId,
				"error":   res.Error,
			}, "Record not found.")
			return false
		}
		log.ErrorWithFields(log.Fields{
			"user_id": userId,
			"error":   res.Error,
		}, "An error occurred while db.Where.")
		return false
	}

	return u.IsSuperuser
}

// ListUserRoles 关联表操作::列出用户所有角色
func ListUserRoles(userId int) (*[]model.Role, bool) {
	rs := make([]model.Role, 0)

	res := db.Model(&model.Role{}).
		Joins("LEFT JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_id=?", userId).
		Find(&rs)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error,
			}, "Record not found.")
			return &rs, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Select.Joins.Where.Find.")
		return nil, false
	}

	return &rs, true
}

// ListUserProducts 关联表操作::列出用户所有产品线
func ListUserProducts(userId int) (*[]model.UserProductMember, bool) {
	ups := make([]model.UserProductMember, 0)

	res := db.Model(&model.Product{}).
		Select("products.id AS id, "+
			"products.name AS name, "+
			"products.alias AS alias, "+
			"products.disabled AS disabled, "+
			"up.is_owner AS is_owner, "+
			"up.joined_at AS joined_at").
		Joins("LEFT JOIN user_products AS up ON products.id = up.product_id").
		Where("user_id=?", userId).
		Find(&ups)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error,
			}, "Record not found.")
			return &ups, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while Ldb.Select.Joins.Where.Find.")
		return nil, false
	}

	return &ups, true
}

// ListUserGroups 关联表操作::列出用户所有组
func ListUserGroups(userId int) (*[]model.UserMember, bool) {
	ugs := make([]model.UserMember, 0)

	res := db.Model(&model.Group{}).
		Select("groups.id AS id, "+
			"groups.name AS name, "+
			"ug.is_owner AS is_owner, "+
			"ug.joined_at AS joined_at").
		Joins("LEFT JOIN user_groups AS ug ON groups.id = ug.group_id").
		Where("user_id=?", userId).
		Find(&ugs)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error,
			}, "Record not found.")
			return &ugs, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Select.Joins.Where.Find.")
		return nil, false
	}

	return &ugs, true
}

// GetUserDepartment 关联表操作::获得用户所在部门
func GetUserDepartment(userId int) (*model.UserDepartmentMember, bool) {
	var uDeptMember = model.UserDepartmentMember{}

	res := db.Model(&model.Department{}).
		Select("departments.id AS id, "+
			"departments.name AS name, "+
			"departments.parent_id AS parent_id, "+
			"ud.is_owner AS is_owner, "+
			"ud.joined_at AS joined_at").
		Joins("LEFT JOIN user_departments AS ud ON departments.id = ud.department_id").
		Where("user_id=?", userId).
		First(&uDeptMember)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error,
			}, "Record not found.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Select.Joins.Where.Find.")
		return nil, false
	}

	return &uDeptMember, true
}

// GetUserDepartment 关联表操作::获得用户所在部门
func MakeUserHandover(uId, tgtId int) error {
	tx := db.Begin()

	// 获得交接用户
	user := &model.User{}
	if res := tx.Where("id=?", uId).Find(user); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id": uId,
			"error":   res.Error,
		}, "Failed to find user.")
		tx.Rollback()
		return res.Error
	}

	// 获得交接目标用户
	tgtUser := &model.User{}
	if res := tx.Where("id=?", tgtId).Find(tgtUser); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"target_user_id": tgtId,
			"error":          res.Error,
		}, "Failed to find handover target user.")
		tx.Rollback()
		return res.Error
	}

	// 交接产品线Owner权限
	res := tx.Model(model.UserProduct{}).
		Where("user_id=? AND is_owner=?", uId, true).
		Updates(map[string]interface{}{"user_id": tgtId, "joined_at": time.Now()})
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":        uId,
			"target_user_id": tgtId,
			"error":          res.Error,
		}, "Failed to handover user's products.")
		tx.Rollback()
		return res.Error
	}

	// 交接组Owner权限
	res = tx.Model(model.UserGroup{}).
		Where("user_id=? AND is_owner=?", uId, true).
		Updates(map[string]interface{}{"user_id": tgtId, "joined_at": time.Now()})
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":        uId,
			"target_user_id": tgtId,
			"error":          res.Error,
		}, "Failed to handover user's groups.")
		tx.Rollback()
		return res.Error
	}

	// 交接角色权限
	res = tx.Model(model.UserRole{}).
		Where("user_id=?", uId).
		Updates(map[string]interface{}{"user_id": tgtId, "joined_at": time.Now()})
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":        uId,
			"target_user_id": tgtId,
			"error":          res.Error,
		}, "Failed to handover user's roles.")
		tx.Rollback()
		return res.Error
	}

	// 删除所在产品线
	if res = tx.Where("user_id=?", uId).Delete(model.UserProduct{}); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id": uId,
			"error":   res.Error,
		}, "Failed to remove user products.")
		tx.Rollback()
		return res.Error
	}

	// 删除所在部门
	if res = tx.Where("user_id=?", uId).Delete(model.UserDepartment{}); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id": uId,
			"error":   res.Error,
		}, "Failed to remove user department.")
		tx.Rollback()
		return res.Error
	}

	// 删除所在组
	if res = tx.Where("user_id=?", uId).Delete(model.UserGroup{}); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id": uId,
			"error":   res.Error,
		}, "Failed to remove user groups.")
		tx.Rollback()
		return res.Error
	}

	// 删除所在角色
	if res = tx.Where("user_id=?", uId).Delete(model.UserRole{}); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id": uId,
			"error":   res.Error,
		}, "Failed to remove user roles.")
		tx.Rollback()
		return res.Error
	}

	// 删除用户
	if res = tx.Where("id=?", uId).Delete(model.User{}); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id": uId,
			"error":   res.Error,
		}, "Failed to remove user.")
		tx.Rollback()
		return res.Error
	}

	tx.Commit()
	return nil
}
