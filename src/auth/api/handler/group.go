package handler

import (
	"database/sql"
	"eago/auth/api/form"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/common/log"
	"eago/common/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// NewGroup 新建组
// @Summary 新建组
// @Tags 组
// @Param token header string true "Token"
// @Param data body model.Group true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","group":{"id":1,"name":"group1"}}"
// @Router /groups [POST]
func NewGroup(c *gin.Context) {
	var gForm model.Group

	// 序列化request body
	if err := c.ShouldBindJSON(&gForm); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'name', 'description' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	group := model.NewGroup(gForm.Name, *gForm.Description)
	if group == nil {
		resp := msg.ErrDatabase.GenResponse("Error in NewGroup.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("group", group)
	resp.Write(c)
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
		resp := msg.WarnInvalidUri.GenResponse("Field 'group_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if ok := model.RemoveGroup(gId); !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.RemoveGroup.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
}

// SetGroup 更新组
// @Summary 更新组
// @Tags 组
// @Param token header string true "Token"
// @Param group_id path string true "组ID"
// @Param data body model.Group true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","group":{"id":1,"name":"group_rename"}}"
// @Router /groups/{group_id} [PUT]
func SetGroup(c *gin.Context) {
	var gFm model.Group

	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'group_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&gFm); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'name', 'description' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	group, ok := model.SetGroup(gId, gFm.Name, *gFm.Description)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.SetGroup.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("group", group)
	resp.Write(c)
}

// ListGroups 列出所有组
// @Summary 列出所有组
// @Tags 组
// @Param token header string true "Token"
// @Param query query string false "过滤条件"
// @Param order_by query string false "排序字段(多个间逗号分割)"
// @Param page query string false "页数"
// @Param page_size query string false "页尺寸"
// @Success 200 {string} string "{"code":0,"groups":[{"id":1,"name":"group1","description":"group1","created_at":"2021-01-21 11:20:29","updated_at":"2021-01-21 11:20:29"}],"message":"Success","page":1,"page_size":50,"pages":1,"total":1}"
// @Router /groups [GET]
func ListGroups(c *gin.Context) {
	var query model.Query

	q := c.GetString("Query")
	if q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query = model.Query{"name LIKE @query OR id LIKE @query": sql.Named("query", likeQuery)}
	}

	paged, ok := model.PagedListGroups(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.PageListGroups.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPagedPayload(paged, "groups")
	resp.Write(c)
}

// AddUser2Group 添加用户至组
// @Summary 添加用户至组
// @Tags 组
// @Param token header string true "Token"
// @Param group_id path string true "组ID"
// @Param data body model.UserGroup true "body"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /groups/{group_id}/users [POST]
func AddUser2Group(c *gin.Context) {
	var uGroup model.UserGroup

	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'group_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&uGroup); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'user_id', 'is_owner' required, and 'user_id' must greater than 0.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if !model.AddGroupUser(uGroup.UserId, gId, *uGroup.IsOwner) {
		resp := msg.ErrDatabase.GenResponse("Error in model.AddGroupUser.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
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
		resp := msg.WarnInvalidUri.GenResponse("Field 'group_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if !model.RemoveGroupUser(userId, gId) {
		resp := msg.ErrDatabase.GenResponse("Error in model.RemoveUser.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
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
		resp := msg.WarnInvalidUri.GenResponse("Field 'group_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

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
	if err := c.ShouldBindJSON(&fm); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'is_owner' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if !model.SetGroupUserIsOwner(userId, gId, *fm.IsOwner) {
		resp := msg.ErrDatabase.GenResponse("Error in model.SetGroupUserIsOwner.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
}

// ListGroupUsers 列出角色所有用户
// @Summary 列出角色所有用户
// @Tags 组
// @Param token header string true "Token"
// @Param group_id path string true "组ID"
// @Success 200 {string} string "{"code":0,"message":"Success","users":[{"id":4,"username":"test2","is_owner":false,"joined_at":"2021-01-20 11:01:16"},{"id":3,"username":"test","is_owner":true,"joined_at":"2021-01-20 11:01:32"}]}"
// @Router /groups/{group_id}/users [GET]
func ListGroupUsers(c *gin.Context) {
	var query = model.Query{}

	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'group_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// TODO(lai.li)
	// 方法中is_owner传值只能是0 or 1，待将来解决
	isOwner, err := strconv.Atoi(c.DefaultQuery("is_owner", "-1"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'is_owner' required, and must integer 0 or 1.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}
	if isOwner >= 0 {
		query["is_owner"] = utils.IntMin(isOwner, 1)
	}

	u, ok := model.ListGroupUsers(gId, query)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.ListGroupUsers.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("users", u)
	resp.Write(c)
}
