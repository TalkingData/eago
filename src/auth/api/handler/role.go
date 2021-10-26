package handler

import (
	"eago/auth/conf/msg"
	"eago/auth/dao"
	"eago/auth/dto"
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"github.com/gin-gonic/gin"
	"strconv"
)

// NewRole 新建角
func NewRole(c *gin.Context) {
	var newRoleFrm dto.NewRole
	// 序列化request body
	if err := c.ShouldBindJSON(&newRoleFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := newRoleFrm.Validate(); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	r, err := dao.NewRole(newRoleFrm.Name, *newRoleFrm.Description)
	if err != nil {
		m := msg.UnknownError.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "role", r)
}

// RemoveRole 删除角色
func RemoveRole(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "role_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 验证数据
	var removeRoleFrm dto.RemoveRole
	if m := removeRoleFrm.Validate(roleId); m != nil {
		// 数据验证未通过
		m.WriteRest(c)
		return
	}

	// 删除
	if dbErr := dao.RemoveRole(roleId); dbErr != nil {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetRole 更新角色
func SetRole(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		m := msg.InvalidUriFailed
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var setRoleFrm dto.SetRole
	// 序列化request body
	if err = c.ShouldBindJSON(&setRoleFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := setRoleFrm.Validate(roleId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	role, err := dao.SetRole(roleId, setRoleFrm.Name, *setRoleFrm.Description)
	if err != nil {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "role", role)
}

// ListRoles 列出所有角色
func ListRoles(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	lrq := dto.ListRolesQuery{}
	if c.ShouldBindQuery(&lrq) == nil {
		_ = lrq.UpdateQuery(query)
	}

	paged, ok := dao.PagedListRoles(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "roles", paged)
}

// AddUser2Role 添加用户至角色
func AddUser2Role(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "role_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var aurFrm dto.AddUser2Role
	// 序列化request body
	if err = c.ShouldBindJSON(&aurFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := aurFrm.Validate(roleId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.AddRoleUser(roleId, aurFrm.UserId) {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// RemoveRoleUser 移除角色中用户
func RemoveRoleUser(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "role_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "user_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var rruFrm dto.RemoveRoleUser
	// 验证数据
	if m := rruFrm.Validate(roleId, userId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.RemoveRoleUser(roleId, userId) {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// ListRoleUsers 列出角色所有用户
func ListRoleUsers(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "role_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	u, ok := dao.ListRoleUsers(roleId)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "users", u)
}
