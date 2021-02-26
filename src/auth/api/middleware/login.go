package middleware

import (
	"eago-auth/api/form"
	"eago-auth/config/msg"
	"eago-common/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Login 登录
// @Summary 登录
// @Tags 登录
// @Param data body form.LoginForm true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","token":"2acff1bc1de905d67c1312aa97699dd70c74ade1ad4efb831462ed5122e7d404"}"
// @Router /login [POST]
func ReadLoginForm() gin.HandlerFunc { // 登录表单预读
	return func(c *gin.Context) {
		log.Info("ReadLoginForm called.")
		defer log.Info("ReadLoginForm end.")

		var lf form.LoginForm

		// 序列化request body获取用户名密码
		if err := c.ShouldBindJSON(&lf); err != nil {
			m := msg.WarnInvalidBody.NewMsg("Field 'username', 'password' required.")
			log.WarnWithFields(log.Fields{
				"error": err.Error(),
			}, m.String())
			c.AbortWithStatusJSON(http.StatusOK, m.GinH())
			return
		}

		c.Set("LoginUser", map[string]string{
			"username": lf.Username, "password": lf.Password,
		})
		c.Next()
	}
}
