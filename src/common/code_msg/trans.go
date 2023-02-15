package code_msg

import (
	"github.com/micro/go-micro/v2/errors"
	"google.golang.org/grpc/status"
)

// TransMicroErr2CodeMsg 尝试将MicroErr转为CodeMsg类型
func TransMicroErr2CodeMsg(microErr error) (*CodeMsg, bool) {
	rErr := errors.FromError(microErr)
	if rErr.Id != defaultMicroErrId {
		return nil, false
	}

	return NewCodeMsg(int(rErr.Code), rErr.Detail), true
}

// TransRpcErr2CodeMsg 尝试将RpcError转为CodeMsg类型
func TransRpcErr2CodeMsg(rpcErr error) (*CodeMsg, bool) {
	sts := status.Convert(rpcErr)
	if sts.Code() < 400 {
		return nil, false
	}
	return NewCodeMsg(int(sts.Code()), sts.Message()), true
}
