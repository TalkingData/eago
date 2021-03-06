package worker

import (
	"context"
	"fmt"
	"time"
)

type TaskFunc func(ctx context.Context, param *Param) error

type Param struct {
	TaskUniqueId    string
	Caller          string
	Timeout         int64
	Arguments       string
	LocalStartTime  time.Time
	RemoteStartTime time.Time
	Log             Logger
}

type Task struct {
	Codename string
	Param    *Param
	Cxt      context.Context
	Cancel   context.CancelFunc
	fn       TaskFunc
	logger   *logger
}

// Run 运行任务
func (t *Task) Run(callback func(ok bool)) {
	err := t.fn(t.Cxt, t.Param)
	if err != nil {
		t.Param.Log.Error("Task \"%s-%s\" failed with returned error \"%s\".", t.Codename, t.Param.TaskUniqueId, err.Error())
		callback(false)
		return
	}

	callback(true)
}

func (t *Task) Info() string {
	return fmt.Sprintf("TaskCodename: %s.", t.Codename)
}
