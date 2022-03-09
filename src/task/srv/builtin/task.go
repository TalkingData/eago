package builtin

import (
	"eago/common/log"
	"eago/task/cli"
	"eago/task/conf"
	"eago/task/dao"
	"eago/task/worker"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// CallTask 调用任务
func CallTask(taskCodename, arguments, caller string, timeout int64) (taskUniqueId string, err error) {
	taskResult := dao.NewResult(taskCodename, caller, arguments, timeout, worker.TASK_INITIALIZATION_STATUS)
	if taskResult == nil {
		err = fmt.Errorf("result object not found")
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An nil object is returned after calling dao.NewResult.")
		return
	}
	partition, partitionErr := taskResult.GetPartition()
	if partitionErr != nil {
		log.ErrorWithFields(log.Fields{
			"error": partitionErr.Error(),
		}, "An error occurred while Result.GetPartition.")
		return "", partitionErr
	}

	cNameSplit := strings.Split(taskCodename, ".")
	if len(cNameSplit) < 2 {
		err = fmt.Errorf("invalid task codename")
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while strings.Split for taskCodename.")
		dao.SetResultStatus(partition, taskResult.Id, worker.TASK_CALL_ERROR_END_STATUS, true)
		return "", err
	}

	// 查找对应模块的worker
	modular := cNameSplit[0]
	wks := cli.WorkerClient.ListByModular(modular)
	if wks == nil || len(wks) < 1 {
		// 找不到模块所属的worker
		dao.SetResultStatus(partition, taskResult.Id, worker.TASK_NO_WORKER_ERROR_END_STATUS, true)
		log.ErrorWithFields(log.Fields{
			"worker": modular,
		}, "Can not kill task, no worker found.")
		return "", fmt.Errorf("no worker found for %s", modular)
	}

	// 生成任务实例唯一ID
	taskUniqueId = TaskUniqueIdEncode(partition, taskResult.Id)

	// 随机找一个Worker
	w := wks[rand.Intn(len(wks))]
	// 调用Worker
	err = cli.WorkerClient.CallTask(w, cNameSplit[1], taskUniqueId, arguments, caller, timeout, taskResult.StartAt.Unix())
	if err != nil {
		// 任务调用错误
		dao.SetResultStatus(partition, taskResult.Id, worker.TASK_CALL_ERROR_END_STATUS, true)
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"result_id": taskResult.Id,
			"error":     err,
		}, "An error occurred while cli.WorkerClient.KillTask.")
		return "", err
	}
	// 填充执行任务的WorkerId
	_ = dao.SetResultWorker(partition, taskResult.Id, w.WorkerId)

	return
}

// TaskUniqueIdEncode 将任务结果Id和分区编码为任务唯一Id
func TaskUniqueIdEncode(partition string, taskResultId int) (taskUniqueId string) {
	return fmt.Sprintf("%s%s%d", partition, conf.TASK_UNIQUE_ID_SEPARATOR, taskResultId)
}

// TaskUniqueIdDecode 将任务唯一Id解码为任务结果Id和分区
func TaskUniqueIdDecode(taskUniqueId string) (partition string, taskResultId int, err error) {
	// 根据分割符拆分任务唯一Id
	split := strings.Split(taskUniqueId, conf.TASK_UNIQUE_ID_SEPARATOR)
	// 拆分后切片长度不是2，则说明任务唯一Id不正确
	if len(split) != 2 {
		err = fmt.Errorf("task unique id invalid")
		log.ErrorWithFields(log.Fields{
			"task_unique_id": taskUniqueId,
			"split_len":      len(split),
			"error":          err,
		}, "An error occurred while strings.Split for taskUniqueId.")
		return "", -1, err
	}

	// 将拆分后切片的第2个元素转为int类型，转换失败也说明任务唯一Id不正确
	taskResultId, err = strconv.Atoi(split[1])
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": taskUniqueId,
			"error":          err,
		}, "An error occurred while for taskResultId.")
		return "", -1, err
	}

	partition = split[0]

	return
}
