package builtin

import (
	"eago/common/log"
	"eago/task/dao"
)

// NewLog 新增一条任务Log
func NewLog(taskUniqueId string, content *string) error {
	// 将任务唯一Id解码为任务结果Id和分区
	p, id, err := TaskUniqueIdDecode(taskUniqueId)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": taskUniqueId,
			"error":          err,
		}, "An error occurred while builtin.TaskUniqueIdDecode.")
		return err
	}

	l := dao.NewLog(p, id, content)
	if l == nil {
		// 创建日志时返回了空对象
		log.WarnWithFields(log.Fields{
			"partition": p,
			"id":        id,
		}, "An nil object is returned after calling model.NewLog, This error will be ignored.")
		return nil
	}

	log.DebugWithFields(log.Fields{
		"partition": p,
		"result_id": id,
		"log_id":    l.Id,
	}, "New log success.")
	return nil
}
