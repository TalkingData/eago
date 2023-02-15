package handler

import (
	authpb "eago/auth/proto"
	"eago/cli"
	"eago/common/api/menu"
	"eago/common/logger"
	"eago/task/biz"
	"eago/task/conf"
	"eago/task/dao"
	taskpb "eago/task/proto"
	workerCli "eago/task/worker/client"
)

type TaskHandler struct {
	dao *dao.Dao

	biz *biz.Biz

	taskCli   taskpb.TaskService
	workerCli *workerCli.WorkerClient

	authCli authpb.AuthService

	menu *menu.Menu

	conf   *conf.Conf
	logger *logger.Logger
}

func NewTaskHandler(dao *dao.Dao, _conf *conf.Conf, _logger *logger.Logger) *TaskHandler {
	return &TaskHandler{
		dao: dao,

		// 生成Biz
		biz: biz.NewBiz(dao, nil, _conf, _logger),

		// 创建Auth客户端
		taskCli: cli.NewTaskClient(_conf.EtcdUsername, _conf.EtcdPassword, _conf.EtcdAddresses),
		// 创建Worker客户端
		workerCli: workerCli.NewWorkerCli(_conf.EtcdAddresses, _logger),

		// 创建Auth客户端
		authCli: cli.NewAuthClient(_conf.EtcdUsername, _conf.EtcdPassword, _conf.EtcdAddresses),

		menu: conf.NewMenu(_conf),

		conf:   _conf,
		logger: _logger,
	}
}

func (th *TaskHandler) GetAuthCli() authpb.AuthService {
	return th.authCli
}
