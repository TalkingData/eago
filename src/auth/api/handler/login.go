package handler

import (
	"bytes"
	"eago/auth/conf"
	"eago/auth/conf/msg"
	"eago/auth/dao"
	"eago/auth/model"
	"eago/auth/srv/builtin"
	"eago/auth/util/sso"
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/common/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Logout 登出
func Logout(c *gin.Context) {
	builtin.RemoveToken(c.GetHeader("Token"))
	w.WriteSuccess(c)
}

// Heartbeat 心跳
func Heartbeat(c *gin.Context) {
	builtin.RenewalToken(c.GetHeader("Token"))
	w.WriteSuccess(c)
}

// GetTokenContent 获得TokenContent
func GetTokenContent(c *gin.Context) {
	tk := c.GetHeader("Token")
	tc, ok := builtin.GetTokenContent(tk)
	if !ok {
		m := "Invalid token or Not login yet."
		log.WarnWithFields(log.Fields{"token": tk}, m)
		w.WriteAnyAndAbort(c, http.StatusForbidden, m)
		return
	}

	w.WriteSuccessPayload(c, "content", tc)
}

// IamLogin 从IAM登录处理
func IamLogin(c *gin.Context) {
	log.Info("IamLogin called.")
	defer log.Info("IamLogin end.")

	loginUser := c.GetStringMapString("LoginUser")

	log.InfoWithFields(log.Fields{
		"username": loginUser["username"],
	}, "CrowdLoginHandler called and got a user.")

	data := []byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, loginUser["username"], loginUser["password"]))

	iamResp, err := http.Post(conf.Conf.IamAddress, "application/json", bytes.NewReader(data))
	if err == nil {
		defer func() {
			_ = iamResp.Body.Close()
		}()

		log.DebugWithFields(log.Fields{
			"response_status_code": iamResp.StatusCode,
			"response_body":        iamResp.Body,
		}, "Got iam response.")

		if iamResp.StatusCode == 200 {
			log.Debug("Iam login success.")
			updateOrCreateUserLastLogin(c, loginUser["username"])
			return
		}
	} else {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "Iam login failed.")
		c.Next()
		return
	}

	log.Error("Iam login failed.")
	c.Next()
}

// CrowdLogin 从Crowd登录处理
func CrowdLogin(c *gin.Context) {
	log.Info("CrowdLogin called.")
	defer log.Info("CrowdLogin end.")

	loginUser := c.GetStringMapString("LoginUser")

	log.InfoWithFields(log.Fields{
		"username": loginUser["username"],
	}, "CrowdLoginHandler called and got a user.")

	// 通过crowd认证
	crowdUser, err := sso.Crowd.Authenticate(
		strings.TrimSpace(loginUser["username"]),
		strings.TrimSpace(loginUser["password"]),
	)
	// 通过crowd认证成功
	if err == nil {
		log.Debug("Crowd login success.")
		// 阻止登录，非启用的用户
		if !crowdUser.Active {
			m := msg.LoginInactiveCrowdUserFailed
			log.WarnWithFields(log.Fields{
				"username": loginUser["username"],
			}, m.String())
			w.WriteAnyAndAbort(c, m.Code(), m.String())
			return
		}

		updateOrCreateUserLastLogin(c, crowdUser.Email)
		return
	}

	log.ErrorWithFields(log.Fields{
		"error": err,
	}, "Crowd login failed.")
	c.Next()
}

// DatabaseLogin 从数据库登录处理
func DatabaseLogin(c *gin.Context) {
	log.Info("DatabaseLogin called.")
	defer log.Info("DatabaseLogin end.")

	loginUser := c.GetStringMapString("LoginUser")

	log.InfoWithFields(log.Fields{
		"username": loginUser["username"],
	}, "DatabaseLoginHandler called and got a user.")

	// 查询该用户在本地数据库中的数据
	user, ok := dao.GetUser(dao.Query{"username=?": loginUser["username"]})
	if !ok {
		// 调用数据库出错
		m := "An error occurred while GetUser, Please contact admin."
		log.ErrorWithFields(log.Fields{
			"username": loginUser["username"],
		}, m)
		w.WriteAnyAndAbort(c, http.StatusInternalServerError, m)
		return
	}

	// 判断用户是否在DB中存在
	if user == nil {
		m := msg.LoginAuthenticationFailed
		log.WarnWithFields(log.Fields{
			"username": loginUser["username"],
		}, m.String())
		w.WriteAnyAndAbort(c, m.Code(), m.String())
		return
	}

	// 判断是否是被禁用的用户
	if user.Disabled {
		m := msg.LoginDisabledUserFailed
		log.WarnWithFields(log.Fields{
			"username": loginUser["username"],
		}, m.String())
		w.WriteAnyAndAbort(c, m.Code(), m.String())
		return
	}

	// 判断是否是密码为空的用户
	if user.Password == "" {
		m := msg.LoginNoPasswordUserFailed
		log.WarnWithFields(log.Fields{
			"username": loginUser["username"],
		}, m.String())
		w.WriteAnyAndAbort(c, m.Code(), m.String())
		return
	}

	// 判断账号密码是否匹配
	saltedPasswd := utils.GenSha256HashCode(loginUser["password"] + conf.Conf.SecretKey)
	if loginUser["username"] == user.Username && saltedPasswd == user.Password {
		dao.SetUserLastLogin(user.Id)
		// 登录成功并返回token
		newTokenResponse(c, user)
		return
	}

	c.Next()
}

// LoginFailed 返回登录失败
func LoginFailed(c *gin.Context) {
	m := msg.LoginAuthenticationFailed
	log.Warn(m.String())
	w.WriteAnyAndAbort(c, m.Code(), m.String())
	return
}

// updateOrCreateUserLastLogin 更新用户最近登录时间或创建用户并填入response
func updateOrCreateUserLastLogin(c *gin.Context, username string) {
	log.InfoWithFields(log.Fields{
		"username": username,
	}, "updateOrCreateUserLastLogin called.")
	defer log.InfoWithFields(log.Fields{
		"username": username,
	}, "updateOrCreateUserLastLogin end.")

	// 查询该用户在本地数据库中的数据
	user, ok := dao.GetUser(dao.Query{"username=?": username})
	if !ok {
		m := msg.LoginUnknownFailed
		log.ErrorWithFields(log.Fields{
			"username": username,
		}, m.String())
		w.WriteAnyAndAbort(c, m.Code(), m.String())
		return
	}

	// 用户在DB中存在
	if user != nil {
		// 判断是否为禁用的账号
		if user.Disabled {
			m := msg.LoginDisabledUserFailed
			log.WarnWithFields(log.Fields{
				"username": username,
			}, m.String())
			w.WriteAnyAndAbort(c, m.Code(), m.String())
			return
		}

		// 如果查找到用户则更新其登录时间
		dao.SetUserLastLogin(user.Id)

		// 登录成功并返回token
		newTokenResponse(c, user)
		return
	}

	// 在本地数据库中没查找到用户则创建一个用户
	user = dao.NewUser(username, username, true)
	if user != nil {
		// 登录成功并返回token
		newTokenResponse(c, user)
		return
	}

	m := msg.LoginUnknownFailed
	log.ErrorWithFields(log.Fields{
		"username": username,
	}, m.String())
	w.WriteAnyAndAbort(c, m.Code(), m.String())
}

// newTokenResponse 生成Token并填入response
func newTokenResponse(c *gin.Context, userObj *model.User) {
	log.InfoWithFields(log.Fields{
		"user_id":  userObj.Id,
		"username": userObj.Username,
	}, "newTokenResponse called.")
	defer log.InfoWithFields(log.Fields{
		"user_id":  userObj.Id,
		"username": userObj.Username,
	}, "newTokenResponse end.")

	// 登录成功并返回token
	tk := builtin.NewToken(userObj)
	if tk == "" {
		m := msg.LoginNewTokenFailed
		log.ErrorWithFields(log.Fields{
			"user_id":  userObj.Id,
			"username": userObj.Username,
		}, m.String())
		w.WriteAnyAndAbort(c, m.Code(), m.String())
		return
	}

	log.DebugWithFields(log.Fields{
		"user_id":  userObj.Id,
		"username": userObj.Username,
		"token":    tk,
	}, "New token success.")
	w.WriteSuccessPayload(c, "token", tk)
	c.Abort()
	return
}
