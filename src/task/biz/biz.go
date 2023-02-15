package biz

import (
	"eago/common/logger"
	"eago/common/redis"
	"eago/task/conf"
	"eago/task/dao"
	workerCli "eago/task/worker/client"
)

type Biz struct {
	dao   *dao.Dao
	redis *redis.RedisTool

	workerCli *workerCli.WorkerClient

	conf   *conf.Conf
	logger *logger.Logger
}

func NewBiz(dao *dao.Dao, redis *redis.RedisTool, conf *conf.Conf, logger *logger.Logger) *Biz {
	return &Biz{
		dao:   dao,
		redis: redis,

		workerCli: workerCli.NewWorkerCli(conf.EtcdAddresses, logger),

		conf:   conf,
		logger: logger,
	}
}
