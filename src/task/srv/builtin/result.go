package builtin

import (
	"eago/common/log"
	"eago/task/cli"
	"eago/task/dao"
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
			"error":          err,
		}, "An error occurred while builtin.TaskUniqueIdDecode.")
		return err
	}

	// 获得任务结果对象
	resObj, ok := dao.GetResult(p, id)
	if !ok {
		// 获得任务时返回了空对象
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"id":        id,
		}, "An nil object is returned after calling dao.GetResult.")
		return errors.New("an nil object is returned after calling dao.GetResult")
	}
	// 找不到数据的处理
	if resObj == nil {
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
		}, "Result object not found.")
		return errors.New("result object not found")
	}

	// 任务不是运行状态的，无法手动结束任务
	if resObj.Status != worker.TASK_RUNNING_STATUS {
		log.ErrorWithFields(log.Fields{
			"partition":   p,
			"result_id":   id,
			"task_status": resObj.Status,
		}, "Can not kill task, it is not in running status.")
		return errors.New("task is not in running status")
	}

	// 找不到Worker直接结束
	wk := cli.WorkerClient.GetWorkerById(resObj.Worker)
	if wk == nil {
		// 找不到任务所属的worker
		dao.SetResultStatus(p, id, worker.TASK_NO_WORKER_ERROR_END_STATUS, true)
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
			"worker":    resObj.Worker,
		}, "Can not kill task, no worker found.")
		return fmt.Errorf("worker not found for %s", resObj.Worker)
	}

	// 调用Worker
	err = cli.WorkerClient.KillTask(wk, taskUniqueId)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
			"error":     err,
		}, "An error occurred while cli.WorkerClient.KillTask.")
		return err
	}

	return nil
}
