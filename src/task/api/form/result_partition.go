package form

import (
	"database/sql"
	cMsg "eago/common/code_msg"
	"eago/common/orm"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

type NewResultPartitionForm struct {
	Partition string `json:"partition" valid:"Required;MinSize(2);MaxSize(10)"`
}

func (f *NewResultPartitionForm) Validate() *cMsg.CodeMsg {
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

type ListResultPartitionsParamsForm struct {
	Query *string `form:"query"`
}

func (pf *ListResultPartitionsParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["partition LIKE @query"] = sql.Named("query", likeQuery)
	}

	return query
}
