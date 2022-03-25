package pagination

import (
	"eago/common/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// ListPageHelper 分页助手中间件
func ListPageHelper() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil {
			page = 1
		}
		c.Set("Page", page)

		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(_DEFAULT_PAGE_SIZE)))
		if err != nil {
			pageSize = _DEFAULT_PAGE_SIZE
		}
		c.Set("PageSize", utils.IntMin(pageSize, _MAX_PAGE_SIZE))

		c.Set("Query", c.Query("query"))

		orderBy := strings.Split(c.DefaultQuery("order_by", "id desc"), ",")
		c.Set("OrderBy", orderBy)

		c.Next()
	}
}
