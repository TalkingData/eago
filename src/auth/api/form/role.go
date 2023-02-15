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

type NewRoleForm struct {
	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9_-]{1,}$/)"`
	Description *string `json:"description" valid:"MinSize(0);MaxSize(500)"`

	dao *dao.Dao
	ctx context.Context
}

func (f *NewRoleForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsRoleExist(f.ctx, orm.Query{"name=?": f.Name}); exist {
		_ = v.SetError("Name", "已有相同名称的角色存在")
	}
}

func (f *NewRoleForm) Validate(ctx context.Context, dao *dao.Dao) *cMsg.CodeMsg {
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

type RemoveRoleForm struct{}

func (*RemoveRoleForm) Validate(ctx context.Context, dao *dao.Dao, roleId uint32) *cMsg.CodeMsg {
	// 验证角色是否存在
	if exist, _ := dao.IsRoleExist(ctx, orm.Query{"id=?": roleId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("角色不存在")
	}

	// 判断角色是否是空角色，不是空角色说明角色有关联，无法删除
	if isEmpty, _ := dao.IsEmptyRole(ctx, roleId); !isEmpty {
		return msg.MsgAssociatedRoleFailed
	}

	return nil
}

type SetRoleForm struct {
	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9_-]{1,}$/)"`
	Description *string `json:"description" valid:"MinSize(0);MaxSize(500)"`

	roleId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *SetRoleForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsRoleExist(f.ctx, orm.Query{"name=?": f.Name, "id<>?": f.roleId}); exist {
		_ = v.SetError("Name", "已有相同名称的角色存在")
	}
}

func (f *SetRoleForm) Validate(ctx context.Context, dao *dao.Dao, roleId uint32) *cMsg.CodeMsg {
	// 验证角色是否存在
	if exist, _ := dao.IsRoleExist(ctx, orm.Query{"id=?": roleId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("角色不存在")
	}

	f.roleId = roleId

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

type PagedListRolesParamsForm struct {
	Query *string `form:"query"`
}

func (pf *PagedListRolesParamsForm) GenQuery() orm.Query {
	query := orm.Query{}
	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(name LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}
	return query
}

type AddUser2RoleForm struct {
	UserId uint32 `json:"user_id" valid:"Required"`

	roleId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *AddUser2RoleForm) Valid(v *validation.Validation) {
	// 验证用户是否存在
	if exist, _ := f.dao.IsUserExist(f.ctx, orm.Query{"id=?": f.UserId}); !exist {
		_ = v.SetError("UserId", "用户不存在")
	}

	// 验证用户是否已经属于该角色
	if ct, _ := f.dao.GetRolesUserCount(f.ctx, orm.Query{"role_id=?": f.roleId, "user_id": f.UserId}); ct > 0 {
		_ = v.SetError("UserId", "用户已经属于该角色")
	}
}

func (f *AddUser2RoleForm) Validate(ctx context.Context, dao *dao.Dao, roleId uint32) *cMsg.CodeMsg {
	// 验证角色是否存在
	if exist, _ := dao.IsRoleExist(ctx, orm.Query{"id=?": roleId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("角色不存在")
	}

	f.roleId = roleId

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

type RemoveRolesUserForm struct{}

func (*RemoveRolesUserForm) Validate(ctx context.Context, dao *dao.Dao, roleId, userId uint32) *cMsg.CodeMsg {
	// 验证用户是否已经属于该角色
	if ct, _ := dao.GetRolesUserCount(ctx, orm.Query{"role_id=?": roleId, "user_id": userId}); ct < 1 {
		return cMsg.MsgNotFoundFailed
	}

	return nil
}
