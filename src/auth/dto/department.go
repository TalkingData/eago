package dto

import (
	"database/sql"
	"eago/auth/conf/msg"
	"eago/auth/dao"
	"eago/common/message"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

// NewDepartment struct
type NewDepartment struct {
	Name     string `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	ParentId *int   `json:"parent_id"`
}

// Valid
func (nd *NewDepartment) Valid(v *validation.Validation) {
	if ct, _ := dao.GetDepartmentCount(dao.Query{"name=?": nd.Name}); ct > 0 {
		_ = v.SetError("Name", "部门名称已存在")
	}
}

// Validate
func (nd *NewDepartment) Validate() *message.Message {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(nd)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// RemoveDepartment struct
type RemoveDepartment struct{}

// Validate
func (*RemoveDepartment) Validate(deptId int) *message.Message {
	// 验证部门是否存在
	if ct, _ := dao.GetDepartmentCount(dao.Query{"id=?": deptId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("部门不存在")
	}

	// 验证部门是存在子树
	if ct, _ := dao.GetDepartmentCount(dao.Query{"parent_id=?": deptId}); ct > 0 {
		return msg.AssociatedDepartmentFailed
	}

	// 验证部门是否有关联存在
	if ct, _ := dao.GetDepartmentUserCount(dao.Query{"department_id=?": deptId}); ct > 0 {

		return msg.AssociatedDepartmentUserFailed
	}

	return nil
}

// SetDepartment struct
type SetDepartment struct {
	departmentId int

	Name     string `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	ParentId *int   `json:"parent_id"`
}

// Valid
func (sd *SetDepartment) Valid(v *validation.Validation) {
	if ct, _ := dao.GetDepartmentCount(dao.Query{"name=?": sd.Name, "id<>?": sd.departmentId}); ct > 0 {
		_ = v.SetError("Name", "部门名称已存在")
	}
}

// Validate
func (sd *SetDepartment) Validate(deptId int) *message.Message {
	// 验证部门是否存在
	if ct, _ := dao.GetDepartmentCount(dao.Query{"id=?": deptId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("部门不存在")
	}

	sd.departmentId = deptId
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(sd)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// ListDepartmentsQuery struct
type ListDepartmentsQuery struct {
	Query *string `form:"query"`
}

// UpdateQuery
func (ldq *ListDepartmentsQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if ldq.Query != nil && *ldq.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *ldq.Query)
		query["name LIKE @query OR id LIKE @query"] = sql.Named("query", likeQuery)
	}

	return nil
}

// AddUser2Department struct
type AddUser2Department struct {
	departmentId int

	UserId  int  `json:"user_id" valid:"Required"`
	IsOwner bool `json:"is_owner" valid:"Required"`
}

// Valid
func (aud *AddUser2Department) Valid(v *validation.Validation) {
	// 验证用户是否存在
	if ct, _ := dao.GetUserCount(dao.Query{"id=?": aud.UserId}); ct < 1 {
		_ = v.SetError("UserId", "用户不存在")
	}

	// 验证用户是否已经属于该部门
	if ct, _ := dao.GetDepartmentUserCount(dao.Query{"user_id": aud.UserId}); ct > 0 {
		_ = v.SetError("UserId", "用户已经属于某部门（用户只能属于一个部门）")
	}
}

// Validate
func (aug *AddUser2Department) Validate(deptId int) *message.Message {
	// 验证组是否存在
	if ct, _ := dao.GetDepartmentCount(dao.Query{"id=?": deptId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("部门不存在")
	}

	aug.departmentId = deptId
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(aug)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// RemoveDepartmentUser struct
type RemoveDepartmentUser struct{}

// Validate
func (*RemoveDepartmentUser) Validate(deptId, userId int) *message.Message {
	// 验证用户是否已经属于该部门
	if ct, _ := dao.GetDepartmentUserCount(dao.Query{"department_id=?": deptId, "user_id": userId}); ct < 1 {
		return msg.NotFoundFailed
	}

	return nil
}

// SetUserIsDepartmentOwner struct
type SetUserIsDepartmentOwner struct {
	IsOwner bool `json:"is_owner" valid:"Required"`
}

// Validate
func (suo *SetUserIsDepartmentOwner) Validate(deptId, userId int) *message.Message {
	// 验证用户是否已经属于该组
	if ct, _ := dao.GetDepartmentUserCount(dao.Query{"department_id=?": deptId, "user_id": userId}); ct < 1 {
		return msg.NotFoundFailed
	}

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(suo)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// ListDepartmentUsersQuery struct
type ListDepartmentUsersQuery struct {
	IsOwner *bool `form:"is_owner"`
}

func (*ListDepartmentUsersQuery) Validate(deptId int) *message.Message {
	// 组不存在
	if ct, _ := dao.GetDepartmentCount(dao.Query{"id=?": deptId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("部门不存在")
	}

	return nil
}

// UpdateQuery
func (ldu *ListDepartmentUsersQuery) UpdateQuery(query dao.Query) error {
	if ldu.IsOwner != nil {
		query["is_owner=?"] = *ldu.IsOwner
	}

	return nil
}
