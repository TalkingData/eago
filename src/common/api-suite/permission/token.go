package permission

import (
	auth "eago/auth/srv/proto"
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/common/tracer"
	"github.com/gin-gonic/gin"
	"net/http"
)

// MustLogin 验证是否登录并装载TokenContent
func MustLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tk auth.Token

		tk.Token = c.GetHeader("Token")
		ctx := tracer.ExtractTraceContext(c)
		tc, err := authCli.GetTokenContent(ctx, &tk)
		if err != nil {
			m := "An error occurred while authCli.GetTokenContent."
			log.WarnWithFields(log.Fields{
				"token": tk.Token,
				"error": err,
			}, m)
			w.WriteAnyAndAbort(c, http.StatusInternalServerError, m)
			return
		} else if !tc.Ok {
			log.WarnWithFields(log.Fields{
				"token": tk.Token,
				"error": err,
			}, InvalidToken.String())
			w.WriteAnyAndAbort(c, InvalidToken.Code(), InvalidToken.String())
			return
		}

		c.Set("TokenContent", map[string]interface{}{
			"UserId":      tc.UserId,
			"Username":    tc.Username,
			"Phone":       tc.Phone,
			"IsSuperuser": tc.IsSuperuser,

			"Department":  tc.Department,
			"Roles":       tc.Roles,
			"Products":    tc.Products,
			"OwnProducts": tc.OwnProducts,
			"Groups":      tc.Groups,
			"OwnGroups":   tc.OwnGroups,
		})
	}
}

// MustLoginWs WebSocket验证是否登录并装载TokenContent
func MustLoginWs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tk auth.Token

		tk.Token = c.DefaultQuery("token", "")
		ctx := tracer.ExtractTraceContext(c)
		tc, err := authCli.GetTokenContent(ctx, &tk)
		if err != nil {
			m := "An error occurred while authCli.GetTokenContent."
			log.WarnWithFields(log.Fields{
				"token": tk.Token,
				"error": err,
			}, m)
			w.WriteAnyAndAbort(c, http.StatusInternalServerError, m)
			return
		} else if !tc.Ok {
			log.WarnWithFields(log.Fields{
				"token": tk.Token,
				"error": err,
			}, InvalidToken.String())
			w.WriteAnyAndAbort(c, InvalidToken.Code(), InvalidToken.String())
			return
		}

		c.Set("TokenContent", map[string]interface{}{
			"UserId":      tc.UserId,
			"Username":    tc.Username,
			"Phone":       tc.Phone,
			"IsSuperuser": tc.IsSuperuser,

			"Department":  tc.Department,
			"Roles":       tc.Roles,
			"Products":    tc.Products,
			"OwnProducts": tc.OwnProducts,
			"Groups":      tc.Groups,
			"OwnGroups":   tc.OwnGroups,
		})
	}
}
