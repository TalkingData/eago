package service

import (
	"eago/auth/biz"
	"eago/auth/conf"
	"eago/auth/dao"
	"eago/common/logger"
	"eago/common/redis"
)

type AuthService struct {
	dao   *dao.Dao
	redis *redis.RedisTool

	biz *biz.Biz

	conf   *conf.Conf
	logger *logger.Logger
}

// NewAuthService 新建Auth服务
func NewAuthService(
	dao *dao.Dao, redis *redis.RedisTool, biz *biz.Biz, conf *conf.Conf, logger *logger.Logger,
) *AuthService {
	return &AuthService{
		dao:   dao,
		redis: redis,

		biz: biz,

		conf:   conf,
		logger: logger,
	}
}
