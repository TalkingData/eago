package main

import (
	"eago/common/logger"
	"eago/common/orm"
	"eago/common/service"
	"eago/task/conf"
	"eago/task/dao"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	task service.EagoSrv

	taskDao *dao.Dao

	taskConf *conf.Conf
	taskLg   *logger.Logger
)

func main() {
	task = NewTaskApi(taskDao, taskConf, taskLg)

	e := make(chan error)
	go func() {
		e <- task.Start()
	}()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case err := <-e:
			if err != nil {
				taskLg.ErrorWithFields(logger.Fields{
					"error": err,
				}, "An error occurred while Start.")
			}
			closeAll()
			return
		case sig := <-quit:
			taskLg.InfoWithFields(logger.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			closeAll()
			return
		}
	}
}

// closeAll
func closeAll() {
	if task != nil {
		task.Stop()
	}
	if taskLg != nil {
		taskLg.Close()
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置
	taskConf = conf.NewConfig()

	// 生成Logger
	lg, err := logger.NewLogger(
		logger.LogLevel(taskConf.LogLevel),
		logger.LogPath(taskConf.LogPath),
		logger.Filename(taskConf.Const.ServiceName, "api"),
	)
	if err != nil {
		fmt.Println("An error occurred while logger.NewLogger, error:", err.Error())
		panic(err)
	}
	taskLg = lg

	taskDao = dao.NewDao(orm.NewMysqlGorm(
		taskConf.MysqlAddress,
		taskConf.MysqlUser,
		taskConf.MysqlPassword,
		taskConf.MysqlDbName,
		orm.MysqlMaxIdleConns(taskConf.MysqlMaxIdleConns),
		orm.MysqlMaxOpenConns(taskConf.MysqlMaxOpenConns),
		orm.UsingOpentracingPlugin(),
	), taskConf, taskLg)
}
