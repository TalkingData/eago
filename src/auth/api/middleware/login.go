package middleware

import (
	"eago/auth/conf/msg"
	"eago/auth/dto"
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"github.com/gin-gonic/gin"
)

// ReadLoginForm 登录表单预读
func ReadLoginForm() gin.HandlerFunc { // 登录表单预读
	return func(c *gin.Context) {
		log.Info("ReadLoginForm called.")
		defer log.Info("ReadLoginForm end.")

		var loginFrm dto.Login
		// 序列化request body获取用户名密码
		if err := c.ShouldBindJSON(&loginFrm); err != nil {
			m := msg.SerializeFailed
			w.WriteAnyAndAbort(c, m.Code(), m.String())
			return
		}
		// 验证数据
		if err := loginFrm.Validate(); err != nil {
			// 数据验证未通过
			m := msg.ValidateFailed
			log.WarnWithFields(m.LogFields())
			w.WriteAnyAndAbort(c, m.Code(), m.String())
			return
		}

		c.Set("LoginUser", map[string]string{
			"username": loginFrm.Username, "password": loginFrm.Password,
		})
	}
}
