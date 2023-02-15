package form

import (
	"context"
	"database/sql"
	cMsg "eago/common/code_msg"
	"eago/common/orm"
	"eago/task/dao"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

type NewScheduleForm struct {
	TaskCodename string  `json:"task_codename" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9._]{1,}$/)"`
	Expression   string  `json:"expression" valid:"Required"`
	Timeout      *int64  `json:"timeout" valid:"Range(0,86400000)"`
	Arguments    string  `json:"arguments" valid:"Required;MinSize(2)"`
	Disabled     *bool   `json:"disabled" gorm:"default:0" valid:"Required"`
	Description  *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
}

func (f *NewScheduleForm) Validate() *cMsg.CodeMsg {
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

type RemoveScheduleForm struct{}

func (*RemoveScheduleForm) Validate(ctx context.Context, dao *dao.Dao, schId uint32) *cMsg.CodeMsg {
	// 验证任务是否存在
	if exist, _ := dao.IsScheduleExist(ctx, orm.Query{"id=?": schId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("计划任务不存在")
	}

	return nil
}

type SetScheduleForm struct {
	TaskCodename string  `json:"task_codename" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9._]{1,}$/)"`
	Expression   string  `json:"expression" valid:"Required"`
	Timeout      *int64  `json:"timeout" valid:"Range(0,86400000)"`
	Arguments    string  `json:"arguments" valid:"Required;MinSize(2)"`
	Disabled     *bool   `json:"disabled" gorm:"default:0" valid:"Required"`
	Description  *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
}

func (f *SetScheduleForm) Validate(ctx context.Context, dao *dao.Dao, schId uint32) *cMsg.CodeMsg {
	// 验证任务是否存在
	if exist, _ := dao.IsScheduleExist(ctx, orm.Query{"id=?": schId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("计划任务不存在")
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

type ListSchedulesParamsForm struct {
	Query    *string `form:"query"`
	Disabled *bool   `form:"disabled"`
}

func (pf *ListSchedulesParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(task_codename LIKE @query OR "+
			"description LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if pf.Disabled != nil {
		query["disabled=?"] = *pf.Disabled
	}

	return query
}
