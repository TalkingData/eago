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

// NewGroup 新建组
func NewGroup(c *gin.Context) {
	var ngFrm dto.NewGroup

	// 序列化request body
	if err := c.ShouldBindJSON(&ngFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := ngFrm.Validate(); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	group, err := dao.NewGroup(ngFrm.Name, *ngFrm.Description)
	if err != nil {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "group", group)
}

// RemoveGroup 删除组
func RemoveGroup(c *gin.Context) {
	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "group_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var rGpFrm dto.RemoveGroup
	// 验证数据
	if m := rGpFrm.Validate(gId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if ok := dao.RemoveGroup(gId); !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetGroup 更新组
func SetGroup(c *gin.Context) {
	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "group_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var sgFrm dto.SetGroup
	// 序列化request body
	if err = c.ShouldBindJSON(&sgFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := sgFrm.Validate(gId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	group, err := dao.SetGroup(gId, sgFrm.Name, *sgFrm.Description)
	if err != nil {
		m := msg.UnknownError.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "group", group)
}

// ListGroups 列出所有组
func ListGroups(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	lgq := dto.ListGroupsQuery{}
	if c.ShouldBindQuery(&lgq) == nil {
		_ = lgq.UpdateQuery(query)
	}

	paged, ok := dao.PagedListGroups(
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

	w.WriteSuccessPayload(c, "groups", paged)
}

// AddUser2Group 添加用户至组
func AddUser2Group(c *gin.Context) {
	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "group_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var augFrm dto.AddUser2Group
	// 序列化request body
	if err = c.ShouldBindJSON(&augFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := augFrm.Validate(gId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.AddGroupUser(gId, augFrm.UserId, augFrm.IsOwner) {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// RemoveGroupUser 移除组中用户
func RemoveGroupUser(c *gin.Context) {
	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "group_id")
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

	var rpuFrm dto.RemoveGroupUser
	// 验证数据
	if m := rpuFrm.Validate(gId, userId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.RemoveGroupUser(gId, userId) {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetUserIsGroupOwner 设置用户是否是组Owner
func SetUserIsGroupOwner(c *gin.Context) {
	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "group_id")
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

	var suoFrm dto.SetUserIsGroupOwner
	// 序列化request body
	if err = c.ShouldBindJSON(&suoFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := suoFrm.Validate(gId, userId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.SetGroupUserIsOwner(gId, userId, suoFrm.IsOwner) {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// ListGroupUsers 列出角色所有用户
func ListGroupUsers(c *gin.Context) {
	gId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "group_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var lguFrm dto.ListGroupUsersQuery
	// 序列化request body
	_ = c.ShouldBindQuery(&lguFrm)
	if m := lguFrm.Validate(gId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	query := dao.Query{}
	// 设置查询filter
	_ = lguFrm.UpdateQuery(query)

	u, ok := dao.ListGroupUsers(gId, query)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "users", u)
}
