package biz

import (
	"context"
	"eago/common/logger"
	"eago/common/orm"
	"eago/common/utils"
	"eago/flow/dto"
	"eago/flow/model"
	"encoding/json"
	"errors"
	"sync"
)

// InstantiateFlow 实例化流程
func (b *Biz) InstantiateFlow(ctx context.Context, flowId uint32, fromData, createdBy string) (uint32, error) {
	b.logger.Debug("biz.InstantiateFlow called.")
	defer b.logger.Debug("biz.InstantiateFlow end.")

	// 获取流程
	flow, err := b.dao.GetFlow(ctx, orm.Query{"id=?": flowId})
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"flow_id": flowId,
			"error":   err,
		}, "An error occurred while dao.GetFlow in biz.InstantiateFlow.")
		return 0, errors.New("get flow error")
	}

	if flow == nil || flow.Id < 1 {
		b.logger.ErrorWithFields(logger.Fields{
			"flow_id": flowId,
		}, "An nil object is returned after calling dao.GetFlow in biz.InstantiateFlow.")
		return 0, errors.New("flow not found")
	}
	b.logger.DebugWithFields(logger.Fields{"flow_id": flow.Id}, "dao.GetFlow success.")

	// 获取流程首节点
	headNode, err := b.dao.GetNode(ctx, orm.Query{"id=?": flow.FirstNodeId})
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"flow_id": flowId,
			"node_id": flow.FirstNodeId,
			"error":   err,
		}, "An error occurred while dao.GetNode in biz.InstantiateFlow.")
		return 0, errors.New("get first node error")
	}

	// 头节点为空则不无法创建流程
	if headNode == nil || headNode.Id < 1 {
		b.logger.ErrorWithFields(logger.Fields{
			"flow_id": flowId,
			"node_id": flow.FirstNodeId,
		}, "An nil object is returned after calling dao.GetNode in biz.InstantiateFlow.")
		return 0, errors.New("first node not found")
	}
	b.logger.DebugWithFields(logger.Fields{"first_node_id": headNode.Id}, "Get first node success.")

	headChain := &model.NodeChain{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 将流程节点转换为链格式
	go func(ws *sync.WaitGroup) {
		b.logger.Debug("Start calling dao.Node2Chain.")
		headChain = b.dao.Node2Chain(ctx, headNode)
		if err = b.dao.GetNodeChain(ctx, headChain); err != nil {
			b.logger.WarnWithFields(logger.Fields{
				"flow_id": flowId,
				"node_id": flow.FirstNodeId,
				"error":   err,
			}, "An error occurred while dao.GetNodeChain in biz.InstantiateFlow, skipped.")
		}

		wg.Done()
		b.logger.Debug("dao.Node2Chain success.")
	}(&wg)

	// 解析form data成为map结构
	b.logger.Debug("Start calling json.Unmarshal(fromData).")
	mapData := make(map[string]interface{})
	if err = json.Unmarshal([]byte(fromData), &mapData); err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"flow_id": flowId,
			"error":   err,
		}, "An error occurred while json.Unmarshal(fromData) in biz.InstantiateFlow.")
		return 0, err
	}
	b.logger.Debug("json.Unmarshal(fromData) success.")

	wg.Wait()

	// 依次获取所有流程节点链审批人
	if err = b.getAssignees(ctx, headChain, mapData); err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"flow_id": flowId,
			"error":   err,
		}, "An error occurred while biz.getAssignees in biz.InstantiateFlow.")
		return 0, err
	}

	chainStr, err := chain2String(headChain)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"flow_id": flowId,
			"node_id": flow.FirstNodeId,
		}, "An error occurred while chain2String in biz.InstantiateFlow.")
		return 0, err
	}

	// 创建流程实例
	inst, err := b.dao.NewInstance(
		ctx,
		flow.FormId,
		dto.InstanceStatusPending,
		b.renderInstanceName(flow.InstanceTitle, mapData),
		fromData,
		chainStr,
		createdBy,
	)

	// 创建流程实例完毕后，流转实例流转至下一步
	_ = b.InstanceNextStep(ctx, inst.Id)
	return inst.Id, nil
}

// renderInstanceName 完成流程名称渲染
func (b *Biz) renderInstanceName(in string, data map[string]interface{}) (out string) {
	b.logger.Info("biz.renderInstanceName called.")
	defer b.logger.Info("biz.renderInstanceName end.")

	res, err := utils.RenderString(in, data)
	if err != nil {
		b.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while chain2String in biz.renderInstanceName.")
		return in
	}

	return res
}

// getter 取值方法
func getter(ac *dto.AssigneeCondition, data map[string]interface{}) (realData interface{}) {
	if ac.Getter == dto.GetterDirect {
		return ac.Data
	}

	return data[ac.Data]
}

func chain2String(head *model.NodeChain) (string, error) {
	ch, err := json.Marshal(head)
	if err != nil {
		return "", err
	}

	return string(ch), nil
}
