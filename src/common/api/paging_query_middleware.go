package api

import (
	"eago/common/global"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func PagingQueryMiddleware(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery(global.GinParamPageKey, "1"))
	if err != nil {
		page = 1
	}
	c.Set(global.GinCtxPageKey, page)

	pageSize, err := strconv.Atoi(c.DefaultQuery(global.GinParamPageSizeKey, strconv.Itoa(global.DefaultPageSize)))
	if err != nil {
		pageSize = global.DefaultPageSize
	}
	// 设置PageSize
	if pageSize < 1 || pageSize > global.MaxPageSize {
		pageSize = global.DefaultPageSize
	}
	c.Set(global.GinCtxPageSizeKey, pageSize)

	c.Set(global.GinCtxQueryKey, c.Query(global.GinParamQueryKey))

	orderBy := strings.Split(c.DefaultQuery(global.GinParamOrderByKey, "id desc"), ",")
	c.Set(global.GinCtxOrderByKey, orderBy)

	c.Next()
}
