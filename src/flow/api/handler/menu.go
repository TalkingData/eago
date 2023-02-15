package handler

import (
	"eago/common/api/ext"
	"github.com/gin-gonic/gin"
)

// ListMenus 根据当前登录用户权限列出菜单
func (h *FlowHandler) ListMenus(c *gin.Context) {
	ext.WriteSuccessPayload(c, "menus", h.menu.ListMenusByContext(c))
}
