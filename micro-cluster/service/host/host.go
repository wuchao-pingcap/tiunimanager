package host

import (
	"context"

	"github.com/pingcap-inc/tiem/library/client"
	"github.com/pingcap-inc/tiem/library/framework"

	hostPb "github.com/pingcap-inc/tiem/micro-cluster/proto"
	dbPb "github.com/pingcap-inc/tiem/micro-metadb/proto"
	"google.golang.org/grpc/codes"
)

type ResourceManager struct{}

func NewResourceManager() *ResourceManager {
	m := new(ResourceManager)
	return m
}

func copyHostToDBReq(src *hostPb.HostInfo, dst *dbPb.DBHostInfoDTO) {
	dst.HostName = src.HostName
	dst.Ip = src.Ip
	dst.UserName = src.UserName
	dst.Passwd = src.Passwd
	dst.Arch = src.Arch
	dst.Os = src.Os
	dst.Kernel = src.Kernel
	dst.FreeCpuCores = src.FreeCpuCores
	dst.FreeMemory = src.FreeMemory
	dst.Spec = src.Spec
	dst.CpuCores = src.CpuCores
	dst.Memory = src.Memory
	dst.Nic = src.Nic
	dst.Region = src.Region
	dst.Az = src.Az
	dst.Rack = src.Rack
	dst.Status = src.Status
	dst.Stat = src.Stat
	dst.Purpose = src.Purpose
	dst.DiskType = src.DiskType
	dst.Reserved = src.Reserved
	for _, disk := range src.Disks {
		dst.Disks = append(dst.Disks, &dbPb.DBDiskDTO{
			Name:     disk.Name,
			Path:     disk.Path,
			Capacity: disk.Capacity,
			Status:   disk.Status,
			Type:     disk.Type,
		})
	}
}

func copyHostFromDBRsp(src *dbPb.DBHostInfoDTO, dst *hostPb.HostInfo) {
	dst.HostId = src.HostId
	dst.HostName = src.HostName
	dst.Ip = src.Ip
	dst.Arch = src.Arch
	dst.Os = src.Os
	dst.Kernel = src.Kernel
	dst.FreeCpuCores = src.FreeCpuCores
	dst.FreeMemory = src.FreeMemory
	dst.Spec = src.Spec
	dst.CpuCores = src.CpuCores
	dst.Memory = src.Memory
	dst.Nic = src.Nic
	dst.Region = src.Region
	dst.Az = src.Az
	dst.Rack = src.Rack
	dst.Status = src.Status
	dst.Stat = src.Stat
	dst.Purpose = src.Purpose
	dst.DiskType = src.DiskType
	dst.CreateAt = src.CreateAt
	dst.Reserved = src.Reserved
	for _, disk := range src.Disks {
		dst.Disks = append(dst.Disks, &hostPb.Disk{
			DiskId:   disk.DiskId,
			Name:     disk.Name,
			Path:     disk.Path,
			Capacity: disk.Capacity,
			Status:   disk.Status,
			Type:     disk.Type,
			UsedBy:   disk.UsedBy,
		})
	}
}

func (m *ResourceManager) ImportHost(ctx context.Context, in *hostPb.ImportHostRequest, out *hostPb.ImportHostResponse) error {
	var req dbPb.DBAddHostRequest
	req.Host = new(dbPb.DBHostInfoDTO)
	copyHostToDBReq(in.Host, req.Host)
	var err error
	rsp, err := client.DBClient.AddHost(ctx, &req)
	if err != nil {
		framework.Log().Errorf("import host %s error, %v", req.Host.Ip, err)
		return err
	}
	out.Rs = new(hostPb.ResponseStatus)
	out.Rs.Code = rsp.Rs.Code
	out.Rs.Message = rsp.Rs.Message
	if rsp.Rs.Code != int32(codes.OK) {
		framework.Log().Warnf("import host failed from db service: %d, %s", rsp.Rs.Code, rsp.Rs.Message)
		return nil
	}
	framework.Log().Infof("import host(%s) succeed from db service: %s", in.Host.Ip, rsp.HostId)
	out.HostId = rsp.HostId

	return nil
}

func (m *ResourceManager) ImportHostsInBatch(ctx context.Context, in *hostPb.ImportHostsInBatchRequest, out *hostPb.ImportHostsInBatchResponse) error {
	var req dbPb.DBAddHostsInBatchRequest
	for _, v := range in.Hosts {
		var host dbPb.DBHostInfoDTO
		copyHostToDBReq(v, &host)
		req.Hosts = append(req.Hosts, &host)
	}
	var err error
	rsp, err := client.DBClient.AddHostsInBatch(ctx, &req)
	if err != nil {
		framework.Log().Errorf("import hosts in batch error, %v", err)
		return err
	}
	out.Rs = new(hostPb.ResponseStatus)
	out.Rs.Code = rsp.Rs.Code
	out.Rs.Message = rsp.Rs.Message
	if rsp.Rs.Code != int32(codes.OK) {
		framework.Log().Warnf("import hosts in batch failed from db service: %d, %s", rsp.Rs.Code, rsp.Rs.Message)
		return nil
	}
	framework.Log().Infof("import %d hosts in batch succeed from db service.", len(rsp.HostIds))
	out.HostIds = rsp.HostIds

	return nil
}

func (m *ResourceManager) RemoveHost(ctx context.Context, in *hostPb.RemoveHostRequest, out *hostPb.RemoveHostResponse) error {
	var req dbPb.DBRemoveHostRequest
	req.HostId = in.HostId
	rsp, err := client.DBClient.RemoveHost(ctx, &req)
	if err != nil {
		framework.Log().Errorf("remove host %s error, %v", req.HostId, err)
		return err
	}
	out.Rs = new(hostPb.ResponseStatus)
	out.Rs.Code = rsp.Rs.Code
	out.Rs.Message = rsp.Rs.Message
	if rsp.Rs.Code != int32(codes.OK) {
		framework.Log().Warnf("remove host %s failed from db service: %d, %s", req.HostId, rsp.Rs.Code, rsp.Rs.Message)
		return nil
	}

	framework.Log().Infof("remove host %s succeed from db service", req.HostId)
	return nil
}

func (m *ResourceManager) RemoveHostsInBatch(ctx context.Context, in *hostPb.RemoveHostsInBatchRequest, out *hostPb.RemoveHostsInBatchResponse) error {
	var req dbPb.DBRemoveHostsInBatchRequest
	req.HostIds = in.HostIds
	rsp, err := client.DBClient.RemoveHostsInBatch(ctx, &req)
	if err != nil {
		framework.Log().Errorf("remove hosts in batch error, %v", err)
		return err
	}
	out.Rs = new(hostPb.ResponseStatus)
	out.Rs.Code = rsp.Rs.Code
	out.Rs.Message = rsp.Rs.Message
	if rsp.Rs.Code != int32(codes.OK) {
		framework.Log().Warnf("remove hosts in batch failed from db service: %d, %s", rsp.Rs.Code, rsp.Rs.Message)
		return nil
	}

	framework.Log().Infof("remove %d hosts succeed from db service", len(req.HostIds))
	return nil
}

func (m *ResourceManager) ListHost(ctx context.Context, in *hostPb.ListHostsRequest, out *hostPb.ListHostsResponse) error {
	var req dbPb.DBListHostsRequest
	req.Purpose = in.Purpose
	req.Status = in.Status
	req.Page = new(dbPb.DBHostPageDTO)
	req.Page.Page = in.PageReq.Page
	req.Page.PageSize = in.PageReq.PageSize
	rsp, err := client.DBClient.ListHost(ctx, &req)
	if err != nil {
		framework.Log().Errorf("list hosts error, %v", err)
		return err
	}
	out.Rs = new(hostPb.ResponseStatus)
	out.PageReq = new(hostPb.PageDTO)
	out.Rs.Code = rsp.Rs.Code
	out.Rs.Message = rsp.Rs.Message

	if rsp.Rs.Code != int32(codes.OK) {
		framework.Log().Warnf("list hosts info from db service failed: %d, %s", rsp.Rs.Code, rsp.Rs.Message)
		return nil
	}

	framework.Log().Infof("list %d hosts info from db service succeed", len(rsp.HostList))
	for _, v := range rsp.HostList {
		var host hostPb.HostInfo
		copyHostFromDBRsp(v, &host)
		out.HostList = append(out.HostList, &host)
	}
	out.PageReq.Page = rsp.Page.Page
	out.PageReq.PageSize = rsp.Page.PageSize
	out.PageReq.Total = rsp.Page.Total
	return nil
}

func (m *ResourceManager) CheckDetails(ctx context.Context, in *hostPb.CheckDetailsRequest, out *hostPb.CheckDetailsResponse) error {
	var req dbPb.DBCheckDetailsRequest
	req.HostId = in.HostId
	rsp, err := client.DBClient.CheckDetails(ctx, &req)
	if err != nil {
		framework.Log().Errorf("check host %s details failed, %v", req.HostId, err)
		return err
	}
	out.Rs = new(hostPb.ResponseStatus)
	out.Rs.Code = rsp.Rs.Code
	out.Rs.Message = rsp.Rs.Message

	if rsp.Rs.Code != int32(codes.OK) {
		framework.Log().Warnf("check host %s details from db service failed: %d, %s", req.HostId, rsp.Rs.Code, rsp.Rs.Message)
		return nil
	}

	framework.Log().Infof("check host %s details from db service succeed", req.HostId)
	out.Details = new(hostPb.HostInfo)
	copyHostFromDBRsp(rsp.Details, out.Details)

	return nil
}

func makeAllocReq(src []*hostPb.AllocationReq) (dst []*dbPb.DBAllocationReq) {
	for _, req := range src {
		dst = append(dst, &dbPb.DBAllocationReq{
			FailureDomain: req.FailureDomain,
			CpuCores:      req.CpuCores,
			Memory:        req.Memory,
			Count:         req.Count,
			Purpose:       req.Purpose,
		})
	}
	return dst
}

func buildDiskFromDB(disk *dbPb.DBDiskDTO) *hostPb.Disk {
	var hd hostPb.Disk
	hd.DiskId = disk.DiskId
	hd.Name = disk.Name
	hd.Path = disk.Path
	hd.Capacity = disk.Capacity
	hd.Status = disk.Status

	return &hd
}

func getAllocRsp(src []*dbPb.DBAllocHostDTO) (dst []*hostPb.AllocHost) {
	for _, rsp := range src {
		dst = append(dst, &hostPb.AllocHost{
			HostName: rsp.HostName,
			Ip:       rsp.Ip,
			UserName: rsp.UserName,
			Passwd:   rsp.Passwd,
			CpuCores: rsp.CpuCores,
			Memory:   rsp.Memory,
			Disk:     buildDiskFromDB(rsp.Disk),
		})
	}
	return dst
}

func (m *ResourceManager) AllocHosts(ctx context.Context, in *hostPb.AllocHostsRequest, out *hostPb.AllocHostResponse) error {
	req := new(dbPb.DBAllocHostsRequest)
	req.PdReq = makeAllocReq(in.PdReq)
	req.TidbReq = makeAllocReq(in.TidbReq)
	req.TikvReq = makeAllocReq(in.TikvReq)

	rsp, err := client.DBClient.AllocHosts(ctx, req)
	if err != nil {
		framework.Log().Errorf("alloc hosts error, %v", err)
		return err
	}
	out.Rs = new(hostPb.ResponseStatus)
	out.Rs.Code = rsp.Rs.Code
	out.Rs.Message = rsp.Rs.Message

	if rsp.Rs.Code != int32(codes.OK) {
		framework.Log().Warnf("alloc hosts from db service failed: %d, %s", rsp.Rs.Code, rsp.Rs.Message)
		return nil
	}

	out.PdHosts = getAllocRsp(rsp.PdHosts)
	out.TidbHosts = getAllocRsp(rsp.TidbHosts)
	out.TikvHosts = getAllocRsp(rsp.TikvHosts)
	return nil
}

func (m *ResourceManager) GetFailureDomain(ctx context.Context, in *hostPb.GetFailureDomainRequest, out *hostPb.GetFailureDomainResponse) error {
	var req dbPb.DBGetFailureDomainRequest
	req.FailureDomainType = in.FailureDomainType
	rsp, err := client.DBClient.GetFailureDomain(ctx, &req)
	if err != nil {
		framework.Log().Errorf("get failure domains error, %v", err)
		return err
	}
	out.Rs = new(hostPb.ResponseStatus)
	out.Rs.Code = rsp.Rs.Code
	out.Rs.Message = rsp.Rs.Message

	if rsp.Rs.Code != int32(codes.OK) {
		framework.Log().Warnf("get failure domains info from db service failed: %d, %s", rsp.Rs.Code, rsp.Rs.Message)
		return nil
	}

	framework.Log().Infof("get failure domain type %d from db service succeed", req.FailureDomainType)
	for _, v := range rsp.FdList {
		out.FdList = append(out.FdList, &hostPb.FailureDomainResource{
			FailureDomain: v.FailureDomain,
			Purpose:       v.Purpose,
			Spec:          v.Spec,
			Count:         v.Count,
		})
	}
	return nil
}

func copyAllocRequirement(src *hostPb.AllocRequirement, dst *dbPb.DBAllocRequirement) {
	dst.Location = new(dbPb.DBLocation)
	dst.Location.Region = src.Location.Region
	dst.Location.Zone = src.Location.Zone
	dst.Location.Rack = src.Location.Rack
	dst.Location.Host = src.Location.Host

	dst.HostExcluded = new(dbPb.DBExcluded)
	if src.HostExcluded != nil {
		dst.HostExcluded.Hosts = append(dst.HostExcluded.Hosts, src.HostExcluded.Hosts...)
	}

	dst.HostFilter = new(dbPb.DBFilter)
	dst.HostFilter.Arch = src.HostFilter.Arch
	dst.HostFilter.Purpose = src.HostFilter.Purpose
	dst.HostFilter.DiskType = src.HostFilter.DiskType

	// copy requirement
	dst.Require = new(dbPb.DBRequirement)
	dst.Require.Exclusive = src.Require.Exclusive
	dst.Require.ComputeReq = new(dbPb.DBComputeRequirement)
	dst.Require.ComputeReq.CpuCores = src.Require.ComputeReq.CpuCores
	dst.Require.ComputeReq.Memory = src.Require.ComputeReq.Memory

	dst.Require.DiskReq = new(dbPb.DBDiskRequirement)
	dst.Require.DiskReq.NeedDisk = src.Require.DiskReq.NeedDisk
	dst.Require.DiskReq.DiskType = src.Require.DiskReq.DiskType
	dst.Require.DiskReq.Capacity = src.Require.DiskReq.Capacity

	for _, portReq := range src.Require.PortReq {
		dst.Require.PortReq = append(dst.Require.PortReq, &dbPb.DBPortRequirement{
			Start:   portReq.Start,
			End:     portReq.End,
			PortCnt: portReq.PortCnt,
		})
	}

	dst.Strategy = src.Strategy
	dst.Count = src.Count
}

func buildDBAllocRequest(src *hostPb.BatchAllocRequest, dst *dbPb.DBBatchAllocRequest) {
	for _, request := range src.BatchRequests {
		var dbReq dbPb.DBAllocRequest
		dbReq.Applicant = new(dbPb.DBApplicant)
		dbReq.Applicant.HolderId = request.Applicant.HolderId
		dbReq.Applicant.RequestId = request.Applicant.RequestId
		for _, require := range request.Requires {
			var dbRequire dbPb.DBAllocRequirement
			copyAllocRequirement(require, &dbRequire)
			dbReq.Requires = append(dbReq.Requires, &dbRequire)
		}
		dst.BatchRequests = append(dst.BatchRequests, &dbReq)
	}
}

func copyResourcesInfoFromRsp(src *dbPb.DBHostResource, dst *hostPb.HostResource) {
	dst.Reqseq = src.Reqseq
	dst.HostId = src.HostId
	dst.HostIp = src.HostIp
	dst.HostName = src.HostName
	dst.Passwd = src.Passwd
	dst.UserName = src.UserName
	dst.ComputeRes = new(hostPb.ComputeRequirement)
	dst.ComputeRes.CpuCores = src.ComputeRes.CpuCores
	dst.ComputeRes.Memory = src.ComputeRes.Memory
	dst.Location = new(hostPb.Location)
	dst.Location.Region = src.Location.Region
	dst.Location.Zone = src.Location.Zone
	dst.Location.Rack = src.Location.Rack
	dst.Location.Host = src.Location.Host
	dst.DiskRes = new(hostPb.DiskResource)
	dst.DiskRes.DiskId = src.DiskRes.DiskId
	dst.DiskRes.DiskName = src.DiskRes.DiskName
	dst.DiskRes.Path = src.DiskRes.Path
	dst.DiskRes.Type = src.DiskRes.Type
	dst.DiskRes.Capacity = src.DiskRes.Capacity
	for _, portRes := range src.PortRes {
		var portResource hostPb.PortResource
		portResource.Start = portRes.Start
		portResource.End = portRes.End
		portResource.Ports = append(portResource.Ports, portRes.Ports...)
		dst.PortRes = append(dst.PortRes, &portResource)
	}
}

func (m *ResourceManager) AllocResourcesInBatch(ctx context.Context, in *hostPb.BatchAllocRequest, out *hostPb.BatchAllocResponse) error {
	var req dbPb.DBBatchAllocRequest
	buildDBAllocRequest(in, &req)
	rsp, err := client.DBClient.AllocResourcesInBatch(ctx, &req)
	if err != nil {
		framework.Log().Errorf("alloc resources error, %v", err)
		return err
	}
	out.Rs = new(hostPb.AllocResponseStatus)
	out.Rs.Code = rsp.Rs.Code
	out.Rs.Message = rsp.Rs.Message

	if rsp.Rs.Code != int32(codes.OK) {
		framework.Log().Warnf("alloc resources from db service failed: %d, %s", rsp.Rs.Code, rsp.Rs.Message)
		return nil
	}

	framework.Log().Infof("alloc resources from db service succeed, holde id: %s, repest id: %s", in.BatchRequests[0].Applicant.HolderId, in.BatchRequests[0].Applicant.RequestId)
	for _, result := range rsp.BatchResults {
		var res hostPb.AllocResponse
		res.Rs = new(hostPb.AllocResponseStatus)
		res.Rs.Code = result.Rs.Code
		res.Rs.Message = result.Rs.Message
		for _, r := range result.Results {
			var hostResource hostPb.HostResource
			copyResourcesInfoFromRsp(r, &hostResource)
			res.Results = append(res.Results, &hostResource)
		}
		out.BatchResults = append(out.BatchResults, &res)
	}
	return nil
}

func buildDBRecycleRequire(src *hostPb.RecycleRequire, dst *dbPb.DBRecycleRequire) {
	dst.RecycleType = src.RecycleType
	dst.HolderId = src.HolderId
	dst.RequestId = src.RequestId
	dst.HostId = src.HostId
	dst.HostIp = src.HostIp
	dst.ComputeReq = new(dbPb.DBComputeRequirement)
	if src.ComputeReq != nil {
		dst.ComputeReq.CpuCores = src.ComputeReq.CpuCores
		dst.ComputeReq.Memory = src.ComputeReq.Memory
	}

	dst.DiskReq = new(dbPb.DBDiskResource)
	if src.DiskReq != nil {
		dst.DiskReq.DiskId = src.DiskReq.DiskId
		dst.DiskReq.DiskName = src.DiskReq.DiskName
		dst.DiskReq.Path = src.DiskReq.Path
		dst.DiskReq.Type = src.DiskReq.Type
		dst.DiskReq.Capacity = src.DiskReq.Capacity
	}

	for _, portReq := range src.PortReq {
		var portResource dbPb.DBPortResource
		portResource.Start = portReq.Start
		portResource.End = portReq.End
		portResource.Ports = append(portResource.Ports, portReq.Ports...)
		dst.PortReq = append(dst.PortReq, &portResource)
	}

}

func (m *ResourceManager) RecycleResources(ctx context.Context, in *hostPb.RecycleRequest, out *hostPb.RecycleResponse) error {
	var req dbPb.DBRecycleRequest
	for _, request := range in.RecycleReqs {
		var r dbPb.DBRecycleRequire
		buildDBRecycleRequire(request, &r)
		req.RecycleReqs = append(req.RecycleReqs, &r)
	}
	rsp, err := client.DBClient.RecycleResources(ctx, &req)
	if err != nil {
		framework.Log().Errorf("recycle resources for type %d error, %v", err, in.RecycleReqs[0].RecycleType)
		return err
	}
	out.Rs = new(hostPb.AllocResponseStatus)
	out.Rs.Code = rsp.Rs.Code
	out.Rs.Message = rsp.Rs.Message

	if rsp.Rs.Code != int32(codes.OK) {
		framework.Log().Warnf("recycle resources from db service failed: %d, %s", rsp.Rs.Code, rsp.Rs.Message)
		return nil
	}

	framework.Log().Infof("recycle resources from db service succeed, recycle type %d, holderId %s, requestId %s", in.RecycleReqs[0].RecycleType, in.RecycleReqs[0].HolderId, in.RecycleReqs[0].RequestId)

	return nil
}