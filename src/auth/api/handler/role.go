package handler

import (
	"database/sql"
	"eago-auth/conf/msg"
	db "eago-auth/database"
	"eago-common/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// NewRole 新建角色
// @Summary 新建角色
// @Tags 角色
// @Param token header string true "Token"
// @Param data body db.Role true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","role":{"id":2,"name":"new_role"}}"
// @Router /roles [POST]
func NewRole(c *gin.Context) {
	var role db.Role

	// 序列化request body
	if err := c.ShouldBindJSON(&role); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'name' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	r := db.RoleModel.New(role.Name)
	if r == nil {
		m := msg.ErrDatabase.NewMsg("Error in db.RoleModel.New.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"role": r})
	c.JSON(http.StatusOK, m.GinH())
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
		m := msg.WarnInvalidUri.NewMsg("Field 'role_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if suc := db.RoleModel.Remove(roleId); !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.RoleModel.Remove.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
}

// SetRole 更新角色
// @Summary 更新角色
// @Tags 角色
// @Param token header string true "Token"
// @Param role_id path string true "角色ID"
// @Param data body db.Role true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","role":{"id":2,"name":"tester"}}"
// @Router /roles/{role_id} [PUT]
func SetRole(c *gin.Context) {
	var roleFm db.Role

	roleId, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'role_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&roleFm); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'name' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	role, suc := db.RoleModel.Set(roleId, roleFm.Name)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.RoleModel.Set.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"role": role})
	c.JSON(http.StatusOK, m.GinH())
}

// ListRoles 列出所有角色
// @Summary 列出所有角色
// @Tags 角色
// @Param token header string true "Token"
// @Param query query string false "过滤条件"
// @Param page query string false "页数"
// @Param page_size query string false "页尺寸"
// @Success 200 {string} string "{"code":0,"message":"Success","page":1,"page_size":50,"pages":1,"roles":[{"id":1,"name":"auth_admin"},{"id":2,"name":"tester"}],"total":2}"
// @Router /roles [GET]
func ListRoles(c *gin.Context) {
	var query db.Query
	q := c.GetString("Query")
	if q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query = db.Query{"name LIKE @query OR id LIKE @query": sql.Named("query", likeQuery)}
	}

	paged, suc := db.RoleModel.PagedList(
		&query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.RoleModel.PageList.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPagedPayload(paged, "roles")
	c.JSON(http.StatusOK, m.GinH())
}

// AddUser2Role 添加用户至角色
// @Summary 添加用户至角色
// @Tags 角色
// @Param token header string true "Token"
// @Param role_id path string true "角色ID"
// @Param data body db.UserRole true "body"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /roles/{role_id}/users [POST]
func AddUser2Role(c *gin.Context) {
	var uRole db.UserRole

	roleId, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'role_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&uRole); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if !db.RoleModel.AddUser(uRole.UserId, roleId) {
		m := msg.ErrDatabase.NewMsg("Error in db.RoleModel.AddUser.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
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
		m := msg.WarnInvalidUri.NewMsg("Field 'role_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if !db.RoleModel.RemoveUser(userId, roleId) {
		m := msg.ErrDatabase.NewMsg("Error in db.RoleModel.RemoveUser.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
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
		m := msg.WarnInvalidUri.NewMsg("Field 'role_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())

		c.JSON(http.StatusOK, m.GinH())
		return
	}

	u, suc := db.RoleModel.ListUsers(roleId)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.RoleModel.ListUsers.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"users": u})
	c.JSON(http.StatusOK, m.GinH())
}
