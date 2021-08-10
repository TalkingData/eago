package middleware

import (
	"context"
	auth "eago/auth/srv/proto"
	"eago/common/log"
	"eago/common/utils"
	"eago/task/cli"
	"eago/task/conf/msg"
	"github.com/gin-gonic/gin"
)

// MustLogin 验证是否登录并装载TokenContent
func MustLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tk auth.Token

		tk.Token = c.GetHeader("Token")
		tc, err := cli.AuthClient.GetTokenContent(context.Background(), &tk)
		if err != nil {
			resp := msg.ErrGetToken.GenResponse()
			log.WarnWithFields(log.Fields{
				"token": tk.Token,
				"error": err.Error(),
			}, resp.String())
			resp.WriteAndAbort(c)
			return
		} else if !tc.Ok {
			resp := msg.WarnPermissionDeny.GenResponse("Invalid token or Not login yet.")
			log.WarnWithFields(log.Fields{
				"token": tk.Token,
			}, resp.String())
			resp.WriteAndAbort(c)
			return
		}

		c.Set("TokenContent", map[string]interface{}{
			"UserId":      tc.UserId,
			"Username":    tc.Username,
			"Phone":       tc.Phone,
			"IsSuperuser": tc.IsSuperuser,

			"Department":  tc.Department,
			"Roles":       tc.Roles,
			"Products":    tc.Products,
			"OwnProducts": tc.OwnProducts,
			"Groups":      tc.Groups,
			"OwnGroups":   tc.OwnGroups,
		})
	}
}

// MustRole 检测当前用户是指定角色
func MustRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isRoleHandler(c, role)
	}
}

// isRoleHandler 是否是指定角色处理
func isRoleHandler(c *gin.Context, role string) {
	tc := c.GetStringMap("TokenContent")
	roles := tc["Roles"].([]string)

	// 超级用户拥有所有权限，跳过判断
	if tc["IsSuperuser"].(bool) {
		return
	}

	// 判断是否是对应角色
	res, err := utils.IsInSlice(roles, role)
	if err != nil {
		resp := msg.ErrUnknown.GenResponse()
		log.ErrorWithFields(log.Fields{
			"user_id":       tc["UserId"].(int32),
			"username":      tc["Username"].(string),
			"is_superuser":  tc["IsSuperuser"].(bool),
			"roles":         tc["Roles"].([]string),
			"role_required": role,
			"error":         err.Error(),
		}, resp.String())
		resp.WriteAndAbort(c)
		return
	} else if !res {
		resp := msg.WarnPermissionDeny.GenResponse("Role: '" + role + "' required.")
		log.WarnWithFields(log.Fields{
			"user_id":       tc["UserId"].(int32),
			"username":      tc["Username"].(string),
			"is_superuser":  tc["IsSuperuser"].(bool),
			"roles":         tc["Roles"].([]string),
			"role_required": role,
		}, resp.String())
		resp.WriteAndAbort(c)
		return
	}

	return
}
