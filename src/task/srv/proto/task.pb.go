// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: proto/task.proto

package task

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Tasks struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tasks []*Task `protobuf:"bytes,1,rep,name=tasks,proto3" json:"tasks,omitempty"`
}

func (x *Tasks) Reset() {
	*x = Tasks{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_task_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tasks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tasks) ProtoMessage() {}

func (x *Tasks) ProtoReflect() protoreflect.Message {
	mi := &file_proto_task_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tasks.ProtoReflect.Descriptor instead.
func (*Tasks) Descriptor() ([]byte, []int) {
	return file_proto_task_proto_rawDescGZIP(), []int{0}
}

func (x *Tasks) GetTasks() []*Task {
	if x != nil {
		return x.Tasks
	}
	return nil
}

type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Codename    string `protobuf:"bytes,2,opt,name=codename,proto3" json:"codename,omitempty"`
	Arguments   string `protobuf:"bytes,3,opt,name=arguments,proto3" json:"arguments,omitempty"`
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_task_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_proto_task_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_proto_task_proto_rawDescGZIP(), []int{1}
}

func (x *Task) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Task) GetCodename() string {
	if x != nil {
		return x.Codename
	}
	return ""
}

func (x *Task) GetArguments() string {
	if x != nil {
		return x.Arguments
	}
	return ""
}

func (x *Task) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type CallTaskReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskCodename string `protobuf:"bytes,1,opt,name=task_codename,json=taskCodename,proto3" json:"task_codename,omitempty"`
	Timeout      int64  `protobuf:"varint,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Arguments    string `protobuf:"bytes,2,opt,name=arguments,proto3" json:"arguments,omitempty"`
	Caller       string `protobuf:"bytes,4,opt,name=caller,proto3" json:"caller,omitempty"`
}

func (x *CallTaskReq) Reset() {
	*x = CallTaskReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_task_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CallTaskReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallTaskReq) ProtoMessage() {}

func (x *CallTaskReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_task_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallTaskReq.ProtoReflect.Descriptor instead.
func (*CallTaskReq) Descriptor() ([]byte, []int) {
	return file_proto_task_proto_rawDescGZIP(), []int{2}
}

func (x *CallTaskReq) GetTaskCodename() string {
	if x != nil {
		return x.TaskCodename
	}
	return ""
}

func (x *CallTaskReq) GetTimeout() int64 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *CallTaskReq) GetArguments() string {
	if x != nil {
		return x.Arguments
	}
	return ""
}

func (x *CallTaskReq) GetCaller() string {
	if x != nil {
		return x.Caller
	}
	return ""
}

type CallTaskRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskUniqueId string `protobuf:"bytes,1,opt,name=task_unique_id,json=taskUniqueId,proto3" json:"task_unique_id,omitempty"`
}

func (x *CallTaskRsp) Reset() {
	*x = CallTaskRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_task_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CallTaskRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallTaskRsp) ProtoMessage() {}

func (x *CallTaskRsp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_task_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallTaskRsp.ProtoReflect.Descriptor instead.
func (*CallTaskRsp) Descriptor() ([]byte, []int) {
	return file_proto_task_proto_rawDescGZIP(), []int{3}
}

func (x *CallTaskRsp) GetTaskUniqueId() string {
	if x != nil {
		return x.TaskUniqueId
	}
	return ""
}

type SetTaskStatusReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskUniqueId string `protobuf:"bytes,1,opt,name=task_unique_id,json=taskUniqueId,proto3" json:"task_unique_id,omitempty"`
	Status       int32  `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *SetTaskStatusReq) Reset() {
	*x = SetTaskStatusReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_task_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetTaskStatusReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetTaskStatusReq) ProtoMessage() {}

func (x *SetTaskStatusReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_task_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetTaskStatusReq.ProtoReflect.Descriptor instead.
func (*SetTaskStatusReq) Descriptor() ([]byte, []int) {
	return file_proto_task_proto_rawDescGZIP(), []int{4}
}

func (x *SetTaskStatusReq) GetTaskUniqueId() string {
	if x != nil {
		return x.TaskUniqueId
	}
	return ""
}

func (x *SetTaskStatusReq) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type AppendTaskLogReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskUniqueId string `protobuf:"bytes,1,opt,name=task_unique_id,json=taskUniqueId,proto3" json:"task_unique_id,omitempty"`
	Content      string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *AppendTaskLogReq) Reset() {
	*x = AppendTaskLogReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_task_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendTaskLogReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendTaskLogReq) ProtoMessage() {}

func (x *AppendTaskLogReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_task_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendTaskLogReq.ProtoReflect.Descriptor instead.
func (*AppendTaskLogReq) Descriptor() ([]byte, []int) {
	return file_proto_task_proto_rawDescGZIP(), []int{5}
}

func (x *AppendTaskLogReq) GetTaskUniqueId() string {
	if x != nil {
		return x.TaskUniqueId
	}
	return ""
}

func (x *AppendTaskLogReq) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type BoolMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *BoolMsg) Reset() {
	*x = BoolMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_task_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoolMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoolMsg) ProtoMessage() {}

func (x *BoolMsg) ProtoReflect() protoreflect.Message {
	mi := &file_proto_task_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoolMsg.ProtoReflect.Descriptor instead.
func (*BoolMsg) Descriptor() ([]byte, []int) {
	return file_proto_task_proto_rawDescGZIP(), []int{6}
}

func (x *BoolMsg) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_task_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_proto_task_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_proto_task_proto_rawDescGZIP(), []int{7}
}

var File_proto_task_proto protoreflect.FileDescriptor

var file_proto_task_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x04, 0x74, 0x61, 0x73, 0x6b, 0x22, 0x29, 0x0a, 0x05, 0x54, 0x61, 0x73, 0x6b,
	0x73, 0x12, 0x20, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x05, 0x74, 0x61,
	0x73, 0x6b, 0x73, 0x22, 0x72, 0x0a, 0x04, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x63,
	0x6f, 0x64, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63,
	0x6f, 0x64, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x72, 0x67, 0x75, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x72, 0x67, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x82, 0x01, 0x0a, 0x0b, 0x43, 0x61, 0x6c, 0x6c,
	0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x61, 0x73, 0x6b, 0x5f,
	0x63, 0x6f, 0x64, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x74, 0x61, 0x73, 0x6b, 0x43, 0x6f, 0x64, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x74,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x72, 0x67, 0x75, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x22, 0x33, 0x0a, 0x0b,
	0x43, 0x61, 0x6c, 0x6c, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x73, 0x70, 0x12, 0x24, 0x0a, 0x0e, 0x74,
	0x61, 0x73, 0x6b, 0x5f, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x61, 0x73, 0x6b, 0x55, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x49,
	0x64, 0x22, 0x50, 0x0a, 0x10, 0x53, 0x65, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x65, 0x71, 0x12, 0x24, 0x0a, 0x0e, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x75, 0x6e,
	0x69, 0x71, 0x75, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74,
	0x61, 0x73, 0x6b, 0x55, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x52, 0x0a, 0x10, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x54, 0x61, 0x73,
	0x6b, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x12, 0x24, 0x0a, 0x0e, 0x74, 0x61, 0x73, 0x6b, 0x5f,
	0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x74, 0x61, 0x73, 0x6b, 0x55, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x19, 0x0a, 0x07, 0x42, 0x6f, 0x6f, 0x6c, 0x4d,
	0x73, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02,
	0x6f, 0x6b, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0xe2, 0x01, 0x0a, 0x0b,
	0x54, 0x61, 0x73, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x27, 0x0a, 0x09, 0x4c,
	0x69, 0x73, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x0b, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0b, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x54, 0x61, 0x73,
	0x6b, 0x73, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x08, 0x43, 0x61, 0x6c, 0x6c, 0x54, 0x61, 0x73, 0x6b,
	0x12, 0x11, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x43, 0x61, 0x6c, 0x6c, 0x54, 0x61, 0x73, 0x6b,
	0x52, 0x65, 0x71, 0x1a, 0x11, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x43, 0x61, 0x6c, 0x6c, 0x54,
	0x61, 0x73, 0x6b, 0x52, 0x73, 0x70, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x0d, 0x53, 0x65, 0x74, 0x54,
	0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x2e, 0x74, 0x61, 0x73, 0x6b,
	0x2e, 0x53, 0x65, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65,
	0x71, 0x1a, 0x0d, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x4d, 0x73, 0x67,
	0x22, 0x00, 0x12, 0x3c, 0x0a, 0x0d, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x54, 0x61, 0x73, 0x6b,
	0x4c, 0x6f, 0x67, 0x12, 0x16, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e,
	0x64, 0x54, 0x61, 0x73, 0x6b, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x1a, 0x0d, 0x2e, 0x74, 0x61,
	0x73, 0x6b, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x4d, 0x73, 0x67, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01,
	0x42, 0x0c, 0x5a, 0x0a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x74, 0x61, 0x73, 0x6b, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_task_proto_rawDescOnce sync.Once
	file_proto_task_proto_rawDescData = file_proto_task_proto_rawDesc
)

func file_proto_task_proto_rawDescGZIP() []byte {
	file_proto_task_proto_rawDescOnce.Do(func() {
		file_proto_task_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_task_proto_rawDescData)
	})
	return file_proto_task_proto_rawDescData
}

var file_proto_task_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_task_proto_goTypes = []interface{}{
	(*Tasks)(nil),            // 0: task.Tasks
	(*Task)(nil),             // 1: task.Task
	(*CallTaskReq)(nil),      // 2: task.CallTaskReq
	(*CallTaskRsp)(nil),      // 3: task.CallTaskRsp
	(*SetTaskStatusReq)(nil), // 4: task.SetTaskStatusReq
	(*AppendTaskLogReq)(nil), // 5: task.AppendTaskLogReq
	(*BoolMsg)(nil),          // 6: task.BoolMsg
	(*Empty)(nil),            // 7: task.Empty
}
var file_proto_task_proto_depIdxs = []int32{
	1, // 0: task.Tasks.tasks:type_name -> task.Task
	7, // 1: task.TaskService.ListTasks:input_type -> task.Empty
	2, // 2: task.TaskService.CallTask:input_type -> task.CallTaskReq
	4, // 3: task.TaskService.SetTaskStatus:input_type -> task.SetTaskStatusReq
	5, // 4: task.TaskService.AppendTaskLog:input_type -> task.AppendTaskLogReq
	0, // 5: task.TaskService.ListTasks:output_type -> task.Tasks
	3, // 6: task.TaskService.CallTask:output_type -> task.CallTaskRsp
	6, // 7: task.TaskService.SetTaskStatus:output_type -> task.BoolMsg
	6, // 8: task.TaskService.AppendTaskLog:output_type -> task.BoolMsg
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_task_proto_init() }
func file_proto_task_proto_init() {
	if File_proto_task_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_task_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tasks); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_task_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Task); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_task_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CallTaskReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_task_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CallTaskRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_task_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetTaskStatusReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_task_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendTaskLogReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_task_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoolMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_task_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_task_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_task_proto_goTypes,
		DependencyIndexes: file_proto_task_proto_depIdxs,
		MessageInfos:      file_proto_task_proto_msgTypes,
	}.Build()
	File_proto_task_proto = out.File
	file_proto_task_proto_rawDesc = nil
	file_proto_task_proto_goTypes = nil
	file_proto_task_proto_depIdxs = nil
}