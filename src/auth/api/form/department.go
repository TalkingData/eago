package form

import (
	"context"
	"database/sql"
	"eago/auth/conf/msg"
	"eago/auth/dao"
	cMsg "eago/common/code_msg"
	"eago/common/orm"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

// NewDepartmentForm struct
type NewDepartmentForm struct {
	Name     string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	ParentId *uint32 `json:"parent_id"`

	dao *dao.Dao
	ctx context.Context
}

func (f *NewDepartmentForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsDepartmentExist(f.ctx, orm.Query{"name=?": f.Name}); exist {
		_ = v.SetError("Name", "部门名称已存在")
	}
}

func (f *NewDepartmentForm) Validate(ctx context.Context, dao *dao.Dao) *cMsg.CodeMsg {
	f.ctx = ctx
	f.dao = dao

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(f)
	if err != nil {
		return cMsg.MsgValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return cMsg.MsgValidateFailed.SetError(valid.Errors)
	}

	return nil
}

type RemoveDepartmentForm struct{}

func (*RemoveDepartmentForm) Validate(ctx context.Context, dao *dao.Dao, deptId uint32) *cMsg.CodeMsg {
	// 验证部门是否存在
	if exist, _ := dao.IsDepartmentExist(ctx, orm.Query{"id=?": deptId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("部门不存在")
	}

	// 验证部门是否存在子树
	if exist, _ := dao.IsDepartmentExist(ctx, orm.Query{"parent_id=?": deptId}); exist {
		return msg.MsgAssociatedDepartmentFailed
	}

	// 验证部门是否有关联存在
	if isEmpty, _ := dao.IsEmptyDepartment(ctx, deptId); !isEmpty {
		return msg.MsgAssociatedDepartmentUserFailed
	}

	return nil
}

type SetDepartmentForm struct {
	Name     string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	ParentId *uint32 `json:"parent_id"`

	departmentId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *SetDepartmentForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsDepartmentExist(f.ctx, orm.Query{"name=?": f.Name, "id<>?": f.departmentId}); exist {
		_ = v.SetError("Name", "部门名称已存在")
	}
}

func (f *SetDepartmentForm) Validate(ctx context.Context, dao *dao.Dao, deptId uint32) *cMsg.CodeMsg {
	// 验证部门是否存在
	if exist, _ := dao.IsDepartmentExist(ctx, orm.Query{"id=?": deptId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("部门不存在")
	}

	f.departmentId = deptId

	f.ctx = ctx
	f.dao = dao

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(f)
	if err != nil {
		return cMsg.MsgValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return cMsg.MsgValidateFailed.SetError(valid.Errors)
	}

	return nil
}

type PagedListDepartmentsParamsForm struct {
	Query *string `form:"query"`
}

func (pf *PagedListDepartmentsParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(name LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	return query
}

type AddUser2DepartmentForm struct {
	UserId  uint32 `json:"user_id" valid:"Required"`
	IsOwner bool   `json:"is_owner" valid:"Required"`

	departmentId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *AddUser2DepartmentForm) Valid(v *validation.Validation) {
	// 验证用户是否存在
	if exist, _ := f.dao.IsUserExist(f.ctx, orm.Query{"id=?": f.UserId}); !exist {
		_ = v.SetError("UserId", "用户不存在")
	}

	// 验证用户是否已经属于该部门
	if ct, _ := f.dao.GetDepartmentUserCount(f.ctx, orm.Query{"user_id": f.UserId}); ct > 0 {
		_ = v.SetError("UserId", "用户已经属于某部门（用户只能属于一个部门）")
	}
}

func (f *AddUser2DepartmentForm) Validate(ctx context.Context, dao *dao.Dao, deptId uint32) *cMsg.CodeMsg {
	// 验证组是否存在
	if exist, _ := dao.IsDepartmentExist(ctx, orm.Query{"id=?": deptId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("部门不存在")
	}

	f.departmentId = deptId

	f.ctx = ctx
	f.dao = dao

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(f)
	if err != nil {
		return cMsg.MsgValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return cMsg.MsgValidateFailed.SetError(valid.Errors)
	}

	return nil
}

type RemoveDepartmentsUserForm struct{}

func (*RemoveDepartmentsUserForm) Validate(ctx context.Context, dao *dao.Dao, deptId, userId uint32) *cMsg.CodeMsg {
	// 验证用户是否已经属于该部门
	if ct, _ := dao.GetDepartmentUserCount(ctx, orm.Query{"department_id=?": deptId, "user_id": userId}); ct < 1 {
		return cMsg.MsgNotFoundFailed
	}

	return nil
}

type SetDepartmentsOwnerForm struct {
	IsOwner bool `json:"is_owner" valid:"Required"`
}

func (f *SetDepartmentsOwnerForm) Validate(ctx context.Context, dao *dao.Dao, deptId, userId uint32) *cMsg.CodeMsg {
	// 验证用户是否已经属于该组
	if ct, _ := dao.GetDepartmentUserCount(ctx, orm.Query{"department_id=?": deptId, "user_id": userId}); ct < 1 {
		return cMsg.MsgNotFoundFailed
	}

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(f)
	if err != nil {
		return cMsg.MsgValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return cMsg.MsgValidateFailed.SetError(valid.Errors)
	}

	return nil
}

type ListDepartmentsUsersParamsForm struct {
	IsOwner *bool `form:"is_owner"`
}

func (*ListDepartmentsUsersParamsForm) Validate(ctx context.Context, dao *dao.Dao, deptId uint32) *cMsg.CodeMsg {
	// 组不存在
	if exist, _ := dao.IsDepartmentExist(ctx, orm.Query{"id=?": deptId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("部门不存在")
	}

	return nil
}

func (f *ListDepartmentsUsersParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	if f.IsOwner != nil {
		query["is_owner=?"] = *f.IsOwner
	}

	return query
}
