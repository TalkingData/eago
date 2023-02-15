package form

import (
	"context"
	"database/sql"
	cMsg "eago/common/code_msg"
	"eago/common/orm"
	"eago/task/conf/msg"
	"eago/task/dao"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

type CallTask struct {
	TaskCodeName string

	Timeout   *int64 `json:"timeout" valid:"Range(0,86400000)"`
	Arguments string `json:"arguments" valid:"Required;MinSize(2)"`
}

func (ct *CallTask) Validate(ctx context.Context, dao *dao.Dao, tId uint32) *cMsg.CodeMsg {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(ct)
	if err != nil {
		return cMsg.MsgValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return cMsg.MsgValidateFailed.SetError(valid.Errors)
	}

	// 验证任务是否存在
	taskObj, err := dao.GetTask(ctx, orm.Query{"id=?": tId, "disabled=?": false})
	if err != nil {
		return cMsg.MsgNotFoundFailed.SetDetail("任务不存在")
	}

	// 没找到任务或者任务是禁用状态
	if taskObj == nil || taskObj.Id < 1 {
		return cMsg.MsgNotFoundFailed.SetDetail("任务不存在或任务不是可用的状态")
	}

	ct.TaskCodeName = taskObj.Codename

	return nil
}

type NewTaskForm struct {
	Category     *int32  `json:"category" valid:"Min(0)"`
	Codename     string  `json:"codename" valid:"Required;MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9._]{1,}$/)"`
	FormalParams string  `json:"formal_params" valid:"Required;MinSize(2)"`
	Disabled     *bool   `json:"disabled" gorm:"default:0" valid:"Required"`
	Description  *string `json:"description" valid:"MinSize(0);MaxSize(500)"`

	dao *dao.Dao
	ctx context.Context
}

func (f *NewTaskForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsTaskExist(f.ctx, orm.Query{"codename=?": f.Codename}); exist {
		_ = v.SetError("Codename", "已有相同代号的任务存在")
	}
}

func (f *NewTaskForm) Validate(ctx context.Context, dao *dao.Dao) *cMsg.CodeMsg {
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

type RemoveTaskForm struct{}

func (*RemoveTaskForm) Validate(ctx context.Context, dao *dao.Dao, tId uint32) *cMsg.CodeMsg {
	// 验证任务是否存在
	taskObj, err := dao.GetTask(ctx, orm.Query{"id=?": tId})
	if err != nil {
		return cMsg.MsgNotFoundFailed.SetDetail("任务不存在")
	}

	if taskObj == nil || taskObj.Id < 1 {
		return cMsg.MsgNotFoundFailed.SetDetail("任务不存在")
	}

	// 验证计划任务是否存在
	if exist, _ := dao.IsScheduleExist(ctx, orm.Query{"codename=?": taskObj.Codename}); exist {
		return msg.MsgAssociatedScheduleFailed
	}

	return nil
}

type SetTaskForm struct {
	Category     *int32  `json:"category" valid:"Min(0)"`
	Codename     string  `json:"codename" valid:"Required;MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9._]{1,}$/)"`
	FormalParams string  `json:"formal_params" valid:"Required;MinSize(2)"`
	Disabled     *bool   `json:"disabled" gorm:"default:0" valid:"Required"`
	Description  *string `json:"description" valid:"MinSize(0);MaxSize(500)"`

	taskId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *SetTaskForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsTaskExist(f.ctx, orm.Query{"codename=?": f.Codename, "id<>?": f.taskId}); exist {
		_ = v.SetError("Codename", "已有相同代号的任务存在")
	}
}

func (f *SetTaskForm) Validate(ctx context.Context, dao *dao.Dao, taskId uint32) *cMsg.CodeMsg {
	// 验证角色是否存在
	if exist, _ := dao.IsTaskExist(ctx, orm.Query{"id=?": taskId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("任务不存在")
	}

	f.taskId = taskId

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

type PagedListTasksParamsForm struct {
	Query    *string `form:"query"`
	Disabled *bool   `form:"disabled"`
}

func (pf *PagedListTasksParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(codename LIKE @query OR "+
			"description LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if pf.Disabled != nil {
		query["disabled=?"] = *pf.Disabled
	}

	return query
}
