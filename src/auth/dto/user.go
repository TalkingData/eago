package dto

import (
	"database/sql"
	"eago/auth/conf/msg"
	"eago/auth/dao"
	"eago/common/message"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

// SetUser struct
type SetUser struct {
	Email string `json:"email" valid:"Required;Email;MinSize(3);MaxSize(100)"`
	Phone string `json:"phone" valid:"Required;Phone;MinSize(8);MaxSize(20)"`
}

// Validate
func (su *SetUser) Validate(userId int) *message.Message {
	// 户不存在
	if ct, _ := dao.GetUserCount(dao.Query{"id=?": userId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("用户不存在")
	}

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(su)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// MakeUserHandover struct
type MakeUserHandover struct{}

// Validate
func (*MakeUserHandover) Validate(frmUserId, tgtUserId int) *message.Message {
	// 原用户不存在
	if ct, _ := dao.GetUserCount(dao.Query{"id=?": frmUserId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("原用户不存在")
	}

	// 目标用户不存在
	if ct, _ := dao.GetUserCount(dao.Query{"id=?": tgtUserId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("目标用户不存在")

	}
	return nil
}

// ListUsersQuery struct
type ListUsersQuery struct {
	Query       *string `form:"query"`
	IsSuperuser *bool   `form:"is_superuser"`
	Disabled    *bool   `form:"disabled"`
}

// UpdateQuery
func (luq *ListUsersQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if luq.Query != nil && *luq.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *luq.Query)
		query["(username LIKE @query OR id LIKE @query OR email LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if luq.Disabled != nil {
		query["disabled=?"] = luq.Disabled
	}
	if luq.IsSuperuser != nil {
		query["is_superuser=?"] = luq.IsSuperuser
	}

	return nil
}
