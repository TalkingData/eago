package middleware

import (
	"eago/auth/conf/msg"
	"eago/auth/srv/dto"
	"eago/auth/srv/local"
	"eago/common/log"
	"eago/common/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// MustLogin 验证是否登录并装载TokenContent
func MustLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		tc, ok := local.GetTokenContent(token)
		if !ok {
			resp := msg.ErrGetToken.GenResponse()
			log.WarnWithFields(log.Fields{
				"token": token,
			}, resp.String())
			resp.WriteAndAbort(c)
			return
		} else if tc == nil {
			resp := msg.WarnPermissionDeny.GenResponse("Not login yet.")
			log.WarnWithFields(log.Fields{
				"token": token,
			}, resp.String())
			resp.WriteAndAbort(c)
			return
		}

		c.Set("TokenContent", map[string]interface{}{
			"UserId":   tc.UserId,
			"Username": tc.Username,
			"Phone":    tc.Phone,

			"UserIsSuperuser": tc.IsSuperuser,

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

// MustCurrUserOrRole 检测指定字段Id=当前用户Id，或当前用户是指定角色
func MustCurrUserOrRole(userIdField, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := isCurrUserHandler(c, userIdField)
		if err != nil {
			resp := msg.WarnInvalidUri.GenResponse("Field '" + userIdField + "' required.")
			log.WarnWithFields(log.Fields{
				"error": err.Error(),
			}, resp.String())
			resp.WriteAndAbort(c)
			return
		} else if res {
			return
		}

		isRoleHandler(c, role)
	}
}

// MustCurrUserInProductOrRole 检测当前用户必须在指定产品线内，或当前用户是指定角色
func MustCurrUserInProductOrRole(prodIdField, role string, isOwner bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := isCurrUserInProductHandler(c, prodIdField, isOwner)
		if err != nil {
			resp := msg.WarnInvalidUri.GenResponse("Field '" + prodIdField + "' required.")
			log.WarnWithFields(log.Fields{
				"error": err.Error(),
			}, resp.String())
			resp.WriteAndAbort(c)
			return
		} else if res {
			return
		}

		isRoleHandler(c, role)
	}
}

// MustCurrUserInGroupOrRole 检测当前用户必须在指定组内，或当前用户是指定角色
func MustCurrUserInGroupOrRole(groupIdField, role string, isOwner bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := isCurrUserInGroupHandler(c, groupIdField, isOwner)
		if err != nil {
			resp := msg.WarnInvalidUri.GenResponse("Field '" + groupIdField + "' required.")
			log.WarnWithFields(log.Fields{
				"error": err.Error(),
			}, resp.String())
			resp.WriteAndAbort(c)
			return
		} else if res {
			return
		}

		isRoleHandler(c, role)
	}
}

// isCurrUserHandler 指定字段Id=当前用户Id处理
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

// isCurrUserInProductHandler 判断当前用户在指定产品线内处理
func isCurrUserInProductHandler(c *gin.Context, prodIdField string, isOwner bool) (bool, error) {
	var prods *[]dto.ProductInToken

	tc := c.GetStringMap("TokenContent")
	if isOwner {
		prods = tc["OwnProducts"].(*[]dto.ProductInToken)
	} else {
		prods = tc["Products"].(*[]dto.ProductInToken)
	}

	prodId, err := strconv.Atoi(c.Param(prodIdField))
	if err != nil {
		return false, err
	}
	for _, p := range *prods {
		if p.Id == prodId {
			return true, nil
		}
	}

	return false, nil
}

// isCurrUserInGroupHandler 判断当前用户在指定组内处理
func isCurrUserInGroupHandler(c *gin.Context, groupIdField string, isOwner bool) (bool, error) {
	var groups *[]dto.GroupInToken

	tc := c.GetStringMap("TokenContent")
	if isOwner {
		groups = tc["OwnGroups"].(*[]dto.GroupInToken)
	} else {
		groups = tc["Groups"].(*[]dto.GroupInToken)
	}

	groupId, err := strconv.Atoi(c.Param(groupIdField))
	if err != nil {
		return false, err
	}
	for _, g := range *groups {
		if g.Id == groupId {
			return true, nil
		}
	}

	return false, nil
}

// isRoleHandler 是否是指定角色处理
func isRoleHandler(c *gin.Context, role string) {
	tc := c.GetStringMap("TokenContent")
	roles := tc["Roles"].(*[]string)

	// 超级用户拥有所有权限，跳过判断
	if tc["UserIsSuperuser"].(bool) {
		return
	}

	// 判断是否是对应角色
	res, err := utils.IsInSlice(*roles, role)
	if err != nil {
		resp := msg.ErrUnknown.GenResponse()
		log.ErrorWithFields(log.Fields{
			"user_id":       tc["UserId"].(int),
			"username":      tc["Username"].(string),
			"is_superuser":  tc["UserIsSuperuser"].(bool),
			"roles":         tc["Roles"].(*[]string),
			"role_required": role,
			"error":         err.Error(),
		}, resp.String())
		resp.WriteAndAbort(c)
		return
	} else if !res {
		resp := msg.WarnPermissionDeny.GenResponse("Role: '" + role + "' required.")
		log.WarnWithFields(log.Fields{
			"user_id":       tc["UserId"].(int),
			"username":      tc["Username"].(string),
			"is_superuser":  tc["UserIsSuperuser"].(bool),
			"roles":         tc["Roles"].(*[]string),
			"role_required": role,
		}, resp.String())
		resp.WriteAndAbort(c)
		return
	}
}
