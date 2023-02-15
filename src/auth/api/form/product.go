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

type NewProductForm struct {
	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z0-9-\u4e00-\u9fa5]+$/)"`
	Alias       string  `json:"alias" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9-]{1,}$/)"`
	Disabled    *bool   `json:"disabled" valid:"Required"`
	Description *string `json:"description" valid:"MinSize(0);MaxSize(500)"`

	dao *dao.Dao
	ctx context.Context
}

func (f *NewProductForm) Valid(v *validation.Validation) {
	nameQuery := orm.Query{"name=@name OR alias=@name": sql.Named("name", f.Name)}
	if exist, _ := f.dao.IsProductExist(f.ctx, nameQuery); exist {
		_ = v.SetError("Name", "名称已存在，或与其他产品线的别名相同")
	}

	aliasQ := orm.Query{"name=@alias OR alias=@alias": sql.Named("alias", f.Alias)}
	if exist, _ := f.dao.IsProductExist(f.ctx, aliasQ); exist {
		_ = v.SetError("Alias", "别名重复，或与其他产品线的名称相同")
	}
}

func (f *NewProductForm) Validate(ctx context.Context, dao *dao.Dao) *cMsg.CodeMsg {
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

type RemoveProductForm struct{}

func (*RemoveProductForm) Validate(ctx context.Context, dao *dao.Dao, prodId uint32) *cMsg.CodeMsg {
	// 验证产品线是否存在
	if exist, _ := dao.IsProductExist(ctx, orm.Query{"id=?": prodId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("产品线不存在")
	}

	// 验证产品线是否有关联存在
	if isEmpty, _ := dao.IsEmptyProduct(ctx, prodId); !isEmpty {
		return msg.MsgAssociatedProductFailed
	}

	return nil
}

type SetProductForm struct {
	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z0-9-\u4e00-\u9fa5]+$/)"`
	Alias       string  `json:"alias" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9-]{1,}$/)"`
	Disabled    *bool   `json:"disabled" valid:"Required"`
	Description *string `json:"description" valid:"MinSize(0);MaxSize(500)"`

	prodId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *SetProductForm) Valid(v *validation.Validation) {
	nameQuery := orm.Query{
		"id<>?":                       f.prodId,
		"(name=@name OR alias=@name)": sql.Named("name", f.Name),
	}
	if exist, _ := f.dao.IsProductExist(f.ctx, nameQuery); exist {
		_ = v.SetError("Name", "产品线名称已存在，或者该名称与其他产品线的别名相同")
	}

	aliasQuery := orm.Query{
		"id<>?":                         f.prodId,
		"(name=@alias OR alias=@alias)": sql.Named("alias", f.Alias),
	}
	if exist, _ := f.dao.IsProductExist(f.ctx, aliasQuery); exist {
		_ = v.SetError("Alias", "产品线别名重复，或者该别名与其他产品线的名称相同")
	}
}

func (f *SetProductForm) Validate(ctx context.Context, dao *dao.Dao, prodId uint32) *cMsg.CodeMsg {
	// 验证产品线是否存在
	if exist, _ := dao.IsProductExist(ctx, orm.Query{"id=?": prodId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("产品线不存在")
	}

	f.prodId = prodId

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

type PagedListProductsParamsForm struct {
	Query    *string `form:"query"`
	Disabled *bool   `form:"disabled"`
}

func (pf *PagedListProductsParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(name LIKE @query OR "+
			"alias LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if pf.Disabled != nil {
		query["disabled=?"] = *pf.Disabled
	}

	return query
}

type AddUser2ProductForm struct {
	UserId  uint32 `json:"user_id" valid:"Required"`
	IsOwner bool   `json:"is_owner" valid:"Required"`

	prodId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *AddUser2ProductForm) Valid(v *validation.Validation) {
	// 验证用户是否存在
	if exist, _ := f.dao.IsUserExist(f.ctx, orm.Query{"id=?": f.UserId}); !exist {
		_ = v.SetError("UserId", "用户不存在")
	}

	// 验证用户是否已经属于该产品线
	if ct, _ := f.dao.GetProductsUserCount(f.ctx, orm.Query{"product_id=?": f.prodId, "user_id": f.UserId}); ct > 0 {
		_ = v.SetError("UserId", "用户已经属于该产品线")
	}
}

func (f *AddUser2ProductForm) Validate(ctx context.Context, dao *dao.Dao, prodId uint32) *cMsg.CodeMsg {
	// 验证产品线是否存在
	if exist, _ := dao.IsProductExist(ctx, orm.Query{"id=?": prodId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("产品线不存在")
	}

	f.prodId = prodId

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

type RemoveProductsUserForm struct{}

func (*RemoveProductsUserForm) Validate(ctx context.Context, dao *dao.Dao, prodId, userId uint32) *cMsg.CodeMsg {
	// 验证用户是否已经属于该产品线
	if ct, _ := dao.GetProductsUserCount(ctx, orm.Query{"product_id=?": prodId, "user_id": userId}); ct < 1 {
		return cMsg.MsgNotFoundFailed
	}

	return nil
}

type SetProductsOwnerForm struct {
	IsOwner bool `json:"is_owner" valid:"Required"`
}

func (f *SetProductsOwnerForm) Validate(ctx context.Context, dao *dao.Dao, prodId, userId uint32) *cMsg.CodeMsg {
	// 验证用户是否已经属于该产品线
	if ct, _ := dao.GetProductsUserCount(ctx, orm.Query{"product_id=?": prodId, "user_id": userId}); ct < 1 {
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

type ListProductsUsersParamsForm struct {
	IsOwner *bool `form:"is_owner"`
}

func (pf *ListProductsUsersParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	if pf.IsOwner != nil {
		query["is_owner=?"] = *pf.IsOwner
	}

	return query
}
