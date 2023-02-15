package form

import (
	"database/sql"
	"eago/common/orm"
	"fmt"
)

type ListResultsParamsForm struct {
	Query  *string `form:"query"`
	Status *int32  `form:"status"`
}

func (pf *ListResultsParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(task_codename LIKE @query OR "+
			"caller LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if pf.Status != nil {
		query["status=?"] = *pf.Status
	}

	return query
}
