package handler

import (
	"database/sql"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/common/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// NewRole 新建角色
// @Summary 新建角色
// @Tags 角色
// @Param token header string true "Token"
// @Param data body model.Role true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","role":{"id":2,"name":"new_role"}}"
// @Router /roles [POST]
func NewRole(c *gin.Context) {
	var role model.Role

	// 序列化request body
	if err := c.ShouldBindJSON(&role); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'name' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	r := model.NewRole(role.Name)
	if r == nil {
		resp := msg.ErrDatabase.GenResponse("Error in model.NewRole.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("role", r)
	resp.Write(c)
}

// RemoveRole 删除角色
// @Summary 删除角色
// @Tags 角色
// @Param token header string true "Token"
// @Param role_id path string true "角色ID"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /roles/{role_id} [DELETE]
func RemoveRole(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'role_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if ok := model.RemoveRole(roleId); !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.RemoveRole.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
}

// SetRole 更新角色
// @Summary 更新角色
// @Tags 角色
// @Param token header string true "Token"
// @Param role_id path string true "角色ID"
// @Param data body model.Role true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","role":{"id":2,"name":"tester"}}"
// @Router /roles/{role_id} [PUT]
func SetRole(c *gin.Context) {
	var roleFm model.Role

	roleId, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'role_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&roleFm); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'name' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	role, ok := model.SetRole(roleId, roleFm.Name)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.SetRole.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("role", role)
	resp.Write(c)
}

// ListRoles 列出所有角色
// @Summary 列出所有角色
// @Tags 角色
// @Param token header string true "Token"
// @Param query query string false "过滤条件"
// @Param order_by query string false "排序字段(多个间逗号分割)"
// @Param page query string false "页数"
// @Param page_size query string false "页尺寸"
// @Success 200 {string} string "{"code":0,"message":"Success","page":1,"page_size":50,"pages":1,"roles":[{"id":1,"name":"auth_admin"},{"id":2,"name":"tester"}],"total":2}"
// @Router /roles [GET]
func ListRoles(c *gin.Context) {
	var query model.Query
	q := c.GetString("Query")
	if q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query = model.Query{"name LIKE @query OR id LIKE @query": sql.Named("query", likeQuery)}
	}

	paged, ok := model.PagedListRoles(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.PageListRoles.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPagedPayload(paged, "roles")
	resp.Write(c)
}

// AddUser2Role 添加用户至角色
// @Summary 添加用户至角色
// @Tags 角色
// @Param token header string true "Token"
// @Param role_id path string true "角色ID"
// @Param data body model.UserRole true "body"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /roles/{role_id}/users [POST]
func AddUser2Role(c *gin.Context) {
	var uRole model.UserRole

	roleId, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'role_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&uRole); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'user_id' required, and it must greater than 0.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if !model.AddRoleUser(uRole.UserId, roleId) {
		resp := msg.ErrDatabase.GenResponse("Error in model.AddRoleUser.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
}

// RemoveRoleUser 移除角色中用户
// @Summary 移除角色中用户
// @Tags 角色
// @Param token header string true "Token"
// @Param role_id path string true "角色ID"
// @Param user_id path string true "用户ID"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /roles/{role_id}/users/{user_id} [DELETE]
func RemoveRoleUser(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'role_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'user_id' required, and it must greater than 0.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if !model.RemoveRoleUser(userId, roleId) {
		resp := msg.ErrDatabase.GenResponse("Error in model.RemoveRoleUser.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
}

// ListRoleUsers 列出角色所有用户
// @Summary 列出角色所有用户
// @Tags 角色
// @Param token header string true "Token"
// @Param role_id path string true "角色ID"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /roles/{role_id}/users [GET]
func ListRoleUsers(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'role_id' required")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())

		resp.Write(c)
		return
	}

	u, ok := model.ListRoleUsers(roleId)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.ListUsers")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("users", u)
	resp.Write(c)
}
