package biz

import (
	"context"
	"eago/common/global"
	"eago/common/logger"
	"eago/common/orm"
	commonpb "eago/common/proto"
	"eago/common/utils"
	"eago/flow/dto"
	"eago/flow/model"
	taskpb "eago/task/proto"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// HandleInstance 处理流程实例
func (b *Biz) HandleInstance(
	ctx context.Context,
	inst *model.Instance, createdBy string, result bool, content *string,
) error {
	b.logger.Info("biz.HandleInstance called.")
	defer b.logger.Info("biz.HandleInstance end.")

	_, _ = b.dao.NewLog(ctx, inst.Id, result, *content, createdBy)

	// 审批被拒，直接结束流程
	if !result {
		b.logger.Info("The HandleInstance result is rejected.")
		_ = b.dao.SetHandleInstance(
			ctx,
			inst.Id,
			dto.InstanceStatusRejectedEnd,
			-1,
			0,
			*inst.FormData,
			"",
			appendPassedAssignees(inst.CurrentAssignees, inst.PassedAssignees),
			createdBy,
		)
		return nil
	}

	// 反序列化Assignees
	currAss := strings.Split(inst.CurrentAssignees, dto.AssigneesSpiltTag)
	// 为CurrentAssignees去除已审批人
	b.logger.Info("Exclude current user from CurrentAssignees.")
	currAss = utils.RemoveStringSliceElement(currAss, createdBy)

	// 结束当前节点的审批
	if inst.AssigneesRequired <= 1 || len(currAss) < 1 {
		b.logger.Info("The HandleInstance go next step, Final Set instance status is dto.InstanceStatusPending.")
		_ = b.dao.SetHandleInstance(
			ctx,
			inst.Id,
			dto.InstanceStatusPending,
			inst.CurrentStep,
			0,
			*inst.FormData,
			strings.Join(currAss, dto.AssigneesSpiltTag),
			appendPassedAssignees(inst.PassedAssignees, createdBy),
			createdBy,
		)
		_ = b.InstanceNextStep(ctx, inst.Id)

		return nil
	}

	b.logger.Info("Final Set instance status is still dto.InstanceStatusRunnin.")
	_ = b.dao.SetHandleInstance(
		ctx,
		inst.Id,
		dto.InstanceStatusRunning,
		inst.CurrentStep,
		inst.AssigneesRequired-1,
		*inst.FormData,
		strings.Join(currAss, dto.AssigneesSpiltTag),
		appendPassedAssignees(inst.PassedAssignees, createdBy),
		createdBy,
	)

	return nil
}

// InstanceNextStep 流转实例流转至下一步
func (b *Biz) InstanceNextStep(ctx context.Context, insId uint32) error {
	b.logger.Info("biz.InstanceNextStep called.")
	defer b.logger.Info("biz.InstanceNextStep end.")

	// 查找流程实例
	inst, err := b.dao.GetInstance(ctx, orm.Query{"id=?": insId, "status=?": dto.InstanceStatusPending})
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"instance_id": insId,
			"error":       err,
		}, "An error occurred while dao.GetInstance in biz.InstanceNextStep.")
		return errors.New("get instance error")
	}

	if inst == nil || inst.Id == 0 {
		b.logger.ErrorWithFields(logger.Fields{
			"instance_id": insId,
		}, "An nil object is returned after calling dao.GetInstance in biz.InstanceNextStep.")
		return errors.New("instance not found")
	}

	// 反序列化审批节点链
	headChain := &model.NodeChain{}
	if err = json.Unmarshal([]byte(*inst.FlowChain), headChain); err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"instance_id": insId,
			"error":       err,
		}, "An error occurred while json.Unmarshal for inst.FlowChain in biz.InstanceNextStep.")
		return err
	}

	// 找到当前审批节点
	currStep := inst.CurrentStep + 1
	currNode := headChain
	for i := int32(0); i < currStep; i++ {
		if currNode.SubNode == nil || currNode.SubNode.Id == 0 {
			_, _ = b.dao.SetInstance(
				ctx,
				insId,
				dto.InstanceStatusApprovedEnd,
				-1,
				0,
				*inst.FlowChain,
				"",
				appendPassedAssignees(inst.CurrentAssignees, inst.PassedAssignees),
				"",
			)

			return nil
		}
		currNode = currNode.SubNode
	}

	// 解析form data成为map结构
	mapData := make(map[string]interface{})
	if err = json.Unmarshal([]byte(*inst.FormData), &mapData); err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"instance_id": insId,
			"error":       err,
		}, "An error occurred while json.Unmarshal for ins.FormData in biz.InstanceNextStep.")
		return err
	}

	// 依次获取所有流程节点链审批人
	if err = b.getAssignees(ctx, currNode, mapData); err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"instance_id": insId,
			"error":       err,
		}, "An error occurred while b.getAssignees in biz.InstanceNextStep.")
		return err
	}

	// assReq int32 标记当前节点需要多少个用户审批
	var assReq int32
	switch currNode.Category {
	case dto.NodeCategoryFirst:
		// 首节点，需要0个用户审批
		assReq = 0
	case dto.NodeCategoryAny:
		// 或签，需要1用户审批
		assReq = 1
	case dto.NodeCategoryAll:
		// 会签，需要全部用户审批
		assReq = int32(len(currNode.Assignees))
	case dto.NodeCategoryInform:
		// 知会，需要0个用户审批
		assReq = 0
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	// 结束之前等待所有并发结束
	defer wg.Wait()

	// 并发：通知审批人
	go func(wg *sync.WaitGroup) {
		b.notifyAssignees(ctx, currNode.Assignees, inst)
		wg.Done()
	}(wg)

	// 并发：调用触发器
	go func(wg *sync.WaitGroup) {
		_ = b.callTriggers(ctx, []int{}, mapData)
		wg.Done()
	}(wg)

	// 知会节点或没有审批人的节点，保存后再次调用Next
	if assReq < 1 || len(currNode.Assignees) < 1 {
		_, _ = b.dao.SetInstance(
			ctx,
			insId,
			dto.InstanceStatusPending,
			currStep,
			assReq,
			*inst.FlowChain,
			"",
			appendPassedAssignees(inst.CurrentAssignees, inst.PassedAssignees),
			"",
		)
		return b.InstanceNextStep(ctx, inst.Id)
	}

	// 找到审批人且不是知会节点，则直接至状态为审批中
	_, _ = b.dao.SetInstance(
		ctx,
		insId,
		dto.InstanceStatusRunning,
		currStep,
		assReq,
		*inst.FlowChain,
		strings.Join(currNode.Assignees, dto.AssigneesSpiltTag),
		appendPassedAssignees(inst.CurrentAssignees, inst.PassedAssignees),
		"",
	)

	return nil
}

// getAssignees 获取指定节点实际审批人
func (b *Biz) getAssignees(ctx context.Context, currNode *model.NodeChain, data map[string]interface{}) error {
	b.logger.Info("biz.getAssignees called.")
	defer b.logger.Info("biz.getAssignees end.")

	ac := dto.AssigneeCondition{}
	// 反序列化AssigneeCondition
	if err := json.Unmarshal([]byte(currNode.AssigneeCondition), &ac); err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while json.Unmarshal for currNode.AssigneeCondition in biz.getAssignees.")
		return err
	}

	b.logger.DebugWithFields(logger.Fields{"assignee_condition": ac}, "Before switch ac.Condition.")
	// 处理具体的Condition
	switch ac.Condition {
	case dto.AssigneeConditionInitiator:
		b.logger.Info("The ac.Condition match to dto.AssigneeConditionInitiator.")
		currNode.Assignees = append(currNode.Assignees, data[dto.InitiatorKeyUsernameKey].(string))

	case dto.AssigneeConditionInitiatorsDepartmentsOwner:
		b.logger.Info("The ac.Condition match to dto.AssigneeConditionInitiatorsDepartmentsOwner.")
		req := commonpb.IdQuery{Value: uint32(data[dto.InitiatorKeyUserId].(int))}
		memUsers, err := b.authCli.ListUsersSameDepartmentUsers(ctx, &req)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while authClient.ListUserDepartmentUsers in biz.getAssignees.")
			return err
		}

		// 将成员添加到审批人列表内
		for _, u := range memUsers.Users {
			// 不是Owner的用户直接跳过
			if !u.IsOwner {
				continue
			}
			currNode.Assignees = append(currNode.Assignees, u.Username)
		}

	case dto.AssigneeConditionInitiatorsParentDepartmentsOwner:
		b.logger.Info("The ac.Condition match to dto.AssigneeConditionInitiatorsParentDepartmentsOwner.")
		// 获得用户所在部门
		dept, err := b.authCli.GetUsersDepartment(
			ctx, &commonpb.IdQuery{Value: uint32(data[dto.InitiatorKeyUserId].(int))},
		)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while authClient.GetUsersDepartment in biz.getAssignees.")
			return err
		}

		// 如果找不到用户所在部门，责审批人置空
		if dept.Id < 1 {
			return nil
		}

		// 获取用户所在部门的父部门成员
		memUsers, err := b.authCli.ListParentDepartmentUsers(ctx, &commonpb.IdQuery{Value: dept.Id})
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while authClient.ListParentDepartmentUsers in biz.getAssignees.")
			return err
		}

		// 将成员添加到审批人列表内
		for _, u := range memUsers.Users {
			// 不是Owner的用户直接跳过
			if !u.IsOwner {
				continue
			}
			currNode.Assignees = append(currNode.Assignees, u.Username)
		}

	case dto.AssigneeConditionSpecifiedUsers:
		b.logger.Info("The ac.Condition match to dto.AssigneeConditionSpecifiedUsers.")
		currNode.Assignees = strings.Split(getter(&ac, data).(string), ",")

	case dto.AssigneeConditionSpecifiedProductOwner:
		b.logger.Info("The ac.Condition match to dto.AssigneeConditionSpecifiedProductOwner.")
		req := commonpb.IdQuery{Value: uint32(getter(&ac, data).(float64))}
		memUsers, err := b.authCli.ListProductsUsers(ctx, &req)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while authClient.ListProductsUsers in biz.getAssignees.")
			return err
		}

		// 将成员添加到审批人列表内
		for _, u := range memUsers.Users {
			// 不是Owner的用户直接跳过
			if !u.IsOwner {
				continue
			}
			currNode.Assignees = append(currNode.Assignees, u.Username)
		}

	case dto.AssigneeConditionSpecifiedGroupOwner:
		b.logger.Info("The ac.Condition match to dto.AssigneeConditionSpecifiedGroupOwner.")
		req := commonpb.IdQuery{Value: uint32(getter(&ac, data).(int))}
		memUsers, err := b.authCli.ListGroupsUsers(ctx, &req)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while authClient.ListGroupsUsers in biz.getAssignees.")
			return err
		}

		// 将成员添加到审批人列表内
		for _, u := range memUsers.Users {
			// 不是Owner的用户直接跳过
			if !u.IsOwner {
				continue
			}
			currNode.Assignees = append(currNode.Assignees, u.Username)
		}

	case dto.AssigneeConditionSpecifiedDepartmentOwner:
		b.logger.Info("The ac.Condition match to dto.AssigneeConditionSpecifiedDepartmentOwner.")
		req := commonpb.IdQuery{Value: uint32(getter(&ac, data).(int))}
		memUsers, err := b.authCli.ListDepartmentsUsers(ctx, &req)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while authClient.ListDepartmentsUsers in biz.getAssignees.")
			return err
		}

		// 将成员添加到审批人列表内
		for _, u := range memUsers.Users {
			// 不是Owner的用户直接跳过
			if !u.IsOwner {
				continue
			}
			currNode.Assignees = append(currNode.Assignees, u.Username)
		}

	case dto.AssigneeConditionSpecifiedRole:
		b.logger.Info("The ac.Condition match to dto.AssigneeConditionSpecifiedRole.")
		req := commonpb.NameQuery{Value: getter(&ac, data).(string)}
		memUsers, err := b.authCli.ListRolesUsers(ctx, &req)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while authClient.ListRolesUsers in biz.getAssignees.")
			return err
		}

		// 将成员添加到审批人列表内
		for _, u := range memUsers.Users {
			currNode.Assignees = append(currNode.Assignees, u.Username)
		}

	}
	return nil
}

// notifyAssignees 通知审批人
func (b *Biz) notifyAssignees(ctx context.Context, assignees []string, ins *model.Instance) {
	b.logger.Info("biz.notifyAssignees called.")
	defer b.logger.Info("biz.notifyAssignees end.")

	if len(assignees) < 1 {
		b.logger.Warn("The len of local.notifyAssignees incoming arguments assignees is zero.")
		return
	}

	b.logger.Info("pub.Publish called in biz.notifyAssignees.")
	err := b.pub.Publish(
		ctx,
		"instance",
		"Notify",
		"message",
		"NewMessage",
		map[string]interface{}{
			"content_type": "textcard",
			"subject":      b.conf.NotifyTitle,
			"to":           assignees,
			"content": map[string]string{
				"description": fmt.Sprintf(
					"<div class=\"gray\">%s</div>"+
						"<div class=\"normal\">ID：%d</div>"+
						"<div class=\"normal\">流程名：%s</div>"+
						"<div class=\"normal\">发起人：%s</div>"+
						"<div class=\"highlight\">请您审批处理</div>",
					time.Now().Format(global.TimestampFormat),
					ins.Id,
					ins.Name,
					ins.CreatedBy,
				),
				"url":    b.conf.NotifyBaseUrl,
				"btntxt": "查看详情",
			},
		},
	)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while pub.Publish in biz.notifyAssignees.")
	}
}

// callTriggers 调用触发器
func (b *Biz) callTriggers(ctx context.Context, tIds []int, data map[string]interface{}) error {
	b.logger.Info("biz.callTriggers called.")
	defer b.logger.Info("biz.callTriggers end.")

	if len(tIds) < 1 {
		b.logger.Info("The len of local.callTriggers incoming arguments is zero.")
		return nil
	}

	b.logger.Info("Finding triggers in biz.notifyAssignees.")
	triggers, err := b.dao.ListTriggers(ctx, orm.Query{"id": tIds})
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"trigger_ids": tIds,
			"error":       err,
		}, "An error occurred while dao.ListTriggers in biz.notifyAssignees.")
	}

	for _, tri := range triggers {
		// 反序列化Trigger内的Arguments
		b.logger.Info("Unmarshal trigger's Arguments in biz.notifyAssignees.")
		trigArgs := make(map[string]interface{})
		err := json.Unmarshal([]byte(tri.Arguments), &trigArgs)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"trigger_id":    tri.Id,
				"task_codename": tri.TaskCodename,
				"error":         err,
			}, "An error occurred while json.Unmarshal for t.Arguments in biz.notifyAssignees.")
		}

		// 将FormData与Trigger内的Arguments合并
		b.logger.Info("Merging trigger's arguments and form data in biz.notifyAssignees.")
		args := utils.MergeMapStringInterface(trigArgs, data)

		// 序列化组合后的Arguments
		b.logger.Info("Marshal Merged arguments in biz.notifyAssignees.")
		arg, err := json.Marshal(args)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"trigger_id":    tri.Id,
				"task_codename": tri.TaskCodename,
				"error":         err,
			}, "An error occurred while json.Marshal for args in biz.notifyAssignees.")
		}

		// 组装调用任务的请求
		b.logger.Info("Loading task.CallTaskReq in biz.notifyAssignees.")
		req := &taskpb.CallTaskReq{
			TaskCodename: tri.TaskCodename,
			Arguments:    arg,
			Timeout:      0,
			Caller:       b.conf.Const.ServiceName + "::local.callTriggers",
		}

		// 调用任务
		b.logger.Info("taskClient.CallTask called in biz.notifyAssignees.")
		rsp, err := b.taskCli.CallTask(ctx, req)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"trigger_id":    tri.Id,
				"task_codename": tri.TaskCodename,
				"arguments":     arg,
				"error":         err,
			}, "An error occurred while taskClient.CallTask in biz.notifyAssignees.")
		}
		b.logger.Info("taskClient.CallTask done.")

		if rsp != nil {
			b.logger.InfoWithFields(logger.Fields{
				"task_unique_id": rsp.TaskUniqueId,
			}, "taskClient.CallTask success.")
		} else {
			b.logger.Warn("cli.TaskClient.CallTask got an nil response.")
		}
	}

	return nil
}

// appendPassedAssignees 追加已审批人
func appendPassedAssignees(currAss, passedAss string) string {
	passedAssignees := utils.MergeStringSlice(
		strings.Split(passedAss, dto.AssigneesSpiltTag),
		strings.Split(currAss, dto.AssigneesSpiltTag),
	)
	return strings.Join(passedAssignees, dto.AssigneesSpiltTag)
}
