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

type NewCategoryForm struct {
	Name string `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`

	dao *dao.Dao
	ctx context.Context
}

func (f *NewCategoryForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsCategoryExist(f.ctx, orm.Query{"name=?": f.Name}); exist {
		_ = v.SetError("Name", "类别名称已存在")
	}
}

func (f *NewCategoryForm) Validate(ctx context.Context, dao *dao.Dao) *cMsg.CodeMsg {
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

type RemoveCategoryForm struct{}

func (*RemoveCategoryForm) Validate(ctx context.Context, dao *dao.Dao, catId uint32) *cMsg.CodeMsg {
	// 验证类别是否存在
	if exist, _ := dao.IsCategoryExist(ctx, orm.Query{"id=?": catId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("类别不存在")
	}

	// 验证类别与流程关联
	if exist, _ := dao.IsFlowExist(ctx, orm.Query{"categories_id=?": catId}); exist {
		return msg.MsgAssociatedCategoryFlowFailed
	}

	return nil
}

type SetCategoryForm struct {
	Name string `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`

	categoryId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *SetCategoryForm) Valid(v *validation.Validation) {
	if exist, _ := f.dao.IsCategoryExist(f.ctx, orm.Query{"name=?": f.Name, "id<>?": f.categoryId}); exist {
		_ = v.SetError("Name", "类别名称已存在")
	}
}

func (f *SetCategoryForm) Validate(ctx context.Context, dao *dao.Dao, catId uint32) *cMsg.CodeMsg {
	f.ctx = ctx
	f.dao = dao

	f.categoryId = catId

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

type ListCategoriesParamsForm struct {
	Query *string `form:"query"`
}

func (pf *ListCategoriesParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(name LIKE @query OR "+
			"id LIKE @query OR "+
			"created_by LIKE @query OR "+
			"updated_by LIKE @query)"] = sql.Named("query", likeQuery)
	}

	return query
}

type ListCategoriesRelations struct {
	Disabled *bool `form:"disabled"`
}

func (*ListCategoriesRelations) Validate(ctx context.Context, dao *dao.Dao, catId uint32) *cMsg.CodeMsg {
	// 验证类别是否存在
	if exist, _ := dao.IsCategoryExist(ctx, orm.Query{"id=?": catId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("类别不存在")
	}

	return nil
}

func (pf *ListCategoriesRelations) GenQuery() orm.Query {
	// TODO: 检测调用这里的方法是否正确
	query := orm.Query{}

	if pf.Disabled != nil {
		query["flows.disabled=?"] = *pf.Disabled
	}

	return query
}
