// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/task.proto

package task

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for TaskService service

func NewTaskServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for TaskService service

type TaskService interface {
	// ListTasks 列出所有任务
	ListTasks(ctx context.Context, in *Empty, opts ...client.CallOption) (*Tasks, error)
	// CallTask 调用任务
	CallTask(ctx context.Context, in *CallTaskReq, opts ...client.CallOption) (*CallTaskRsp, error)
	// SetTaskStatus 设置任务状态
	SetTaskStatus(ctx context.Context, in *SetTaskStatusReq, opts ...client.CallOption) (*BoolMsg, error)
	// AppendTaskLog 追加任务日志
	AppendTaskLog(ctx context.Context, opts ...client.CallOption) (TaskService_AppendTaskLogService, error)
}

type taskService struct {
	c    client.Client
	name string
}

func NewTaskService(name string, c client.Client) TaskService {
	return &taskService{
		c:    c,
		name: name,
	}
}

func (c *taskService) ListTasks(ctx context.Context, in *Empty, opts ...client.CallOption) (*Tasks, error) {
	req := c.c.NewRequest(c.name, "TaskService.ListTasks", in)
	out := new(Tasks)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskService) CallTask(ctx context.Context, in *CallTaskReq, opts ...client.CallOption) (*CallTaskRsp, error) {
	req := c.c.NewRequest(c.name, "TaskService.CallTask", in)
	out := new(CallTaskRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskService) SetTaskStatus(ctx context.Context, in *SetTaskStatusReq, opts ...client.CallOption) (*BoolMsg, error) {
	req := c.c.NewRequest(c.name, "TaskService.SetTaskStatus", in)
	out := new(BoolMsg)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskService) AppendTaskLog(ctx context.Context, opts ...client.CallOption) (TaskService_AppendTaskLogService, error) {
	req := c.c.NewRequest(c.name, "TaskService.AppendTaskLog", &AppendTaskLogReq{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &taskServiceAppendTaskLog{stream}, nil
}

type TaskService_AppendTaskLogService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*AppendTaskLogReq) error
	Recv() (*BoolMsg, error)
}

type taskServiceAppendTaskLog struct {
	stream client.Stream
}

func (x *taskServiceAppendTaskLog) Close() error {
	return x.stream.Close()
}

func (x *taskServiceAppendTaskLog) Context() context.Context {
	return x.stream.Context()
}

func (x *taskServiceAppendTaskLog) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *taskServiceAppendTaskLog) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *taskServiceAppendTaskLog) Send(m *AppendTaskLogReq) error {
	return x.stream.Send(m)
}

func (x *taskServiceAppendTaskLog) Recv() (*BoolMsg, error) {
	m := new(BoolMsg)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for TaskService service

type TaskServiceHandler interface {
	// ListTasks 列出所有任务
	ListTasks(context.Context, *Empty, *Tasks) error
	// CallTask 调用任务
	CallTask(context.Context, *CallTaskReq, *CallTaskRsp) error
	// SetTaskStatus 设置任务状态
	SetTaskStatus(context.Context, *SetTaskStatusReq, *BoolMsg) error
	// AppendTaskLog 追加任务日志
	AppendTaskLog(context.Context, TaskService_AppendTaskLogStream) error
}

func RegisterTaskServiceHandler(s server.Server, hdlr TaskServiceHandler, opts ...server.HandlerOption) error {
	type taskService interface {
		ListTasks(ctx context.Context, in *Empty, out *Tasks) error
		CallTask(ctx context.Context, in *CallTaskReq, out *CallTaskRsp) error
		SetTaskStatus(ctx context.Context, in *SetTaskStatusReq, out *BoolMsg) error
		AppendTaskLog(ctx context.Context, stream server.Stream) error
	}
	type TaskService struct {
		taskService
	}
	h := &taskServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&TaskService{h}, opts...))
}

type taskServiceHandler struct {
	TaskServiceHandler
}

func (h *taskServiceHandler) ListTasks(ctx context.Context, in *Empty, out *Tasks) error {
	return h.TaskServiceHandler.ListTasks(ctx, in, out)
}

func (h *taskServiceHandler) CallTask(ctx context.Context, in *CallTaskReq, out *CallTaskRsp) error {
	return h.TaskServiceHandler.CallTask(ctx, in, out)
}

func (h *taskServiceHandler) SetTaskStatus(ctx context.Context, in *SetTaskStatusReq, out *BoolMsg) error {
	return h.TaskServiceHandler.SetTaskStatus(ctx, in, out)
}

func (h *taskServiceHandler) AppendTaskLog(ctx context.Context, stream server.Stream) error {
	return h.TaskServiceHandler.AppendTaskLog(ctx, &taskServiceAppendTaskLogStream{stream})
}

type TaskService_AppendTaskLogStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*BoolMsg) error
	Recv() (*AppendTaskLogReq, error)
}

type taskServiceAppendTaskLogStream struct {
	stream server.Stream
}

func (x *taskServiceAppendTaskLogStream) Close() error {
	return x.stream.Close()
}

func (x *taskServiceAppendTaskLogStream) Context() context.Context {
	return x.stream.Context()
}

func (x *taskServiceAppendTaskLogStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *taskServiceAppendTaskLogStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *taskServiceAppendTaskLogStream) Send(m *BoolMsg) error {
	return x.stream.Send(m)
}

func (x *taskServiceAppendTaskLogStream) Recv() (*AppendTaskLogReq, error) {
	m := new(AppendTaskLogReq)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
