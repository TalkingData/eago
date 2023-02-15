package handler

import (
	"eago/auth/biz"
	"eago/auth/conf"
	"eago/auth/dao"
	authpb "eago/auth/proto"
	"eago/cli"
	"eago/common/api/menu"
	"eago/common/broker"
	"eago/common/logger"
	"eago/common/redis"
	"fmt"
	"github.com/jda/go-crowd"
)

type AuthHandler struct {
	dao   *dao.Dao
	redis *redis.RedisTool

	biz *biz.Biz

	authCli authpb.AuthService

	menu *menu.Menu

	crowdCli *crowd.Crowd

	conf   *conf.Conf
	logger *logger.Logger
}

func NewAuthHandler(_dao *dao.Dao, redis *redis.RedisTool, _conf *conf.Conf, _logger *logger.Logger) *AuthHandler {
	// 生成Broker
	_broker, err := broker.NewKafkaBroker(_conf.KafkaAddresses)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init KafkaBroker: %v\n", err))
	}

	// 生成Publisher
	_pub := broker.NewPublisher(
		_broker,
		broker.ServiceName(_conf.Const.ServiceName),
		broker.Logger(_logger),
	)

	// 生成crowdCli
	_crowdCli, err := crowd.New(_conf.CrowdAppName, _conf.CrowdAppPassword, _conf.CrowdAddress)
	if err != nil {
		_logger.WarnWithFields(logger.Fields{
			"crowd_address":  _conf.CrowdAddress,
			"crowd_app_name": _conf.CrowdAppName,
		}, "An error occurred while crowd.New in NewAuthHandler, skipped it.")
	}

	return &AuthHandler{
		dao:   _dao,
		redis: redis,

		// 生成Biz
		biz: biz.NewBiz(_dao, redis, _pub, _conf, _logger),

		// 创建Auth客户端
		authCli: cli.NewAuthClient(_conf.EtcdUsername, _conf.EtcdPassword, _conf.EtcdAddresses),

		menu: conf.NewMenu(_conf),

		crowdCli: &_crowdCli,

		conf:   _conf,
		logger: _logger,
	}
}

func (ah *AuthHandler) GetAuthCli() authpb.AuthService {
	return ah.authCli
}
