package dto

import (
	"database/sql"
	"eago/common/message"
	"eago/task/conf/msg"
	"eago/task/dao"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

// NewResultPartition struct
type NewResultPartition struct {
	Partition string `json:"partition" valid:"Required;MinSize(2);MaxSize(10)"`
}

func (rp *NewResultPartition) Validate() *message.Message {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(rp)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// ListResultPartitionsQuery struct
type ListResultPartitionsQuery struct {
	Query *string `form:"query"`
}

func (q *ListResultPartitionsQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["partition LIKE @query"] = sql.Named("query", likeQuery)
	}

	return nil
}
