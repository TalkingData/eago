package biz

import (
	"context"
	"eago/common/logger"
)

// MakeUserHandover 用户交接
func (b *Biz) MakeUserHandover(ctx context.Context, srcUserId, tgtUserId uint32) error {
	b.logger.InfoWithFields(logger.Fields{
		"src_user_id": srcUserId,
		"tgt_user_id": tgtUserId,
	}, "Biz.MakeUserHandover called.")
	defer b.logger.Info("Biz.MakeUserHandover end.")

	// 执行交接
	srcUser, tgtUser, err := b.dao.MakeUserHandover(ctx, srcUserId, tgtUserId)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"src_user_id": srcUserId,
			"tgt_user_id": tgtUserId,
			"error":       err,
		}, "An error occurred while dao.MakeUserHandover in Biz.MakeUserHandover.")
		return err
	}

	if srcUser == nil || tgtUser == nil {
		b.logger.ErrorWithFields(logger.Fields{
			"src_user_id": srcUserId,
			"tgt_user_id": tgtUserId,
			"error":       err,
		}, "Got a nil srcUser or tgtUser object in Biz.MakeUserHandover.")
		return err
	}

	bd := map[string]interface{}{
		// 交接用户
		"from": map[string]interface{}{
			"id":       srcUser.Id,
			"username": srcUser.Username,
			"email":    srcUser.Email,
			"phone":    srcUser.Phone,
		},
		// 交接给
		"to": map[string]interface{}{
			"id":       tgtUser.Id,
			"username": tgtUser.Username,
			"email":    tgtUser.Email,
			"phone":    tgtUser.Phone,
		},
	}

	// 发送消息用户交接消息
	// Topic：eago-auth.topic.user.MakeUserHandover
	err = b.pub.Publish(ctx, "user", b.conf.Const.ServiceName, "user", "MakeUserHandover", bd)
	if err != nil {
		b.logger.WarnWithFields(logger.Fields{
			"src_user_id":  srcUserId,
			"src_username": srcUser.Username,
			"tgt_user_id":  tgtUserId,
			"tgt_username": tgtUser.Username,
			"error":        err,
		}, "Failed when broker.Publisher.Publish in Biz.MakeUserHandover.")
	}

	return nil
}
