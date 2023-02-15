package biz

import (
	"context"
	"eago/common/logger"
	"eago/common/utils"
	"eago/task/dto"
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

// CallTask 调用任务
func (b *Biz) CallTask(
	ctx context.Context,
	taskCodename, arguments, caller string,
	timeout int64,
) (taskUniqueId string, err error) {
	resObj, err := b.dao.NewResult(ctx, taskCodename, caller, arguments, timeout, dto.TaskResultStatusInitialization)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"task_codename": taskCodename,
			"caller":        caller,
			"arguments":     arguments,
			"error":         err,
		}, "An error occurred while dao.NewResult in biz.CallTask.")
		return "", err
	}
	if resObj == nil || resObj.Id < 1 {
		err = fmt.Errorf("result object not found")
		b.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An nil object is returned after calling dao.NewResult in biz.CallTask.")
		return
	}

	part, err := resObj.GetPartition(b.conf.Const.TaskResultPartitionTsFormat)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"error": err.Error(),
		}, "An error occurred while Result.GetPartition in biz.CallTask.")
		return "", err
	}

	cNameSplit := strings.Split(taskCodename, ".")
	if len(cNameSplit) < 2 {
		err = fmt.Errorf("invalid task codename")
		b.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while strings.Split for taskCodename.")
		_ = b.dao.SetResultStatus(ctx, part, resObj.Id, dto.TaskResultStatusCallErrEnd, true)
		return "", err
	}

	// 查找对应模块的worker
	modular := cNameSplit[0]
	wks := b.workerCli.ListByModular(ctx, modular)
	if wks == nil || len(wks) < 1 {
		// 找不到模块所属的worker
		_ = b.dao.SetResultStatus(ctx, part, resObj.Id, dto.TaskResultStatusNoWorkerErrEnd, true)
		b.logger.ErrorWithFields(logger.Fields{
			"worker": modular,
		}, "Can not kill task, no worker found.")
		return "", fmt.Errorf("no worker found for %s", modular)
	}

	// 生成任务实例唯一ID
	taskUniqueId = b.TaskUniqueIdEncode(part, resObj.Id)

	// 随机找一个Worker
	wk := wks[rand.Intn(len(wks))]

	// 调用Worker
	if err = b.workerCli.CallTask(
		b.NewSrvTokenWithCtx(ctx, wk.Address),
		wk,
		cNameSplit[1],
		taskUniqueId,
		arguments,
		caller,
		timeout,
		resObj.StartAt.Unix(),
	); err != nil {
		// 任务调用错误
		_ = b.dao.SetResultStatus(ctx, part, resObj.Id, dto.TaskResultStatusCallErrEnd, true)
		b.logger.ErrorWithFields(logger.Fields{
			"partition": part,
			"result_id": resObj.Id,
			"error":     err,
		}, "An error occurred while workerCli.CallTask in biz.CallTask.")
		return "", err
	}
	// 填充执行任务的WorkerId
	_ = b.dao.SetResultWorker(ctx, part, resObj.Id, wk.WorkerId)

	return
}

// KillTask 结束任务
func (b *Biz) KillTask(ctx context.Context, taskUniqueId string) error {
	// 将任务唯一Id解码为任务结果Id和分区
	part, resId, err := b.TaskUniqueIdDecode(taskUniqueId)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"task_unique_id": taskUniqueId,
			"error":          err,
		}, "An error occurred while biz.TaskUniqueIdDecode in biz.KillTask.")
		return err
	}

	// 获得任务结果对象
	resObj, err := b.dao.GetResult(ctx, part, resId)
	if err != nil {
		// 获得任务时返回了空对象
		b.logger.ErrorWithFields(logger.Fields{
			"partition": part,
			"id":        resId,
			"error":     err,
		}, "An error occurred while dao.GetResult in biz.KillTask.")
		return err
	}
	// 找不到数据的处理
	if resObj == nil || resObj.Id < 1 {
		b.logger.ErrorWithFields(logger.Fields{
			"partition": part,
			"result_id": resId,
		}, "Result object not found.")
		return errors.New("result object not found")
	}

	// 任务不是运行状态的，无法手动结束任务
	if resObj.Status != dto.TaskResultStatusRunning {
		b.logger.ErrorWithFields(logger.Fields{
			"partition":   part,
			"result_id":   resId,
			"task_status": resObj.Status,
		}, "Can not kill task, it is not in running status.")
		return errors.New("task is not in running status")
	}

	// 找不到Worker直接结束
	wk := b.workerCli.GetWorkerById(ctx, resObj.Worker)
	if wk == nil {
		// 找不到任务所属的worker
		_ = b.dao.SetResultStatus(ctx, part, resId, dto.TaskResultStatusNoWorkerErrEnd, true)
		b.logger.ErrorWithFields(logger.Fields{
			"partition": part,
			"result_id": resId,
			"worker":    resObj.Worker,
		}, "Can not kill task, no worker found.")
		return fmt.Errorf("worker not found for %s", resObj.Worker)
	}

	// 调用Worker
	err = b.workerCli.KillTask(b.NewSrvTokenWithCtx(ctx, wk.Address), wk, taskUniqueId)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"partition": part,
			"result_id": resId,
			"error":     err,
		}, "An error occurred while workerCli.KillTask in biz.KillTask.")
		return err
	}

	return nil
}

// TaskUniqueIdEncode 将任务结果Id和分区编码为任务唯一Id
func (b *Biz) TaskUniqueIdEncode(partition string, taskResultId uint32) (taskUniqueId string) {
	return fmt.Sprintf("%s%s%d", partition, b.conf.Const.TaskUniqueIdSeparator, taskResultId)
}

// TaskUniqueIdDecode 将任务唯一Id解码为任务结果Id和分区
func (b *Biz) TaskUniqueIdDecode(taskUniqueId string) (partition string, taskResultId uint32, err error) {
	// 根据分割符拆分任务唯一Id
	split := strings.Split(taskUniqueId, b.conf.Const.TaskUniqueIdSeparator)
	// 拆分后切片长度不是2，则说明任务唯一Id不正确
	if len(split) != 2 {
		err = fmt.Errorf("task unique id invalid")
		b.logger.ErrorWithFields(logger.Fields{
			"task_unique_id": taskUniqueId,
			"split_len":      len(split),
			"error":          err,
		}, "An error occurred while strings.Split for taskUniqueId in biz.TaskUniqueIdDecode.")
		return "", 0, err
	}

	// 将拆分后切片的第2个元素转为int类型，转换失败也说明任务唯一Id不正确
	taskResultId, err = utils.Str2Uint32(split[1])
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"task_unique_id": taskUniqueId,
			"error":          err,
		}, "An error occurred while utils.Str2Uint32 for taskResultId in biz.TaskUniqueIdDecode.")
		return "", 0, err
	}

	partition = split[0]
	return
}
