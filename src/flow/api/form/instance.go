package form

import (
	"context"
	"database/sql"
	cMsg "eago/common/code_msg"
	"eago/common/orm"
	"eago/common/utils"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"eago/flow/dto"
	"eago/flow/model"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
	"strings"
)

// HandleInstanceForm struct 处理流程实例的数据结构
type HandleInstanceForm struct {
	CreatedBy string
	Instance  *model.Instance

	FormData *string `json:"form_data" valid:"MinSize(2)"`
	Result   *bool   `json:"result" valid:"Required"`
	Content  *string `json:"content" valid:"MinSize(0)"`
}

func (f *HandleInstanceForm) Validate(ctx context.Context, dao *dao.Dao, instId uint32, currUname string) *cMsg.CodeMsg {
	// 验证实例是否存在
	q := orm.Query{"id=?": instId, "status=?": dto.InstanceStatusRunning}
	instObj, err := dao.GetInstance(ctx, q)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetDetail("查找流程时失败")
		return m
	}

	if instObj == nil || instObj.Id < 1 {
		return cMsg.MsgNotFoundFailed.SetDetail("流程实例不存在或状态不为流转中")
	}

	f.Instance = instObj
	// 装载CreatedBy
	f.CreatedBy = currUname
	// 一般验证
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

	// 准备验证审批权限
	// 反序列化Assignees
	assignees := strings.Split(instObj.CurrentAssignees, dto.AssigneesSpiltTag)

	// 判断当前用户是否有权限审批
	res, err := utils.IsInSlice(assignees, f.CreatedBy)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetDetail("判断当前用户是否有权限审批时失败")
		return m
	}

	// 无审批权限的情况
	if !res {
		return msg.MsgHandleInstancePermDenyErr
	}

	return nil
}

type PagedListInstancesParamsForm struct {
	Query  *string `form:"query"`
	Status *int    `form:"status"`
}

func (pf *PagedListInstancesParamsForm) GenDefaultQuery(currUname string) orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(name LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if pf.Status != nil {
		query["status=?"] = *pf.Status
	}

	return query
}

func (pf *PagedListInstancesParamsForm) GenMyInstancesQuery(currUname string) orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(name LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if pf.Status != nil {
		query["status=?"] = *pf.Status
	}

	// 筛选当前用户发起的流程
	query["created_by=?"] = currUname

	return query
}

func (pf *PagedListInstancesParamsForm) GenTodoInstancesQuery(currUname string) orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(name LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if pf.Status != nil {
		query["status=?"] = *pf.Status
	}

	// 筛选当前审批人包含当前用户的流程
	likeQuery := fmt.Sprintf("%%%s%%", currUname)
	query["(current_assignees LIKE @query)"] = sql.Named("query", likeQuery)

	return query
}

func (pf *PagedListInstancesParamsForm) GenDoneInstancesQuery(currUname string) orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(name LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if pf.Status != nil {
		query["status=?"] = *pf.Status
	}

	// 已审批人包含当前用户的流程
	likeQuery := fmt.Sprintf("%%%s%%", currUname)
	query["(passed_assignees LIKE @query)"] = sql.Named("query", likeQuery)

	return query
}
