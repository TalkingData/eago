package handler

import (
	authpb "eago/auth/proto"
	"eago/cli"
	"eago/common/api/menu"
	"eago/common/logger"
	"eago/flow/biz"
	"eago/flow/conf"
	"eago/flow/dao"
)

type FlowHandler struct {
	dao *dao.Dao

	biz *biz.Biz

	authCli authpb.AuthService

	menu *menu.Menu

	conf   *conf.Conf
	logger *logger.Logger
}

func NewFlowHandler(dao *dao.Dao, biz *biz.Biz, _conf *conf.Conf, _logger *logger.Logger) *FlowHandler {
	return &FlowHandler{
		dao: dao,

		biz: biz,

		// 创建Auth客户端
		authCli: cli.NewAuthClient(_conf.EtcdUsername, _conf.EtcdPassword, _conf.EtcdAddresses),

		menu: conf.NewMenu(_conf),

		conf:   _conf,
		logger: _logger,
	}
}

func (h *FlowHandler) GetAuthCli() authpb.AuthService {
	return h.authCli
}
