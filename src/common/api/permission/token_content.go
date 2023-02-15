package permission

import (
	authpb "eago/auth/proto"
	"github.com/gin-gonic/gin"
)

func MustGetTokenContent(c *gin.Context) *authpb.TokenContent {
	return c.MustGet(defaultTokenContentGinCtxKey).(*authpb.TokenContent)
}

func GetTokenContent(c *gin.Context) (*authpb.TokenContent, bool) {
	val, ok := c.Get(defaultTokenContentGinCtxKey)
	if !ok || val == nil {
		return nil, false
	}

	tContent := val.(*authpb.TokenContent)
	if tContent.UserId < 1 {
		return nil, false
	}

	return tContent, true
}
