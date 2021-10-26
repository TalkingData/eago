package handler

import (
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"eago/flow/dto"
	"github.com/gin-gonic/gin"
	"strconv"
)

// NewNode 新建节点
func NewNode(c *gin.Context) {
	var nFrm dto.NewNode

	// 序列化request body
	if err := c.ShouldBindJSON(&nFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := nFrm.Validate(); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 新建
	tc := c.GetStringMap("TokenContent")
	n := dao.NewNode(
		nFrm.Name,
		nFrm.ParentId,
		nFrm.Category,
		nFrm.EntryCondition,
		nFrm.AssigneeCondition,
		nFrm.VisibleFields,
		nFrm.EditableFields,
		tc["Username"].(string),
	)
	if n == nil {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "node", n)
}

// RemoveNode 删除节点
func RemoveNode(c *gin.Context) {
	nId, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "node_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var rnFrm dto.RemoveNode
	// 验证数据
	if m := rnFrm.Validate(nId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if ok := dao.RemoveNode(nId); !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetNode 更新节点
func SetNode(c *gin.Context) {
	nId, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "node_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var nFrm dto.SetNode
	// 序列化request body
	if err = c.ShouldBindJSON(&nFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := nFrm.Validate(nId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	n, ok := dao.SetNode(
		nId,
		nFrm.Name,
		nFrm.ParentId,
		nFrm.Category,
		nFrm.EntryCondition,
		nFrm.AssigneeCondition,
		nFrm.VisibleFields,
		nFrm.EditableFields,
		tc["Username"].(string),
	)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "node", n)
}

// GetNodeChain 列出指定节点链
func GetNodeChain(c *gin.Context) {
	nID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "node_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 查找根节点
	node, ok := dao.GetNode(dao.Query{"id=?": nID})
	if !ok {
		m := msg.UnknownError
		log.ErrorWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 找不到根节点则直接返回空
	if node == nil {
		w.WriteSuccessPayload(c, "chain", make(map[string]interface{}))
		return
	}

	// 将根节点转化为链结构
	root := dao.Node2Chain(node)
	if ok = dao.GetNodeChain(root); !ok {
		m := msg.UnknownError
		log.WarnWithFields(log.Fields{"node_id": nID}, m.String())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "chain", root)
}

// ListNodes 列出所有节点
func ListNodes(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	ltq := dto.ListNodesQuery{}
	if c.ShouldBindQuery(&ltq) == nil {
		_ = ltq.UpdateQuery(query)
	}

	paged, ok := dao.PagedListNodes(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while model.PagedListNodes.")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "nodes", paged)
}

// AddTrigger2Node 添加触发器至节点
func AddTrigger2Node(c *gin.Context) {
	nID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "node_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var atnFrm dto.AddTrigger2Node
	// 序列化request body
	if err = c.ShouldBindJSON(&atnFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := atnFrm.Validate(nID); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	if !dao.AddNodeTrigger(nID, atnFrm.TriggerId, tc["Username"].(string)) {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// RemoveNodeTrigger 移除节点中触发器
func RemoveNodeTrigger(c *gin.Context) {
	nID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "node_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	tId, err := strconv.Atoi(c.Param("trigger_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "trigger_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var rntFrm dto.RemoveNodeTrigger
	// 验证数据
	if m := rntFrm.Validate(nID, tId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.RemoveNodeTrigger(nID, tId) {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// ListNodeTriggers 列出节点中所有触发器
func ListNodeTriggers(c *gin.Context) {
	nID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "node_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	u, ok := dao.ListNodeTriggers(nID)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "triggers", u)
}
