package local

import (
	"bytes"
	"eago/common/log"
	"eago/task/cli"
	"eago/task/conf"
	"eago/task/model"
	"eago/task/worker"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

// CallTask 调用任务
func CallTask(taskCodename, arguments, caller string, timeout int64) (task_unique_id string, err error) {
	taskResult := model.NewResult(taskCodename, caller, arguments, timeout, worker.TASK_INITIALIZATION_STATUS)
	if taskResult == nil {
		err = fmt.Errorf("Got an empty result object.")
		log.ErrorWithFields(log.Fields{
			"error": err.Error(),
		}, "Error in model.NewResult.")
		return
	}
	partition, partitionErr := taskResult.GetPartition()
	if partitionErr != nil {
		log.ErrorWithFields(log.Fields{
			"error": partitionErr.Error(),
		}, "Error in model.Result.GetPartition.")
		return "", partitionErr
	}

	cNameSplit := strings.Split(taskCodename, ".")
	if len(cNameSplit) < 2 {
		err := fmt.Errorf("Invalid task codename.")
		log.ErrorWithFields(log.Fields{
			"error": err.Error(),
		}, "Error in taskCodename strings.Split.")
		model.SetResultStatus(partition, taskResult.Id, worker.TASK_CALL_ERROR_END_STATUS, true)
		return "", err
	}

	// 查找对应模块的worker
	modular := cNameSplit[0]
	wks := cli.WorkerClient.ListByModular(modular)
	if wks == nil || len(wks) < 1 {
		// 找不到模块所属的worker
		model.SetResultStatus(partition, taskResult.Id, worker.TASK_NO_WORKER_ERROR_END_STATUS, true)
		return "", fmt.Errorf("No worker for %s", modular)
	}

	// 生成任务实例唯一ID
	task_unique_id = TaskUniqueIdEncode(partition, taskResult.Id)

	data, _ := json.Marshal(worker.CallTaskReq{
		TaskCodename: cNameSplit[1],

		TaskUniqueId: task_unique_id,
		Arguments:    arguments,
		Timeout:      timeout,
		Caller:       caller,

		Timestamp: taskResult.StartAt.Unix(),
	})

	// 随机找一个Worker
	w := wks[rand.Intn(len(wks))]
	// 调用Worker
	uri := fmt.Sprintf("http://%s/call", w.Address)
	req, _ := http.NewRequest("POST", uri, bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err.Error(),
		}, "Error in call task client.Do.")
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := ioutil.ReadAll(resp.Body)

	wkResp := &worker.WorkerResponse{}
	if err := json.Unmarshal(body, wkResp); err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err.Error(),
		}, "Error in call task json.Unmarshal.")
		return "", err
	}

	if wkResp.Code != 0 {
		log.ErrorWithFields(log.Fields{
			"status_code": wkResp.Code,
		}, "Error in call task, worker status is not 0.")
		return "", err
	}

	// 填充执行任务的WorkerId
	_ = model.SetResultWorker(partition, taskResult.Id, w.WorkerId)

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
		err = fmt.Errorf("Task unique id invalid.")
		log.ErrorWithFields(log.Fields{
			"task_unique_id": taskUniqueId,
			"split_len":      len(split),
			"error":          err.Error(),
		}, "Error in TaskUniqueIdDecode.")
		return "", -1, err
	}

	// 将拆分后切片的第2个元素转为int类型，转换失败也说明任务唯一Id不正确
	taskResultId, err = strconv.Atoi(split[1])
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": taskUniqueId,
			"error":          err.Error(),
		}, "Error in TaskUniqueIdDecode.")
		return "", -1, err
	}
	partition = split[0]
	return
}
