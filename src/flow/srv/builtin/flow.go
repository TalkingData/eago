package builtin

import (
	"eago/common/log"
	"eago/common/utils"
	"eago/flow/conf"
	"eago/flow/dao"
	"eago/flow/model"
	"encoding/json"
	"errors"
	"sync"
)

// InstantiateFlow 实例化流程
func InstantiateFlow(flowId int, fromData, createdBy string) (int, error) {
	log.Info("builtin.InstantiateFlow called.")
	defer log.Info("builtin.InstantiateFlow end.")

	// 获取流程
	flow, ok := dao.GetFlow(dao.Query{"id=?": flowId})
	if !ok {
		log.ErrorWithFields(log.Fields{
			"flow_id": flowId,
		}, "An error occurred while dao.GetFlow.")
		return 0, errors.New("get flow error")
	}
	if flow == nil {
		log.ErrorWithFields(log.Fields{
			"flow_id": flowId,
		}, "An nil object is returned after calling dao.GetFlow.")
		return 0, errors.New("flow not found")
	}
	log.InfoWithFields(log.Fields{"flow_id": flow.Id}, "dao.GetFlow success.")

	// 获取流程首节点
	headNode, ok := dao.GetNode(dao.Query{"id=?": flow.FirstNodeId})
	if !ok {
		log.ErrorWithFields(log.Fields{
			"flow_id": flowId,
			"node_id": flow.FirstNodeId,
		}, "An error occurred while dao.GetNode.")
		return 0, errors.New("get first node error")
	}
	// 头节点为空则不无法创建流程
	if headNode == nil {
		log.ErrorWithFields(log.Fields{
			"flow_id": flowId,
			"node_id": flow.FirstNodeId,
		}, "An nil object is returned after calling dao.GetNode.")
		return 0, errors.New("first node not found")
	}
	log.InfoWithFields(log.Fields{"first_node_id": headNode.Id}, "Get first node success.")

	headChain := &model.NodeChain{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 将流程节点转换为链格式
	go func(ws *sync.WaitGroup) {
		log.Info("Start calling dao.Node2Chain.")
		headChain = dao.Node2Chain(headNode)
		dao.GetNodeChain(headChain)
		wg.Done()
		log.Info("dao.Node2Chain success.")
	}(&wg)

	// 解析form data成为map结构
	log.Info("Start calling json.Unmarshal(fromData).")
	mapData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(fromData), &mapData); err != nil {
		return 0, err
	}
	log.Info("json.Unmarshal(fromData) success.")

	wg.Wait()

	// 依次获取所有流程节点链审批人
	if err := getAssignees(headChain, mapData); err != nil {
		return 0, err
	}

	chainStr, err := chain2String(headChain)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"flow_id": flowId,
			"node_id": flow.FirstNodeId,
		}, "An error occurred while InstantiateFlow, in chain2String json.Marshal(head).")
		return 0, err
	}
	// 创建流程实例
	ins := dao.NewInstance(
		flow.FormId,
		conf.INSTANCE_PENDING_STATUS,
		renderInstanceName(flow.Name, mapData),
		fromData,
		chainStr,
		createdBy,
	)

	// 创建流程实例完毕后，流转实例流转至下一步
	_ = InstanceNextStep(ins.Id)
	return ins.Id, nil
}

// getter 取值方法
func getter(ac *model.AssigneeCondition, data map[string]interface{}) (realData interface{}) {
	if ac.Getter == "direct" {
		return ac.Data
	}

	return data[ac.Data]
}

// renderInstanceName 完成流程名称渲染
func renderInstanceName(in string, data map[string]interface{}) (out string) {
	log.Info("builtin.renderInstanceName called.")
	defer log.Info("builtin.renderInstanceName end.")

	res, err := utils.RenderString(in, data)
	if err != nil {
		log.WarnWithFields(log.Fields{
			"error": err,
		}, "Warning when renderInstanceName called, in utils.RenderString.")
		return in
	}

	return res
}

// chain2String
func chain2String(head *model.NodeChain) (string, error) {
	ch, err := json.Marshal(head)
	if err != nil {
		return "", err
	}
	return string(ch), nil
}
