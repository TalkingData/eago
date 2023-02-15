package permission

import (
	authpb "eago/auth/proto"
	"eago/common/api/ext"
	cMsg "eago/common/code_msg"
	"eago/common/tracer"
	"github.com/gin-gonic/gin"
)

func isRoleHandler(c *gin.Context, role string) {
	ok, err := IsRole(c, role)
	if err != nil {
		m := cMsg.MsgCheckRoleErr.SetError(err)
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}

	if !ok {
		m := cMsg.MsgUserNotRoleFailed.SetDetail(role)
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}
}

func getTokenContentHandler(c *gin.Context, token string, authCli authpb.AuthService) {
	if len(token) < 1 {
		m := cMsg.MsgInvalidTokenFailed
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
	}

	tc, err := authCli.GetTokenContent(tracer.ExtractTraceCtxFromGin(c), &authpb.Token{Value: token})
	if err != nil {
		m, ok := cMsg.TransMicroErr2CodeMsg(err)
		// 如果err可以被转换为code_msg可直接用其返回，否则创建一个新的code_msg返回
		if !ok || m == nil {
			m = msgAuthClientErr.SetError(err)
		}
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}

	// 获得了无效的Token
	if tc == nil || tc.UserId < 1 {
		m := cMsg.MsgInvalidUriFailed.SetError(err)
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}

	// 成功，将TokenContent写入gin.Context
	c.Set(defaultTokenContentGinCtxKey, tc)
}
