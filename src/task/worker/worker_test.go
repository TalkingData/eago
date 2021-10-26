package worker

import (
	"context"
	"eago/common/log"
	"eago/task/conf"
	"fmt"
	"runtime"
	"testing"
)

func TestWorker(t *testing.T) {
	wk := NewWorker(
		EtcdAddresses(conf.Conf.EtcdAddresses),
		EtcdUsername(conf.Conf.EtcdUsername),
		EtcdPassword(conf.Conf.EtcdPassword),
		ServiceName("test"),
	)

	wk.RegTask("test_worker", foo)
	// 启动Worker
	if err := wk.Start(); err != nil {
		log.Error(err.Error())
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

func foo(ctx context.Context, param *Param) error {
	defer param.Log.Info("test_worker ended.")

	param.Log.Info("Got test_worker call, and there is foo.")
	param.Log.Info("Your arguments is: %s", param.Arguments)

	return nil
}
