package handler

import (
	"eago/auth/cli"
	"eago/auth/conf/msg"
	"eago/auth/dao"
	"eago/auth/dto"
	auth "eago/auth/srv/proto"
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/common/tracer"
	"github.com/gin-gonic/gin"
	"strconv"
)

// SetUser 更新用户
func SetUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "user_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var setUserFrm dto.SetUser
	// 序列化request body
	if err = c.ShouldBindJSON(&setUserFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := setUserFrm.Validate(userId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	u, ok := dao.SetUser(userId, setUserFrm.Email, setUserFrm.Phone)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "user", u)
}

// ListUsers 列出所有用户
func ListUsers(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	luQ := dto.ListUsersQuery{}
	if c.ShouldBindQuery(&luQ) == nil {
		_ = luQ.UpdateQuery(query)
	}

	paged, ok := dao.PagedListUsers(
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

	w.WriteSuccessPayload(c, "users", paged)
}

// ListUserRoles 列出用户所在角色
func ListUserRoles(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "user_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	roles, ok := dao.ListUserRoles(userId)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "roles", roles)
}

// ListUserProducts 列出用户所在产品线
func ListUserProducts(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "user_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	prods, ok := dao.ListUserProducts(userId)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "products", prods)
}

// ListUserGroups 列出用户所在组
func ListUserGroups(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "user_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	gps, ok := dao.ListUserGroups(userId)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "groups", gps)
}

// GetUserDepartment 获得用户所在部门
func GetUserDepartment(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "user_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	dept, ok := dao.GetUserDepartment(userId)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if dept != nil {
		w.WriteSuccessPayload(c, "department", dept)
		return
	}
	w.WriteSuccessPayload(c, "department", "{}")
}

// MakeUserHandover 用户交接
func MakeUserHandover(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "user_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	tgtUserId, err := strconv.Atoi(c.Param("target_user_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "target_user_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var handoverFrm dto.MakeUserHandover
	// 验证数据
	if m := handoverFrm.Validate(userId, tgtUserId); m != nil {
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	ctx := tracer.ExtractTraceContext(c)
	hdReq := &auth.HandoverRequest{UserId: int32(userId), TargetUserId: int32(tgtUserId)}
	res, err := cli.AuthClient.MakeUserHandover(ctx, hdReq)
	if err != nil {
		m := msg.HandoverUnknownError.SetError(err)
		log.ErrorWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	if !res.Ok {
		m := msg.HandoverUnknownError
		log.ErrorWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}
