package dto

import (
	"database/sql"
	"eago/task/dao"
	"fmt"
)

// ListResultsQuery struct
type ListResultsQuery struct {
	Query  *string `form:"query"`
	Status *int    `form:"status"`
}

// UpdateQuery
func (q *ListResultsQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["(task_codename LIKE @query OR caller LIKE @query OR id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if q.Status != nil {
		query["status=?"] = *q.Status
	}

	return nil
}
