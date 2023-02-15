package main

import (
	"eago/common/logger"
	"eago/common/orm"
	"eago/common/service"
	"eago/flow/conf"
	"eago/flow/dao"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	flow service.EagoSrv

	flowDao *dao.Dao

	flowConf *conf.Conf
	flowLg   *logger.Logger
)

func main() {
	flow = NewFlowApi(flowDao, flowConf, flowLg)

	e := make(chan error)
	go func() {
		e <- flow.Start()
	}()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case err := <-e:
			if err != nil {
				flowLg.ErrorWithFields(logger.Fields{
					"error": err,
				}, "An error occurred while Start.")
			}
			closeAll()
			return
		case sig := <-quit:
			flowLg.InfoWithFields(logger.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			closeAll()
			return
		}
	}
}

// closeAll
func closeAll() {
	if flow != nil {
		flow.Stop()
	}
	if flowLg != nil {
		flowLg.Close()
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置
	flowConf = conf.NewConfig()

	// 生成Logger
	lg, err := logger.NewLogger(
		logger.LogLevel(flowConf.LogLevel),
		logger.LogPath(flowConf.LogPath),
		logger.Filename(flowConf.Const.ServiceName, "api"),
	)
	if err != nil {
		fmt.Println("An error occurred while logger.NewLogger, error:", err.Error())
		panic(err)
	}
	flowLg = lg

	flowDao = dao.NewDao(orm.NewMysqlGorm(
		flowConf.MysqlAddress,
		flowConf.MysqlUser,
		flowConf.MysqlPassword,
		flowConf.MysqlDbName,
		orm.MysqlMaxIdleConns(flowConf.MysqlMaxIdleConns),
		orm.MysqlMaxOpenConns(flowConf.MysqlMaxOpenConns),
		orm.UsingOpentracingPlugin(),
	), flowConf, flowLg)
}
