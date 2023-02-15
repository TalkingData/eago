package handler

import (
	"eago/auth/conf/msg"
	"eago/common/api/ext"
	cMsg "eago/common/code_msg"
	"eago/common/tracer"
	"github.com/gin-gonic/gin"
)

// GetTokenContent 获得TokenContent
func (ah *AuthHandler) GetTokenContent(c *gin.Context) {
	tk := c.GetHeader("Token")

	tc, err := ah.biz.GetTokenContent(tracer.ExtractTraceCtxFromGin(c), tk)
	if err != nil {
		m := msg.MsgGetTokenContentFailed.SetError(err)
		logF := m.ToLoggerFields()
		logF["token"] = tk
		ah.logger.WarnWithFields(logF, m.GetMsg())
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}
	if tc == nil {
		m := cMsg.MsgInvalidTokenFailed
		logF := m.ToLoggerFields()
		logF["token"] = tk
		ah.logger.WarnWithFields(logF, m.GetMsg())
		ext.WriteAnyAndAbort(c, m.GetCode(), m.GetMsg())
		return
	}

	ext.WriteSuccessPayload(c, "content", tc)
}
