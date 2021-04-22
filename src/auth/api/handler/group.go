package handler

import (
	"database/sql"
	"eago-auth/api/form"
	"eago-auth/conf/msg"
	db "eago-auth/database"
	"eago-common/log"
	"eago-common/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// NewGroup 新建组
// @Summary 新建组
// @Tags 组
// @Param token header string true "Token"
// @Param data body db.Group true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","group":{"id":1,"name":"group1"}}"
// @Router /groups [POST]
func NewGroup(c *gin.Context) {
	var gForm db.Group

	// 序列化request body
	if err := c.ShouldBindJSON(&gForm); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'name', 'description' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	group := db.GroupModel.New(gForm.Name, gForm.Description)
	if group == nil {
		m := msg.ErrDatabase.NewMsg("Error in db.GroupModel.New.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"group": group})
	c.JSON(http.StatusOK, m.GinH())
}

// RemoveGroup 删除组
// @Summary 删除组
// @Tags 组
// @Param token header string true "Token"
// @Param group_id path string true "组ID"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /groups/{group_id} [DELETE]
func RemoveGroup(c *gin.Context) {
	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'group_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if suc := db.GroupModel.Remove(gId); !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.GroupModel.Remove.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
}

// SetGroup 更新组
// @Summary 更新组
// @Tags 组
// @Param token header string true "Token"
// @Param group_id path string true "组ID"
// @Param data body db.Group true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","group":{"id":1,"name":"group_rename"}}"
// @Router /groups/{group_id} [PUT]
func SetGroup(c *gin.Context) {
	var gFm db.Group

	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'group_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&gFm); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'name', 'description' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	group, suc := db.GroupModel.Set(gId, gFm.Name, gFm.Description)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.GroupModel.Set.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"group": group})
	c.JSON(http.StatusOK, m.GinH())
}

// ListGroups 列出所有组
// @Summary 列出所有组
// @Tags 组
// @Param token header string true "Token"
// @Param query query string false "过滤条件"
// @Param page query string false "页数"
// @Param page_size query string false "页尺寸"
// @Success 200 {string} string "{"code":0,"groups":[{"id":1,"name":"group1","description":"group1","created_at":"2021-01-21 11:20:29","updated_at":"2021-01-21 11:20:29"}],"message":"Success","page":1,"page_size":50,"pages":1,"total":1}"
// @Router /groups [GET]
func ListGroups(c *gin.Context) {
	var query db.Query

	q := c.GetString("Query")
	if q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query = db.Query{"name LIKE @query OR alias LIKE @query id LIKE @query": sql.Named("query", likeQuery)}
	}

	paged, suc := db.GroupModel.PagedList(
		&query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.GroupModel.PageList.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPagedPayload(paged, "groups")
	c.JSON(http.StatusOK, m.GinH())
}

// AddUser2Group 添加用户至组
// @Summary 添加用户至组
// @Tags 组
// @Param token header string true "Token"
// @Param group_id path string true "组ID"
// @Param data body db.UserGroup true "body"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /groups/{group_id}/users [POST]
func AddUser2Group(c *gin.Context) {
	var uGroup db.UserGroup

	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'group_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&uGroup); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'user_id', 'is_owner' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if !db.GroupModel.AddUser(uGroup.UserId, gId, *uGroup.IsOwner) {
		m := msg.ErrDatabase.NewMsg("Error in db.GroupModel.AddUser.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
}

// RemoveGroupUser 移除组中用户
// @Summary 移除组中用户
// @Tags 组
// @Param token header string true "Token"
// @Param group_id path string true "组ID"
// @Param user_id path string true "用户ID"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /groups/{group_id}/users/{user_id} [DELETE]
func RemoveGroupUser(c *gin.Context) {
	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'group_id' required.")
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

	if !db.GroupModel.RemoveUser(userId, gId) {
		m := msg.ErrDatabase.NewMsg("Error in db.GroupModel.RemoveUser.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
}

// SetUserIsGroupOwner 设置用户是否是组Owner
// @Summary 设置用户是否是组Owner
// @Tags 组
// @Param token header string true "Token"
// @Param product_id path string true "组ID"
// @Param user_id path string true "用户ID"
// @Param data body form.IsOwnerForm true "body"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /groups/{group_id}/users/{user_id} [PUT]
func SetUserIsGroupOwner(c *gin.Context) {
	var fm = form.IsOwnerForm{}

	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'group_id' required.")
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

	// 序列化request body
	if err := c.ShouldBindJSON(&fm); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'is_owner' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if !db.GroupModel.SetUserIsOwner(userId, gId, *fm.IsOwner) {
		m := msg.ErrDatabase.NewMsg("Error in db.GroupModel.SetUserIsOwner.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
}

// ListGroupUsers 列出角色所有用户
// @Summary 列出角色所有用户
// @Tags 组
// @Param token header string true "Token"
// @Param group_id path string true "组ID"
// @Success 200 {string} string "{"code":0,"message":"Success","users":[{"id":4,"username":"test2","is_owner":false,"joined_at":"2021-01-20 11:01:16"},{"id":3,"username":"test","is_owner":true,"joined_at":"2021-01-20 11:01:32"}]}"
// @Router /groups/{group_id}/users [GET]
func ListGroupUsers(c *gin.Context) {
	var query = db.Query{}

	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'group_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// TODO(lai.li)
	// 方法中is_owner传值只能是0 or 1，待将来解决
	isOwner, err := strconv.Atoi(c.DefaultQuery("is_owner", "-1"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'is_owner' required, and must integer 0 or 1.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}
	if isOwner >= 0 {
		query["is_owner"] = tools.IntMin(isOwner, 1)
	}

	u, suc := db.GroupModel.ListUsers(gId, &query)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.GroupModel.ListUsers.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"users": u})
	c.JSON(http.StatusOK, m.GinH())
}
