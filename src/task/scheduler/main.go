package main

import (
	"context"
	"eago/common/logger"
	"eago/common/service"
	"eago/task/conf"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	schSrv service.EagoSrv

	schConf *conf.Conf
	schLg   *logger.Logger
)

func main() {
	schSrv = NewScheduler(
		context.Background(),
		EtcdAddresses(schConf.EtcdAddresses),
		EtcdUsername(schConf.EtcdUsername),
		EtcdPassword(schConf.EtcdPassword),
		RegisterTtl(schConf.SchedulerRegisterTtl),
		TaskRpcRegisterKey(schConf.Const.RpcRegisterKey),
		Logger(schLg),
	)

	e := make(chan error)
	go func() {
		e <- schSrv.Start()
	}()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case err := <-e:
			if err != nil {
				schLg.ErrorWithFields(logger.Fields{
					"error": err,
				}, "An error occurred while Start.")
			}
			closeAll()
			return
		case sig := <-quit:
			schLg.InfoWithFields(logger.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			closeAll()
			return
		}
	}
}

// closeAll
func closeAll() {
	if schSrv != nil {
		schSrv.Stop()
	}
	if schLg != nil {
		schLg.Close()
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置
	schConf = conf.NewConfig()

	// 生成Logger
	lg, err := logger.NewLogger(
		logger.LogLevel(schConf.LogLevel),
		logger.LogPath(schConf.LogPath),
		logger.Filename(schConf.Const.ServiceName, "scheduler"),
	)
	if err != nil {
		fmt.Println("An error occurred while logger.NewLogger, error:", err.Error())
		panic(err)
	}
	schLg = lg
}
