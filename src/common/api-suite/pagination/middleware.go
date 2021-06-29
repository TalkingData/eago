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

		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(DEFAULT_PAGE_SIZE)))
		if err != nil {
			pageSize = DEFAULT_PAGE_SIZE
		}
		c.Set("PageSize", utils.IntMin(pageSize, MAX_PAGE_SIZE))

		c.Set("Query", c.Query("query"))

		orderBy := strings.Split(c.DefaultQuery("order_by", "-id"), ",")
		c.Set("OrderBy", orderBy)

		c.Next()
	}
}
