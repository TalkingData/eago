package permission

import (
	authpb "eago/auth/proto"
	"eago/common/api/ext"
	"eago/common/utils"
	"github.com/gin-gonic/gin"
)

// IsCurrUser 指定字段Id=当前用户Id处理
func IsCurrUser(c *gin.Context, userIdField string) (bool, error) {
	val, exists := c.Get(defaultTokenContentGinCtxKey)
	if !exists || val == nil {
		return false, nil
	}
	tc := val.(*authpb.TokenContent)

	userId, err := ext.ParamInt(c, userIdField)
	if err != nil {
		return false, err
	}

	return int(tc.UserId) == userId, nil
}

// IsCurrUserInProduct 判断当前用户在指定产品线内处理
func IsCurrUserInProduct(c *gin.Context, prodIdField string, isOwner bool) (bool, error) {
	val, exists := c.Get(defaultTokenContentGinCtxKey)
	if !exists || val == nil {
		return false, nil
	}
	tc := val.(*authpb.TokenContent)

	prods := []*authpb.Product{}
	if isOwner {
		prods = tc.OwnProducts
	} else {
		prods = tc.Products
	}

	prodId, err := ext.ParamUint32(c, prodIdField)
	if err != nil {
		return false, err
	}

	for _, p := range prods {
		if p.Id == prodId {
			return true, nil
		}
	}

	return false, nil
}

// IsCurrUserInGroup 判断当前用户在指定组内处理
func IsCurrUserInGroup(c *gin.Context, groupIdField string, isOwner bool) (bool, error) {
	val, exists := c.Get(defaultTokenContentGinCtxKey)
	if !exists || val == nil {
		return false, nil
	}
	tc := val.(*authpb.TokenContent)

	groups := []*authpb.Group{}
	if isOwner {
		groups = tc.OwnGroups
	} else {
		groups = tc.Groups
	}

	groupId, err := ext.ParamUint32(c, groupIdField)
	if err != nil {
		return false, err
	}

	for _, g := range groups {
		if g.Id == groupId {
			return true, nil
		}
	}

	return false, nil
}

// IsRole 判断当前用户的角色，是否是指定角色的处理
func IsRole(c *gin.Context, roleRequired string) (bool, error) {
	val, exists := c.Get(defaultTokenContentGinCtxKey)
	if !exists || val == nil {
		return false, nil
	}
	tc := val.(*authpb.TokenContent)

	// 超级用户拥有所有权限，跳过判断
	if tc.IsSuperuser {
		return true, nil
	}

	// 判断是否是对应角色
	return utils.IsInSlice(tc.Roles, roleRequired)
}

// IsRoleIn 判断当前用户的角色，是否是指定角色数组中的任意一个的处理
func IsRoleIn(c *gin.Context, rolesRequired []string) (res bool, err error) {
	val, exists := c.Get(defaultTokenContentGinCtxKey)
	if !exists || val == nil {
		return false, nil
	}
	tc := val.(*authpb.TokenContent)

	// 超级用户拥有所有权限，跳过判断
	if tc.IsSuperuser {
		return true, nil
	}

	// 判断是否是对应角色之一
	for _, roleReq := range rolesRequired {
		res, err = utils.IsInSlice(tc.Roles, roleReq)
		if err != nil {
			return false, err
		}
		if res {
			break
		}
	}

	return
}
