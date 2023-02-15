package biz

import (
	"context"
	"eago/common/global"
	"eago/common/logger"
	"eago/common/utils"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
	"time"
)

// VerifyAndUnregisterSrvToken 验证SrvToken，验证成功后会注销token
func (b *Biz) VerifyAndUnregisterSrvToken(ctx context.Context, srvToken string) bool {
	ok := b.isValidSrvToken(ctx, srvToken)
	if ok {
		b.unregisterSrvToken(ctx, srvToken)
	}

	return ok
}

// NewSrvTokenWithCtx 新建SrvToken并绑定到context
func (b *Biz) NewSrvTokenWithCtx(ctx context.Context, tokenContent string) context.Context {
	srvToken := genSrvToken(tokenContent)
	b.registerSrvToken(ctx, srvToken)
	return bindSrvToken2Ctx(ctx, srvToken)
}

// registerSrvToken 注册SrvToken
func (b *Biz) registerSrvToken(ctx context.Context, srvToken string) {
	if err := b.redis.Set(ctx, genSrvTokenKey(srvToken), "", b.conf.SrvTokenTtlSecs); err != nil {
		b.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while registerSrvToken.")
	}
}

// unregisterSrvToken 注销Srv的Token
func (b *Biz) unregisterSrvToken(ctx context.Context, srvToken string) {
	if err := b.redis.Del(ctx, genSrvTokenKey(srvToken)); err != nil {
		b.logger.WarnWithFields(logger.Fields{
			"srv_token": srvToken,
			"error":     err,
		}, "An error occurred while biz.unregisterSrvToken, Skipped it.")
	}
}

// isValidSrvToken 判断Srv是否有效
func (b *Biz) isValidSrvToken(ctx context.Context, srvToken string) bool {
	// 如果Key不存在，肯定没有注册srv
	return b.redis.Exist(ctx, genSrvTokenKey(srvToken))
}

// genSrvToken 生成SrvToken
func genSrvToken(tokenContent string) string {
	// 计算srv token值
	return utils.GenSha256HashCode(
		fmt.Sprintf("%s%s%s", tokenContent, time.Now().Format(global.TimestampFormat), uuid.NewV4().String()),
	)
}

// bindSrvToken2Ctx 绑定SrvToken到context
func bindSrvToken2Ctx(ctx context.Context, srvToken string) context.Context {
	return metadata.NewOutgoingContext(
		ctx,
		metadata.New(map[string]string{"srv_token": srvToken}),
	)
}

// genSrvTokenKey 生成Srv的Token key
func genSrvTokenKey(srvToken string) string {
	return fmt.Sprintf("srv_token/%s", srvToken)
}
