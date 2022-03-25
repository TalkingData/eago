package dto

import (
	"database/sql"
	"eago/common/log"
	"eago/common/message"
	"eago/common/utils"
	"eago/flow/conf"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"eago/flow/model"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
	"strings"
)

// HandleInstance struct 处理流程实例的数据结构
type HandleInstance struct {
	CreatedBy string
	Instance  *model.Instance

	FormData *string `json:"form_data" valid:"MinSize(2)"`
	Result   *bool   `json:"result" valid:"Required"`
	Content  *string `json:"content" valid:"MinSize(0)"`
}

// Validate
func (hi *HandleInstance) Validate(iId int, currUname string) *message.Message {
	// 验证实例是否存在
	q := dao.Query{"id=?": iId, "status=?": conf.INSTANCE_STATUS_RUNNING}
	insObj, err := dao.GetInstance(q)
	if err != nil {
		m := msg.UndefinedError.SetDetail("查找流程时失败")
		log.Error(log.Fields{
			"query": q,
			"error": err,
		}, m.String())
		return m
	}
	if insObj == nil || insObj.Id == 0 {
		return msg.NotFoundFailed.SetDetail("流程实例不存在或状态不为流转中")
	}

	hi.Instance = insObj
	// 装载CreatedBy
	hi.CreatedBy = currUname
	// 一般验证
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(hi)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	// 准备验证审批权限
	// 反序列化Assignees
	assignees := strings.Split(insObj.CurrentAssignees, conf.ASSIGNEES_SPILT_TAG)

	// 判断当前用户是否有权限审批
	res, err := utils.IsInSlice(assignees, hi.CreatedBy)
	if err != nil {
		m := msg.UndefinedError.SetDetail("判断当前用户是否有权限审批时失败")
		log.Error(log.Fields{
			"instance_id": insObj.Id,
			"username":    hi.CreatedBy,
			"error":       err,
		}, m.String())
		return m
	}

	// 无审批权限的情况
	if !res {
		return msg.HandleInstancePermDenyError
	}

	return nil
}

// ListInstancesQuery struct
type ListInstancesQuery struct {
	Query  *string `form:"query"`
	Status *int    `form:"status"`
}

// UpdateDefaultQuery
func (q *ListInstancesQuery) UpdateDefaultQuery(query dao.Query, currUname string) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["(name LIKE @query OR id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if q.Status != nil {
		query["status=?"] = *q.Status
	}

	return nil
}

// UpdateMyInstancesQuery
func (q *ListInstancesQuery) UpdateMyInstancesQuery(query dao.Query, currUname string) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["(name LIKE @query OR id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if q.Status != nil {
		query["status=?"] = *q.Status
	}

	// 筛选当前用户发起的流程
	query["created_by=?"] = currUname

	return nil
}

// UpdateTodoInstancesQuery
func (q *ListInstancesQuery) UpdateTodoInstancesQuery(query dao.Query, currUname string) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["(name LIKE @query OR id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if q.Status != nil {
		query["status=?"] = *q.Status
	}

	// 筛选当前审批人包含当前用户的流程
	likeQuery := fmt.Sprintf("%%%s%%", currUname)
	query["(current_assignees LIKE @query)"] = sql.Named("query", likeQuery)

	return nil
}

// UpdateDoneInstancesQuery
func (q *ListInstancesQuery) UpdateDoneInstancesQuery(query dao.Query, currUname string) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["(name LIKE @query OR id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if q.Status != nil {
		query["status=?"] = *q.Status
	}

	// 已审批人包含当前用户的流程
	likeQuery := fmt.Sprintf("%%%s%%", currUname)
	query["(passed_assignees LIKE @query)"] = sql.Named("query", likeQuery)

	return nil
}
