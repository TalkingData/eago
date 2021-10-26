package permission

import (
	auth "eago/auth/srv/proto"
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/common/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

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
			m := InvalidParams.SetError(err, userIdField)
			log.WarnWithFields(m.LogFields())
			w.WriteAnyAndAbort(c, m.Code(), m.String())
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
			m := InvalidParams.SetError(err, prodIdField)
			log.WarnWithFields(m.LogFields())
			w.WriteAnyAndAbort(c, m.Code(), m.String())
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
			m := InvalidParams.SetError(err, groupIdField)
			log.WarnWithFields(m.LogFields())
			w.WriteAnyAndAbort(c, m.Code(), m.String())
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
	currUserId := int(tc["UserId"].(int32))

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
	var prods []*auth.Product

	tc := c.GetStringMap("TokenContent")
	if isOwner {
		prods = tc["OwnProducts"].([]*auth.Product)
	} else {
		prods = tc["Products"].([]*auth.Product)
	}

	prodId, err := strconv.Atoi(c.Param(prodIdField))
	if err != nil {
		return false, err
	}
	for _, p := range prods {
		if int(p.Id) == prodId {
			return true, nil
		}
	}

	return false, nil
}

// isCurrUserInGroupHandler 判断当前用户在指定组内处理
func isCurrUserInGroupHandler(c *gin.Context, groupIdField string, isOwner bool) (bool, error) {
	var groups []*auth.Group

	tc := c.GetStringMap("TokenContent")
	if isOwner {
		groups = tc["OwnGroups"].([]*auth.Group)
	} else {
		groups = tc["Groups"].([]*auth.Group)
	}

	groupId, err := strconv.Atoi(c.Param(groupIdField))
	if err != nil {
		return false, err
	}
	for _, g := range groups {
		if int(g.Id) == groupId {
			return true, nil
		}
	}

	return false, nil
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
		log.ErrorWithFields(log.Fields{
			"user_id":       tc["UserId"].(int32),
			"username":      tc["Username"].(string),
			"is_superuser":  tc["IsSuperuser"].(bool),
			"roles":         tc["Roles"].([]string),
			"role_required": role,
			"error":         err,
		}, CheckRoleError.String())

		w.WriteAnyAndAbort(c, CheckRoleError.Code(), CheckRoleError.String())
		return
	} else if !res {
		m := UserNotRole.SetDetail(role)
		log.WarnWithFields(log.Fields{
			"user_id":       tc["UserId"].(int32),
			"username":      tc["Username"].(string),
			"is_superuser":  tc["IsSuperuser"].(bool),
			"roles":         tc["Roles"].([]string),
			"role_required": role,
		}, m.String())
		w.WriteAnyAndAbort(c, m.Code(), m.String())
		return
	}

	return
}
