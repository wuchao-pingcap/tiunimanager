/******************************************************************************
 * Copyright (c)  2021 PingCAP, Inc.                                          *
 * Licensed under the Apache License, Version 2.0 (the "License");            *
 * you may not use this file except in compliance with the License.           *
 * You may obtain a copy of the License at                                    *
 *                                                                            *
 * http://www.apache.org/licenses/LICENSE-2.0                                 *
 *                                                                            *
 * Unless required by applicable law or agreed to in writing, software        *
 * distributed under the License is distributed on an "AS IS" BASIS,          *
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.   *
 * See the License for the specific language governing permissions and        *
 * limitations under the License.                                             *
 *                                                                            *
 ******************************************************************************/

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pingcap-inc/tiem/message"
	"github.com/pingcap-inc/tiem/message/cluster"
	management2 "github.com/pingcap-inc/tiem/micro-api/controller/cluster/management"
	"github.com/pingcap-inc/tiem/micro-cluster/cluster/backuprestore"
	changeFeedManager "github.com/pingcap-inc/tiem/micro-cluster/cluster/changefeed"
	clusterManager "github.com/pingcap-inc/tiem/micro-cluster/cluster/management"
	"github.com/pingcap-inc/tiem/micro-cluster/datatransfer/importexport"

	"github.com/pingcap-inc/tiem/workflow"
	"net/http"
	"strconv"
	"time"

	"github.com/pingcap-inc/tiem/library/util/convert"

	"github.com/pingcap-inc/tiem/library/thirdparty/metrics"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/pingcap-inc/tiem/library/client/cluster/clusterpb"
	"github.com/pingcap-inc/tiem/library/common"
	"github.com/pingcap-inc/tiem/micro-cluster/service/resource"
	"github.com/pingcap-inc/tiem/micro-cluster/service/user/adapt"
	user "github.com/pingcap-inc/tiem/micro-cluster/service/user/application"

	"github.com/pingcap-inc/tiem/library/framework"
	"github.com/pingcap-inc/tiem/micro-cluster/service/cluster/domain"
	userDomain "github.com/pingcap-inc/tiem/micro-cluster/service/user/domain"

	log "github.com/sirupsen/logrus"
)

var TiEMClusterServiceName = "go.micro.tiem.cluster"

var SuccessResponseStatus = &clusterpb.ResponseStatusDTO{Code: 0}
var BizErrorResponseStatus = &clusterpb.ResponseStatusDTO{Code: 500}

type ClusterServiceHandler struct {
	resourceManager   *resource.ResourceManager
	authManager       *user.AuthManager
	tenantManager     *user.TenantManager
	userManager       *user.UserManager
	changeFeedManager *changeFeedManager.Manager
	clusterManager    *clusterManager.Manager
}

func handleRequest(ctx context.Context, req *clusterpb.RpcRequest, resp *clusterpb.RpcResponse, requestBody interface{}) bool {
	err := json.Unmarshal([]byte(req.GetRequest()), &requestBody)
	if err != nil {
		errMsg := fmt.Sprintf("unmarshal request error, request = %s, err = %s", req.GetRequest(), err.Error())
		handleResponse(ctx, resp, framework.NewTiEMErrorf(common.TIEM_UNMARSHAL_ERROR, errMsg), nil, nil)
		return false
	} else {
		return true
	}
}

func handleResponse(ctx context.Context, resp *clusterpb.RpcResponse, err error, responseData interface{}, page *clusterpb.RpcPage) {
	if err == nil {
		data, getDataError := json.Marshal(responseData)
		if getDataError != nil {
			// deal with err uniformly later
			err = framework.WrapError(common.TIEM_MARSHAL_ERROR, fmt.Sprintf("marshal request data error, data = %v", responseData), getDataError)
		} else {
			// handle data and page
			resp.Code = int32(common.TIEM_SUCCESS)
			resp.Response = string(data)
			if page != nil {
				resp.Page = page
			}
			return
		}
	}

	if err != nil {
		if _, ok := err.(framework.TiEMError); !ok {
			err = framework.WrapError(common.TIEM_UNRECOGNIZED_ERROR, "", err)
		}
		framework.LogWithContext(ctx).Error(err.Error())
		resp.Code = int32(err.(framework.TiEMError).GetCode())
		resp.Message = err.(framework.TiEMError).GetMsg()
	}
}

func (handler *ClusterServiceHandler) CreateChangeFeedTask(ctx context.Context, req *clusterpb.RpcRequest, resp *clusterpb.RpcResponse) error {
	request := cluster.CreateChangeFeedTaskReq{}

	if handleRequest(ctx, req, resp, request) {
		result, err := handler.changeFeedManager.Create(ctx, request)
		handleResponse(ctx, resp, err, result, nil)
	}

	return nil
}

func (handler *ClusterServiceHandler) PauseChangeFeedTask(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	panic("implement me")
}

func (handler *ClusterServiceHandler) ResumeChangeFeedTask(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	panic("implement me")
}

func (handler *ClusterServiceHandler) DeleteChangeFeedTask(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	panic("implement me")
}

func (handler *ClusterServiceHandler) UpdateChangeFeedTask(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	panic("implement me")
}

func (handler *ClusterServiceHandler) QueryChangeFeedTasks(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	panic("implement me")
}

func NewClusterServiceHandler(fw *framework.BaseFramework) *ClusterServiceHandler {
	handler := new(ClusterServiceHandler)
	handler.SetResourceManager(resource.NewResourceManager())
	handler.userManager = user.NewUserManager(adapt.MicroMetaDbRepo{})
	handler.tenantManager = user.NewTenantManager(adapt.MicroMetaDbRepo{})
	handler.authManager = user.NewAuthManager(handler.userManager, adapt.MicroMetaDbRepo{})
	handler.changeFeedManager = changeFeedManager.NewManager()
	handler.clusterManager = clusterManager.NewClusterManager()

	domain.InitFlowMap()
	return handler
}

func (handler *ClusterServiceHandler) SetResourceManager(resourceManager *resource.ResourceManager) {
	handler.resourceManager = resourceManager
}

func (handler *ClusterServiceHandler) ResourceManager() *resource.ResourceManager {
	return handler.resourceManager
}

func getLoggerWithContext(ctx context.Context) *log.Entry {
	return framework.LogWithContext(ctx)
}

func handleMetrics(start time.Time, funcName string, code int) {
	duration := time.Since(start)
	framework.Current.GetMetrics().MicroDurationHistogramMetric.With(prometheus.Labels{
		metrics.ServiceLabel: framework.Current.GetServiceMeta().ServiceName.ServerName(),
		metrics.MethodLabel:  funcName,
		metrics.CodeLabel:    strconv.Itoa(code)}).
		Observe(duration.Seconds())
	framework.Current.GetMetrics().MicroRequestsCounterMetric.With(prometheus.Labels{
		metrics.ServiceLabel: framework.Current.GetServiceMeta().ServiceName.ServerName(),
		metrics.MethodLabel:  funcName,
		metrics.CodeLabel:    strconv.Itoa(code)}).
		Inc()
}

func (c ClusterServiceHandler) CreateCluster(ctx context.Context, req *clusterpb.RpcRequest, resp *clusterpb.RpcResponse) (err error) {
	start := time.Now()
	defer handleMetrics(start, "CreateCluster", int(resp.GetCode()))

	request := cluster.CreateClusterReq{}

	if handleRequest(ctx, req, resp, request) {
		result, err := c.clusterManager.CreateCluster(ctx, request)

		handleResponse(ctx, resp, err, result, nil)
	}

	return nil
}

func (handler *ClusterServiceHandler) ScaleOutCluster(ctx context.Context, req *clusterpb.RpcRequest, resp *clusterpb.RpcResponse) error {
	start := time.Now()
	defer handleMetrics(start, "ScaleOutCluster", int(resp.GetCode()))

	request := cluster.ScaleOutClusterReq{}

	if handleRequest(ctx, req, resp, request) {
		result, err := handler.clusterManager.ScaleOut(ctx, request)

		handleResponse(ctx, resp, err, result, nil)
	}

	return nil
}

func (handler *ClusterServiceHandler) ScaleInCluster(ctx context.Context, req *clusterpb.RpcRequest, resp *clusterpb.RpcResponse) error {
	start := time.Now()
	defer handleMetrics(start, "ScaleInCluster", int(resp.GetCode()))

	request := cluster.ScaleInClusterReq{}

	if handleRequest(ctx, req, resp, request) {
		result, err := handler.clusterManager.ScaleIn(ctx, request)

		handleResponse(ctx, resp, err, result, nil)
	}

	return nil
}

func (handler *ClusterServiceHandler) CloneCluster(ctx context.Context, req *clusterpb.RpcRequest, resp *clusterpb.RpcResponse) error {
	start := time.Now()
	defer handleMetrics(start, "CloneCluster", int(resp.GetCode()))

	request := cluster.CloneClusterReq{}

	if handleRequest(ctx, req, resp, request) {
		result, err := handler.clusterManager.Clone(ctx, &request)

		handleResponse(ctx, resp, err, *result, nil)
	}

	return nil
}

func (c ClusterServiceHandler) TakeoverClusters(ctx context.Context, req *clusterpb.ClusterTakeoverReqDTO, resp *clusterpb.ClusterTakeoverRespDTO) (err error) {
	framework.LogWithContext(ctx).Info("takeover clusters")
	clusters, err := domain.TakeoverClusters(ctx, req.Operator, req)
	if err != nil {
		framework.LogWithContext(ctx).Info(err)
		return nil
	} else {
		resp.RespStatus = SuccessResponseStatus
		resp.Clusters = make([]*clusterpb.ClusterDisplayDTO, len(clusters))

		return nil
	}
}

func (c ClusterServiceHandler) QueryCluster(ctx context.Context, req *clusterpb.RpcRequest, resp *clusterpb.RpcResponse) (err error) {
	framework.LogWithContext(ctx).Info("query cluster")
	request := &management2.QueryReq{}
	err = json.Unmarshal([]byte(req.Request), request)
	if err != nil {
		resp.Code = int32(common.TIEM_PARAMETER_INVALID)
		resp.Message = err.Error()
		return
	}
	clusters, total, err := domain.ListCluster(ctx, request)
	if err != nil {
		resp.Code = int32(err.(framework.TiEMError).GetCode())
		resp.Message = err.(framework.TiEMError).GetMsg()
	} else {
		response := make([]management2.ClusterDisplayInfo, 0)

		for _, cluster := range clusters {
			response = append(response, cluster.ExtractDisplayInfo())
		}
		body, err := json.Marshal(response)
		if err != nil {
			resp.Code = int32(common.TIEM_PARAMETER_INVALID)
			resp.Message = err.Error()
		} else {
			resp.Code = int32(common.TIEM_SUCCESS)
			resp.Response = string(body)
			resp.Page = &clusterpb.RpcPage{
				Page:     int32(request.Page),
				PageSize: int32(request.PageSize),
				Total:    int32(total),
			}
		}
	}
	return
}

func (c ClusterServiceHandler) DeleteCluster(ctx context.Context, req *clusterpb.RpcRequest, resp *clusterpb.RpcResponse) (err error) {
	start := time.Now()
	defer handleMetrics(start, "DeleteCluster", int(resp.GetCode()))

	request := cluster.DeleteClusterReq{}

	if handleRequest(ctx, req, resp, request) {
		result, err := c.clusterManager.DeleteCluster(ctx, request)
		handleResponse(ctx, resp, err, result, nil)
	}

	return nil
}

func (c ClusterServiceHandler) RestartCluster(ctx context.Context, req *clusterpb.RpcRequest, resp *clusterpb.RpcResponse) (err error) {
	start := time.Now()
	defer handleMetrics(start, "RestartCluster", int(resp.GetCode()))

	request := cluster.RestartClusterReq{}

	if handleRequest(ctx, req, resp, request) {
		result, err := c.clusterManager.RestartCluster(ctx, request)
		handleResponse(ctx, resp, err, result, nil)
	}

	return nil
}

func (c ClusterServiceHandler) StopCluster(ctx context.Context, req *clusterpb.RpcRequest, resp *clusterpb.RpcResponse) (err error) {
	start := time.Now()
	defer handleMetrics(start, "StopCluster", int(resp.GetCode()))

	request := cluster.StopClusterReq{}

	if handleRequest(ctx, req, resp, request) {
		result, err := c.clusterManager.StopCluster(ctx, request)
		handleResponse(ctx, resp, err, result, nil)
	}

	return nil
}

func (c ClusterServiceHandler) DetailCluster(ctx context.Context, req *clusterpb.RpcRequest, resp *clusterpb.RpcResponse) (err error) {
	framework.LogWithContext(ctx).Info("detail cluster")

	request := &management2.DetailReq{}
	err = json.Unmarshal([]byte(req.Request), request)
	if err != nil {
		resp.Code = int32(common.TIEM_PARAMETER_INVALID)
		resp.Message = err.Error()
		return
	}

	cluster, err := domain.GetClusterDetail(ctx, request.ClusterID)

	if err != nil {
		resp.Code = int32(err.(framework.TiEMError).GetCode())
		resp.Message = err.(framework.TiEMError).GetMsg()
	} else {
		response, err := domain.ExtractClusterInfo(cluster)
		if err != nil {
			resp.Code = int32(common.TIEM_PARAMETER_INVALID)
			resp.Message = err.Error()
		} else {
			resp.Code = int32(common.TIEM_SUCCESS)
			resp.Response = response
		}
	}
	return
}

func (c ClusterServiceHandler) ExportData(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	start := time.Now()
	defer handleMetrics(start, "ExportData", int(response.GetCode()))
	framework.LogWithContext(ctx).Info("export data")
	exportReq := message.DataExportReq{}

	if handleRequest(ctx, request, response, exportReq) {
		manager := importexport.GetImportExportService()
		result, err := manager.ExportData(ctx, &exportReq)
		handleResponse(ctx, response, err, *result, nil)
	}

	return nil
}

func (c ClusterServiceHandler) ImportData(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	start := time.Now()
	defer handleMetrics(start, "ImportData", int(response.GetCode()))
	framework.LogWithContext(ctx).Info("import data")
	importReq := message.DataImportReq{}

	if handleRequest(ctx, request, response, importReq) {
		manager := importexport.GetImportExportService()
		result, err := manager.ImportData(ctx, &importReq)
		handleResponse(ctx, response, err, *result, nil)
	}

	return nil
}

func (c ClusterServiceHandler) QueryDataTransport(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	start := time.Now()
	defer handleMetrics(start, "QueryDataTransport", int(response.GetCode()))
	framework.LogWithContext(ctx).Info("query data transport")
	queryReq := message.QueryDataImportExportRecordsReq{}

	if handleRequest(ctx, request, response, queryReq) {
		manager := importexport.GetImportExportService()
		result, page, err := manager.QueryDataTransportRecords(ctx, &queryReq)
		handleResponse(ctx, response, err, *result, &clusterpb.RpcPage{
			Page:     int32(page.Page),
			PageSize: int32(page.PageSize),
			Total:    int32(page.Total),
		})
	}

	return nil
}

func (c ClusterServiceHandler) DeleteDataTransportRecord(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	start := time.Now()
	defer handleMetrics(start, "DeleteDataTransportRecord", int(response.GetCode()))
	framework.LogWithContext(ctx).Info("delete data transport record")
	deleteReq := message.DeleteImportExportRecordReq{}

	if handleRequest(ctx, request, response, deleteReq) {
		manager := importexport.GetImportExportService()
		result, err := manager.DeleteDataTransportRecord(ctx, &deleteReq)
		handleResponse(ctx, response, err, *result, nil)
	}

	return nil
}

func (c ClusterServiceHandler) CreateBackup(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	start := time.Now()
	defer handleMetrics(start, "CreateBackup", int(response.GetCode()))
	framework.LogWithContext(ctx).Info("create backup")
	backupReq := cluster.BackupClusterDataReq{}

	if handleRequest(ctx, request, response, backupReq) {
		manager := backuprestore.GetBRService()
		result, err := manager.BackupCluster(ctx, &backupReq)
		handleResponse(ctx, response, err, *result, nil)
	}

	return nil
}

func (c ClusterServiceHandler) RecoverCluster(ctx context.Context, req *clusterpb.RecoverRequest, resp *clusterpb.RecoverResponse) (err error) {
	start := time.Now()
	defer handleMetrics(start, "RecoverCluster", int(resp.GetRespStatus().GetCode()))
	if err = domain.RecoverPreCheck(ctx, req); err != nil {
		getLoggerWithContext(ctx).Errorf("recover cluster pre check failed, %s", err.Error())
		resp.RespStatus = &clusterpb.ResponseStatusDTO{Code: int32(common.TIEM_RECOVER_PARAM_INVALID), Message: err.Error()}
		return nil
	}

	clusterAggregation, err := domain.Recover(ctx, req.GetOperator(), req.GetCluster(), req.GetCommonDemand(), req.GetDemands())
	if err != nil {
		getLoggerWithContext(ctx).Error(err)
		resp.RespStatus = &clusterpb.ResponseStatusDTO{Code: int32(common.TIEM_RECOVER_PROCESS_FAILED), Message: common.TIEM_RECOVER_PROCESS_FAILED.Explain()}
	} else {
		resp.RespStatus = SuccessResponseStatus
		resp.ClusterId = clusterAggregation.Cluster.Id
		resp.BaseInfo = clusterAggregation.ExtractBaseInfoDTO()
		resp.ClusterStatus = clusterAggregation.ExtractStatusDTO()
	}
	return nil
}

func (c ClusterServiceHandler) DeleteBackupRecords(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	start := time.Now()
	defer handleMetrics(start, "DeleteBackupRecord", int(response.GetCode()))
	framework.LogWithContext(ctx).Info("delete backup records")
	deleteReq := cluster.DeleteBackupDataReq{}

	if handleRequest(ctx, request, response, deleteReq) {
		manager := backuprestore.GetBRService()
		result, err := manager.DeleteBackupRecords(ctx, &deleteReq)
		handleResponse(ctx, response, err, *result, nil)
	}

	return nil
}

func (c ClusterServiceHandler) SaveBackupStrategy(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	start := time.Now()
	defer handleMetrics(start, "SaveBackupStrategy", int(response.GetCode()))
	framework.LogWithContext(ctx).Info("save backup strategy")
	saveReq := cluster.SaveBackupStrategyReq{}

	if handleRequest(ctx, request, response, saveReq) {
		manager := backuprestore.GetBRService()
		result, err := manager.SaveBackupStrategy(ctx, &saveReq)
		handleResponse(ctx, response, err, *result, nil)
	}

	return nil
}

func (c ClusterServiceHandler) GetBackupStrategy(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	start := time.Now()
	defer handleMetrics(start, "GetBackupStrategy", int(response.GetCode()))
	framework.LogWithContext(ctx).Info("get backup strategy")
	getReq := cluster.GetBackupStrategyReq{}

	if handleRequest(ctx, request, response, getReq) {
		manager := backuprestore.GetBRService()
		result, err := manager.GetBackupStrategy(ctx, &getReq)
		handleResponse(ctx, response, err, *result, nil)
	}

	return nil
}

func (c ClusterServiceHandler) QueryBackupRecords(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) (err error) {
	start := time.Now()
	defer handleMetrics(start, "QueryBackupRecords", int(response.GetCode()))
	framework.LogWithContext(ctx).Info("query backup records")
	queryReq := cluster.QueryBackupRecordsReq{}

	if handleRequest(ctx, request, response, queryReq) {
		manager := backuprestore.GetBRService()
		result, page, err := manager.QueryClusterBackupRecords(ctx, &queryReq)
		handleResponse(ctx, response, err, *result, &clusterpb.RpcPage{
			Page:     int32(page.Page),
			PageSize: int32(page.PageSize),
			Total:    int32(page.Total),
		})
	}

	return nil
}

func (c ClusterServiceHandler) QueryParameters(ctx context.Context, request *clusterpb.QueryClusterParametersRequest, response *clusterpb.QueryClusterParametersResponse) (err error) {

	content, err := domain.GetParameters(ctx, request.Operator, request.ClusterId)

	if err != nil {
		framework.LogWithContext(ctx).Info(err)
		return nil
	} else {
		response.Status = SuccessResponseStatus

		response.ClusterId = request.ClusterId
		response.ParametersJson = content
		return nil
	}
}

func (c ClusterServiceHandler) SaveParameters(ctx context.Context, request *clusterpb.SaveClusterParametersRequest, response *clusterpb.SaveClusterParametersResponse) (err error) {

	//clusterAggregation, err := domain.ModifyParameters(ctx, request.Operator, request.ClusterId, request.ParametersJson)
	//
	//if err != nil {
	//	framework.LogWithContext(ctx).Info(err)
	//	return nil
	//} else {
	//	response.ChangeFeedStatus = SuccessResponseStatus
	//	response.DisplayInfo = &clusterpb.DisplayStatusDTO{
	//		InProcessFlowId: int32(clusterAggregation.CurrentWorkFlow.Id),
	//	}
	//	return nil
	//}
	return err
}

func (c ClusterServiceHandler) DescribeDashboard(ctx context.Context, request *clusterpb.DescribeDashboardRequest, response *clusterpb.DescribeDashboardResponse) (err error) {
	start := time.Now()
	defer handleMetrics(start, "DescribeDashboard", int(response.GetStatus().GetCode()))
	info, err := domain.DescribeDashboard(ctx, request.Operator, request.ClusterId)
	if err != nil {
		getLoggerWithContext(ctx).Error(err)
		response.Status = &clusterpb.ResponseStatusDTO{Code: int32(common.TIEM_DASHBOARD_NOT_FOUND), Message: common.TIEM_DASHBOARD_NOT_FOUND.Explain()}
	} else {
		response.Status = SuccessResponseStatus
		response.ClusterId = info.ClusterId
		response.Url = info.Url
		response.Token = info.Token
	}

	return nil
}

func (c ClusterServiceHandler) DescribeMonitor(ctx context.Context, request *clusterpb.DescribeMonitorRequest, response *clusterpb.DescribeMonitorResponse) (err error) {
	start := time.Now()
	defer handleMetrics(start, "DescribeMonitor", int(response.GetStatus().GetCode()))
	monitor, err := domain.DescribeMonitor(ctx, request.Operator, request.ClusterId)
	if err != nil {
		getLoggerWithContext(ctx).Error(err)
		response.Status = &clusterpb.ResponseStatusDTO{Code: int32(common.TIEM_MONITOR_NOT_FOUND), Message: common.TIEM_MONITOR_NOT_FOUND.Explain()}
	} else {
		response.Status = SuccessResponseStatus
		response.ClusterId = monitor.ClusterId
		response.AlertUrl = monitor.AlertUrl
		response.GrafanaUrl = monitor.GrafanaUrl
	}

	return nil
}

func (c ClusterServiceHandler) ListFlows(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	framework.LogWithContext(ctx).Info("list flows")
	reqData := request.GetRequest()

	listReq := &message.QueryWorkFlowsReq{}
	err := json.Unmarshal([]byte(reqData), listReq)
	if err != nil {
		framework.LogWithContext(ctx).Errorf("json unmarshal reuqest failed %s", err.Error())
		handleResponse(ctx, response, framework.NewTiEMError(common.TIEM_PARAMETER_INVALID, err.Error()), nil, nil)
		return nil
	}

	manager := workflow.GetWorkFlowService()
	flows, total, err := manager.ListWorkFlows(ctx, listReq.BizID, listReq.FlowName, listReq.Status, listReq.Page, listReq.PageSize)
	if err != nil {
		framework.LogWithContext(ctx).Errorf("call workflow manager list flows failed %s", err.Error())
		handleResponse(ctx, response, framework.NewTiEMError(common.TIEM_LIST_WORKFLOW_FAILED, err.Error()), nil, nil)
		return nil
	}

	listResp := message.QueryWorkFlowsResp{
		WorkFlows: flows,
	}
	data, err := json.Marshal(listResp)
	if err != nil {
		framework.LogWithContext(ctx).Errorf("json marshal response failed %s", err.Error())
		handleResponse(ctx, response, framework.NewTiEMError(common.TIEM_LIST_WORKFLOW_FAILED, err.Error()), nil, nil)
	} else {
		response.Code = int32(common.TIEM_SUCCESS)
		response.Response = string(data)
		response.Page = &clusterpb.RpcPage{
			Page:     int32(listReq.Page),
			PageSize: int32(listReq.PageSize),
			Total:    int32(total),
		}
	}

	return nil
}

func (c *ClusterServiceHandler) DetailFlow(ctx context.Context, request *clusterpb.RpcRequest, response *clusterpb.RpcResponse) error {
	framework.LogWithContext(ctx).Info("detail flow")
	reqData := request.GetRequest()

	detailReq := &message.QueryWorkFlowDetailReq{}
	err := json.Unmarshal([]byte(reqData), detailReq)
	if err != nil {
		handleResponse(ctx, response, framework.SimpleError(common.TIEM_PARAMETER_INVALID), nil, nil)
		return nil
	}

	manager := workflow.GetWorkFlowService()
	flowDetail, err := manager.DetailWorkFlow(ctx, detailReq.WorkFlowID)
	if err != nil {
		handleResponse(ctx, response, framework.NewTiEMError(common.TIEM_DETAIL_WORKFLOW_FAILED, err.Error()), nil, nil)
		return nil
	}

	detailResp := message.QueryWorkFlowDetailResp{
		Info:      flowDetail.Flow,
		NodeInfo:  flowDetail.Nodes,
		NodeNames: flowDetail.NodeNames,
	}

	data, err := json.Marshal(detailResp)
	if err != nil {
		handleResponse(ctx, response, framework.NewTiEMError(common.TIEM_DETAIL_WORKFLOW_FAILED, err.Error()), nil, nil)
	} else {
		response.Code = int32(common.TIEM_SUCCESS)
		response.Response = string(data)
	}

	return nil
}

var ManageSuccessResponseStatus = &clusterpb.ManagerResponseStatus{
	Code: 0,
}

func (p *ClusterServiceHandler) Login(ctx context.Context, req *clusterpb.LoginRequest, resp *clusterpb.LoginResponse) error {
	log := framework.LogWithContext(ctx).WithField("fp", "ClusterServiceHandler.Login")
	log.Debug("req:", req)
	token, err := p.authManager.Login(ctx, req.GetAccountName(), req.GetPassword())

	if err != nil {
		resp.Status = &clusterpb.ManagerResponseStatus{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		resp.Status.Message = err.Error()
		log.Error("resp:", resp)
	} else {
		resp.Status = ManageSuccessResponseStatus
		resp.TokenString = token
		log.Debug("resp:", resp)
	}
	return nil

}

func (p *ClusterServiceHandler) Logout(ctx context.Context, req *clusterpb.LogoutRequest, resp *clusterpb.LogoutResponse) error {
	accountName, err := p.authManager.Logout(ctx, req.TokenString)
	if err != nil {
		resp.Status = &clusterpb.ManagerResponseStatus{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		resp.Status.Message = err.Error()
	} else {
		resp.Status = ManageSuccessResponseStatus
		resp.AccountName = accountName
	}
	return nil

}

func (p *ClusterServiceHandler) VerifyIdentity(ctx context.Context, req *clusterpb.VerifyIdentityRequest, resp *clusterpb.VerifyIdentityResponse) error {
	tenantId, accountId, accountName, err := p.authManager.Accessible(ctx, req.GetAuthType(), req.GetPath(), req.GetTokenString())

	if err != nil {
		if _, ok := err.(*userDomain.UnauthorizedError); ok {
			resp.Status = &clusterpb.ManagerResponseStatus{
				Code:    http.StatusUnauthorized,
				Message: "未登录或登录失效，请重试",
			}
		} else if _, ok := err.(*userDomain.ForbiddenError); ok {
			resp.Status = &clusterpb.ManagerResponseStatus{
				Code:    http.StatusForbidden,
				Message: "无权限",
			}
		} else {
			resp.Status = &clusterpb.ManagerResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
	} else {
		resp.Status = ManageSuccessResponseStatus
		resp.TenantId = tenantId
		resp.AccountId = accountId
		resp.AccountName = accountName
	}

	return nil
}

func (c *ClusterServiceHandler) CreateParamGroup(ctx context.Context, req *clusterpb.CreateParamGroupRequest, resp *clusterpb.CreateParamGroupResponse) error {
	start := time.Now()
	defer handleMetrics(start, "CreateParamGroup", int(resp.GetRespStatus().GetCode()))
	return domain.CreateParamGroup(ctx, req, resp)
}

func (c *ClusterServiceHandler) UpdateParamGroup(ctx context.Context, req *clusterpb.UpdateParamGroupRequest, resp *clusterpb.UpdateParamGroupResponse) error {
	start := time.Now()
	defer handleMetrics(start, "UpdateParamGroup", int(resp.GetRespStatus().GetCode()))
	return domain.UpdateParamGroup(ctx, req, resp)
}

func (c *ClusterServiceHandler) DeleteParamGroup(ctx context.Context, req *clusterpb.DeleteParamGroupRequest, resp *clusterpb.DeleteParamGroupResponse) error {
	start := time.Now()
	defer handleMetrics(start, "DeleteParamGroup", int(resp.GetRespStatus().GetCode()))
	return domain.DeleteParamGroup(ctx, req, resp)
}

func (c *ClusterServiceHandler) ListParamGroup(ctx context.Context, req *clusterpb.ListParamGroupRequest, resp *clusterpb.ListParamGroupResponse) error {
	start := time.Now()
	defer handleMetrics(start, "ListParamGroup", int(resp.GetRespStatus().GetCode()))
	return domain.ListParamGroup(ctx, req, resp)
}

func (c *ClusterServiceHandler) DetailParamGroup(ctx context.Context, req *clusterpb.DetailParamGroupRequest, resp *clusterpb.DetailParamGroupResponse) error {
	start := time.Now()
	defer handleMetrics(start, "DetailParamGroup", int(resp.GetRespStatus().GetCode()))
	return domain.DetailParamGroup(ctx, req, resp)
}

func (c *ClusterServiceHandler) ApplyParamGroup(ctx context.Context, req *clusterpb.ApplyParamGroupRequest, resp *clusterpb.ApplyParamGroupResponse) error {
	start := time.Now()
	defer handleMetrics(start, "ApplyParamGroup", int(resp.GetRespStatus().GetCode()))

	// query params list by param group id
	pgDetail := &clusterpb.DetailParamGroupResponse{}
	err := domain.DetailParamGroup(ctx, &clusterpb.DetailParamGroupRequest{ParamGroupId: req.ParamGroupId}, pgDetail)
	if err != nil {
		framework.LogWithContext(ctx).Errorf("apply param group err: %v", err)
		return err
	}

	params := make([]*domain.ApplyParam, len(pgDetail.ParamGroup.Params))
	err = convert.ConvertObj(pgDetail.ParamGroup.Params, &params)
	if err != nil {
		framework.LogWithContext(ctx).Errorf("apply param group convert obj err: %v", err)
		resp.RespStatus = &clusterpb.ResponseStatusDTO{Code: int32(common.TIEM_CONVERT_OBJ_FAILED), Message: common.TIEM_CONVERT_OBJ_FAILED.Explain() + err.Error()}
		return err
	}
	// Convert the default value of the parameter group to the real value of the modified parameter
	for i, param := range pgDetail.ParamGroup.Params {
		params[i].RealValue.Cluster = param.DefaultValue
	}
	modifyParam := &domain.ModifyParam{NeedReboot: req.NeedReboot, Params: params}
	clusterAggregation, err := domain.ModifyParameters(ctx, req.Operator, req.ClusterId, modifyParam)
	if err != nil {
		framework.LogWithContext(ctx).Errorf("apply param group modify parameters err: %v", err)
		resp.RespStatus = &clusterpb.ResponseStatusDTO{Code: int32(common.TIEM_MODIFY_PARAM_FAILED), Message: common.TIEM_MODIFY_PARAM_FAILED.Explain() + err.Error()}
		return nil
	} else {
		resp.DisplayInfo = &clusterpb.DisplayStatusDTO{
			StatusCode:      strconv.Itoa(int(clusterAggregation.CurrentWorkFlow.Status)),
			StatusName:      clusterAggregation.CurrentWorkFlow.Status.Display(),
			InProcessFlowId: int32(clusterAggregation.CurrentWorkFlow.Id),
			CreateTime:      clusterAggregation.CurrentWorkFlow.CreateTime.Unix(),
			UpdateTime:      clusterAggregation.CurrentWorkFlow.UpdateTime.Unix(),
		}
		return domain.ApplyParamGroup(ctx, req, resp)
	}
}

func (c *ClusterServiceHandler) CopyParamGroup(ctx context.Context, req *clusterpb.CopyParamGroupRequest, resp *clusterpb.CopyParamGroupResponse) error {
	start := time.Now()
	defer handleMetrics(start, "CopyParamGroup", int(resp.GetRespStatus().GetCode()))
	return domain.CopyParamGroup(ctx, req, resp)
}

func (c *ClusterServiceHandler) ListClusterParams(ctx context.Context, req *clusterpb.ListClusterParamsRequest, resp *clusterpb.ListClusterParamsResponse) error {
	start := time.Now()
	defer handleMetrics(start, "ListClusterParams", int(resp.GetRespStatus().GetCode()))
	return domain.ListClusterParams(ctx, req, resp)
}

func (c *ClusterServiceHandler) UpdateClusterParams(ctx context.Context, req *clusterpb.UpdateClusterParamsRequest, resp *clusterpb.UpdateClusterParamsResponse) error {
	start := time.Now()
	defer handleMetrics(start, "UpdateClusterParams", int(resp.GetRespStatus().GetCode()))

	params := make([]*domain.ApplyParam, len(req.Params))
	err := convert.ConvertObj(req.Params, &params)
	if err != nil {
		framework.LogWithContext(ctx).Errorf("update cluster params convert obj err: %v", err)
		resp.RespStatus = &clusterpb.ResponseStatusDTO{Code: int32(common.TIEM_CONVERT_OBJ_FAILED), Message: common.TIEM_CONVERT_OBJ_FAILED.Explain() + err.Error()}
		return err
	}

	modifyParam := &domain.ModifyParam{NeedReboot: req.NeedReboot, Params: params}
	clusterAggregation, err := domain.ModifyParameters(ctx, req.Operator, req.ClusterId, modifyParam)
	if err != nil {
		framework.LogWithContext(ctx).Errorf("update cluster params modify parameters err: %v", err)
		resp.RespStatus = &clusterpb.ResponseStatusDTO{Code: int32(common.TIEM_MODIFY_PARAM_FAILED), Message: common.TIEM_MODIFY_PARAM_FAILED.Explain() + err.Error()}
		return nil
	} else {
		resp.DisplayInfo = &clusterpb.DisplayStatusDTO{
			StatusCode:      strconv.Itoa(int(clusterAggregation.CurrentWorkFlow.Status)),
			StatusName:      clusterAggregation.CurrentWorkFlow.Status.Display(),
			InProcessFlowId: int32(clusterAggregation.CurrentWorkFlow.Id),
			CreateTime:      clusterAggregation.CurrentWorkFlow.CreateTime.Unix(),
			UpdateTime:      clusterAggregation.CurrentWorkFlow.UpdateTime.Unix(),
		}
		return domain.UpdateClusterParams(ctx, req, resp)
	}
}

func (c *ClusterServiceHandler) InspectClusterParams(ctx context.Context, req *clusterpb.InspectClusterParamsRequest, resp *clusterpb.InspectClusterParamsResponse) error {
	start := time.Now()
	defer handleMetrics(start, "InspectClusterParams", int(resp.GetRespStatus().GetCode()))
	return domain.InspectClusterParams(ctx, req, resp)
}

func (clusterManager *ClusterServiceHandler) ImportHost(ctx context.Context, in *clusterpb.ImportHostRequest, out *clusterpb.ImportHostResponse) error {
	return clusterManager.resourceManager.ImportHost(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) ImportHostsInBatch(ctx context.Context, in *clusterpb.ImportHostsInBatchRequest, out *clusterpb.ImportHostsInBatchResponse) error {
	return clusterManager.resourceManager.ImportHostsInBatch(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) RemoveHost(ctx context.Context, in *clusterpb.RemoveHostRequest, out *clusterpb.RemoveHostResponse) error {
	return clusterManager.resourceManager.RemoveHost(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) RemoveHostsInBatch(ctx context.Context, in *clusterpb.RemoveHostsInBatchRequest, out *clusterpb.RemoveHostsInBatchResponse) error {
	return clusterManager.resourceManager.RemoveHostsInBatch(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) ListHost(ctx context.Context, in *clusterpb.ListHostsRequest, out *clusterpb.ListHostsResponse) error {
	return clusterManager.resourceManager.ListHost(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) CheckDetails(ctx context.Context, in *clusterpb.CheckDetailsRequest, out *clusterpb.CheckDetailsResponse) error {
	return clusterManager.resourceManager.CheckDetails(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) AllocHosts(ctx context.Context, in *clusterpb.AllocHostsRequest, out *clusterpb.AllocHostResponse) error {
	return clusterManager.resourceManager.AllocHosts(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) GetFailureDomain(ctx context.Context, in *clusterpb.GetFailureDomainRequest, out *clusterpb.GetFailureDomainResponse) error {
	return clusterManager.resourceManager.GetFailureDomain(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) AllocResourcesInBatch(ctx context.Context, in *clusterpb.BatchAllocRequest, out *clusterpb.BatchAllocResponse) error {
	return clusterManager.resourceManager.AllocResourcesInBatch(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) RecycleResources(ctx context.Context, in *clusterpb.RecycleRequest, out *clusterpb.RecycleResponse) error {
	return clusterManager.resourceManager.RecycleResources(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) UpdateHostStatus(ctx context.Context, in *clusterpb.UpdateHostStatusRequest, out *clusterpb.UpdateHostStatusResponse) error {
	return clusterManager.resourceManager.UpdateHostStatus(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) ReserveHost(ctx context.Context, in *clusterpb.ReserveHostRequest, out *clusterpb.ReserveHostResponse) error {
	return clusterManager.resourceManager.ReserveHost(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) GetHierarchy(ctx context.Context, in *clusterpb.GetHierarchyRequest, out *clusterpb.GetHierarchyResponse) error {
	return clusterManager.resourceManager.GetHierarchy(ctx, in, out)
}

func (clusterManager *ClusterServiceHandler) GetStocks(ctx context.Context, in *clusterpb.GetStocksRequest, out *clusterpb.GetStocksResponse) error {
	return clusterManager.resourceManager.GetStocks(ctx, in, out)
}
