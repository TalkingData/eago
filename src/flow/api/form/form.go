package form

import (
	"context"
	"database/sql"
	cMsg "eago/common/code_msg"
	"eago/common/orm"
	"eago/flow/dao"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

type NewFormForm struct {
	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	Disabled    *bool   `json:"disabled" valid:"Required"`
	Description *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
	Body        *string `json:"body" valid:"MinSize(2)"`

	dao *dao.Dao
	ctx context.Context
}

func (f *NewFormForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsFlowExist(f.ctx, orm.Query{"name=?": f.Name}); exist {
		_ = v.SetError("Name", "表单名称已存在")
	}
}

func (f *NewFormForm) Validate(ctx context.Context, dao *dao.Dao) *cMsg.CodeMsg {
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

type SetFormForm struct {
	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	Disabled    *bool   `json:"disabled" valid:"Required"`
	Description *string `json:"description" valid:"MinSize(0);MaxSize(500)"`

	formId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *SetFormForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsFormExist(f.ctx, orm.Query{"name=?": f.Name, "id<>?": f.formId}); exist {
		_ = v.SetError("Name", "表单名称已存在")
	}
}

func (f *SetFormForm) Validate(ctx context.Context, dao *dao.Dao, frmId uint32) *cMsg.CodeMsg {
	f.ctx = ctx
	f.dao = dao

	f.formId = frmId

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

type GetFormForm struct{}

func (*GetFormForm) Validate(ctx context.Context, dao *dao.Dao, frmId uint32) *cMsg.CodeMsg {
	// 验证表单是否存在
	if exist, _ := dao.IsFormExist(ctx, orm.Query{"id=?": frmId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("表单不存在")
	}

	return nil
}

type PagedListFormsParamsForm struct {
	Query    *string `form:"query"`
	Disabled *bool   `form:"disabled"`
}

func (pf *PagedListFormsParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(id LIKE @query OR "+
			"name LIKE @query OR "+
			"description LIKE @query OR "+
			"created_by LIKE @query OR "+
			"updated_by LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if pf.Disabled != nil {
		query["disabled=?"] = *pf.Disabled
	}

	return query
}

type ListFormRelationsForm struct{}

func (*ListFormRelationsForm) Validate(ctx context.Context, dao *dao.Dao, frmId uint32) *cMsg.CodeMsg {
	// 验证表单是否存在
	if exist, _ := dao.IsFormExist(ctx, orm.Query{"id=?": frmId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("表单不存在")
	}

	return nil
}
