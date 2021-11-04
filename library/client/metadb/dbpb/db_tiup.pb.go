// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.2
// source: db_tiup.proto

package dbpb

import (
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

type TiupTaskStatus int32

const (
	TiupTaskStatus_Init       TiupTaskStatus = 0
	TiupTaskStatus_Processing TiupTaskStatus = 1
	TiupTaskStatus_Finished   TiupTaskStatus = 2
	TiupTaskStatus_Error      TiupTaskStatus = 3
)

// Enum value maps for TiupTaskStatus.
var (
	TiupTaskStatus_name = map[int32]string{
		0: "Init",
		1: "Processing",
		2: "Finished",
		3: "Error",
	}
	TiupTaskStatus_value = map[string]int32{
		"Init":       0,
		"Processing": 1,
		"Finished":   2,
		"Error":      3,
	}
)

func (x TiupTaskStatus) Enum() *TiupTaskStatus {
	p := new(TiupTaskStatus)
	*p = x
	return p
}

func (x TiupTaskStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TiupTaskStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_db_tiup_proto_enumTypes[0].Descriptor()
}

func (TiupTaskStatus) Type() protoreflect.EnumType {
	return &file_db_tiup_proto_enumTypes[0]
}

func (x TiupTaskStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TiupTaskStatus.Descriptor instead.
func (TiupTaskStatus) EnumDescriptor() ([]byte, []int) {
	return file_db_tiup_proto_rawDescGZIP(), []int{0}
}

type TiupTaskType int32

const (
	TiupTaskType_Deploy    TiupTaskType = 0
	TiupTaskType_Start     TiupTaskType = 1
	TiupTaskType_Destroy   TiupTaskType = 2
	TiupTaskType_List      TiupTaskType = 3
	TiupTaskType_Dumpling  TiupTaskType = 4
	TiupTaskType_Lightning TiupTaskType = 5
	TiupTaskType_Backup    TiupTaskType = 6
	TiupTaskType_Restore   TiupTaskType = 7
	TiupTaskType_Restart   TiupTaskType = 8
	TiupTaskType_Stop      TiupTaskType = 9
)

// Enum value maps for TiupTaskType.
var (
	TiupTaskType_name = map[int32]string{
		0: "Deploy",
		1: "Start",
		2: "Destroy",
		3: "List",
		4: "Dumpling",
		5: "Lightning",
		6: "Backup",
		7: "Restore",
		8: "Restart",
		9: "Stop",
	}
	TiupTaskType_value = map[string]int32{
		"Deploy":    0,
		"Start":     1,
		"Destroy":   2,
		"List":      3,
		"Dumpling":  4,
		"Lightning": 5,
		"Backup":    6,
		"Restore":   7,
		"Restart":   8,
		"Stop":      9,
	}
)

func (x TiupTaskType) Enum() *TiupTaskType {
	p := new(TiupTaskType)
	*p = x
	return p
}

func (x TiupTaskType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TiupTaskType) Descriptor() protoreflect.EnumDescriptor {
	return file_db_tiup_proto_enumTypes[1].Descriptor()
}

func (TiupTaskType) Type() protoreflect.EnumType {
	return &file_db_tiup_proto_enumTypes[1]
}

func (x TiupTaskType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TiupTaskType.Descriptor instead.
func (TiupTaskType) EnumDescriptor() ([]byte, []int) {
	return file_db_tiup_proto_rawDescGZIP(), []int{1}
}

type TiupTask struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        uint64         `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	CreatedAt string         `protobuf:"bytes,2,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	UpdatedAt string         `protobuf:"bytes,3,opt,name=UpdatedAt,proto3" json:"UpdatedAt,omitempty"`
	DeletedAt string         `protobuf:"bytes,4,opt,name=DeletedAt,proto3" json:"DeletedAt,omitempty"`
	Type      TiupTaskType   `protobuf:"varint,5,opt,name=Type,proto3,enum=TiupTaskType" json:"Type,omitempty"`
	Status    TiupTaskStatus `protobuf:"varint,6,opt,name=Status,proto3,enum=TiupTaskStatus" json:"Status,omitempty"`
	ErrorStr  string         `protobuf:"bytes,7,opt,name=ErrorStr,proto3" json:"ErrorStr,omitempty"`
	BizID     uint64         `protobuf:"varint,8,opt,name=BizID,proto3" json:"BizID,omitempty"`
}

func (x *TiupTask) Reset() {
	*x = TiupTask{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_tiup_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TiupTask) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TiupTask) ProtoMessage() {}

func (x *TiupTask) ProtoReflect() protoreflect.Message {
	mi := &file_db_tiup_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TiupTask.ProtoReflect.Descriptor instead.
func (*TiupTask) Descriptor() ([]byte, []int) {
	return file_db_tiup_proto_rawDescGZIP(), []int{0}
}

func (x *TiupTask) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *TiupTask) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *TiupTask) GetUpdatedAt() string {
	if x != nil {
		return x.UpdatedAt
	}
	return ""
}

func (x *TiupTask) GetDeletedAt() string {
	if x != nil {
		return x.DeletedAt
	}
	return ""
}

func (x *TiupTask) GetType() TiupTaskType {
	if x != nil {
		return x.Type
	}
	return TiupTaskType_Deploy
}

func (x *TiupTask) GetStatus() TiupTaskStatus {
	if x != nil {
		return x.Status
	}
	return TiupTaskStatus_Init
}

func (x *TiupTask) GetErrorStr() string {
	if x != nil {
		return x.ErrorStr
	}
	return ""
}

func (x *TiupTask) GetBizID() uint64 {
	if x != nil {
		return x.BizID
	}
	return 0
}

type CreateTiupTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type TiupTaskType `protobuf:"varint,1,opt,name=type,proto3,enum=TiupTaskType" json:"type,omitempty"`
	// zero means bizID has no effect
	BizID uint64 `protobuf:"varint,2,opt,name=bizID,proto3" json:"bizID,omitempty"`
}

func (x *CreateTiupTaskRequest) Reset() {
	*x = CreateTiupTaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_tiup_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTiupTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTiupTaskRequest) ProtoMessage() {}

func (x *CreateTiupTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_db_tiup_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTiupTaskRequest.ProtoReflect.Descriptor instead.
func (*CreateTiupTaskRequest) Descriptor() ([]byte, []int) {
	return file_db_tiup_proto_rawDescGZIP(), []int{1}
}

func (x *CreateTiupTaskRequest) GetType() TiupTaskType {
	if x != nil {
		return x.Type
	}
	return TiupTaskType_Deploy
}

func (x *CreateTiupTaskRequest) GetBizID() uint64 {
	if x != nil {
		return x.BizID
	}
	return 0
}

type CreateTiupTaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// zero means ok
	ErrCode int32  `protobuf:"varint,1,opt,name=errCode,proto3" json:"errCode,omitempty"`
	ErrStr  string `protobuf:"bytes,2,opt,name=errStr,proto3" json:"errStr,omitempty"`
	Id      uint64 `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CreateTiupTaskResponse) Reset() {
	*x = CreateTiupTaskResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_tiup_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTiupTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTiupTaskResponse) ProtoMessage() {}

func (x *CreateTiupTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_db_tiup_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTiupTaskResponse.ProtoReflect.Descriptor instead.
func (*CreateTiupTaskResponse) Descriptor() ([]byte, []int) {
	return file_db_tiup_proto_rawDescGZIP(), []int{2}
}

func (x *CreateTiupTaskResponse) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

func (x *CreateTiupTaskResponse) GetErrStr() string {
	if x != nil {
		return x.ErrStr
	}
	return ""
}

func (x *CreateTiupTaskResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type UpdateTiupTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint64         `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Status TiupTaskStatus `protobuf:"varint,2,opt,name=status,proto3,enum=TiupTaskStatus" json:"status,omitempty"`
	ErrStr string         `protobuf:"bytes,3,opt,name=errStr,proto3" json:"errStr,omitempty"`
}

func (x *UpdateTiupTaskRequest) Reset() {
	*x = UpdateTiupTaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_tiup_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateTiupTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTiupTaskRequest) ProtoMessage() {}

func (x *UpdateTiupTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_db_tiup_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTiupTaskRequest.ProtoReflect.Descriptor instead.
func (*UpdateTiupTaskRequest) Descriptor() ([]byte, []int) {
	return file_db_tiup_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateTiupTaskRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateTiupTaskRequest) GetStatus() TiupTaskStatus {
	if x != nil {
		return x.Status
	}
	return TiupTaskStatus_Init
}

func (x *UpdateTiupTaskRequest) GetErrStr() string {
	if x != nil {
		return x.ErrStr
	}
	return ""
}

type UpdateTiupTaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// zero means ok
	ErrCode int32  `protobuf:"varint,1,opt,name=errCode,proto3" json:"errCode,omitempty"`
	ErrStr  string `protobuf:"bytes,2,opt,name=errStr,proto3" json:"errStr,omitempty"`
}

func (x *UpdateTiupTaskResponse) Reset() {
	*x = UpdateTiupTaskResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_tiup_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateTiupTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTiupTaskResponse) ProtoMessage() {}

func (x *UpdateTiupTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_db_tiup_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTiupTaskResponse.ProtoReflect.Descriptor instead.
func (*UpdateTiupTaskResponse) Descriptor() ([]byte, []int) {
	return file_db_tiup_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateTiupTaskResponse) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

func (x *UpdateTiupTaskResponse) GetErrStr() string {
	if x != nil {
		return x.ErrStr
	}
	return ""
}

type FindTiupTaskByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *FindTiupTaskByIDRequest) Reset() {
	*x = FindTiupTaskByIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_tiup_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindTiupTaskByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindTiupTaskByIDRequest) ProtoMessage() {}

func (x *FindTiupTaskByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_db_tiup_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindTiupTaskByIDRequest.ProtoReflect.Descriptor instead.
func (*FindTiupTaskByIDRequest) Descriptor() ([]byte, []int) {
	return file_db_tiup_proto_rawDescGZIP(), []int{5}
}

func (x *FindTiupTaskByIDRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type FindTiupTaskByIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// zero means ok
	ErrCode  int32     `protobuf:"varint,1,opt,name=errCode,proto3" json:"errCode,omitempty"`
	ErrStr   string    `protobuf:"bytes,2,opt,name=errStr,proto3" json:"errStr,omitempty"`
	TiupTask *TiupTask `protobuf:"bytes,3,opt,name=tiupTask,proto3" json:"tiupTask,omitempty"`
}

func (x *FindTiupTaskByIDResponse) Reset() {
	*x = FindTiupTaskByIDResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_tiup_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindTiupTaskByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindTiupTaskByIDResponse) ProtoMessage() {}

func (x *FindTiupTaskByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_db_tiup_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindTiupTaskByIDResponse.ProtoReflect.Descriptor instead.
func (*FindTiupTaskByIDResponse) Descriptor() ([]byte, []int) {
	return file_db_tiup_proto_rawDescGZIP(), []int{6}
}

func (x *FindTiupTaskByIDResponse) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

func (x *FindTiupTaskByIDResponse) GetErrStr() string {
	if x != nil {
		return x.ErrStr
	}
	return ""
}

func (x *FindTiupTaskByIDResponse) GetTiupTask() *TiupTask {
	if x != nil {
		return x.TiupTask
	}
	return nil
}

type GetTiupTaskStatusByBizIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BizID uint64 `protobuf:"varint,1,opt,name=bizID,proto3" json:"bizID,omitempty"`
}

func (x *GetTiupTaskStatusByBizIDRequest) Reset() {
	*x = GetTiupTaskStatusByBizIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_tiup_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTiupTaskStatusByBizIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTiupTaskStatusByBizIDRequest) ProtoMessage() {}

func (x *GetTiupTaskStatusByBizIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_db_tiup_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTiupTaskStatusByBizIDRequest.ProtoReflect.Descriptor instead.
func (*GetTiupTaskStatusByBizIDRequest) Descriptor() ([]byte, []int) {
	return file_db_tiup_proto_rawDescGZIP(), []int{7}
}

func (x *GetTiupTaskStatusByBizIDRequest) GetBizID() uint64 {
	if x != nil {
		return x.BizID
	}
	return 0
}

type GetTiupTaskStatusByBizIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// zero means ok
	ErrCode    int32          `protobuf:"varint,1,opt,name=errCode,proto3" json:"errCode,omitempty"`
	ErrStr     string         `protobuf:"bytes,2,opt,name=errStr,proto3" json:"errStr,omitempty"`
	Stat       TiupTaskStatus `protobuf:"varint,3,opt,name=stat,proto3,enum=TiupTaskStatus" json:"stat,omitempty"`
	StatErrStr string         `protobuf:"bytes,4,opt,name=statErrStr,proto3" json:"statErrStr,omitempty"`
}

func (x *GetTiupTaskStatusByBizIDResponse) Reset() {
	*x = GetTiupTaskStatusByBizIDResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_tiup_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTiupTaskStatusByBizIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTiupTaskStatusByBizIDResponse) ProtoMessage() {}

func (x *GetTiupTaskStatusByBizIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_db_tiup_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTiupTaskStatusByBizIDResponse.ProtoReflect.Descriptor instead.
func (*GetTiupTaskStatusByBizIDResponse) Descriptor() ([]byte, []int) {
	return file_db_tiup_proto_rawDescGZIP(), []int{8}
}

func (x *GetTiupTaskStatusByBizIDResponse) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

func (x *GetTiupTaskStatusByBizIDResponse) GetErrStr() string {
	if x != nil {
		return x.ErrStr
	}
	return ""
}

func (x *GetTiupTaskStatusByBizIDResponse) GetStat() TiupTaskStatus {
	if x != nil {
		return x.Stat
	}
	return TiupTaskStatus_Init
}

func (x *GetTiupTaskStatusByBizIDResponse) GetStatErrStr() string {
	if x != nil {
		return x.StatErrStr
	}
	return ""
}

var File_db_tiup_proto protoreflect.FileDescriptor

var file_db_tiup_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x64, 0x62, 0x5f, 0x74, 0x69, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xf2, 0x01, 0x0a, 0x08, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x0e, 0x0a, 0x02,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x21, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x27, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x54, 0x69, 0x75, 0x70,
	0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x53, 0x74, 0x72, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x53, 0x74, 0x72, 0x12, 0x14,
	0x0a, 0x05, 0x42, 0x69, 0x7a, 0x49, 0x44, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x42,
	0x69, 0x7a, 0x49, 0x44, 0x22, 0x50, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69,
	0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x54, 0x69,
	0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x62, 0x69, 0x7a, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x05, 0x62, 0x69, 0x7a, 0x49, 0x44, 0x22, 0x5a, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72,
	0x72, 0x53, 0x74, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x53,
	0x74, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x68, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x75, 0x70,
	0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x27, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x54, 0x69,
	0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x53, 0x74, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x53, 0x74, 0x72, 0x22, 0x4a, 0x0a, 0x16,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x53, 0x74, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x65, 0x72, 0x72, 0x53, 0x74, 0x72, 0x22, 0x29, 0x0a, 0x17, 0x46, 0x69, 0x6e, 0x64,
	0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x73, 0x0a, 0x18, 0x46, 0x69, 0x6e, 0x64, 0x54, 0x69, 0x75, 0x70, 0x54,
	0x61, 0x73, 0x6b, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72,
	0x53, 0x74, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x53, 0x74,
	0x72, 0x12, 0x25, 0x0a, 0x08, 0x74, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x08,
	0x74, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x22, 0x37, 0x0a, 0x1f, 0x47, 0x65, 0x74, 0x54,
	0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x79, 0x42,
	0x69, 0x7a, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x62,
	0x69, 0x7a, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x62, 0x69, 0x7a, 0x49,
	0x44, 0x22, 0x99, 0x01, 0x0a, 0x20, 0x47, 0x65, 0x74, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73,
	0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x79, 0x42, 0x69, 0x7a, 0x49, 0x44, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x53, 0x74, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x65, 0x72, 0x72, 0x53, 0x74, 0x72, 0x12, 0x23, 0x0a, 0x04, 0x73, 0x74, 0x61, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73,
	0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x04, 0x73, 0x74, 0x61, 0x74, 0x12, 0x1e, 0x0a,
	0x0a, 0x73, 0x74, 0x61, 0x74, 0x45, 0x72, 0x72, 0x53, 0x74, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x45, 0x72, 0x72, 0x53, 0x74, 0x72, 0x2a, 0x43, 0x0a,
	0x0e, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x08, 0x0a, 0x04, 0x49, 0x6e, 0x69, 0x74, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x46, 0x69, 0x6e,
	0x69, 0x73, 0x68, 0x65, 0x64, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x10, 0x03, 0x2a, 0x89, 0x01, 0x0a, 0x0c, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x10, 0x00, 0x12,
	0x09, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x72, 0x74, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x65,
	0x73, 0x74, 0x72, 0x6f, 0x79, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x10,
	0x03, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x75, 0x6d, 0x70, 0x6c, 0x69, 0x6e, 0x67, 0x10, 0x04, 0x12,
	0x0d, 0x0a, 0x09, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x6e, 0x69, 0x6e, 0x67, 0x10, 0x05, 0x12, 0x0a,
	0x0a, 0x06, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x10, 0x06, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x65,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x10, 0x07, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x10, 0x08, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x74, 0x6f, 0x70, 0x10, 0x09, 0x42, 0x11,
	0x5a, 0x0f, 0x2e, 0x2f, 0x2e, 0x2e, 0x2f, 0x64, 0x62, 0x70, 0x62, 0x2f, 0x3b, 0x64, 0x62, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_db_tiup_proto_rawDescOnce sync.Once
	file_db_tiup_proto_rawDescData = file_db_tiup_proto_rawDesc
)

func file_db_tiup_proto_rawDescGZIP() []byte {
	file_db_tiup_proto_rawDescOnce.Do(func() {
		file_db_tiup_proto_rawDescData = protoimpl.X.CompressGZIP(file_db_tiup_proto_rawDescData)
	})
	return file_db_tiup_proto_rawDescData
}

var file_db_tiup_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_db_tiup_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_db_tiup_proto_goTypes = []interface{}{
	(TiupTaskStatus)(0),                      // 0: TiupTaskStatus
	(TiupTaskType)(0),                        // 1: TiupTaskType
	(*TiupTask)(nil),                         // 2: TiupTask
	(*CreateTiupTaskRequest)(nil),            // 3: CreateTiupTaskRequest
	(*CreateTiupTaskResponse)(nil),           // 4: CreateTiupTaskResponse
	(*UpdateTiupTaskRequest)(nil),            // 5: UpdateTiupTaskRequest
	(*UpdateTiupTaskResponse)(nil),           // 6: UpdateTiupTaskResponse
	(*FindTiupTaskByIDRequest)(nil),          // 7: FindTiupTaskByIDRequest
	(*FindTiupTaskByIDResponse)(nil),         // 8: FindTiupTaskByIDResponse
	(*GetTiupTaskStatusByBizIDRequest)(nil),  // 9: GetTiupTaskStatusByBizIDRequest
	(*GetTiupTaskStatusByBizIDResponse)(nil), // 10: GetTiupTaskStatusByBizIDResponse
}
var file_db_tiup_proto_depIdxs = []int32{
	1, // 0: TiupTask.Type:type_name -> TiupTaskType
	0, // 1: TiupTask.Status:type_name -> TiupTaskStatus
	1, // 2: CreateTiupTaskRequest.type:type_name -> TiupTaskType
	0, // 3: UpdateTiupTaskRequest.status:type_name -> TiupTaskStatus
	2, // 4: FindTiupTaskByIDResponse.tiupTask:type_name -> TiupTask
	0, // 5: GetTiupTaskStatusByBizIDResponse.stat:type_name -> TiupTaskStatus
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_db_tiup_proto_init() }
func file_db_tiup_proto_init() {
	if File_db_tiup_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_db_tiup_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TiupTask); i {
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
		file_db_tiup_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTiupTaskRequest); i {
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
		file_db_tiup_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTiupTaskResponse); i {
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
		file_db_tiup_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateTiupTaskRequest); i {
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
		file_db_tiup_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateTiupTaskResponse); i {
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
		file_db_tiup_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindTiupTaskByIDRequest); i {
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
		file_db_tiup_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindTiupTaskByIDResponse); i {
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
		file_db_tiup_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTiupTaskStatusByBizIDRequest); i {
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
		file_db_tiup_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTiupTaskStatusByBizIDResponse); i {
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
			RawDescriptor: file_db_tiup_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_db_tiup_proto_goTypes,
		DependencyIndexes: file_db_tiup_proto_depIdxs,
		EnumInfos:         file_db_tiup_proto_enumTypes,
		MessageInfos:      file_db_tiup_proto_msgTypes,
	}.Build()
	File_db_tiup_proto = out.File
	file_db_tiup_proto_rawDesc = nil
	file_db_tiup_proto_goTypes = nil
	file_db_tiup_proto_depIdxs = nil
}
