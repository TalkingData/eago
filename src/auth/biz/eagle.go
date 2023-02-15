package biz

import (
	"context"
	"eago/auth/dto"
	"eago/auth/model"
	"eago/common/logger"
	"eago/common/orm"
	"encoding/json"
	"errors"
	"fmt"
)

// GetUserObjectFromEagleToken 本地服务::根据EagleToken获取用户对象
func (b *Biz) GetUserObjectFromEagleToken(ctx context.Context, tkStr string) (*model.User, error) {
	tkKey := fmt.Sprintf("%s_%s", b.conf.Const.EagleTokenKeyPrefix, tkStr)

	if !b.redis.Exist(ctx, tkKey) {
		// 找不到eagle token，退出
		return nil, fmt.Errorf("eagle token not found")
	}

	// 获取eagle token content
	tkContent, err := b.redis.DirectGet(ctx, tkKey)
	if err != nil {
		return nil, err
	}

	// 反序列化eagle token
	tkObj := new(dto.EagleToken)
	// 反序列化失败，退出
	if err = json.Unmarshal([]byte(tkContent), tkObj); err != nil {
		return nil, err
	}

	// 查询该用户在本地数据库中的数据
	userObj, err := b.dao.GetUser(ctx, orm.Query{"username=?": tkObj.Username, "disabled=?": 0})
	if err != nil {
		// 调用数据库出错
		b.logger.ErrorWithFields(logger.Fields{
			"username": tkObj.Username,
			"error":    err,
		}, "An error occurred while dao.GetUser in Biz.GetUserObjectFromEagleToken.")
		return nil, fmt.Errorf("find user error")
	}
	// 找不到用户
	if userObj == nil || userObj.Id < 1 {
		b.logger.ErrorWithFields(logger.Fields{
			"username": tkObj.Username,
		}, "An nil object is returned after calling dao.GetUser in Biz.GetUserObjectFromEagleToken.")
		return nil, errors.New("got an nil user")
	}

	return userObj, nil
}
