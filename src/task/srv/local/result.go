package local

import (
	"eago/common/log"
	"eago/task/cli"
	"eago/task/conf/msg"
	"eago/task/model"
	"eago/task/worker"
	"errors"
	"fmt"
)

// KillTask 结束任务
func KillTask(taskUniqueId string) error {
	// 将任务唯一Id解码为任务结果Id和分区
	p, id, err := TaskUniqueIdDecode(taskUniqueId)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": taskUniqueId,
			"error":          err.Error(),
		}, "Error in local.TaskUniqueIdDecode.")
		return err
	}

	// 获得任务结果对象
	resObj, ok := model.GetResult(p, id)
	if !ok {
		m := msg.ErrDatabase.SetDetail("Error when TaskService.AppendTaskLog called, in model.GetResult.")
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
		}, m.String())
		return m.Error()
	}
	// 找不到数据的处理
	if resObj == nil {
		m := msg.WarnNotFound.SetDetail("Result object not found.")
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
		}, m.String())
		return m.Error()
	}

	// 任务不是运行状态的，无法手动结束任务
	if resObj.Status != worker.TASK_RUNNING_STATUS {
		err := errors.New("Task not in running status.")
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
			"error":     err.Error(),
		}, "Error in resObj.Status.")
		return err
	}

	// 找不到Worker直接结束
	wk := cli.WorkerClient.GetWorkerById(resObj.Worker)
	if wk == nil {
		// 找不到任务所属的worker
		model.SetResultStatus(p, id, worker.TASK_NO_WORKER_ERROR_END_STATUS, true)
		return fmt.Errorf("No worker for %s", resObj.Worker)
	}

	// 调用Worker
	err = cli.WorkerClient.KillTask(wk, taskUniqueId)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
			"error":     err.Error(),
		}, "Error in cli.WorkerClient.KillTask.")
		return err
	}

	return nil
}
