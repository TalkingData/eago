package handler

import (
	"eago-auth/config/msg"
	db "eago-auth/database"
	"eago-auth/srv"
	"eago-auth/util/sso"
	"eago-common/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// GetTokenContent 获得TokenContent
// @Summary 获得TokenContent
// @Tags 登录
// @Param token header string true "Token"
// @Success 200 {string} string "{"code":0,"content":{"groups":[],"is_superuser":false,"own_groups":[],"own_products":[{"id":2,"name":"scm","alias":"scm","disabled":false}],"products":[],"roles":["auth_admin","tester"],"user_id":3,"username":"test"},"message":"Success"}"
// @Router /token/content [GET]
func GetTokenContent(c *gin.Context) {
	tc := c.GetStringMap("TokenContent")
	m := msg.Success.NewMsg()
	m.SetPayload(&gin.H{"content": gin.H{
		"user_id":  tc["UserId"].(int),
		"username": tc["Username"].(string),

		"is_superuser": tc["IsSuperuser"].(bool),
		"roles":        tc["Roles"].(*[]string),
		"products":     tc["Products"].(*[]srv.ProductInToken),
		"own_products": tc["OwnProducts"].(*[]srv.ProductInToken),
		"groups":       tc["Groups"].(*[]srv.GroupInToken),
		"own_groups":   tc["OwnGroups"].(*[]srv.GroupInToken),
	}})
	c.JSON(http.StatusOK, m.GinH())
}

// Heartbeat 心跳
// @Summary 心跳（续期Token）
// @Tags 登录
// @Param token header string true "Token"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /heartbeat [POST]
func Heartbeat(c *gin.Context) {
	srv.RenewalToken(c.GetHeader("token"))
	c.JSON(http.StatusOK, msg.Success.NewMsg().GinH())
}

// Logout 登出
// @Summary 登出（销毁Token）
// @Tags 登录
// @Param token header string true "Token"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /logout [DELETE]
func Logout(c *gin.Context) {
	srv.DeleteToken(c.GetHeader("token"))
	c.JSON(http.StatusOK, msg.Success.NewMsg().GinH())
}

// 从Crowd登录处理
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
			m := msg.WarnLoginFailed.NewMsg("Inactive crowd user.")
			log.WarnWithFields(log.Fields{
				"username": loginUser["username"],
			}, m.String())
			c.AbortWithStatusJSON(http.StatusOK, m.GinH())
			return
		}

		c.AbortWithStatusJSON(
			http.StatusOK,
			updateOrCreateUserLastLogin(crowdUser.Email),
		)
		return
	}

	log.ErrorWithFields(log.Fields{
		"error": err.Error(),
	}, "Crowd login failed.")
	c.Next()
}

// 从数据库登录处理
func DatabaseLogin(c *gin.Context) {
	log.Info("DatabaseLogin called.")
	defer log.Info("DatabaseLogin end.")

	loginUser := c.GetStringMapString("LoginUser")

	log.InfoWithFields(log.Fields{
		"username": loginUser["username"],
	}, "DatabaseLoginHandler called and got a user.")

	// 查询该用户在本地数据库中的数据
	user, suc := db.UserModel.Get(&db.Query{"username=?": loginUser["username"]})
	if !suc {
		// 调用数据库出错
		m := msg.ErrDatabase.NewMsg("Error in db.UserModel.GetUser.")
		log.ErrorWithFields(log.Fields{
			"username": loginUser["username"],
		}, m.String())
		c.AbortWithStatusJSON(http.StatusOK, m.GinH())
		return
	}

	// 判断用户是否在DB中存在
	if user == nil {
		m := msg.WarnLoginFailed.NewMsg("Not exist user.")
		log.WarnWithFields(log.Fields{
			"username": loginUser["username"],
		}, m.String())
		c.AbortWithStatusJSON(http.StatusOK, m.GinH())
		return
	}

	// 判断是否是被禁用的用户
	if user.Disabled {
		m := msg.WarnLoginFailed.NewMsg("Disabled database user.")
		log.WarnWithFields(log.Fields{
			"username": loginUser["username"],
		}, m.String())
		c.AbortWithStatusJSON(http.StatusOK, m.GinH())
		return
	}

	// 判断是否是密码为空的用户
	if user.Password == "" {
		m := msg.WarnLoginFailed.NewMsg("No password user.")
		log.WarnWithFields(log.Fields{
			"username": loginUser["username"],
		}, m.String())
		c.AbortWithStatusJSON(http.StatusOK, m.GinH())
		return
	}

	// 判断账号密码是否匹配
	if loginUser["username"] == user.Username && loginUser["password"] == user.Password {
		db.UserModel.SetLastLogin(&db.Query{"id=?": user.Id})
		// 登录成功并返回token
		c.AbortWithStatusJSON(http.StatusOK, newTokenResponse(user))
		return
	}

	c.Next()
}

// 返回登录失败
func LoginFailed(c *gin.Context) {
	m := msg.WarnLoginFailed.NewMsg("Check username and password.")
	log.Warn(m.String())
	c.JSON(http.StatusOK, m.GinH())
}

// 更新用户最近登录时间或创建用户并填入response
func updateOrCreateUserLastLogin(username string) *gin.H {
	log.InfoWithFields(log.Fields{
		"username": username,
	}, "updateOrCreateUserLastLogin called.")
	defer log.InfoWithFields(log.Fields{
		"username": username,
	}, "updateOrCreateUserLastLogin end.")

	// 查询该用户在本地数据库中的数据
	user, suc := db.UserModel.Get(&db.Query{"username=?": username})
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.UserModel.GetUser.")
		log.ErrorWithFields(log.Fields{
			"username": username,
		}, m.String())
		return m.GinH()
	}

	// 用户在DB中存在
	if user != nil {
		// 判断是否为禁用的账号
		if user.Disabled {
			m := msg.WarnLoginFailed.NewMsg("Disabled database user.")
			log.WarnWithFields(log.Fields{
				"username": username,
			}, m.String())
			return m.GinH()
		}

		// 如果查找到用户则更新其登录时间
		db.UserModel.SetLastLogin(&db.Query{"id=?": user.Id})

		// 登录成功并返回token
		return newTokenResponse(user)
	}

	// 在本地数据库中没查找到用户则创建一个用户
	user = db.UserModel.New(username, username, true)
	if user != nil {
		// 登录成功并返回token
		return newTokenResponse(user)
	}

	m := msg.ErrUnknown.NewMsg()
	log.ErrorWithFields(log.Fields{
		"username": username,
	}, m.String())
	return m.GinH()
}

// 生成Token并填入response
func newTokenResponse(userObj *db.User) *gin.H {
	log.InfoWithFields(log.Fields{
		"user_id":  userObj.Id,
		"username": userObj.Username,
	}, "newTokenResponse called.")
	defer log.InfoWithFields(log.Fields{
		"user_id":  userObj.Id,
		"username": userObj.Username,
	}, "newTokenResponse end.")

	// 登录成功并返回token
	tk := srv.NewToken(userObj)
	if tk == "" {
		m := msg.ErrGenToken.NewMsg()
		log.ErrorWithFields(log.Fields{
			"user_id":  userObj.Id,
			"username": userObj.Username,
		}, m.String())
		return m.GinH()
	}

	m := msg.Success.NewMsg()
	m.SetPayload(&gin.H{"token": tk})
	log.DebugWithFields(log.Fields{
		"user_id":  userObj.Id,
		"username": userObj.Username,
		"token":    tk,
	}, m.String())
	return m.GinH()
}
