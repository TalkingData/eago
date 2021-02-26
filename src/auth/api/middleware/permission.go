package middleware

import (
	"eago-auth/config/msg"
	"eago-auth/srv"
	"eago-common/log"
	"eago-common/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func MustLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		tc, suc := srv.GetTokenContent(token)
		if !suc {
			m := msg.ErrGenToken.NewMsg()
			log.WarnWithFields(log.Fields{
				"token": token,
			}, m.String())
			c.AbortWithStatusJSON(http.StatusOK, m.GinH())
			return
		} else if tc == nil {
			m := msg.WarnPermissionDeny.NewMsg("Not login yet.")
			log.WarnWithFields(log.Fields{
				"token": token,
			}, m.String())
			c.AbortWithStatusJSON(http.StatusOK, m.GinH())
			return
		}

		c.Set("TokenContent", map[string]interface{}{
			"UserId":      tc.UserId,
			"Username":    tc.Username,
			"IsSuperuser": tc.IsSuperuser,
			"Roles":       tc.Roles,
			"Products":    tc.Products,
			"OwnProducts": tc.OwnProducts,
			"Groups":      tc.Groups,
			"OwnGroups":   tc.OwnGroups,
		})
		c.Next()
	}
}

// 检测当前用户是指定角色
func MustRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isRoleHandler(c, role) {
			return
		}
		c.Next()
	}
}

// 检测指定字段Id=当前用户Id或当前用户是指定角色
func MustCurrUserOrRole(userIdField string, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := isCurrUserHandler(c, userIdField)
		if err != nil {
			m := msg.WarnInvalidUri.NewMsg("Field '" + userIdField + "' required.")
			log.WarnWithFields(log.Fields{
				"error": err.Error(),
			}, m.String())
			c.AbortWithStatusJSON(http.StatusOK, m.GinH())
			return
		} else if res || isRoleHandler(c, role) {
			c.Next()
			return
		}
		return
	}
}

// 指定字段Id=当前用户Id处理
func isCurrUserHandler(c *gin.Context, userIdField string) (bool, error) {
	tc := c.GetStringMap("TokenContent")
	currUserId := tc["UserId"].(int)

	userId, err := strconv.Atoi(c.Param(userIdField))
	if err != nil {
		return false, err
	} else if currUserId == userId {
		return true, nil
	}

	return false, nil
}

// 是否是指定角色处理
func isRoleHandler(c *gin.Context, role string) bool {
	tc := c.GetStringMap("TokenContent")
	roles := tc["Roles"].(*[]string)

	// 超级用户拥有所有权限，跳过判断
	if tc["IsSuperuser"].(bool) {
		c.Next()
		return true
	}

	// 判断是否是对应角色
	res, err := tools.IsInSlice(*roles, role)
	if err != nil {
		m := msg.ErrUnknown.NewMsg()
		log.ErrorWithFields(log.Fields{
			"user_id":       tc["UserId"].(int),
			"username":      tc["Username"].(string),
			"is_superuser":  tc["IsSuperuser"].(bool),
			"roles":         tc["Roles"].(*[]string),
			"role_required": role,
			"error":         err.Error(),
		}, m.String())
		c.AbortWithStatusJSON(http.StatusOK, m.GinH())
		return false
	} else if !res {
		m := msg.WarnPermissionDeny.NewMsg("Role: '" + role + "' required.")
		log.WarnWithFields(log.Fields{
			"user_id":       tc["UserId"].(int),
			"username":      tc["Username"].(string),
			"is_superuser":  tc["IsSuperuser"].(bool),
			"roles":         tc["Roles"].(*[]string),
			"role_required": role,
		}, m.String())
		c.AbortWithStatusJSON(http.StatusOK, m.GinH())
		return false
	}

	return true
}
