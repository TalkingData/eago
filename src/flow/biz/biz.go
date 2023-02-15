package biz

import (
	authpb "eago/auth/proto"
	"eago/cli"
	"eago/common/broker"
	"eago/common/logger"
	"eago/flow/conf"
	"eago/flow/dao"
	taskpb "eago/task/proto"
)

type Biz struct {
	dao *dao.Dao

	pub broker.Publisher

	authCli authpb.AuthService
	taskCli taskpb.TaskService

	conf   *conf.Conf
	logger *logger.Logger
}

func NewBiz(dao *dao.Dao, pub broker.Publisher, _conf *conf.Conf, logger *logger.Logger) *Biz {
	return &Biz{
		dao: dao,

		pub: pub,

		// 创建Auth客户端
		authCli: cli.NewAuthClient(_conf.EtcdUsername, _conf.EtcdPassword, _conf.EtcdAddresses),
		taskCli: cli.NewTaskClient(_conf.EtcdUsername, _conf.EtcdPassword, _conf.EtcdAddresses),

		conf:   _conf,
		logger: logger,
	}
}
