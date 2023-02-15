package form

import (
	"context"
	cMsg "eago/common/code_msg"
	"eago/common/orm"
	"eago/flow/dao"
	"github.com/beego/beego/v2/core/validation"
)

type NewLogForm struct {
	Result  bool    `json:"result" valid:"Required"`
	Content *string `json:"content" valid:"MinSize(0);MaxSize(500)"`
}

func (f *NewLogForm) Validate(ctx context.Context, dao *dao.Dao, instId uint32) *cMsg.CodeMsg {
	// 验证流程实例是否存在
	if exist, _ := dao.IsInstanceExist(ctx, orm.Query{"id=?": instId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("流程实例不存在")
	}

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(f)
	if err != nil {
		return cMsg.MsgValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return cMsg.MsgValidateFailed.SetError(valid.Errors)
	}

	return nil
}

type ListLogForm struct{}

func (*ListLogForm) Validate(ctx context.Context, dao *dao.Dao, instId uint32) *cMsg.CodeMsg {
	// 验证流程实例是否存在
	if exist, _ := dao.IsInstanceExist(ctx, orm.Query{"id=?": instId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("流程实例不存在")
	}

	return nil
}
