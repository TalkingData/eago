package menu

import (
	"eago/common/api/permission"
	"github.com/gin-gonic/gin"
)

const (
	permKeyIsRole   = "is_role"
	permKeyIsRoleIn = "is_role_in"
)

type itemPerm struct {
	key string
	val interface{}
}

func NewPermIsRole(roleRequired string) *itemPerm {
	return &itemPerm{
		key: permKeyIsRole,
		val: roleRequired,
	}
}

func NewPermIsRoleIn(rolesRequired []string) *itemPerm {
	return &itemPerm{
		key: permKeyIsRoleIn,
		val: rolesRequired,
	}
}

func (perm *itemPerm) hasPerm(c *gin.Context) (bool, error) {
	switch perm.key {
	case permKeyIsRole:
		return permission.IsRole(c, perm.val.(string))
	case permKeyIsRoleIn:
		return permission.IsRoleIn(c, perm.val.([]string))
	default:
		return false, nil
	}
}
