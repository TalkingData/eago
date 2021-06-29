package local

import (
	"eago/common/log"
	"eago/task/conf/msg"
	"eago/task/model"
	"errors"
)

// NewLog 新增一条任务Log
func NewLog(taskUniqueId string, content *string) error {
	// 将任务唯一Id解码为任务结果Id和分区
	p, id, err := TaskUniqueIdDecode(taskUniqueId)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": taskUniqueId,
			"error":          err.Error(),
		}, "Error in local.TaskUniqueIdDecode.")
		return err
	}

	l := model.NewLog(p, id, content)
	if l == nil {
		m := msg.ErrDatabase.SetDetail("Error in model.NewLog.")
		err := errors.New(m.String())
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"id":        id,
		}, m.String())
		return err
	}

	log.DebugWithFields(log.Fields{
		"partition": p,
		"result_id": id,
		"log_id":    l.Id,
	}, "TaskLog success.")
	return nil
}
