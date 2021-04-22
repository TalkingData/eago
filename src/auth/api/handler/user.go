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

// SetUser 更新用户
// @Summary 更新用户
// @Tags 用户
// @Param token header string true "Token"
// @Param user_id path string true "用户ID"
// @Param data body db.User true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","user":{"id":1,"username":"username","email":"email","phone":"phone","is_superuser":false,"disabled":false,"last_login":null,"created_at":"2021-01-08 10:57:27","updated_at":"2021-01-26 11:31:24"}}"
// @Router /users/{user_id} [PUT]
func SetUser(c *gin.Context) {
	var userFm db.User

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())

		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&userFm); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'email', 'phone' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	u, suc := db.UserModel.Set(userId, userFm.Email, userFm.Phone)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.UserModel.Set.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"user": u})
	c.JSON(http.StatusOK, m.GinH())
}

// ListUsers 列出所有用户
// @Summary 列出所有用户
// @Tags 用户
// @Param token header string true "Token"
// @Param query query string false "过滤条件"
// @Param page query string false "页数"
// @Param page_size query string false "页尺寸"
// @Success 200 {string} string "{"code":0,"message":"Success","page":1,"page_size":50,"pages":1,"users":[{"id":1,"username":"username","email":"email","phone":"phone","is_superuser":false,"disabled":false,"last_login":null,"created_at":"2021-01-08 10:57:27","updated_at":"2021-01-26 11:31:24"}],"total":1}"
// @Router /users [GET]
func ListUsers(c *gin.Context) {
	var query db.Query

	if q := c.GetString("Query"); q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query = db.Query{"username LIKE @query OR id LIKE @query": sql.Named("query", likeQuery)}
	}

	paged, suc := db.UserModel.PagedList(
		&query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.UserModel.PageList.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPagedPayload(paged, "users")
	c.JSON(http.StatusOK, m.GinH())
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
		m := msg.WarnInvalidUri.NewMsg("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())

		c.JSON(http.StatusOK, m.GinH())
		return
	}

	roles, suc := db.UserModel.ListRoles(userId)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.UserModel.ListRoles.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"roles": roles})
	c.JSON(http.StatusOK, m.GinH())
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
		m := msg.WarnInvalidUri.NewMsg("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())

		c.JSON(http.StatusOK, m.GinH())
		return
	}

	prods, suc := db.UserModel.ListProducts(userId)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.UserModel.ListProducts.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"products": prods})
	c.JSON(http.StatusOK, m.GinH())
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
		m := msg.WarnInvalidUri.NewMsg("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())

		c.JSON(http.StatusOK, m.GinH())
		return
	}

	gps, suc := db.UserModel.ListGroups(userId)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.UserModel.ListGroups.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"groups": gps})
	c.JSON(http.StatusOK, m.GinH())
}

// GetUserDepartment 获得用户所在部门
// @Summary 获得用户所在部门（一个用户只能在一个部门）
// @Tags 用户
// @Param token header string true "Token"
// @Param user_id path string true "用户ID"
// @Success 200 {string} string "{"code":0,"departments":{"id":2,"name":"dept2","is_owner":true,"joined_at":"2021-01-21 11:20:56"},"message":"Success"}"
// @Router /users/{user_id}/department [GET]
func GetUserDepartment(c *gin.Context) {
	var pld = gin.H{}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())

		c.JSON(http.StatusOK, m.GinH())
		return
	}

	dept, suc := db.UserModel.GetDepartment(userId)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.UserModel.GetDepartment.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if dept != nil {
		pld["department"] = dept
	} else {
		pld["department"] = map[string]string{}
	}

	m := msg.Success.NewMsg().SetPayload(&pld)
	c.JSON(http.StatusOK, m.GinH())
}
