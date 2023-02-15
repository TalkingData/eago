package form

import (
	"context"
	"database/sql"
	"eago/auth/dao"
	cMsg "eago/common/code_msg"
	"eago/common/orm"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

type SetUserForm struct {
	Email string `json:"email" valid:"Required;Email;MinSize(3);MaxSize(100)"`
	Phone string `json:"phone" valid:"Required;Phone;MinSize(8);MaxSize(20)"`
}

func (f *SetUserForm) Validate(ctx context.Context, dao *dao.Dao, userId uint32) *cMsg.CodeMsg {
	// 户不存在
	if ct, _ := dao.GetUserCount(ctx, orm.Query{"id=?": userId}); ct < 1 {
		return cMsg.MsgNotFoundFailed.SetDetail("用户不存在")
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

type PagedListUsersParamsForm struct {
	Query       *string `form:"query"`
	IsSuperuser *bool   `form:"is_superuser"`
	Disabled    *bool   `form:"disabled"`
}

func (pf *PagedListUsersParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(username LIKE @query OR "+
			"id LIKE @query OR "+
			"email LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if pf.Disabled != nil {
		query["disabled=?"] = pf.Disabled
	}
	if pf.IsSuperuser != nil {
		query["is_superuser=?"] = pf.IsSuperuser
	}

	return query
}

type MakeUserHandoverForm struct{}

func (*MakeUserHandoverForm) Validate(ctx context.Context, dao *dao.Dao, frmUserId, tgtUserId uint32) *cMsg.CodeMsg {
	// 原用户不存在
	if ct, _ := dao.GetUserCount(ctx, orm.Query{"id=?": frmUserId}); ct < 1 {
		return cMsg.MsgNotFoundFailed.SetDetail("原用户不存在")
	}

	// 目标用户不存在
	if ct, _ := dao.GetUserCount(ctx, orm.Query{"id=?": tgtUserId}); ct < 1 {
		return cMsg.MsgNotFoundFailed.SetDetail("目标用户不存在")

	}
	return nil
}
