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

// SetUser 更新用户
// @Summary 更新用户
// @Tags 用户
// @Param token header string true "Token"
// @Param user_id path string true "用户ID"
// @Param data body model.User true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","user":{"id":1,"username":"username","email":"email","phone":"phone","is_superuser":false,"disabled":false,"last_login":null,"created_at":"2021-01-08 10:57:27","updated_at":"2021-01-26 11:31:24"}}"
// @Router /users/{user_id} [PUT]
func SetUser(c *gin.Context) {
	var userFm model.User

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())

		resp.Write(c)
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&userFm); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'email', 'phone' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	u, ok := model.SetUser(userId, userFm.Email, userFm.Phone)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.SetUser.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("user", u)
	resp.Write(c)
}

// ListUsers 列出所有用户
// @Summary 列出所有用户
// @Tags 用户
// @Param token header string true "Token"
// @Param query query string false "过滤条件"
// @Param order_by query string false "排序字段(多个间逗号分割)"
// @Param page query string false "页数"
// @Param page_size query string false "页尺寸"
// @Success 200 {string} string "{"code":0,"message":"Success","page":1,"page_size":50,"pages":1,"users":[{"id":1,"username":"username","email":"email","phone":"phone","is_superuser":false,"disabled":false,"last_login":null,"created_at":"2021-01-08 10:57:27","updated_at":"2021-01-26 11:31:24"}],"total":1}"
// @Router /users [GET]
func ListUsers(c *gin.Context) {
	var query model.Query

	if q := c.GetString("Query"); q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query = model.Query{"username LIKE @query OR id LIKE @query": sql.Named("query", likeQuery)}
	}

	paged, ok := model.PagedListUsers(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.PageListUsers.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPagedPayload(paged, "users")
	resp.Write(c)
}

// ListUserRoles 列出用户所在角色
// @Summary 列出用户所在角色
// @Tags 用户
// @Param token header string true "Token"
// @Param user_id path string true "用户ID"
// @Success 200 {string} string "{"code":0,"message":"Success","roles":[{"id":1,"name":"auth_admin"},{"id":2,"name":"tester"}]}"
// @Router /users/{user_id}/roles [GET]
func ListUserRoles(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())

		resp.Write(c)
		return
	}

	roles, ok := model.ListUserRoles(userId)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.ListUserRoles.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("roles", roles)
	resp.Write(c)
}

// ListUserProducts 列出用户所在产品线
// @Summary 列出用户所在产品线
// @Tags 用户
// @Param token header string true "Token"
// @Param user_id path string true "用户ID"
// @Success 200 {string} string "{"code":0,"message":"Success","products":[{"id":3,"name":"SCM","is_owner":false,"joined_at":"0001-01-01 00:00:00","alias":"scm","disable":false}]}"
// @Router /users/{user_id}/products [GET]
func ListUserProducts(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'user_id' required")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())

		resp.Write(c)
		return
	}

	prods, ok := model.ListUserProducts(userId)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.ListUserProducts")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("products", prods)
	resp.Write(c)
}

// ListUserGroups 列出用户所在组
// @Summary 列出用户所在组
// @Tags 用户
// @Param token header string true "Token"
// @Param user_id path string true "用户ID"
// @Success 200 {string} string "{"code":0,"groups":[{"id":2,"name":"group2","is_owner":true,"joined_at":"2021-01-21 11:20:56"}],"message":"Success"}"
// @Router /users/{user_id}/groups [GET]
func ListUserGroups(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'user_id' required")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())

		resp.Write(c)
		return
	}

	gps, ok := model.ListUserGroups(userId)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.ListUserGroups")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("groups", gps)
	resp.Write(c)
}

// GetUserDepartment 获得用户所在部门
// @Summary 获得用户所在部门（一个用户只能在一个部门）
// @Tags 用户
// @Param token header string true "Token"
// @Param user_id path string true "用户ID"
// @Success 200 {string} string "{"code":0,"departments":{"id":2,"name":"dept2","is_owner":true,"joined_at":"2021-01-21 11:20:56"},"message":"Success"}"
// @Router /users/{user_id}/department [GET]
func GetUserDepartment(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'user_id' required")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())

		resp.Write(c)
		return
	}

	dept, ok := model.GetUserDepartment(userId)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.GetUserDepartment")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	if dept != nil {
		resp.SetPayload("department", dept)
	} else {
		resp.SetPayload("department", "{}")
	}
	resp.Write(c)
}
