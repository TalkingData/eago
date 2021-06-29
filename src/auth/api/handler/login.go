package handler

import (
	"bytes"
	"eago/auth/conf"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/auth/srv/dto"
	"eago/auth/srv/local"
	"eago/auth/util/sso"
	"eago/common/log"
	"eago/common/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Login 登录
// @Summary 登录
// @Tags 登录
// @Param data body form.LoginForm true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","token":"2acff1bc1de905d67c1312aa97699dd70c74ade1ad4efb831462ed5122e7d404"}"
// @Router /login [POST]
func Login() {
	// 登录Handler见router.go
	// 此处仅生成swag文档
}

// Logout 登出
// @Summary 登出（销毁Token）
// @Tags 登录
// @Param token header string true "Token"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /logout [DELETE]
func Logout(c *gin.Context) {
	local.RemoveToken(c.GetHeader("Token"))
	msg.Success.GenResponse().Write(c)
}

// Heartbeat 心跳
// @Summary 心跳（续期Token）
// @Tags 登录
// @Param token header string true "Token"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /heartbeat [POST]
func Heartbeat(c *gin.Context) {
	local.RenewalToken(c.GetHeader("Token"))
	msg.Success.GenResponse().Write(c)
}

// GetTokenContent 获得TokenContent
// @Summary 获得TokenContent
// @Tags 登录
// @Param token header string true "Token"
// @Success 200 {string} string "{"code":0,"content":{"groups":[],"is_superuser":false,"own_groups":[],"own_products":[{"id":2,"name":"scm","alias":"scm","disabled":false}],"products":[],"roles":["auth_admin","tester"],"user_id":3,"username":"test"},"message":"Success"}"
// @Router /token/content [GET]
func GetTokenContent(c *gin.Context) {
	tc := c.GetStringMap("TokenContent")

	content := make(map[string]interface{})
	content["user_id"] = tc["UserId"]
	content["username"] = tc["Username"].(string)
	content["phone"] = tc["Phone"].(string)

	content["is_superuser"] = tc["UserIsSuperuser"].(bool)

	content["department"] = tc["Department"].(*[]string)
	content["roles"] = tc["Roles"].(*[]string)
	content["products"] = tc["Products"].(*[]dto.ProductInToken)
	content["own_products"] = tc["OwnProducts"].(*[]dto.ProductInToken)
	content["groups"] = tc["Groups"].(*[]dto.GroupInToken)
	content["own_groups"] = tc["OwnGroups"].(*[]dto.GroupInToken)

	resp := msg.Success.GenResponse()
	resp.SetPayload("content", content)
	resp.Write(c)
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

	iamResp, err := http.Post(conf.Config.IamAddress, "application/json", bytes.NewReader(data))
	if err == nil {
		defer func() {
			_ = iamResp.Body.Close()
		}()

		log.DebugWithFields(log.Fields{
			"response_status_code": iamResp.StatusCode,
			"response_body":        iamResp.Body,
		}, "Got iam response")

		if iamResp.StatusCode == 200 {
			log.Debug("Iam login success.")
			updateOrCreateUserLastLogin(c, loginUser["username"])
			return
		}
	} else {
		log.ErrorWithFields(log.Fields{
			"error": err.Error(),
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
			resp := msg.WarnLoginFailed.GenResponse("Inactive crowd user.")
			log.WarnWithFields(log.Fields{
				"username": loginUser["username"],
			}, resp.String())
			resp.WriteAndAbort(c)
			return
		}

		updateOrCreateUserLastLogin(c, crowdUser.Email)
		return
	}

	log.ErrorWithFields(log.Fields{
		"error": err.Error(),
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
	user, ok := model.GetUser(model.Query{"username=?": loginUser["username"]})
	if !ok {
		// 调用数据库出错
		resp := msg.ErrDatabase.GenResponse("Error when GetUser.")
		log.ErrorWithFields(log.Fields{
			"username": loginUser["username"],
		}, resp.String())
		resp.WriteAndAbort(c)
		return
	}

	// 判断用户是否在DB中存在
	if user == nil {
		resp := msg.WarnLoginFailed.GenResponse("Not exist user.")
		log.WarnWithFields(log.Fields{
			"username": loginUser["username"],
		}, resp.String())
		resp.WriteAndAbort(c)
		return
	}

	// 判断是否是被禁用的用户
	if user.Disabled {
		resp := msg.WarnLoginFailed.GenResponse("Disabled database user.")
		log.WarnWithFields(log.Fields{
			"username": loginUser["username"],
		}, resp.String())
		resp.WriteAndAbort(c)
		return
	}

	// 判断是否是密码为空的用户
	if user.Password == "" {
		resp := msg.WarnLoginFailed.GenResponse("No password user.")
		log.WarnWithFields(log.Fields{
			"username": loginUser["username"],
		}, resp.String())
		resp.WriteAndAbort(c)
		return
	}

	// 判断账号密码是否匹配
	saltedPasswd := utils.GenSha256HashCode(loginUser["password"] + conf.Config.SecretKey)
	if loginUser["username"] == user.Username && saltedPasswd == user.Password {
		model.SetUserLastLogin(user.Id)
		// 登录成功并返回token
		newTokenResponse(c, user)
		return
	}

	c.Next()
}

// LoginFailed 返回登录失败
func LoginFailed(c *gin.Context) {
	resp := msg.WarnLoginFailed.GenResponse("Check username and password.")
	log.Warn(resp.String())
	resp.Write(c)
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
	user, ok := model.GetUser(model.Query{"username=?": username})
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error when GetUser.")
		log.ErrorWithFields(log.Fields{
			"username": username,
		}, resp.String())
		resp.WriteAndAbort(c)
	}

	// 用户在DB中存在
	if user != nil {
		// 判断是否为禁用的账号
		if user.Disabled {
			resp := msg.WarnLoginFailed.GenResponse("Disabled database user.")
			log.WarnWithFields(log.Fields{
				"username": username,
			}, resp.String())
			resp.WriteAndAbort(c)
		}

		// 如果查找到用户则更新其登录时间
		model.SetUserLastLogin(user.Id)

		// 登录成功并返回token
		newTokenResponse(c, user)
		return
	}

	// 在本地数据库中没查找到用户则创建一个用户
	user = model.NewUser(username, username, true)
	if user != nil {
		// 登录成功并返回token
		newTokenResponse(c, user)
		return
	}

	resp := msg.ErrUnknown.GenResponse()
	log.ErrorWithFields(log.Fields{
		"username": username,
	}, resp.String())
	resp.WriteAndAbort(c)
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
	tk := local.NewToken(userObj)
	if tk == "" {
		resp := msg.ErrGenToken.GenResponse()
		log.ErrorWithFields(log.Fields{
			"user_id":  userObj.Id,
			"username": userObj.Username,
		}, resp.String())
		resp.WriteAndAbort(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("token", tk)
	log.DebugWithFields(log.Fields{
		"user_id":  userObj.Id,
		"username": userObj.Username,
		"token":    tk,
	}, resp.String())
	resp.WriteAndAbort(c)
	return
}
