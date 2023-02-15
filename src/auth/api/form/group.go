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

type NewGroupForm struct {
	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9_-]{1,}$/)"`
	Description *string `json:"description" valid:"MinSize(0)"`

	dao *dao.Dao
	ctx context.Context
}

func (f *NewGroupForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsGroupExist(f.ctx, orm.Query{"name=?": f.Name}); exist {
		_ = v.SetError("Name", "组名称已存在")
	}
}

func (f *NewGroupForm) Validate(ctx context.Context, dao *dao.Dao) *cMsg.CodeMsg {
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

type RemoveGroupForm struct{}

func (*RemoveGroupForm) Validate(ctx context.Context, dao *dao.Dao, gId uint32) *cMsg.CodeMsg {
	// 验证组是否存在
	if exist, _ := dao.IsGroupExist(ctx, orm.Query{"id=?": gId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("组不存在")
	}

	// 验证组是否有关联存在
	if isEmpty, _ := dao.IsEmptyGroup(ctx, gId); !isEmpty {
		return msg.MsgAssociatedGroupFailed
	}

	return nil
}

type SetGroupForm struct {
	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9_-]{1,}$/)"`
	Description *string `json:"description" valid:"MinSize(0)"`

	groupId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *SetGroupForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsGroupExist(f.ctx, orm.Query{"name=?": f.Name, "id<>?": f.groupId}); exist {
		_ = v.SetError("Name", "组名称已存在")
	}
}

func (f *SetGroupForm) Validate(ctx context.Context, dao *dao.Dao, gId uint32) *cMsg.CodeMsg {
	f.groupId = gId

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

type PagedListGroupsParamsForm struct {
	Query *string `form:"query"`
}

func (pf *PagedListGroupsParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(name LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	return query
}

type AddUser2GroupForm struct {
	UserId  uint32 `json:"user_id" valid:"Required"`
	IsOwner bool   `json:"is_owner" valid:"Required"`

	groupId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *AddUser2GroupForm) Valid(v *validation.Validation) {
	// 验证用户是否存在
	if exist, _ := f.dao.IsUserExist(f.ctx, orm.Query{"id=?": f.UserId}); !exist {
		_ = v.SetError("UserId", "用户不存在")
	}

	// 验证用户是否已经属于该组
	if ct, _ := f.dao.GetGroupsUserCount(f.ctx, orm.Query{"group_id=?": f.groupId, "user_id": f.UserId}); ct > 0 {
		_ = v.SetError("UserId", "用户已经属于该组")
	}
}

func (f *AddUser2GroupForm) Validate(ctx context.Context, dao *dao.Dao, gId uint32) *cMsg.CodeMsg {
	// 验证组是否存在
	if exist, _ := dao.IsGroupExist(ctx, orm.Query{"id=?": gId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("组不存在")
	}

	f.groupId = gId

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

type RemoveGroupsUserForm struct{}

func (*RemoveGroupsUserForm) Validate(ctx context.Context, dao *dao.Dao, gId, userId uint32) *cMsg.CodeMsg {
	// 验证用户是否已经属于该组
	if ct, _ := dao.GetGroupsUserCount(ctx, orm.Query{"group_id=?": gId, "user_id": userId}); ct < 1 {
		return cMsg.MsgNotFoundFailed
	}

	return nil
}

type SetGroupsOwnerForm struct {
	IsOwner bool `json:"is_owner" valid:"Required"`
}

func (f *SetGroupsOwnerForm) Validate(ctx context.Context, dao *dao.Dao, groupId, userId uint32) *cMsg.CodeMsg {
	// 验证用户是否已经属于该组
	if ct, _ := dao.GetGroupsUserCount(ctx, orm.Query{"group_id=?": groupId, "user_id": userId}); ct < 1 {
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

type ListGroupsUsersParamsForm struct {
	IsOwner *bool `form:"is_owner"`
}

func (*ListGroupsUsersParamsForm) Validate(ctx context.Context, dao *dao.Dao, groupId uint32) *cMsg.CodeMsg {
	// 组不存在
	if exist, _ := dao.IsGroupExist(ctx, orm.Query{"id=?": groupId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("组不存在")
	}

	return nil
}

func (pf *ListGroupsUsersParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	if pf.IsOwner != nil {
		query["is_owner=?"] = *pf.IsOwner
	}

	return query
}
