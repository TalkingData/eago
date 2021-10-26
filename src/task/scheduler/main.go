package main

import (
	"eago/common/log"
	"eago/task/conf"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	scheduler := NewScheduler(
		EtcdAddresses(conf.Conf.EtcdAddresses),
		EtcdUsername(conf.Conf.EtcdUsername),
		EtcdPassword(conf.Conf.EtcdPassword),
	)

	go scheduler.Start()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case sig := <-quit:
			log.InfoWithFields(log.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			scheduler.Stop()
			return
		}
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 加载日志设置
	err := log.InitLog(
		conf.Conf.LogPath,
		conf.SERVICE_NAME,
		conf.Conf.LogLevel,
	)
	if err != nil {
		fmt.Println("Failed to init logging, error:", err.Error())
		panic(err)
	}
}
