package service

import (
	"eago/common/logger"
	"eago/common/redis"
	"eago/task/biz"
	"eago/task/conf"
	"eago/task/dao"
)

type TaskService struct {
	dao   *dao.Dao
	redis *redis.RedisTool

	biz *biz.Biz

	conf   *conf.Conf
	logger *logger.Logger
}

// NewTaskService 新建Task服务
func NewTaskService(
	dao *dao.Dao, redis *redis.RedisTool, biz *biz.Biz, conf *conf.Conf, logger *logger.Logger,
) *TaskService {
	return &TaskService{
		dao:   dao,
		redis: redis,

		biz: biz,

		conf:   conf,
		logger: logger,
	}
}
