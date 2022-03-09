package builtin

import (
	"context"
	auth "eago/auth/srv/proto"
	"eago/common/log"
	"eago/common/utils"
	"eago/flow/cli"
	"eago/flow/conf"
	"eago/flow/dao"
	"eago/flow/dto"
	"eago/flow/model"
	task "eago/task/srv/proto"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// HandleInstance 处理流程实例
func HandleInstance(ins *model.Instance, hdIns *dto.HandleInstance) error {
	log.Info("builtin.HandleInstance called.")
	defer log.Info("builtin.HandleInstance end.")

	_ = dao.NewLog(ins.Id, *hdIns.Result, *hdIns.Content, hdIns.CreatedBy)

	// 审批被拒，直接结束流程
	if !*hdIns.Result {
		log.Info("The HandleInstance result is rejected.")
		dao.SetHandleInstance(
			ins.Id,
			conf.INSTANCE_STATUS_REJECTED_END,
			-1,
			0,
			*ins.FormData,
			"",
			appendPassedAssignees(ins.CurrentAssignees, ins.PassedAssignees),
			hdIns.CreatedBy,
		)
		return nil
	}

	// 反序列化Assignees
	currAss := strings.Split(ins.CurrentAssignees, conf.ASSIGNEES_SPILT_TAG)
	// 为CurrentAssignees去除已审批人
	log.Info("Exclude current user from CurrentAssignees.")
	currAss = utils.RemoveStringSliceElement(currAss, hdIns.CreatedBy)

	// 结束当前节点的审批
	if ins.AssigneesRequired <= 1 || len(currAss) < 1 {
		log.Info("The HandleInstance go next step, Final Set instance status is model.INSTANCE_STATUS_PENDING.")
		dao.SetHandleInstance(
			ins.Id,
			conf.INSTANCE_STATUS_PENDING,
			ins.CurrentStep,
			0,
			*ins.FormData,
			strings.Join(currAss, conf.ASSIGNEES_SPILT_TAG),
			appendPassedAssignees(ins.PassedAssignees, hdIns.CreatedBy),
			hdIns.CreatedBy,
		)
		_ = InstanceNextStep(ins.Id)
		return nil
	}

	log.Info("Final Set instance status is still model.INSTANCE_STATUS_RUNNING.")
	dao.SetHandleInstance(
		ins.Id,
		conf.INSTANCE_STATUS_RUNNING,
		ins.CurrentStep,
		ins.AssigneesRequired-1,
		*ins.FormData,
		strings.Join(currAss, conf.ASSIGNEES_SPILT_TAG),
		appendPassedAssignees(ins.PassedAssignees, hdIns.CreatedBy),
		hdIns.CreatedBy,
	)

	return nil
}

// InstanceNextStep 流转实例流转至下一步
func InstanceNextStep(insId int) error {
	log.Info("builtin.InstanceNextStep called.")
	defer log.Info("builtin.InstanceNextStep end.")

	// 查找流程实例
	ins, err := dao.GetInstance(dao.Query{"id=?": insId, "status=?": conf.INSTANCE_STATUS_PENDING})
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"instance_id": insId,
			"error":       err,
		}, "An error occurred while dao.GetInstance.")
		return errors.New("get instance error")
	}
	if ins == nil || ins.Id == 0 {
		log.ErrorWithFields(log.Fields{
			"instance_id": insId,
		}, "An nil object is returned after calling dao.GetInstance.")
		return errors.New("instance not found")
	}

	// 反序列化审批节点链
	headChain := &model.NodeChain{}
	if err := json.Unmarshal([]byte(*ins.FlowChain), headChain); err != nil {
		log.ErrorWithFields(log.Fields{
			"instance_id": insId,
			"error":       err,
		}, "An error occurred while json.Unmarshal for ins.FlowChain.")
		return err
	}

	// 找到当前审批节点
	currStep := ins.CurrentStep + 1
	currNode := headChain
	for i := 0; i < currStep; i++ {
		if currNode.SubNode == nil || currNode.SubNode.Id == 0 {
			dao.SetInstance(
				insId,
				conf.INSTANCE_STATUS_APPROVED_END,
				-1,
				0,
				*ins.FlowChain,
				"",
				appendPassedAssignees(ins.CurrentAssignees, ins.PassedAssignees),
				"",
			)

			return nil
		}
		currNode = currNode.SubNode
	}

	// 解析form data成为map结构
	mapData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(*ins.FormData), &mapData); err != nil {
		log.ErrorWithFields(log.Fields{
			"instance_id": insId,
			"error":       err,
		}, "An error occurred while json.Unmarshal for ins.FormData.")
		return err
	}

	// 依次获取所有流程节点链审批人
	if err := getAssignees(currNode, mapData); err != nil {
		log.ErrorWithFields(log.Fields{
			"instance_id": insId,
			"error":       err,
		}, "An error occurred while getAssignees.")
		return err
	}

	// assReq int 标记当前节点需要多少个用户审批
	var assReq int
	switch currNode.Category {
	case conf.NODE_CATEGORY_FIRST:
		// 首节点，需要0个用户审批
		assReq = 0
	case conf.NODE_CATEGORY_ANY:
		// 或签，需要1用户审批
		assReq = 1
	case conf.NODE_CATEGORY_ALL:
		// 会签，需要全部用户审批
		assReq = len(currNode.Assignees)
	case conf.NODE_CATEGORY_INFORM:
		// 知会，需要0个用户审批
		assReq = 0
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	// 结束之前等待所有并发结束
	defer wg.Wait()

	// 并发：通知审批人
	go func(wg *sync.WaitGroup) {
		notifyAssignees(currNode.Assignees, ins)
		wg.Done()
	}(wg)

	// 并发：调用触发器
	go func(wg *sync.WaitGroup) {
		callTriggers([]int{}, mapData)
		wg.Done()
	}(wg)

	// 知会节点或没有审批人的节点，保存后再次调用Next
	if assReq < 1 || len(currNode.Assignees) < 1 {
		dao.SetInstance(
			insId,
			conf.INSTANCE_STATUS_PENDING,
			currStep,
			assReq,
			*ins.FlowChain,
			"",
			appendPassedAssignees(ins.CurrentAssignees, ins.PassedAssignees),
			"",
		)
		return InstanceNextStep(ins.Id)
	}

	// 找到审批人且不是知会节点，则直接至状态为审批中
	dao.SetInstance(
		insId,
		conf.INSTANCE_STATUS_RUNNING,
		currStep,
		assReq,
		*ins.FlowChain,
		strings.Join(currNode.Assignees, conf.ASSIGNEES_SPILT_TAG),
		appendPassedAssignees(ins.CurrentAssignees, ins.PassedAssignees),
		"",
	)

	return nil
}

// getAssignees 获取指定节点实际审批人
func getAssignees(currNode *model.NodeChain, data map[string]interface{}) error {
	log.Info("builtin.getAssignees called.")
	defer log.Info("builtin.getAssignees end.")

	ac := model.AssigneeCondition{}
	// 反序列化AssigneeCondition
	if err := json.Unmarshal([]byte(currNode.AssigneeCondition), &ac); err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while json.Unmarshal for currNode.AssigneeCondition.")
		return err
	}

	ctx := context.Background()

	log.DebugWithFields(log.Fields{"assignee_condition": ac}, "Before switch ac.Condition.")
	// 处理具体的Condition
	switch ac.Condition {
	case conf.AC_INITIATOR:
		log.Info("The ac.Condition match to model.AC_INITIATOR.")
		currNode.Assignees = append(currNode.Assignees, data[conf.INITIATOR_USERNAME_KEY].(string))

	case conf.AC_INITIATORS_DEPARTMENTS_OWNER:
		log.Info("The ac.Condition match to model.AC_INITIATORS_DEPARTMENTS_OWNER.")
		req := auth.IdQuery{Id: int32(data[conf.INITIATOR_USER_ID_KEY].(int))}
		memUsers, err := cli.AuthClient.ListUserDepartmentUsers(ctx, &req)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"error": err,
			}, "An error occurred while cli.AuthClient.ListUserDepartmentUsers.")
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

	case conf.AC_INITIATORS_PARENT_DEPARTMENTS_OWNER:
		log.Info("The ac.Condition match to model.AC_INITIATORS_PARENT_DEPARTMENTS_OWNER.")
		// 获得用户所在部门
		dept, err := cli.AuthClient.GetUserDepartment(ctx, &auth.IdQuery{Id: int32(data[conf.INITIATOR_USER_ID_KEY].(int))})
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"error": err,
			}, "An error occurred while cli.AuthClient.GetUserDepartment.")
			return err
		}

		// 如果找不到用户所在部门，责审批人置空
		if dept.Id < 1 {
			return nil
		}

		// 获取用户所在部门的父部门成员
		memUsers, err := cli.AuthClient.ListParentDepartmentUsers(ctx, &auth.IdQuery{Id: dept.Id})
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"error": err,
			}, "An error occurred while cli.AuthClient.ListParentDepartmentUsers.")
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

	case conf.AC_SPECIFIED_USERS:
		log.Info("The ac.Condition match to model.AC_SPECIFIED_USERS.")
		currNode.Assignees = strings.Split(getter(&ac, data).(string), ",")

	case conf.AC_SPECIFIED_PRODUCT_OWNER:
		log.Info("The ac.Condition match to model.AC_SPECIFIED_PRODUCT_OWNER.")
		req := auth.IdQuery{Id: int32(getter(&ac, data).(float64))}
		memUsers, err := cli.AuthClient.ListProductUsers(ctx, &req)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"error": err,
			}, "An error occurred while cli.AuthClient.ListProductUsers.")
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

	case conf.AC_SPECIFIED_GROUP_OWNER:
		log.Info("The ac.Condition match to model.AC_SPECIFIED_GROUP_OWNER.")
		req := auth.IdQuery{Id: int32(getter(&ac, data).(int))}
		memUsers, err := cli.AuthClient.ListGroupUsers(ctx, &req)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"error": err,
			}, "An error occurred while cli.AuthClient.ListGroupUsers.")
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

	case conf.AC_SPECIFIED_DEPARTMENT_OWNER:
		log.Info("The ac.Condition match to model.AC_SPECIFIED_DEPARTMENT_OWNER.")
		req := auth.IdQuery{Id: int32(getter(&ac, data).(int))}
		memUsers, err := cli.AuthClient.ListDepartmentUsers(ctx, &req)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"error": err,
			}, "An error occurred while cli.AuthClient.ListDepartmentUsers.")
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

	case conf.AC_SPECIFIED_ROLE:
		log.Info("The ac.Condition match to model.AC_SPECIFIED_ROLE.")
		req := auth.NameQuery{Name: getter(&ac, data).(string)}
		memUsers, err := cli.AuthClient.ListRoleUsers(ctx, &req)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"error": err,
			}, "An error occurred while cli.AuthClient.ListRoleUsers.")
			return err
		}

		// 将成员添加到审批人列表内
		for _, u := range memUsers.Users {
			currNode.Assignees = append(currNode.Assignees, u.Username)
		}

	}
	return nil
}

// callTriggers 调用触发器
func callTriggers(tIds []int, data map[string]interface{}) bool {
	log.Info("builtin.callTriggers called.")
	defer log.Info("builtin.callTriggers end.")

	if len(tIds) < 1 {
		log.Info("The len of local.callTriggers incoming arguments is zero.")
		return true
	}

	log.Info("Finding triggers.")
	triggers, ok := dao.ListTriggers(dao.Query{"id": tIds})
	if !ok {
		log.ErrorWithFields(log.Fields{
			"trigger_ids": tIds,
		}, "An error occurred while dao.ListTriggers.")
	}

	for _, t := range triggers {
		ctx := context.Background()

		// 反序列化Trigger内的Arguments
		log.Info("Unmarshal trigger's Arguments.")
		trigArgs := make(map[string]interface{})
		err := json.Unmarshal([]byte(t.Arguments), &trigArgs)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"trigger_id":    t.Id,
				"task_codename": t.TaskCodename,
				"error":         err,
			}, "An error occurred while json.Unmarshal for t.Arguments.")
		}

		// 将FormData与Trigger内的Arguments合并
		log.Info("Merging trigger's arguments and form data.")
		args := utils.MergeMapStringInterface(trigArgs, data)

		// 序列化组合后的Arguments
		log.Info("Marshal Merged arguments.")
		arg, err := json.Marshal(args)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"trigger_id":    t.Id,
				"task_codename": t.TaskCodename,
				"error":         err,
			}, "An error occurred while json.Marshal for args.")
		}

		// 组装调用任务的请求
		log.Info("Loading task.CallTaskReq.")
		req := &task.CallTaskReq{
			TaskCodename: t.TaskCodename,
			Arguments:    arg,
			Timeout:      0,
			Caller:       conf.RPC_REGISTER_KEY + "::local.callTriggers",
		}

		// 调用任务
		log.Info("Call cli.TaskClient.CallTask.")
		rsp, err := cli.TaskClient.CallTask(ctx, req)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"trigger_id":    t.Id,
				"task_codename": t.TaskCodename,
				"arguments":     arg,
				"error":         err,
			}, "An error occurred while cli.TaskClient.CallTask.")
		}
		log.Info("Call cli.TaskClient.CallTask done.")

		if rsp != nil {
			log.InfoWithFields(log.Fields{"task_unique_id": rsp.TaskUniqueId}, "cli.TaskClient.CallTask success.")
		} else {
			log.Warn("cli.TaskClient.CallTask got an nil response.")
		}

	}

	return true
}

// appendPassedAssignees 追加已审批人
func appendPassedAssignees(currAss, passedAss string) string {
	log.Info("builtin.appendPassedAssignees called.")
	defer log.Info("builtin.appendPassedAssignees end.")

	passedAssignees := utils.MergeStringSlice(
		strings.Split(passedAss, conf.ASSIGNEES_SPILT_TAG),
		strings.Split(currAss, conf.ASSIGNEES_SPILT_TAG),
	)
	return strings.Join(passedAssignees, conf.ASSIGNEES_SPILT_TAG)
}

// 通知审批人
func notifyAssignees(assignees []string, ins *model.Instance) {
	log.Info("builtin.notifyAssignees called.")
	defer log.Info("builtin.notifyAssignees end.")

	if len(assignees) < 1 {
		log.Info("The len of local.notifyAssignees incoming arguments assignees is zero.")
		return
	}

	log.Info("Loading arguments.")
	args := make(map[string]interface{})
	args["content_type"] = "textcard"
	args["subject"] = conf.Conf.NotifyTitle
	args["to"] = assignees
	args["content"] = map[string]string{
		"description": fmt.Sprintf(
			"<div class=\"gray\">%s</div>"+
				"<div class=\"normal\">ID：%d</div>"+
				"<div class=\"normal\">流程名：%s</div>"+
				"<div class=\"normal\">发起人：%s</div>"+
				"<div class=\"highlight\">请您审批处理</div>",
			time.Now().Format(conf.TIMESTAMP_FORMAT),
			ins.Id,
			ins.Name,
			ins.CreatedBy,
		),
		"url":    conf.Conf.NotifyBaseUrl,
		"btntxt": "查看详情",
	}

	// 序列化组合后的Arguments
	log.Info("Marshal arguments.")
	arg, err := json.Marshal(args)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"arguments": args,
			"error":     err,
		}, "An error occurred while json.Marshal for args.")
	}

	// 组装调用任务的请求
	log.Info("Loading task.CallTaskReq.")
	req := &task.CallTaskReq{
		TaskCodename: "notify.send_wework",
		Arguments:    arg,
		Timeout:      0,
		Caller:       conf.RPC_REGISTER_KEY + "::local.notifyAssignees",
	}

	// 调用任务
	log.Info("Call cli.TaskClient.CallTask.")
	ctx := context.Background()
	rsp, err := cli.TaskClient.CallTask(ctx, req)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_codename": "notify.send_wework",
			"arguments":     "",
			"error":         err,
		}, "An error occurred while cli.TaskClient.CallTask.")
	}
	log.Info("Call cli.TaskClient.CallTask done.")

	if rsp != nil {
		log.InfoWithFields(log.Fields{"task_unique_id": rsp.TaskUniqueId}, "cli.TaskClient.CallTask success.")
	} else {
		log.Warn("cli.TaskClient.CallTask got an nil response.")
	}

	return
}
