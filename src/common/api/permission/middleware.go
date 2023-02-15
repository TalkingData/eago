package permission

import (
	authpb "eago/auth/proto"
	"eago/common/api/ext"
	cMsg "eago/common/code_msg"
	"github.com/gin-gonic/gin"
	"strings"
)

// MustRole 检测当前用户是指定角色
func MustRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isRoleHandler(c, role)
	}
}

// MustRoleIn 检测当前用户属于指定角色数组中的任意义一个
func MustRoleIn(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsRoleIn(c, roles)
		if err != nil {
			m := cMsg.MsgCheckRoleErr.SetError(err)
			ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
			return
		}

		if !ok {
			m := cMsg.MsgUserNotRoleFailed.SetDetail(strings.Join(roles, ","))
			ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
			return
		}
	}
}

// MustCurrUserOrRole 检测指定字段Id=当前用户Id，或当前用户是指定角色
func MustCurrUserOrRole(userIdField, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsCurrUser(c, userIdField)
		if err != nil {
			m := cMsg.MsgInvalidUriFailed.SetError(err)
			ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
			return
		}
		if ok {
			return
		}

		isRoleHandler(c, role)
	}
}

// MustCurrUserInProductOrRole 检测当前用户必须在指定产品线内，或当前用户是指定角色
func MustCurrUserInProductOrRole(prodIdField, role string, isOwner bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsCurrUserInProduct(c, prodIdField, isOwner)
		if err != nil {
			m := cMsg.MsgInvalidUriFailed.SetError(err)
			ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
			return
		}
		if ok {
			return
		}

		isRoleHandler(c, role)
	}
}

// MustCurrUserInGroupOrRole 检测当前用户必须在指定组内，或当前用户是指定角色
func MustCurrUserInGroupOrRole(groupIdField, role string, isOwner bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsCurrUserInGroup(c, groupIdField, isOwner)
		if err != nil {
			m := cMsg.MsgInvalidUriFailed.SetError(err)
			ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
			return
		}
		if ok {
			return
		}

		isRoleHandler(c, role)
	}
}

// MustLogin 验证是否登录并装载TokenContent
func MustLogin(authCli authpb.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		getTokenContentHandler(c, c.GetHeader(defaultTokenKey), authCli)
	}
}

// MustLoginWs WebSocket验证是否登录并装载TokenContent
func MustLoginWs(authCli authpb.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		getTokenContentHandler(c, c.DefaultQuery(defaultTokenKey, ""), authCli)
	}
}
