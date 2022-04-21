package dto

import (
	"database/sql"
	"eago/auth/conf/msg"
	"eago/auth/dao"
	"eago/common/message"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

// NewProduct struct
type NewProduct struct {
	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z0-9-\u4e00-\u9fa5]+$/)"`
	Alias       string  `json:"alias" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9-]{1,}$/)"`
	Disabled    *bool   `json:"disabled" valid:"Required"`
	Description *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
}

func (np *NewProduct) Valid(v *validation.Validation) {
	nameQ := dao.Query{"name=@name OR alias=@name": sql.Named("name", np.Name)}
	if ct, _ := dao.GetProductCount(nameQ); ct > 0 {
		_ = v.SetError("Name", "名称已存在，或与其他产品线的别名相同")
	}

	aliasQ := dao.Query{"name=@alias OR alias=@alias": sql.Named("alias", np.Alias)}
	if ct, _ := dao.GetProductCount(aliasQ); ct > 0 {
		_ = v.SetError("Alias", "别名重复，或与其他产品线的名称相同")
	}
}

func (np *NewProduct) Validate() *message.Message {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(np)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// RemoveProduct struct
type RemoveProduct struct{}

func (*RemoveProduct) Validate(prodId int) *message.Message {
	// 验证产品线是否存在
	if ct, _ := dao.GetProductCount(dao.Query{"id=?": prodId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("产品线不存在")
	}

	// 验证产品线是否有关联存在
	if ct, _ := dao.GetProductUserCount(dao.Query{"product_id=?": prodId}); ct > 0 {
		return msg.AssociatedProductFailed
	}

	return nil
}

// SetProduct struct
type SetProduct struct {
	prodId int

	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z0-9-\u4e00-\u9fa5]+$/)"`
	Alias       string  `json:"alias" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9-]{1,}$/)"`
	Disabled    *bool   `json:"disabled" valid:"Required"`
	Description *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
}

func (sp *SetProduct) Valid(v *validation.Validation) {
	nameQ := dao.Query{
		"id<>?":                       sp.prodId,
		"(name=@name OR alias=@name)": sql.Named("name", sp.Name),
	}
	if ct, _ := dao.GetProductCount(nameQ); ct > 0 {
		_ = v.SetError("Name", "产品线名称已存在，或者该名称与其他产品线的别名相同")
	}

	aliasQ := dao.Query{
		"id<>?":                         sp.prodId,
		"(name=@alias OR alias=@alias)": sql.Named("alias", sp.Alias),
	}
	if ct, _ := dao.GetProductCount(aliasQ); ct > 0 {
		_ = v.SetError("Alias", "产品线别名重复，或者该别名与其他产品线的名称相同")
	}
}

func (sp *SetProduct) Validate(prodId int) *message.Message {
	// 验证产品线是否存在
	if ct, _ := dao.GetProductCount(dao.Query{"id=?": prodId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("产品线不存在")
	}

	sp.prodId = prodId
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(sp)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// PagedListProductsQuery struct
type PagedListProductsQuery struct {
	Query    *string `form:"query"`
	Disabled *bool   `form:"disabled"`
}

func (lpq *PagedListProductsQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if lpq.Query != nil && *lpq.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *lpq.Query)
		query["(name LIKE @query OR "+
			"alias LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if lpq.Disabled != nil {
		query["disabled=?"] = *lpq.Disabled
	}

	return nil
}

// AddUser2Product struct
type AddUser2Product struct {
	prodId int

	UserId  int  `json:"user_id" valid:"Required"`
	IsOwner bool `json:"is_owner" valid:"Required"`
}

func (apr *AddUser2Product) Valid(v *validation.Validation) {
	// 验证用户是否存在
	if ct, _ := dao.GetUserCount(dao.Query{"id=?": apr.UserId}); ct < 1 {
		_ = v.SetError("UserId", "用户不存在")
	}

	// 验证用户是否已经属于该产品线
	if ct, _ := dao.GetProductUserCount(dao.Query{"product_id=?": apr.prodId, "user_id": apr.UserId}); ct > 0 {
		_ = v.SetError("UserId", "用户已经属于该产品线")
	}
}

func (apr *AddUser2Product) Validate(prodId int) *message.Message {
	// 验证产品线是否存在
	if ct, _ := dao.GetProductCount(dao.Query{"id=?": prodId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("产品线不存在")
	}

	apr.prodId = prodId
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(apr)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// RemoveProductUser struct
type RemoveProductUser struct{}

func (*RemoveProductUser) Validate(prodId, userId int) *message.Message {
	// 验证用户是否已经属于该产品线
	if ct, _ := dao.GetProductUserCount(dao.Query{"product_id=?": prodId, "user_id": userId}); ct < 1 {
		return msg.NotFoundFailed
	}

	return nil
}

// SetUserIsProductOwner struct
type SetUserIsProductOwner struct {
	IsOwner bool `json:"is_owner" valid:"Required"`
}

func (suo *SetUserIsProductOwner) Validate(prodId, userId int) *message.Message {
	// 验证用户是否已经属于该产品线
	if ct, _ := dao.GetProductUserCount(dao.Query{"product_id=?": prodId, "user_id": userId}); ct < 1 {
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

// ListProductUsersQuery struct
type ListProductUsersQuery struct {
	IsOwner *bool `form:"is_owner"`
}

func (*ListProductUsersQuery) Validate(prodId int) *message.Message {
	// 产品线不存在
	if ct, _ := dao.GetProductCount(dao.Query{"id=?": prodId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("产品线不存在")
	}

	return nil
}

func (lpu *ListProductUsersQuery) UpdateQuery(query dao.Query) error {
	if lpu.IsOwner != nil {
		query["is_owner=?"] = *lpu.IsOwner
	}

	return nil
}
