package dto

import (
	"eago/common/message"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"github.com/beego/beego/v2/core/validation"
)

// NewLog struct
type NewLog struct {
	Result  bool    `json:"result" valid:"Required"`
	Content *string `json:"content" valid:"MinSize(0);MaxSize(500)"`
}

func (n *NewLog) Validate(insId int) *message.Message {
	// 验证流程实例是否存在
	if ct, _ := dao.GetInstancesCount(dao.Query{"id=?": insId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("流程实例不存在")
	}

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(n)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// ListLog struct
type ListLog struct{}

func (*ListLog) Validate(insId int) *message.Message {
	// 验证流程实例是否存在
	if ct, _ := dao.GetInstancesCount(dao.Query{"id=?": insId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("流程实例不存在")
	}

	return nil
}
