package handler

import (
	"eago/auth/api/form"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/common/api/ext"
	cMsg "eago/common/code_msg"
	"eago/common/logger"
	"eago/common/orm"
	"eago/common/tracer"
	"eago/common/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

// Logout 登出
func (ah *AuthHandler) Logout(c *gin.Context) {
	ah.biz.RemoveToken(tracer.ExtractTraceCtxFromGin(c), c.GetHeader("Token"))
	ext.WriteSuccess(c)
}

// Heartbeat 心跳
func (ah *AuthHandler) Heartbeat(c *gin.Context) {
	ah.biz.RenewalToken(tracer.ExtractTraceCtxFromGin(c), c.GetHeader("Token"))
	ext.WriteSuccess(c)
}

// CrowdLogin 从Crowd登录处理
func (ah *AuthHandler) CrowdLogin(c *gin.Context) {
	ah.logger.Info("authHandler.CrowdLogin called.")
	defer ah.logger.Info("authHandler.CrowdLogin end.")

	loginUser := c.GetStringMapString("LoginUser")

	ah.logger.DebugWithFields(logger.Fields{
		"username": loginUser["username"],
	}, "CrowdLogin called and got a user.")

	// 通过crowd认证
	crowdUser, err := ah.crowdCli.Authenticate(loginUser["username"], loginUser["password"])
	// 通过crowd认证成功
	if err == nil {
		// 阻止登录，非启用的用户
		if !crowdUser.Active {
			m := msg.MsgLoginInactiveCrowdUserFailed
			ah.logger.WarnWithFields(logger.Fields{
				"username": loginUser["username"],
			}, m.GetMsg())
			ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
			return
		}

		ah.logger.DebugWithFields(logger.Fields{
			"username": loginUser["username"],
		}, "Crowd login success.")
		ah.updateOrCreateUserLastLogin(c, crowdUser.Email)
		return
	}

	// TODO: 下方switch代码为临时解决方案，后续维护crowd模块解决
	// crowd错误类型完整介绍见：https://developer.atlassian.com/server/crowd/using-the-crowd-rest-apis/
	// 对于Crowd client返回的错误类型，做出相应处理
	switch err.Error() {
	case "INVALID_USER_AUTHENTICATION":
		// 认证失败：返回密码错误
		m := msg.MsgLoginAuthenticationFailed
		ah.logger.WarnWithFields(logger.Fields{
			"username": loginUser["username"],
			"error":    err,
		}, m.GetMsg())
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return

	case "INACTIVE_ACCOUNT":
		// 用户处于被禁用状态：返回当前用户在Crowd中是禁用状态
		m := msg.MsgLoginInactiveCrowdUserFailed
		ah.logger.WarnWithFields(logger.Fields{
			"username": loginUser["username"],
			"error":    err,
		}, m.GetMsg())
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}

	ah.logger.WarnWithFields(logger.Fields{
		"username": loginUser["username"],
		"error":    err,
	}, "Crowd login failed.")
	c.Next()
}

// DatabaseLogin 从数据库登录处理
func (ah *AuthHandler) DatabaseLogin(c *gin.Context) {
	ah.logger.Info("authHandler.DatabaseLogin called.")
	defer ah.logger.Info("authHandler.DatabaseLogin end.")

	loginUser := c.GetStringMapString("LoginUser")

	ah.logger.InfoWithFields(logger.Fields{
		"username": loginUser["username"],
	}, "DatabaseLoginHandler called and got a user.")

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 查询该用户在本地数据库中的数据
	user, err := ah.dao.GetUser(ctx, orm.Query{"username=?": loginUser["username"]})
	if err != nil {
		// 调用数据库出错
		m := msg.MsgAuthDaoErr.SetError(err)
		logF := m.ToLoggerFields()
		logF["username"] = loginUser["username"]
		ah.logger.ErrorWithFields(logF, "An error occurred while dao.GetUser in authHandler.DatabaseLogin.")
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}
	// 判断用户是否在DB中存在
	if user == nil {
		m := msg.MsgLoginAuthenticationFailed
		logF := m.ToLoggerFields()
		logF["username"] = loginUser["username"]
		ah.logger.WarnWithFields(logF, "Got an nil user from dao.GetUser in authHandler.DatabaseLogin.")
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}

	// 判断是否是被禁用的用户
	if user.Disabled {
		m := msg.MsgLoginDisabledUserFailed
		logF := m.ToLoggerFields()
		logF["username"] = loginUser["username"]
		ah.logger.WarnWithFields(logF, "User is disabled.")
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}

	// 判断是否是密码为空的用户
	if len(user.Password) < 1 {
		m := msg.MsgLoginNoPasswordUserFailed
		logF := m.ToLoggerFields()
		logF["username"] = loginUser["username"]
		ah.logger.WarnWithFields(logF, "No password user.")
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}

	// 判断账号密码是否匹配
	saltedPasswd := utils.GenSha256HashCode(loginUser["password"] + ah.conf.SecretKey)
	if loginUser["username"] == user.Username && saltedPasswd == user.Password {
		if err = ah.dao.SetUserLastLogin(ctx, user.Id); err != nil {
			ah.logger.WarnWithFields(logger.Fields{
				"username": loginUser["username"],
				"error":    err,
			}, "An error occurred while dao.SetUserLastLogin in authHandler.DatabaseLogin. but skipped.")
		}
		// 登录成功并返回token
		ah.newTokenResponse(c, user)
		return
	}

	c.Next()
}

// LoginFailed 返回登录失败
func (ah *AuthHandler) LoginFailed(c *gin.Context) {
	m := msg.MsgLoginAuthenticationFailed
	ah.logger.Warn(m.GetMsg())
	ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
	return
}

// updateOrCreateUserLastLogin 更新用户最近登录时间或创建用户并填入response
func (ah *AuthHandler) updateOrCreateUserLastLogin(c *gin.Context, username string) {
	ah.logger.InfoWithFields(logger.Fields{
		"username": username,
	}, "authHandler.updateOrCreateUserLastLogin called.")
	defer ah.logger.InfoWithFields(logger.Fields{
		"username": username,
	}, "authHandler.updateOrCreateUserLastLogin end.")

	ctx := tracer.ExtractTraceCtxFromGin(c)
	// 查询该用户在本地数据库中的数据
	user, err := ah.dao.GetUser(ctx, orm.Query{"username=?": username})
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(logger.Fields{
			"username": username,
		}, m.GetMsg())
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}

	// 用户在DB中存在
	if user != nil {
		// 判断是否为禁用的账号
		if user.Disabled {
			m := msg.MsgLoginDisabledUserFailed
			ah.logger.WarnWithFields(logger.Fields{
				"username": username,
			}, m.GetMsg())
			ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
			return
		}

		// 如果查找到用户则更新其登录时间
		if err = ah.dao.SetUserLastLogin(ctx, user.Id); err != nil {
			ah.logger.WarnWithFields(logger.Fields{
				"username": username,
			}, "An error occurred while dao.SetUserLastLogin in authHandler.updateOrCreateUserLastLogin. but skipped.")
		}

		// 登录成功并返回token
		ah.newTokenResponse(c, user)
		return
	}

	// 在本地数据库中没查找到用户则创建一个用户
	user, err = ah.dao.NewUser(ctx, username, username, true)
	if user != nil {
		// 登录成功并返回token
		ah.newTokenResponse(c, user)
		return
	}

	m := msg.MsgLoginUnknownFailed
	ah.logger.ErrorWithFields(logger.Fields{
		"username": username,
	}, m.GetMsg())
	ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
}

// newTokenResponse 生成Token并填入response
func (ah *AuthHandler) newTokenResponse(c *gin.Context, userObj *model.User) {
	ah.logger.InfoWithFields(logger.Fields{
		"user_id":  userObj.Id,
		"username": userObj.Username,
	}, "authHandler.newTokenResponse called.")
	defer ah.logger.InfoWithFields(logger.Fields{
		"user_id":  userObj.Id,
		"username": userObj.Username,
	}, "authHandler.newTokenResponse end.")

	// 登录成功并返回token
	tk := ah.biz.NewToken(tracer.ExtractTraceCtxFromGin(c), userObj)
	if tk == "" {
		m := msg.MsgLoginNewTokenFailed
		logF := m.ToLoggerFields()
		ah.logger.ErrorWithFields(logF, "An error occurred while biz.NewToken in authHandler.newTokenResponse.")
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}

	ah.logger.DebugWithFields(logger.Fields{
		"user_id":  userObj.Id,
		"username": userObj.Username,
		"token":    tk,
	}, "New token success.")
	ext.WriteSuccessPayload(c, "token", tk)
	c.Abort()
	return
}

func (ah *AuthHandler) ReadLoginForm(c *gin.Context) {
	frm := new(form.LoginForm)
	// 序列化request body获取用户名密码
	if err := c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields())
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}
	// 验证数据
	if err := frm.Validate(); err != nil {
		// 数据验证未通过
		m := cMsg.MsgValidateFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields())
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}

	c.Set("LoginUser", map[string]string{
		"username": strings.ToLower(strings.TrimSpace(frm.Username)),
		"password": strings.TrimSpace(frm.Password),
	})
}
