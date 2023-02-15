package handler

import (
	"eago/common/api/ext"
	"eago/common/logger"
	"eago/common/tracer"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginFromEagle 根据Eagle的token转为本系统的token
func (ah *AuthHandler) LoginFromEagle(c *gin.Context) {
	eagleTk := c.GetHeader("eagle_token")

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 根据eagle token获取用户对象
	userObj, err := ah.biz.GetUserObjectFromEagleToken(ctx, eagleTk)
	if err != nil {
		ah.logger.WarnWithFields(logger.Fields{
			"eagle_token": eagleTk,
			"error":       err,
		}, "An error occurred while builtin.GetUserObjectFromEagleToken, Please contact admin.")
		ext.WriteAnyAndAbort(c, http.StatusForbidden, err.Error())
		return
	}

	// 判断用户是否在DB中存在，并且不是被禁用的用户，不存在则退出
	if userObj == nil {
		m := "user not found or disabled"
		ah.logger.WarnWithFields(logger.Fields{
			"eagle_token": eagleTk,
		}, m)
		ext.WriteAnyAndAbort(c, http.StatusForbidden, m)
		return
	}

	// 生成token
	if token := ah.biz.NewToken(ctx, userObj); len(token) > 0 {
		ah.logger.DebugWithFields(logger.Fields{
			"user_id":  userObj.Id,
			"username": userObj.Username,
			"token":    token,
		}, "New token success.")
		ext.WriteSuccessPayload(c, "token", token)
		return
	}

	ext.WriteAnyAndAbort(c, http.StatusInternalServerError, "login from eagle failed, got nil token.")
}
