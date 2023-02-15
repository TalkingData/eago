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

type InstantiateFlowForm struct {
	FormId   uint32
	Name     string
	FormData *string `json:"form_data" valid:"MinSize(2)"`

	flowId uint32
}

func (f *InstantiateFlowForm) Validate(ctx context.Context, dao *dao.Dao, flowId uint32) *cMsg.CodeMsg {
	flow, err := dao.GetFlow(ctx, orm.Query{"id=?": flowId})
	if err != nil {
		return msg.MsgFlowDaoErr
	}

	if flow == nil || flow.Id < 1 {
		return cMsg.MsgNotFoundFailed.SetDetail("流程不存在")
	}

	if *flow.Disabled {
		return cMsg.MsgNotFoundFailed.SetDetail("无法发起一个禁用的流程")
	}

	f.Name = flow.Name
	f.FormId = flow.Id

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

type NewFlowForm struct {
	Name          string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	InstanceTitle string  `json:"instance_title" valid:"Required;MinSize(3);MaxSize(200)"`
	CategoriesId  *uint32 `json:"categories_id"`
	Disabled      *bool   `json:"disabled" valid:"Required"`
	Description   *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
	FormId        uint32  `json:"form_id" valid:"Required"`
	FirstNodeId   uint32  `json:"first_node_id" valid:"Required"`

	dao *dao.Dao
	ctx context.Context
}

func (f *NewFlowForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsFlowExist(f.ctx, orm.Query{"name=?": f.Name}); exist {
		_ = v.SetError("Name", "流程名称已存在")
	}

	if exist, _ := f.dao.IsFormExist(f.ctx, orm.Query{"id=?": f.FormId, "disabled=?": false}); !exist {
		_ = v.SetError("FormId", "找不到所选表单，请确定该表单存在并且不是禁用状态")
	}
}

func (f *NewFlowForm) Validate(ctx context.Context, dao *dao.Dao) *cMsg.CodeMsg {
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

type RemoveFlowForm struct{}

func (*RemoveFlowForm) Validate(ctx context.Context, dao *dao.Dao, flowId uint32) *cMsg.CodeMsg {
	// 验证流程是否存在
	if exist, _ := dao.IsFlowExist(ctx, orm.Query{"id=?": flowId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("流程不存在")
	}

	return nil
}

type SetFlowForm struct {
	Name          string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	InstanceTitle string  `json:"instance_title" valid:"Required;MinSize(3);MaxSize(200)"`
	CategoriesId  *uint32 `json:"categories_id"`
	Disabled      *bool   `json:"disabled" valid:"Required"`
	Description   *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
	FormId        uint32  `json:"form_id" valid:"Required"`
	FirstNodeId   uint32  `json:"first_node_id" valid:"Required"`

	flowId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *SetFlowForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsFlowExist(f.ctx, orm.Query{"name=?": f.Name, "id<>?": f.flowId}); exist {
		_ = v.SetError("Name", "流程名称已存在")
	}

	if exist, _ := f.dao.IsFormExist(f.ctx, orm.Query{"id=?": f.FormId, "disabled=?": false}); !exist {
		_ = v.SetError("FormId", "找不到所选表单，请确定该表单存在并且不是禁用状态")
	}
}

func (f *SetFlowForm) Validate(ctx context.Context, dao *dao.Dao, flowId uint32) *cMsg.CodeMsg {
	f.ctx = ctx
	f.dao = dao

	f.flowId = flowId

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

type PagedListFlowsParamsForm struct {
	Query        *string `form:"query"`
	Disabled     *bool   `form:"disabled"`
	CategoriesId *int    `form:"categories_id"`
}

func (pf *PagedListFlowsParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(flows.id LIKE @query OR "+
			"flows.name LIKE @query OR "+
			"flows.description LIKE @query OR "+
			"flows.created_by LIKE @query OR "+
			"flows.updated_by LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if pf.Disabled != nil {
		query["flows.disabled=?"] = *pf.Disabled
	}

	if pf.CategoriesId != nil {
		query["flows.categories_id=?"] = *pf.CategoriesId
	}

	return query
}
