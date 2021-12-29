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
 ******************************************************************************/

package management

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pingcap-inc/tiem/common/constants"
	"github.com/pingcap-inc/tiem/common/structs"
	"github.com/pingcap-inc/tiem/library/secondparty"
	"github.com/pingcap-inc/tiem/message/cluster"
	"github.com/pingcap-inc/tiem/micro-cluster/cluster/management/handler"
	"github.com/pingcap-inc/tiem/models"
	"github.com/pingcap-inc/tiem/models/cluster/management"
	"github.com/pingcap-inc/tiem/models/common"
	wfModel "github.com/pingcap-inc/tiem/models/workflow"
	"github.com/pingcap-inc/tiem/test/mockmodels/mockclustermanagement"
	mock_secondparty_v2 "github.com/pingcap-inc/tiem/test/mocksecondparty_v2"
	mock_workflow_service "github.com/pingcap-inc/tiem/test/mockworkflow"
	"github.com/pingcap-inc/tiem/workflow"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var emptyNode = func(task *wfModel.WorkFlowNode, context *workflow.FlowContext) error {
	return nil
}

func getEmptyFlow(name string) *workflow.WorkFlowDefine {
	return &workflow.WorkFlowDefine{
		FlowName: name,
		TaskNodes: map[string]*workflow.NodeDefine{
			"start": {"start", "done", "fail", workflow.SyncFuncNode, emptyNode},
			"done":  {"end", "", "", workflow.SyncFuncNode, emptyNode},
			"fail":  {"fail", "", "", workflow.SyncFuncNode, emptyNode},
		},
	}
}

func TestAsyncMaintenance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workflow.GetWorkFlowService().RegisterWorkFlow(context.TODO(), "testFlow", getEmptyFlow("testFlow"))
	clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
	models.SetClusterReaderWriter(clusterRW)
	clusterRW.EXPECT().SetMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	t.Run("normal", func(t *testing.T) {
		workflowService := mock_workflow_service.NewMockWorkFlowService(ctrl)
		workflow.MockWorkFlowService(workflowService)
		defer workflow.MockWorkFlowService(workflow.NewWorkFlowManager())
		workflowService.EXPECT().CreateWorkFlow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&workflow.WorkFlowAggregation{
			Flow:    &wfModel.WorkFlow{Entity: common.Entity{ID: "flow01"}},
			Context: workflow.FlowContext{Context: context.TODO(), FlowData: make(map[string]interface{})},
		}, nil).AnyTimes()
		workflowService.EXPECT().AsyncStart(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		meta := &handler.ClusterMeta{
			Cluster: &management.Cluster{Entity: common.Entity{ID: "cluster01"}},
		}
		data := make(map[string]interface{})
		data["key"] = "value"
		got, err := asyncMaintenance(context.TODO(), meta, constants.ClusterMaintenanceNone, "testFlow", data)
		assert.NoError(t, err)
		assert.Equal(t, got, "flow01")
	})

	t.Run("create workflow fail", func(t *testing.T) {
		clusterRW.EXPECT().ClearMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		workflowService := mock_workflow_service.NewMockWorkFlowService(ctrl)
		workflow.MockWorkFlowService(workflowService)
		defer workflow.MockWorkFlowService(workflow.NewWorkFlowManager())
		workflowService.EXPECT().CreateWorkFlow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&workflow.WorkFlowAggregation{
			Flow:    &wfModel.WorkFlow{Entity: common.Entity{ID: "flow01"}},
			Context: workflow.FlowContext{Context: context.TODO(), FlowData: make(map[string]interface{})},
		}, fmt.Errorf("create workflow error")).AnyTimes()
		workflowService.EXPECT().AsyncStart(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		meta := &handler.ClusterMeta{
			Cluster: &management.Cluster{Entity: common.Entity{ID: "cluster01"}},
		}
		data := make(map[string]interface{})
		data["key"] = "value"
		_, err := asyncMaintenance(context.TODO(), meta, constants.ClusterMaintenanceNone, "testFlow", data)
		assert.Error(t, err)
	})

	t.Run("start workflow fail", func(t *testing.T) {
		clusterRW.EXPECT().ClearMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		workflowService := mock_workflow_service.NewMockWorkFlowService(ctrl)
		workflow.MockWorkFlowService(workflowService)
		defer workflow.MockWorkFlowService(workflow.NewWorkFlowManager())
		workflowService.EXPECT().CreateWorkFlow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&workflow.WorkFlowAggregation{
			Flow:    &wfModel.WorkFlow{Entity: common.Entity{ID: "flow01"}},
			Context: workflow.FlowContext{Context: context.TODO(), FlowData: make(map[string]interface{})},
		}, nil).AnyTimes()
		workflowService.EXPECT().AsyncStart(gomock.Any(), gomock.Any()).Return(fmt.Errorf("start workflow error")).AnyTimes()
		meta := &handler.ClusterMeta{
			Cluster: &management.Cluster{Entity: common.Entity{ID: "cluster01"}},
		}
		data := make(map[string]interface{})
		data["key"] = "value"
		_, err := asyncMaintenance(context.TODO(), meta, constants.ClusterMaintenanceNone, "testFlow", data)
		assert.Error(t, err)
	})
}

func TestManager_ScaleOut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workflow.GetWorkFlowService().RegisterWorkFlow(context.TODO(), constants.FlowScaleOutCluster, getEmptyFlow(constants.FlowScaleOutCluster))
	manager := &Manager{}
	clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
	models.SetClusterReaderWriter(clusterRW)

	clusterRW.EXPECT().GetMeta(gomock.Any(), "111").Return(&management.Cluster{}, []*management.ClusterInstance{
		{},
		{},
	}, nil).AnyTimes()

	clusterRW.EXPECT().SetMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	workflowService := mock_workflow_service.NewMockWorkFlowService(ctrl)
	workflow.MockWorkFlowService(workflowService)
	defer workflow.MockWorkFlowService(workflow.NewWorkFlowManager())
	workflowService.EXPECT().CreateWorkFlow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&workflow.WorkFlowAggregation{
		Flow:    &wfModel.WorkFlow{Entity: common.Entity{ID: "flow01"}},
		Context: workflow.FlowContext{Context: context.TODO(), FlowData: make(map[string]interface{})},
	}, nil).AnyTimes()
	workflowService.EXPECT().AsyncStart(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockTiup := mock_secondparty_v2.NewMockSecondPartyService(ctrl)
	secondparty.Manager = mockTiup
	mockTiup.EXPECT().ClusterComponentCtl(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any()).Return("{\"enable-placement-rules\": \"true\"}", nil)

	t.Run("normal", func(t *testing.T) {
		_, err := manager.ScaleOut(context.TODO(), cluster.ScaleOutClusterReq{
			ClusterID: "111",
			ClusterResourceInfo: structs.ClusterResourceInfo{
				InstanceResource: []structs.ClusterResourceParameterCompute{
					{Type: "TiDB", Count: 1, Resource: []structs.ClusterResourceParameterComputeResource{
						{Zone: "Test_Zone1", DiskType: "SATA", DiskCapacity: 0, Spec: "4C8G", Count: 1},
					}},
					{Type: "TiKV", Count: 1, Resource: []structs.ClusterResourceParameterComputeResource{
						{Zone: "Test_Zone1", DiskType: "SATA", DiskCapacity: 0, Spec: "4C8G", Count: 1},
					}},
					{Type: "PD", Count: 1, Resource: []structs.ClusterResourceParameterComputeResource{
						{Zone: "Test_Zone1", DiskType: "SATA", DiskCapacity: 0, Spec: "4C8G", Count: 1},
					}},
					{Type: "TiFlash", Count: 1, Resource: []structs.ClusterResourceParameterComputeResource{
						{Zone: "Test_Zone1", DiskType: "SATA", DiskCapacity: 0, Spec: "4C8G", Count: 1},
					}},
				},
			},
		})
		assert.NoError(t, err)
	})

	t.Run("not found cluster", func(t *testing.T) {
		clusterRW.EXPECT().GetMeta(gomock.Any(), "112").Return(&management.Cluster{}, []*management.ClusterInstance{
			{},
			{},
		}, fmt.Errorf("not found cluster"))
		_, err := manager.ScaleOut(context.TODO(), cluster.ScaleOutClusterReq{
			ClusterID:           "112",
			ClusterResourceInfo: structs.ClusterResourceInfo{},
		})
		assert.Error(t, err)
	})

	t.Run("check fail", func(t *testing.T) {
		mockTiup.EXPECT().ClusterComponentCtl(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any()).Return("{\"enable-placement-rules\": \"false\"}", nil).AnyTimes()
		_, err := manager.ScaleOut(context.TODO(), cluster.ScaleOutClusterReq{
			ClusterID: "111",
			ClusterResourceInfo: structs.ClusterResourceInfo{
				InstanceResource: []structs.ClusterResourceParameterCompute{
					{Type: "TiFlash", Count: 1, Resource: []structs.ClusterResourceParameterComputeResource{
						{Zone: "Test_Zone1", DiskType: "SATA", DiskCapacity: 0, Spec: "4C8G", Count: 1},
					}},
				},
			},
		})
		assert.Error(t, err)
	})
}

func TestManager_ScaleIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workflow.GetWorkFlowService().RegisterWorkFlow(context.TODO(), constants.FlowScaleInCluster, getEmptyFlow(constants.FlowScaleInCluster))
	manager := &Manager{}
	clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
	models.SetClusterReaderWriter(clusterRW)

	clusterRW.EXPECT().GetMeta(gomock.Any(), "111").Return(&management.Cluster{}, []*management.ClusterInstance{
		{Entity: common.Entity{ID: "instance01"}, Type: "TiDB"},
		{Entity: common.Entity{ID: "instance02"}, Type: "TiFlash"},
	}, nil).AnyTimes()

	clusterRW.EXPECT().SetMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	workflowService := mock_workflow_service.NewMockWorkFlowService(ctrl)
	workflow.MockWorkFlowService(workflowService)
	defer workflow.MockWorkFlowService(workflow.NewWorkFlowManager())
	workflowService.EXPECT().CreateWorkFlow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&workflow.WorkFlowAggregation{
		Flow:    &wfModel.WorkFlow{Entity: common.Entity{ID: "flow01"}},
		Context: workflow.FlowContext{Context: context.TODO(), FlowData: make(map[string]interface{})},
	}, nil).AnyTimes()
	workflowService.EXPECT().AsyncStart(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	t.Run("normal", func(t *testing.T) {
		_, err := manager.ScaleIn(context.TODO(), cluster.ScaleInClusterReq{
			ClusterID:  "111",
			InstanceID: "instance01",
		})
		assert.NoError(t, err)
	})
	t.Run("not found cluster", func(t *testing.T) {
		clusterRW.EXPECT().GetMeta(gomock.Any(), "112").Return(&management.Cluster{}, []*management.ClusterInstance{
			{Entity: common.Entity{ID: "instance01"}, Type: "TiDB"},
			{Entity: common.Entity{ID: "instance02"}, Type: "PD"},
		}, fmt.Errorf("not found cluster")).AnyTimes()
		_, err := manager.ScaleIn(context.TODO(), cluster.ScaleInClusterReq{
			ClusterID:  "112",
			InstanceID: "instance01",
		})
		assert.Error(t, err)
	})
	t.Run("check fail", func(t *testing.T) {
		_, err := manager.ScaleIn(context.TODO(), cluster.ScaleInClusterReq{
			ClusterID:  "111",
			InstanceID: "instance02",
		})
		assert.Error(t, err)
	})
}

func TestManager_Clone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workflow.GetWorkFlowService().RegisterWorkFlow(context.TODO(), constants.FlowScaleInCluster, getEmptyFlow(constants.FlowScaleInCluster))
	manager := &Manager{}
	clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
	models.SetClusterReaderWriter(clusterRW)

	clusterRW.EXPECT().GetMeta(gomock.Any(), "111").Return(&management.Cluster{}, []*management.ClusterInstance{
		{Entity: common.Entity{ID: "instance01"}, Type: "TiDB"},
		{Entity: common.Entity{ID: "instance02"}, Type: "TiFlash"},
	}, nil).AnyTimes()
	clusterRW.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, nil)
	clusterRW.EXPECT().CreateRelation(gomock.Any(), gomock.Any()).Return(nil)

	clusterRW.EXPECT().SetMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	workflowService := mock_workflow_service.NewMockWorkFlowService(ctrl)
	workflow.MockWorkFlowService(workflowService)
	defer workflow.MockWorkFlowService(workflow.NewWorkFlowManager())
	workflowService.EXPECT().CreateWorkFlow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&workflow.WorkFlowAggregation{
		Flow:    &wfModel.WorkFlow{Entity: common.Entity{ID: "flow01"}},
		Context: workflow.FlowContext{Context: context.TODO(), FlowData: make(map[string]interface{})},
	}, nil).AnyTimes()
	workflowService.EXPECT().AsyncStart(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	t.Run("normal", func(t *testing.T) {
		_, err := manager.Clone(context.TODO(), cluster.CloneClusterReq{
			SourceClusterID: "111",
			CreateClusterParameter: structs.CreateClusterParameter{
				Name:       "testCluster",
				DBUser:     "user01",
				DBPassword: "password01",
				Type:       "TiDB",
				Version:    "v5.0.0",
			},
		})
		assert.NoError(t, err)
	})
}

func TestManager_CreateCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workflow.GetWorkFlowService().RegisterWorkFlow(context.TODO(), constants.FlowCreateCluster, getEmptyFlow(constants.FlowCreateCluster))
	manager := Manager{}
	clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
	models.SetClusterReaderWriter(clusterRW)

	clusterRW.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	clusterRW.EXPECT().SetMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	workflowService := mock_workflow_service.NewMockWorkFlowService(ctrl)
	workflow.MockWorkFlowService(workflowService)
	defer workflow.MockWorkFlowService(workflow.NewWorkFlowManager())
	workflowService.EXPECT().CreateWorkFlow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&workflow.WorkFlowAggregation{
		Flow:    &wfModel.WorkFlow{Entity: common.Entity{ID: "flow01"}},
		Context: workflow.FlowContext{Context: context.TODO(), FlowData: make(map[string]interface{})},
	}, nil).AnyTimes()
	workflowService.EXPECT().AsyncStart(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	t.Run("normal", func(t *testing.T) {
		_, err := manager.CreateCluster(context.TODO(), cluster.CreateClusterReq{
			ResourceParameter: structs.ClusterResourceInfo{
				InstanceResource: []structs.ClusterResourceParameterCompute{
					{Type: "TiDB", Count: 1, Resource: []structs.ClusterResourceParameterComputeResource{
						{Zone: "Test_Zone1", DiskType: "SATA", DiskCapacity: 0, Spec: "4C8G", Count: 1},
					}},
					{Type: "TiKV", Count: 1, Resource: []structs.ClusterResourceParameterComputeResource{
						{Zone: "Test_Zone1", DiskType: "SATA", DiskCapacity: 0, Spec: "4C8G", Count: 1},
					}},
					{Type: "PD", Count: 1, Resource: []structs.ClusterResourceParameterComputeResource{
						{Zone: "Test_Zone1", DiskType: "SATA", DiskCapacity: 0, Spec: "4C8G", Count: 1},
					}},
				},
			},
		})
		assert.NoError(t, err)
	})
	t.Run("no computes", func(t *testing.T) {
		_, err := manager.CreateCluster(context.TODO(), cluster.CreateClusterReq{})
		assert.Error(t, err)
	})
}

func TestManager_StopCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workflow.GetWorkFlowService().RegisterWorkFlow(context.TODO(), constants.FlowStopCluster, getEmptyFlow(constants.FlowStopCluster))
	manager := Manager{}
	t.Run("normal", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)

		clusterRW.EXPECT().GetMeta(gomock.Any(), "111").Return(&management.Cluster{}, []*management.ClusterInstance{
			{},
			{},
		}, nil)

		clusterRW.EXPECT().SetMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		workflowService := mock_workflow_service.NewMockWorkFlowService(ctrl)
		workflow.MockWorkFlowService(workflowService)
		defer workflow.MockWorkFlowService(workflow.NewWorkFlowManager())
		workflowService.EXPECT().CreateWorkFlow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&workflow.WorkFlowAggregation{
			Flow:    &wfModel.WorkFlow{Entity: common.Entity{ID: "flow01"}},
			Context: workflow.FlowContext{Context: context.TODO(), FlowData: make(map[string]interface{})},
		}, nil).AnyTimes()
		workflowService.EXPECT().AsyncStart(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		_, err := manager.StopCluster(context.TODO(), cluster.StopClusterReq{
			ClusterID: "111",
		})
		assert.NoError(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)
		clusterRW.EXPECT().GetMeta(gomock.Any(), "111").Return(nil, nil, errors.New(""))

		_, err := manager.StopCluster(context.TODO(), cluster.StopClusterReq{
			ClusterID: "111",
		})
		assert.Error(t, err)
	})
	t.Run("status", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)
		clusterRW.EXPECT().GetMeta(gomock.Any(), "111").Return(&management.Cluster{}, []*management.ClusterInstance{
			{},
			{},
		}, nil)
		clusterRW.EXPECT().SetMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New(""))

		_, err := manager.StopCluster(context.TODO(), cluster.StopClusterReq{
			ClusterID: "111",
		})
		assert.Error(t, err)
	})
}

func TestManager_RestartCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workflow.GetWorkFlowService().RegisterWorkFlow(context.TODO(), constants.FlowRestartCluster, getEmptyFlow(constants.FlowRestartCluster))
	manager := Manager{}
	t.Run("normal", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)

		clusterRW.EXPECT().GetMeta(gomock.Any(), "111").Return(&management.Cluster{}, []*management.ClusterInstance{
			{},
			{},
		}, nil)

		clusterRW.EXPECT().SetMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		workflowService := mock_workflow_service.NewMockWorkFlowService(ctrl)
		workflow.MockWorkFlowService(workflowService)
		defer workflow.MockWorkFlowService(workflow.NewWorkFlowManager())
		workflowService.EXPECT().CreateWorkFlow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&workflow.WorkFlowAggregation{
			Flow:    &wfModel.WorkFlow{Entity: common.Entity{ID: "flow01"}},
			Context: workflow.FlowContext{Context: context.TODO(), FlowData: make(map[string]interface{})},
		}, nil).AnyTimes()
		workflowService.EXPECT().AsyncStart(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		_, err := manager.StopCluster(context.TODO(), cluster.StopClusterReq{
			ClusterID: "111",
		})
		assert.NoError(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)
		clusterRW.EXPECT().GetMeta(gomock.Any(), "111").Return(nil, nil, errors.New(""))

		_, err := manager.StopCluster(context.TODO(), cluster.StopClusterReq{
			ClusterID: "111",
		})
		assert.Error(t, err)
	})
	t.Run("status", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)
		clusterRW.EXPECT().GetMeta(gomock.Any(), "111").Return(&management.Cluster{}, []*management.ClusterInstance{
			{},
			{},
		}, nil)
		clusterRW.EXPECT().SetMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New(""))

		_, err := manager.StopCluster(context.TODO(), cluster.StopClusterReq{
			ClusterID: "111",
		})
		assert.Error(t, err)
	})
}

func TestManager_DeleteCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workflow.GetWorkFlowService().RegisterWorkFlow(context.TODO(), constants.FlowDeleteCluster, getEmptyFlow(constants.FlowDeleteCluster))
	manager := Manager{}
	t.Run("normal", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)

		clusterRW.EXPECT().GetMeta(gomock.Any(), "111").Return(&management.Cluster{}, []*management.ClusterInstance{
			{},
			{},
		}, nil)

		clusterRW.EXPECT().SetMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		workflowService := mock_workflow_service.NewMockWorkFlowService(ctrl)
		workflow.MockWorkFlowService(workflowService)
		defer workflow.MockWorkFlowService(workflow.NewWorkFlowManager())
		workflowService.EXPECT().CreateWorkFlow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&workflow.WorkFlowAggregation{
			Flow:    &wfModel.WorkFlow{Entity: common.Entity{ID: "flow01"}},
			Context: workflow.FlowContext{Context: context.TODO(), FlowData: make(map[string]interface{})},
		}, nil).AnyTimes()
		workflowService.EXPECT().AsyncStart(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		_, err := manager.StopCluster(context.TODO(), cluster.StopClusterReq{
			ClusterID: "111",
		})
		assert.NoError(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)
		clusterRW.EXPECT().GetMeta(gomock.Any(), "111").Return(nil, nil, errors.New(""))

		_, err := manager.StopCluster(context.TODO(), cluster.StopClusterReq{
			ClusterID: "111",
		})
		assert.Error(t, err)
	})
	t.Run("status", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)
		clusterRW.EXPECT().GetMeta(gomock.Any(), "111").Return(&management.Cluster{}, []*management.ClusterInstance{
			{},
			{},
		}, nil)
		clusterRW.EXPECT().SetMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New(""))

		_, err := manager.DeleteCluster(context.TODO(), cluster.DeleteClusterReq{
			ClusterID: "111",
		})
		assert.Error(t, err)
	})
}

func TestManager_DetailCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
	models.SetClusterReaderWriter(clusterRW)

	clusterRW.EXPECT().GetMeta(gomock.Any(), gomock.Any()).Return(&management.Cluster{
		Entity: common.Entity{
			ID:        "2145635758",
			TenantId:  "324567",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:              "koojdafij",
		DBUser:            "kodjsfn",
		DBPassword:        "mypassword",
		Type:              "TiDB",
		Version:           "v5.0.0",
		Tags:              []string{"111", "333"},
		OwnerId:           "436534636u",
		ParameterGroupID:  "352467890",
		Copies:            4,
		Region:            "Region1",
		CpuArchitecture:   "x86_64",
		MaintenanceStatus: constants.ClusterMaintenanceCreating,
	}, []*management.ClusterInstance{
		{
			Entity: common.Entity{
				Status: string(constants.ClusterInstanceRunning),
			},
			Zone:     "zone1",
			CpuCores: 4,
			Memory:   8,
			Type:     "TiDB",
			Version:  "v5.0.0",
			Ports:    []int32{10001, 10002, 10003, 10004},
			HostIP:   []string{"127.0.0.1"},
		},
		{
			Entity: common.Entity{
				Status: string(constants.ClusterInstanceRunning),
			},
			Zone:     "zone1",
			CpuCores: 3,
			Memory:   7,
			Type:     "TiDB",
			Version:  "v5.0.0",
			Ports:    []int32{10001, 10002, 10003, 10004},
			HostIP:   []string{"127.0.0.1"},
		},
		{
			Entity: common.Entity{
				Status: string(constants.ClusterInstanceRunning),
			},
			Zone:     "zone1",
			CpuCores: 4,
			Memory:   8,
			Type:     "TiKV",
			Version:  "v5.0.0",
			Ports:    []int32{20001, 20002, 20003, 20004},
			HostIP:   []string{"127.0.0.2"},
		},
		{
			Entity: common.Entity{
				Status: string(constants.ClusterInstanceRunning),
			},
			Zone:     "zone1",
			CpuCores: 3,
			Memory:   7,
			Type:     "TiKV",
			Version:  "v5.0.0",
			Ports:    []int32{20001, 20002, 20003, 20004},
			HostIP:   []string{"127.0.0.2"},
		},
		{
			Entity: common.Entity{
				Status: string(constants.ClusterInstanceRunning),
			},
			Zone:     "zone1",
			CpuCores: 4,
			Memory:   8,
			Type:     "PD",
			Version:  "v5.0.0",
			Ports:    []int32{30001, 30002, 30003, 30004},
			HostIP:   []string{"127.0.0.3"},
		},
		{
			Entity: common.Entity{
				Status: string(constants.ClusterInstanceRunning),
			},
			Zone:     "zone1",
			CpuCores: 3,
			Memory:   7,
			Type:     "PD",
			Version:  "v5.0.0",
			Ports:    []int32{30001, 30002, 30003, 30004},
			HostIP:   []string{"127.0.0.3"},
		},
	}, nil)
	manager := Manager{}
	got, err := manager.DetailCluster(context.TODO(), cluster.QueryClusterDetailReq{
		ClusterID: "111",
	})
	assert.NoError(t, err)
	assert.Equal(t, got.Info.Region, "Region1")
	assert.Equal(t, len(got.ClusterTopologyInfo.Topology), 6)
}

func TestManager_RestoreNewCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workflow.GetWorkFlowService().RegisterWorkFlow(context.TODO(), constants.FlowRestoreNewCluster, getEmptyFlow(constants.FlowRestoreNewCluster))
	manager := Manager{}
	clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
	models.SetClusterReaderWriter(clusterRW)

	clusterRW.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	clusterRW.EXPECT().SetMaintenanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	workflowService := mock_workflow_service.NewMockWorkFlowService(ctrl)
	workflow.MockWorkFlowService(workflowService)
	defer workflow.MockWorkFlowService(workflow.NewWorkFlowManager())
	workflowService.EXPECT().CreateWorkFlow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&workflow.WorkFlowAggregation{
		Flow:    &wfModel.WorkFlow{Entity: common.Entity{ID: "flow01"}},
		Context: workflow.FlowContext{Context: context.TODO(), FlowData: make(map[string]interface{})},
	}, nil).AnyTimes()
	workflowService.EXPECT().AsyncStart(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	t.Run("normal", func(t *testing.T) {
		_, err := manager.RestoreNewCluster(context.TODO(), cluster.RestoreNewClusterReq{
			ResourceParameter: structs.ClusterResourceInfo{
				InstanceResource: []structs.ClusterResourceParameterCompute{
					{Type: "TiDB", Count: 1, Resource: []structs.ClusterResourceParameterComputeResource{
						{Zone: "Test_Zone1", DiskType: "SATA", DiskCapacity: 0, Spec: "4C8G", Count: 1},
					}},
					{Type: "TiKV", Count: 1, Resource: []structs.ClusterResourceParameterComputeResource{
						{Zone: "Test_Zone1", DiskType: "SATA", DiskCapacity: 0, Spec: "4C8G", Count: 1},
					}},
					{Type: "PD", Count: 1, Resource: []structs.ClusterResourceParameterComputeResource{
						{Zone: "Test_Zone1", DiskType: "SATA", DiskCapacity: 0, Spec: "4C8G", Count: 1},
					}},
				},
			},
			BackupID: "backup123",
		})
		assert.NoError(t, err)
	})
	t.Run("no computes", func(t *testing.T) {
		_, err := manager.RestoreNewCluster(context.TODO(), cluster.RestoreNewClusterReq{})
		assert.Error(t, err)
	})
}

func TestManager_GetMonitorInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("normal", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)

		clusterRW.EXPECT().GetMeta(gomock.Any(), gomock.Any()).Return(&management.Cluster{
			Entity: common.Entity{
				ID: "2145635758",
			},
		}, []*management.ClusterInstance{
			{
				Entity: common.Entity{
					Status: string(constants.ClusterInstanceRunning),
				},
				Zone:     "zone1",
				CpuCores: 4,
				Memory:   8,
				Type:     "Grafana",
				Version:  "v5.0.0",
				Ports:    []int32{50001, 50002},
				HostIP:   []string{"127.0.0.5"},
			},
			{
				Entity: common.Entity{
					Status: string(constants.ClusterInstanceRunning),
				},
				Zone:     "zone1",
				CpuCores: 4,
				Memory:   8,
				Type:     "AlertManger",
				Version:  "v5.0.0",
				Ports:    []int32{60001, 60002},
				HostIP:   []string{"127.0.0.6"},
			},
		}, nil)
		manager := Manager{}
		got, err := manager.GetMonitorInfo(context.TODO(), cluster.QueryMonitorInfoReq{
			ClusterID: "2145635758",
		})
		assert.NoError(t, err)
		assert.Equal(t, got.ClusterID, "2145635758")
		assert.Equal(t, got.AlertUrl, "http://127.0.0.6:60001")
		assert.Equal(t, got.GrafanaUrl, "http://127.0.0.5:50001")
	})

	t.Run("err ip", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)

		clusterRW.EXPECT().GetMeta(gomock.Any(), gomock.Any()).Return(&management.Cluster{
			Entity: common.Entity{
				ID: "2145635758",
			},
		}, []*management.ClusterInstance{
			{
				Entity: common.Entity{
					Status: string(constants.ClusterInstanceRunning),
				},
				Zone:     "zone1",
				CpuCores: 4,
				Memory:   8,
				Type:     "Grafana",
				Version:  "v5.0.0",
				Ports:    []int32{50001, 50002},
				HostIP:   []string{"127.0.0.5"},
			},
		}, nil)
		manager := Manager{}
		_, err := manager.GetMonitorInfo(context.TODO(), cluster.QueryMonitorInfoReq{
			ClusterID: "2145635758",
		})
		assert.Error(t, err)
	})

	t.Run("grafana err port", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)

		clusterRW.EXPECT().GetMeta(gomock.Any(), gomock.Any()).Return(&management.Cluster{
			Entity: common.Entity{
				ID: "2145635758",
			},
		}, []*management.ClusterInstance{
			{
				Entity: common.Entity{
					Status: string(constants.ClusterInstanceRunning),
				},
				Zone:     "zone1",
				CpuCores: 4,
				Memory:   8,
				Type:     "Grafana",
				Version:  "v5.0.0",
				Ports:    []int32{-50001, -50002},
				HostIP:   []string{"127.0.0.5"},
			},
			{
				Entity: common.Entity{
					Status: string(constants.ClusterInstanceRunning),
				},
				Zone:     "zone1",
				CpuCores: 4,
				Memory:   8,
				Type:     "AlertManger",
				Version:  "v5.0.0",
				Ports:    []int32{60001, 60002},
				HostIP:   []string{"127.0.0.6"},
			},
		}, nil)
		manager := Manager{}
		_, err := manager.GetMonitorInfo(context.TODO(), cluster.QueryMonitorInfoReq{
			ClusterID: "2145635758",
		})
		assert.Error(t, err)
	})

	t.Run("alert err port", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)

		clusterRW.EXPECT().GetMeta(gomock.Any(), gomock.Any()).Return(&management.Cluster{
			Entity: common.Entity{
				ID: "2145635758",
			},
		}, []*management.ClusterInstance{
			{
				Entity: common.Entity{
					Status: string(constants.ClusterInstanceRunning),
				},
				Zone:     "zone1",
				CpuCores: 4,
				Memory:   8,
				Type:     "Grafana",
				Version:  "v5.0.0",
				Ports:    []int32{50001, 50002},
				HostIP:   []string{"127.0.0.5"},
			},
			{
				Entity: common.Entity{
					Status: string(constants.ClusterInstanceRunning),
				},
				Zone:     "zone1",
				CpuCores: 4,
				Memory:   8,
				Type:     "AlertManger",
				Version:  "v5.0.0",
				Ports:    []int32{-60001, -60002},
				HostIP:   []string{"127.0.0.6"},
			},
		}, nil)
		manager := Manager{}
		_, err := manager.GetMonitorInfo(context.TODO(), cluster.QueryMonitorInfoReq{
			ClusterID: "2145635758",
		})
		assert.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		clusterRW := mockclustermanagement.NewMockReaderWriter(ctrl)
		models.SetClusterReaderWriter(clusterRW)

		clusterRW.EXPECT().GetMeta(gomock.Any(), gomock.Any()).Return(&management.Cluster{}, []*management.ClusterInstance{
			{},
			{},
		}, fmt.Errorf("not found"))
		manager := Manager{}
		_, err := manager.GetMonitorInfo(context.TODO(), cluster.QueryMonitorInfoReq{
			ClusterID: "2145635758",
		})
		assert.Error(t, err)
	})
}

func TestNewClusterManager(t *testing.T) {
	got := NewClusterManager()
	assert.NotNil(t, got)
}
