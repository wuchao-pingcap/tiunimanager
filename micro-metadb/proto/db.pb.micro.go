// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: db.proto

package db

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/asim/go-micro/v3/api"
	client "github.com/asim/go-micro/v3/client"
	server "github.com/asim/go-micro/v3/server"
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

// Api Endpoints for TiCPDBService service

func NewTiCPDBServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for TiCPDBService service

type TiCPDBService interface {
	// Auth Module
	FindTenant(ctx context.Context, in *DBFindTenantRequest, opts ...client.CallOption) (*DBFindTenantResponse, error)
	FindAccount(ctx context.Context, in *DBFindAccountRequest, opts ...client.CallOption) (*DBFindAccountResponse, error)
	SaveToken(ctx context.Context, in *DBSaveTokenRequest, opts ...client.CallOption) (*DBSaveTokenResponse, error)
	FindToken(ctx context.Context, in *DBFindTokenRequest, opts ...client.CallOption) (*DBFindTokenResponse, error)
	FindRolesByPermission(ctx context.Context, in *DBFindRolesByPermissionRequest, opts ...client.CallOption) (*DBFindRolesByPermissionResponse, error)
	// Host Module
	AddHost(ctx context.Context, in *DBAddHostRequest, opts ...client.CallOption) (*DBAddHostResponse, error)
	RemoveHost(ctx context.Context, in *DBRemoveHostRequest, opts ...client.CallOption) (*DBRemoveHostResponse, error)
	ListHost(ctx context.Context, in *DBListHostsRequest, opts ...client.CallOption) (*DBListHostsResponse, error)
	CheckDetails(ctx context.Context, in *DBCheckDetailsRequest, opts ...client.CallOption) (*DBCheckDetailsResponse, error)
	AllocHosts(ctx context.Context, in *DBAllocHostsRequest, opts ...client.CallOption) (*DBAllocHostResponse, error)
	AddCluster(ctx context.Context, in *DBCreateClusterRequest, opts ...client.CallOption) (*DBCreateClusterResponse, error)
	FindCluster(ctx context.Context, in *DBFindClusterRequest, opts ...client.CallOption) (*DBFindClusterResponse, error)
	UpdateTiUPConfig(ctx context.Context, in *DBUpdateTiUPConfigRequest, opts ...client.CallOption) (*DBUpdateTiUPConfigResponse, error)
	ListCluster(ctx context.Context, in *DBListClusterRequest, opts ...client.CallOption) (*DBListClusterResponse, error)
	// Tiup Task
	CreateTiupTask(ctx context.Context, in *CreateTiupTaskRequest, opts ...client.CallOption) (*CreateTiupTaskResponse, error)
	UpdateTiupTask(ctx context.Context, in *UpdateTiupTaskRequest, opts ...client.CallOption) (*UpdateTiupTaskResponse, error)
	FindTiupTaskByID(ctx context.Context, in *FindTiupTaskByIDRequest, opts ...client.CallOption) (*FindTiupTaskByIDResponse, error)
	GetTiupTaskStatusByBizID(ctx context.Context, in *GetTiupTaskStatusByBizIDRequest, opts ...client.CallOption) (*GetTiupTaskStatusByBizIDResponse, error)
}

type tiCPDBService struct {
	c    client.Client
	name string
}

func NewTiCPDBService(name string, c client.Client) TiCPDBService {
	return &tiCPDBService{
		c:    c,
		name: name,
	}
}

func (c *tiCPDBService) FindTenant(ctx context.Context, in *DBFindTenantRequest, opts ...client.CallOption) (*DBFindTenantResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.FindTenant", in)
	out := new(DBFindTenantResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) FindAccount(ctx context.Context, in *DBFindAccountRequest, opts ...client.CallOption) (*DBFindAccountResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.FindAccount", in)
	out := new(DBFindAccountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) SaveToken(ctx context.Context, in *DBSaveTokenRequest, opts ...client.CallOption) (*DBSaveTokenResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.SaveToken", in)
	out := new(DBSaveTokenResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) FindToken(ctx context.Context, in *DBFindTokenRequest, opts ...client.CallOption) (*DBFindTokenResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.FindToken", in)
	out := new(DBFindTokenResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) FindRolesByPermission(ctx context.Context, in *DBFindRolesByPermissionRequest, opts ...client.CallOption) (*DBFindRolesByPermissionResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.FindRolesByPermission", in)
	out := new(DBFindRolesByPermissionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) AddHost(ctx context.Context, in *DBAddHostRequest, opts ...client.CallOption) (*DBAddHostResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.AddHost", in)
	out := new(DBAddHostResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) RemoveHost(ctx context.Context, in *DBRemoveHostRequest, opts ...client.CallOption) (*DBRemoveHostResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.RemoveHost", in)
	out := new(DBRemoveHostResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) ListHost(ctx context.Context, in *DBListHostsRequest, opts ...client.CallOption) (*DBListHostsResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.ListHost", in)
	out := new(DBListHostsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) CheckDetails(ctx context.Context, in *DBCheckDetailsRequest, opts ...client.CallOption) (*DBCheckDetailsResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.CheckDetails", in)
	out := new(DBCheckDetailsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) AllocHosts(ctx context.Context, in *DBAllocHostsRequest, opts ...client.CallOption) (*DBAllocHostResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.AllocHosts", in)
	out := new(DBAllocHostResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) AddCluster(ctx context.Context, in *DBCreateClusterRequest, opts ...client.CallOption) (*DBCreateClusterResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.AddCluster", in)
	out := new(DBCreateClusterResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) FindCluster(ctx context.Context, in *DBFindClusterRequest, opts ...client.CallOption) (*DBFindClusterResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.FindCluster", in)
	out := new(DBFindClusterResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) UpdateTiUPConfig(ctx context.Context, in *DBUpdateTiUPConfigRequest, opts ...client.CallOption) (*DBUpdateTiUPConfigResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.UpdateTiUPConfig", in)
	out := new(DBUpdateTiUPConfigResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) ListCluster(ctx context.Context, in *DBListClusterRequest, opts ...client.CallOption) (*DBListClusterResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.ListCluster", in)
	out := new(DBListClusterResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) CreateTiupTask(ctx context.Context, in *CreateTiupTaskRequest, opts ...client.CallOption) (*CreateTiupTaskResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.CreateTiupTask", in)
	out := new(CreateTiupTaskResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) UpdateTiupTask(ctx context.Context, in *UpdateTiupTaskRequest, opts ...client.CallOption) (*UpdateTiupTaskResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.UpdateTiupTask", in)
	out := new(UpdateTiupTaskResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) FindTiupTaskByID(ctx context.Context, in *FindTiupTaskByIDRequest, opts ...client.CallOption) (*FindTiupTaskByIDResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.FindTiupTaskByID", in)
	out := new(FindTiupTaskByIDResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tiCPDBService) GetTiupTaskStatusByBizID(ctx context.Context, in *GetTiupTaskStatusByBizIDRequest, opts ...client.CallOption) (*GetTiupTaskStatusByBizIDResponse, error) {
	req := c.c.NewRequest(c.name, "TiCPDBService.GetTiupTaskStatusByBizID", in)
	out := new(GetTiupTaskStatusByBizIDResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for TiCPDBService service

type TiCPDBServiceHandler interface {
	// Auth Module
	FindTenant(context.Context, *DBFindTenantRequest, *DBFindTenantResponse) error
	FindAccount(context.Context, *DBFindAccountRequest, *DBFindAccountResponse) error
	SaveToken(context.Context, *DBSaveTokenRequest, *DBSaveTokenResponse) error
	FindToken(context.Context, *DBFindTokenRequest, *DBFindTokenResponse) error
	FindRolesByPermission(context.Context, *DBFindRolesByPermissionRequest, *DBFindRolesByPermissionResponse) error
	// Host Module
	AddHost(context.Context, *DBAddHostRequest, *DBAddHostResponse) error
	RemoveHost(context.Context, *DBRemoveHostRequest, *DBRemoveHostResponse) error
	ListHost(context.Context, *DBListHostsRequest, *DBListHostsResponse) error
	CheckDetails(context.Context, *DBCheckDetailsRequest, *DBCheckDetailsResponse) error
	AllocHosts(context.Context, *DBAllocHostsRequest, *DBAllocHostResponse) error
	AddCluster(context.Context, *DBCreateClusterRequest, *DBCreateClusterResponse) error
	FindCluster(context.Context, *DBFindClusterRequest, *DBFindClusterResponse) error
	UpdateTiUPConfig(context.Context, *DBUpdateTiUPConfigRequest, *DBUpdateTiUPConfigResponse) error
	ListCluster(context.Context, *DBListClusterRequest, *DBListClusterResponse) error
	// Tiup Task
	CreateTiupTask(context.Context, *CreateTiupTaskRequest, *CreateTiupTaskResponse) error
	UpdateTiupTask(context.Context, *UpdateTiupTaskRequest, *UpdateTiupTaskResponse) error
	FindTiupTaskByID(context.Context, *FindTiupTaskByIDRequest, *FindTiupTaskByIDResponse) error
	GetTiupTaskStatusByBizID(context.Context, *GetTiupTaskStatusByBizIDRequest, *GetTiupTaskStatusByBizIDResponse) error
}

func RegisterTiCPDBServiceHandler(s server.Server, hdlr TiCPDBServiceHandler, opts ...server.HandlerOption) error {
	type tiCPDBService interface {
		FindTenant(ctx context.Context, in *DBFindTenantRequest, out *DBFindTenantResponse) error
		FindAccount(ctx context.Context, in *DBFindAccountRequest, out *DBFindAccountResponse) error
		SaveToken(ctx context.Context, in *DBSaveTokenRequest, out *DBSaveTokenResponse) error
		FindToken(ctx context.Context, in *DBFindTokenRequest, out *DBFindTokenResponse) error
		FindRolesByPermission(ctx context.Context, in *DBFindRolesByPermissionRequest, out *DBFindRolesByPermissionResponse) error
		AddHost(ctx context.Context, in *DBAddHostRequest, out *DBAddHostResponse) error
		RemoveHost(ctx context.Context, in *DBRemoveHostRequest, out *DBRemoveHostResponse) error
		ListHost(ctx context.Context, in *DBListHostsRequest, out *DBListHostsResponse) error
		CheckDetails(ctx context.Context, in *DBCheckDetailsRequest, out *DBCheckDetailsResponse) error
		AllocHosts(ctx context.Context, in *DBAllocHostsRequest, out *DBAllocHostResponse) error
		AddCluster(ctx context.Context, in *DBCreateClusterRequest, out *DBCreateClusterResponse) error
		FindCluster(ctx context.Context, in *DBFindClusterRequest, out *DBFindClusterResponse) error
		UpdateTiUPConfig(ctx context.Context, in *DBUpdateTiUPConfigRequest, out *DBUpdateTiUPConfigResponse) error
		ListCluster(ctx context.Context, in *DBListClusterRequest, out *DBListClusterResponse) error
		CreateTiupTask(ctx context.Context, in *CreateTiupTaskRequest, out *CreateTiupTaskResponse) error
		UpdateTiupTask(ctx context.Context, in *UpdateTiupTaskRequest, out *UpdateTiupTaskResponse) error
		FindTiupTaskByID(ctx context.Context, in *FindTiupTaskByIDRequest, out *FindTiupTaskByIDResponse) error
		GetTiupTaskStatusByBizID(ctx context.Context, in *GetTiupTaskStatusByBizIDRequest, out *GetTiupTaskStatusByBizIDResponse) error
	}
	type TiCPDBService struct {
		tiCPDBService
	}
	h := &tiCPDBServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&TiCPDBService{h}, opts...))
}

type tiCPDBServiceHandler struct {
	TiCPDBServiceHandler
}

func (h *tiCPDBServiceHandler) FindTenant(ctx context.Context, in *DBFindTenantRequest, out *DBFindTenantResponse) error {
	return h.TiCPDBServiceHandler.FindTenant(ctx, in, out)
}

func (h *tiCPDBServiceHandler) FindAccount(ctx context.Context, in *DBFindAccountRequest, out *DBFindAccountResponse) error {
	return h.TiCPDBServiceHandler.FindAccount(ctx, in, out)
}

func (h *tiCPDBServiceHandler) SaveToken(ctx context.Context, in *DBSaveTokenRequest, out *DBSaveTokenResponse) error {
	return h.TiCPDBServiceHandler.SaveToken(ctx, in, out)
}

func (h *tiCPDBServiceHandler) FindToken(ctx context.Context, in *DBFindTokenRequest, out *DBFindTokenResponse) error {
	return h.TiCPDBServiceHandler.FindToken(ctx, in, out)
}

func (h *tiCPDBServiceHandler) FindRolesByPermission(ctx context.Context, in *DBFindRolesByPermissionRequest, out *DBFindRolesByPermissionResponse) error {
	return h.TiCPDBServiceHandler.FindRolesByPermission(ctx, in, out)
}

func (h *tiCPDBServiceHandler) AddHost(ctx context.Context, in *DBAddHostRequest, out *DBAddHostResponse) error {
	return h.TiCPDBServiceHandler.AddHost(ctx, in, out)
}

func (h *tiCPDBServiceHandler) RemoveHost(ctx context.Context, in *DBRemoveHostRequest, out *DBRemoveHostResponse) error {
	return h.TiCPDBServiceHandler.RemoveHost(ctx, in, out)
}

func (h *tiCPDBServiceHandler) ListHost(ctx context.Context, in *DBListHostsRequest, out *DBListHostsResponse) error {
	return h.TiCPDBServiceHandler.ListHost(ctx, in, out)
}

func (h *tiCPDBServiceHandler) CheckDetails(ctx context.Context, in *DBCheckDetailsRequest, out *DBCheckDetailsResponse) error {
	return h.TiCPDBServiceHandler.CheckDetails(ctx, in, out)
}

func (h *tiCPDBServiceHandler) AllocHosts(ctx context.Context, in *DBAllocHostsRequest, out *DBAllocHostResponse) error {
	return h.TiCPDBServiceHandler.AllocHosts(ctx, in, out)
}

func (h *tiCPDBServiceHandler) AddCluster(ctx context.Context, in *DBCreateClusterRequest, out *DBCreateClusterResponse) error {
	return h.TiCPDBServiceHandler.AddCluster(ctx, in, out)
}

func (h *tiCPDBServiceHandler) FindCluster(ctx context.Context, in *DBFindClusterRequest, out *DBFindClusterResponse) error {
	return h.TiCPDBServiceHandler.FindCluster(ctx, in, out)
}

func (h *tiCPDBServiceHandler) UpdateTiUPConfig(ctx context.Context, in *DBUpdateTiUPConfigRequest, out *DBUpdateTiUPConfigResponse) error {
	return h.TiCPDBServiceHandler.UpdateTiUPConfig(ctx, in, out)
}

func (h *tiCPDBServiceHandler) ListCluster(ctx context.Context, in *DBListClusterRequest, out *DBListClusterResponse) error {
	return h.TiCPDBServiceHandler.ListCluster(ctx, in, out)
}

func (h *tiCPDBServiceHandler) CreateTiupTask(ctx context.Context, in *CreateTiupTaskRequest, out *CreateTiupTaskResponse) error {
	return h.TiCPDBServiceHandler.CreateTiupTask(ctx, in, out)
}

func (h *tiCPDBServiceHandler) UpdateTiupTask(ctx context.Context, in *UpdateTiupTaskRequest, out *UpdateTiupTaskResponse) error {
	return h.TiCPDBServiceHandler.UpdateTiupTask(ctx, in, out)
}

func (h *tiCPDBServiceHandler) FindTiupTaskByID(ctx context.Context, in *FindTiupTaskByIDRequest, out *FindTiupTaskByIDResponse) error {
	return h.TiCPDBServiceHandler.FindTiupTaskByID(ctx, in, out)
}

func (h *tiCPDBServiceHandler) GetTiupTaskStatusByBizID(ctx context.Context, in *GetTiupTaskStatusByBizIDRequest, out *GetTiupTaskStatusByBizIDResponse) error {
	return h.TiCPDBServiceHandler.GetTiupTaskStatusByBizID(ctx, in, out)
}