package form

import (
	"context"
	"database/sql"
	cMsg "eago/common/code_msg"
	"eago/common/orm"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

type NewTriggerForm struct {
	Name         string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z0-9-\u4e00-\u9fa5]+$/)"`
	Description  *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
	TaskCodename string  `json:"task_codename" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9._]{1,}$/)"`
	Arguments    string  `json:"arguments" valid:"Required;MinSize(2);MaxSize(4000)"`

	dao *dao.Dao
	ctx context.Context
}

func (f *NewTriggerForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsTriggerExist(f.ctx, orm.Query{"name=?": f.Name}); exist {
		_ = v.SetError("Name", "触发器名称已存在")
	}
}

func (f *NewTriggerForm) Validate(ctx context.Context, dao *dao.Dao) *cMsg.CodeMsg {
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

type RemoveTriggerForm struct{}

func (*RemoveTriggerForm) Validate(ctx context.Context, dao *dao.Dao, triId uint32) *cMsg.CodeMsg {
	// 验证触发器是否存在
	if exist, _ := dao.IsTriggerExist(ctx, orm.Query{"id=?": triId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("触发器不存在")
	}

	// 验证触发器是否有关联存在
	if ct, _ := dao.GetNodesTriggerCount(ctx, orm.Query{"trigger_id=?": triId}); ct > 0 {
		return msg.MsgAssociatedTriggerNodeFailed
	}

	return nil
}

type SetTriggerForm struct {
	Name         string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z0-9-\u4e00-\u9fa5]+$/)"`
	Description  *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
	TaskCodename string  `json:"task_codename" valid:"Required;MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9._]{1,}$/)"`
	Arguments    string  `json:"arguments" valid:"Required;MinSize(2);MaxSize(4000)"`

	triggerId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *SetTriggerForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsTriggerExist(f.ctx, orm.Query{"name=?": f.Name, "id<>?": f.triggerId}); exist {
		_ = v.SetError("Name", "触发器名称已存在")
	}
}

func (f *SetTriggerForm) Validate(ctx context.Context, dao *dao.Dao, triId uint32) *cMsg.CodeMsg {
	f.triggerId = triId

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

type PagedListTriggersParamsForm struct {
	Query *string `form:"query"`
}

func (pf *PagedListTriggersParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(name LIKE @query OR "+
			"description LIKE @query OR "+
			"task_codename LIKE @query OR "+
			"id LIKE @query OR "+
			"created_by LIKE @query OR "+
			"updated_by LIKE @query)"] = sql.Named("query", likeQuery)
	}

	return query
}

type ListTriggersNodesForm struct{}

func (*ListTriggersNodesForm) Validate(ctx context.Context, dao *dao.Dao, triId uint32) *cMsg.CodeMsg {
	// 验证触发器是否存在
	if exist, _ := dao.IsTriggerExist(ctx, orm.Query{"id=?": triId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("触发器不存在")
	}

	return nil
}
