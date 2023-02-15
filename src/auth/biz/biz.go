package biz

import (
	"eago/auth/conf"
	"eago/auth/dao"
	"eago/common/broker"
	"eago/common/logger"
	"eago/common/redis"
)

type Biz struct {
	dao   *dao.Dao
	redis *redis.RedisTool
	pub   broker.Publisher

	conf   *conf.Conf
	logger *logger.Logger
}

func NewBiz(dao *dao.Dao, redis *redis.RedisTool, pub broker.Publisher, conf *conf.Conf, logger *logger.Logger) *Biz {
	return &Biz{
		dao:   dao,
		redis: redis,
		pub:   pub,

		conf:   conf,
		logger: logger,
	}
}
