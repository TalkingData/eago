package worker

import (
	"context"
	"fmt"
	"time"
)

type TaskFunc func(ctx context.Context, param *Param) error

// Param struct
type Param struct {
	TaskUniqueId    string
	Caller          string
	Timeout         int32
	Arguments       string
	LocalStartTime  time.Time
	RemoteStartTime time.Time
	Log             Logger
}

// Task struct
type Task struct {
	Codename string
	Param    *Param
	Cxt      context.Context
	Cancel   context.CancelFunc
	fn       TaskFunc
	logger   *logger
}

// Run 运行任务
func (t *Task) Run(callback func(err error)) {
	err := t.fn(t.Cxt, t.Param)
	if t.Cxt.Err() != nil {
		t.Param.Log.Error("Task '%s-%s' failed with context error: '%s'.", t.Codename, t.Param.TaskUniqueId, t.Cxt.Err())
		callback(t.Cxt.Err())
		return
	}

	if err != nil {
		t.Param.Log.Error("Task '%s-%s' failed with returned error: '%s'.", t.Codename, t.Param.TaskUniqueId, err.Error())
	}
	callback(err)
}

func (t *Task) Info() string {
	return fmt.Sprintf("TaskCodename: %s.", t.Codename)
}
