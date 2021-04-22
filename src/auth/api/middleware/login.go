package middleware

import (
	"eago-auth/api/form"
	"eago-auth/conf/msg"
	"eago-common/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ReadLoginForm 登录表单预读
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
	}
}
