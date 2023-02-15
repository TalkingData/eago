package service

import (
	"context"
	"eago/auth/conf/msg"
	authpb "eago/auth/proto"
	cMsg "eago/common/code_msg"
	"eago/common/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

// VerifyToken 验证Token
func (authSrv *AuthService) VerifyToken(ctx context.Context, req *authpb.Token, _ *emptypb.Empty) error {
	authSrv.logger.DebugWithFields(logger.Fields{
		"token": "token",
	}, "authSrv.VerifyToken called.")
	defer authSrv.logger.Debug("authSrv.VerifyToken end.")

	if !authSrv.biz.VerifyToken(ctx, req.Value) {
		return cMsg.MsgInvalidTokenFailed.ToMicroErr()
	}
	return nil
}

// GetTokenContent 获取Token内容
func (authSrv *AuthService) GetTokenContent(
	ctx context.Context, req *authpb.Token, rsp *authpb.TokenContent,
) error {
	authSrv.logger.DebugWithFields(logger.Fields{
		"token": "token",
	}, "authSrv.GetTokenContent called.")
	defer authSrv.logger.Debug("authSrv.GetTokenContent end.")

	tc, err := authSrv.biz.GetTokenContent(ctx, req.Value)
	if err != nil {
		return msg.MsgGetTokenContentFailed.ToMicroErr()
	}
	if tc == nil {
		return cMsg.MsgInvalidTokenFailed.ToMicroErr()
	}

	tc.Trans2AuthPb(rsp)
	return nil
}
