package pagination

import (
	"eago-common/tools"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

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
		c.Set("PageSize", tools.IntMin(pageSize, MAX_PAGE_SIZE))

		c.Set("Query", c.Query("query"))
		c.Set("OrderBy", strings.Split(c.Query("order_by"), ","))

		c.Next()
	}
}
