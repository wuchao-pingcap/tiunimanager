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

package domain

import (
	ctx "context"
	"errors"
	"testing"
	"time"

	"github.com/pingcap-inc/tiem/library/common/resource-type"
	"github.com/pingcap-inc/tiem/test/mocksecondparty"

	"github.com/golang/mock/gomock"
	"github.com/pingcap-inc/tiem/library/client"
	"github.com/pingcap-inc/tiem/library/client/cluster/clusterpb"
	"github.com/pingcap-inc/tiem/library/client/metadb/dbpb"
	"github.com/pingcap-inc/tiem/library/secondparty"
	"github.com/pingcap-inc/tiem/micro-metadb/service"
	"github.com/pingcap-inc/tiem/test/mockdb"
	"github.com/pingcap/tiup/pkg/cluster/spec"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestSaveBackupStrategyPreCheck_case1(t *testing.T) {
	err := SaveBackupStrategyPreCheck(&clusterpb.OperatorDTO{}, &clusterpb.BackupStrategy{
		Period:     "14:00-16:00",
		BackupDate: "Monday,Tuesday,Thursday,Sunday",
	})
	assert.NoError(t, err)
}

func TestSaveBackupStrategyPreCheck_case3(t *testing.T) {
	err := SaveBackupStrategyPreCheck(&clusterpb.OperatorDTO{}, &clusterpb.BackupStrategy{
		Period:     "14:0016:00",
		BackupDate: "Monday,Tuesday,Thursday,Sunday",
	})
	assert.NotNil(t, err)
}

func TestSaveBackupStrategyPreCheck_case4(t *testing.T) {
	err := SaveBackupStrategyPreCheck(&clusterpb.OperatorDTO{}, &clusterpb.BackupStrategy{
		Period:     "18-16",
		BackupDate: "Monday,Tuesday,Thursday,Sunday",
	})
	assert.NotNil(t, err)
}

func TestSaveBackupStrategyPreCheck_case5(t *testing.T) {
	err := SaveBackupStrategyPreCheck(&clusterpb.OperatorDTO{}, &clusterpb.BackupStrategy{
		Period:     "14:00-16:00",
		BackupDate: "Monday,weekday,Thursday,Sunday",
	})
	assert.NotNil(t, err)
}

func TestBackup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().SaveBackupRecord(gomock.Any(), gomock.Any()).Return(&dbpb.DBSaveBackupRecordResponse{}, nil)
	client.DBClient = mockClient

	_, err := Backup(ctx.Background(), &clusterpb.OperatorDTO{
		Id:       "123",
		Name:     "123",
		TenantId: "123",
	}, "test-abc", "", "", BackupModeManual, "")

	assert.NoError(t, err)
}

func TestRecoverPreCheck(t *testing.T) {
	request := &clusterpb.RecoverRequest{
		Operator: &clusterpb.OperatorDTO{
			Id:   "123",
			Name: "123",
		},
		Cluster: &clusterpb.ClusterBaseInfoDTO{
			RecoverInfo: &clusterpb.RecoverInfoDTO{
				BackupRecordId:  123,
				SourceClusterId: "test-tidb",
			},
		},
	}

	err := RecoverPreCheck(context.TODO(), request)

	assert.NoError(t, err)
}

func TestRecover(t *testing.T) {
	_, err := Recover(ctx.Background(), &clusterpb.OperatorDTO{
		Id:       "123",
		Name:     "123",
		TenantId: "123",
	}, &clusterpb.ClusterBaseInfoDTO{
		ClusterName:    "test-tidb",
		ClusterVersion: &clusterpb.ClusterVersionDTO{Code: "v4.0.12", Name: "v4.0.12"},
		ClusterType:    &clusterpb.ClusterTypeDTO{Code: "TiDB", Name: "TiDB"},
	}, &clusterpb.ClusterCommonDemandDTO{Region: "111", Exclusive: false, CpuArchitecture: string(resource.X86_64)}, nil)

	assert.NoError(t, err)
}

func TestDeleteBackup_case1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().QueryBackupRecords(gomock.Any(), gomock.Any()).Return(&dbpb.DBQueryBackupRecordResponse{}, nil)
	mockClient.EXPECT().DeleteBackupRecord(gomock.Any(), gomock.Any()).Return(&dbpb.DBDeleteBackupRecordResponse{}, nil)
	client.DBClient = mockClient

	err := DeleteBackup(ctx.Background(), &clusterpb.OperatorDTO{
		Id:       "123",
		Name:     "123",
		TenantId: "123",
	}, "test-abc", 123)

	assert.NoError(t, err)
}

func TestDeleteBackup_case2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().QueryBackupRecords(gomock.Any(), gomock.Any()).Return(&dbpb.DBQueryBackupRecordResponse{}, errors.New("failed"))
	client.DBClient = mockClient

	err := DeleteBackup(ctx.Background(), &clusterpb.OperatorDTO{
		Id:       "123",
		Name:     "123",
		TenantId: "123",
	}, "test-abc", 123)

	assert.NotNil(t, err)
}

func TestDeleteBackup_case3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().QueryBackupRecords(gomock.Any(), gomock.Any()).Return(&dbpb.DBQueryBackupRecordResponse{
		Status: service.BizErrResponseStatus,
	}, nil)
	client.DBClient = mockClient

	err := DeleteBackup(ctx.Background(), &clusterpb.OperatorDTO{
		Id:       "123",
		Name:     "123",
		TenantId: "123",
	}, "test-abc", 123)

	assert.NotNil(t, err)
}

func TestDeleteBackup_case4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().QueryBackupRecords(gomock.Any(), gomock.Any()).Return(&dbpb.DBQueryBackupRecordResponse{}, nil)
	mockClient.EXPECT().DeleteBackupRecord(gomock.Any(), gomock.Any()).Return(&dbpb.DBDeleteBackupRecordResponse{}, errors.New("failed"))
	client.DBClient = mockClient

	err := DeleteBackup(ctx.Background(), &clusterpb.OperatorDTO{
		Id:       "123",
		Name:     "123",
		TenantId: "123",
	}, "test-abc", 123)

	assert.NotNil(t, err)
}

func TestDeleteBackup_case5(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().QueryBackupRecords(gomock.Any(), gomock.Any()).Return(&dbpb.DBQueryBackupRecordResponse{}, nil)
	mockClient.EXPECT().DeleteBackupRecord(gomock.Any(), gomock.Any()).Return(&dbpb.DBDeleteBackupRecordResponse{
		Status: service.BizErrResponseStatus,
	}, nil)
	client.DBClient = mockClient

	err := DeleteBackup(ctx.Background(), &clusterpb.OperatorDTO{
		Id:       "123",
		Name:     "123",
		TenantId: "123",
	}, "test-abc", 123)

	assert.NotNil(t, err)
}

func TestDeleteBackups_case1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().ListBackupRecords(gomock.Any(), gomock.Any()).Return(&dbpb.DBListBackupRecordsResponse{
		BackupRecords: []*dbpb.DBDBBackupRecordDisplayDTO{
			{
				BackupRecord: &dbpb.DBBackupRecordDTO{
					Id: 111,
				},
			},
		},
	}, nil)
	mockClient.EXPECT().DeleteBackupRecord(gomock.Any(), gomock.Any()).Return(&dbpb.DBDeleteBackupRecordResponse{}, nil)
	client.DBClient = mockClient

	err := DeleteBackups(ctx.Background(), &clusterpb.OperatorDTO{
		Id:       "123",
		Name:     "123",
		TenantId: "123",
	}, "test-abc", string(BackupModeAuto))

	assert.NoError(t, err)
}

func TestSaveBackupStrategy_case1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().SaveBackupStrategy(gomock.Any(), gomock.Any()).Return(&dbpb.DBSaveBackupStrategyResponse{}, nil)
	client.DBClient = mockClient

	err := SaveBackupStrategy(ctx.Background(), &clusterpb.OperatorDTO{
		Id:       "123",
		Name:     "123",
		TenantId: "123",
	}, &clusterpb.BackupStrategy{
		ClusterId:  "test-abc",
		BackupDate: "Monday,Sunday",
		Period:     "12:00-13:00",
	})
	assert.NoError(t, err)
}

func TestSaveBackupStrategy_case2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().SaveBackupStrategy(gomock.Any(), gomock.Any()).Return(&dbpb.DBSaveBackupStrategyResponse{}, errors.New("failed"))
	client.DBClient = mockClient

	err := SaveBackupStrategy(ctx.Background(), &clusterpb.OperatorDTO{
		Id:       "123",
		Name:     "123",
		TenantId: "123",
	}, &clusterpb.BackupStrategy{
		ClusterId:  "test-abc",
		BackupDate: "Monday,Sunday",
		Period:     "12:00-13:00",
	})
	assert.NotNil(t, err)
}

func TestSaveBackupStrategy_case3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().SaveBackupStrategy(gomock.Any(), gomock.Any()).Return(&dbpb.DBSaveBackupStrategyResponse{
		Status: service.BizErrResponseStatus,
	}, nil)
	client.DBClient = mockClient

	err := SaveBackupStrategy(ctx.Background(), &clusterpb.OperatorDTO{
		Id:       "123",
		Name:     "123",
		TenantId: "123",
	}, &clusterpb.BackupStrategy{
		ClusterId:  "test-abc",
		BackupDate: "Monday,Sunday",
		Period:     "12:00-13:00",
	})
	assert.NotNil(t, err)
}

func TestQueryBackupStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().QueryBackupStrategy(gomock.Any(), gomock.Any()).Return(&dbpb.DBQueryBackupStrategyResponse{
		Strategy: &dbpb.DBBackupStrategyDTO{
			TenantId:   "123",
			OperatorId: "123",
			ClusterId:  "test-abc",
			BackupDate: "Monday,Sunday",
			StartHour:  12,
			EndHour:    13,
		},
	}, nil)
	client.DBClient = mockClient

	_, err := QueryBackupStrategy(ctx.Background(), &clusterpb.OperatorDTO{
		Id:       "123",
		Name:     "123",
		TenantId: "123",
	}, "test-abc")

	assert.NoError(t, err)
}

func Test_calculateNextBackupTime_case1(t *testing.T) {
	now := time.Date(2021, 9, 17, 11, 0, 0, 0, time.Local)
	weekdayStr := "Monday,Tuesday,Thursday"
	startHour := 10
	nextBackupTime, err := calculateNextBackupTime(now, weekdayStr, startHour)
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2021, 9, 20, 10, 0, 0, 0, time.Local), nextBackupTime)
}

func Test_calculateNextBackupTime_case2(t *testing.T) {
	now := time.Date(2021, 9, 17, 11, 0, 0, 0, time.Local)
	weekdayStr := "Monday,Tuesday,Thursday,Sunday"
	startHour := 10
	nextBackupTime, err := calculateNextBackupTime(now, weekdayStr, startHour)
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2021, 9, 19, 10, 0, 0, 0, time.Local), nextBackupTime)
}

func Test_calculateNextBackupTime_case3(t *testing.T) {
	now := time.Date(2021, 9, 17, 11, 0, 0, 0, time.Local)
	weekdayStr := "Tuesday,Friday"
	startHour := 11
	nextBackupTime, err := calculateNextBackupTime(now, weekdayStr, startHour)
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2021, 9, 21, 11, 0, 0, 0, time.Local), nextBackupTime)
}

func Test_calculateNextBackupTime_case4(t *testing.T) {
	now := time.Date(2021, 9, 17, 11, 0, 0, 0, time.Local)
	weekdayStr := "Tuesday,Friday"
	startHour := 12
	nextBackupTime, err := calculateNextBackupTime(now, weekdayStr, startHour)
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2021, 9, 17, 12, 0, 0, 0, time.Local), nextBackupTime)
}

func Test_calculateNextBackupTime_case5(t *testing.T) {
	now := time.Date(2021, 9, 17, 11, 0, 0, 0, time.Local)
	weekdayStr := "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday"
	startHour := 12
	nextBackupTime, err := calculateNextBackupTime(now, weekdayStr, startHour)
	assert.NoError(t, err)
	t.Log(nextBackupTime)
	assert.Equal(t, time.Date(2021, 9, 17, 12, 0, 0, 0, time.Local), nextBackupTime)
}

func Test_calculateNextBackupTime_case6(t *testing.T) {
	now := time.Date(2021, 9, 17, 11, 0, 0, 0, time.Local)
	weekdayStr := "Tuesday,Weekday"
	startHour := 12
	_, err := calculateNextBackupTime(now, weekdayStr, startHour)
	assert.NotNil(t, err)
}

func Test_convertBrStorageType(t *testing.T) {
	result, _ := convertBrStorageType(string(StorageTypeS3))
	assert.Equal(t, secondparty.StorageTypeS3, result)
	result, _ = convertBrStorageType(string(StorageTypeLocal))
	assert.Equal(t, secondparty.StorageTypeLocal, result)
	_, err := convertBrStorageType("data")
	assert.NotNil(t, err)
}

func Test_updateBackupRecord_case1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().UpdateBackupRecord(gomock.Any(), gomock.Any()).Return(&dbpb.DBUpdateBackupRecordResponse{}, nil)
	mockClient.EXPECT().FindTiupTaskByID(gomock.Any(), gomock.Any()).Return(&dbpb.FindTiupTaskByIDResponse{TiupTask: &dbpb.TiupTask{
		Status: dbpb.TiupTaskStatus_Finished,
	}}, nil)
	client.DBClient = mockClient

	task := &TaskEntity{}
	flowCtx := NewFlowContext(ctx.TODO())
	flowCtx.SetData(contextClusterKey, &ClusterAggregation{
		LastBackupRecord: &BackupRecord{
			Id:   123,
			Size: 1000,
		},
	})
	flowCtx.SetData("backupTaskId", uint64(123))

	ret := updateBackupRecord(task, flowCtx)

	assert.Equal(t, true, ret)
}

func Test_updateBackupRecord_case2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockdb.NewMockTiEMDBService(ctrl)
	mockClient.EXPECT().UpdateBackupRecord(gomock.Any(), gomock.Any()).Return(&dbpb.DBUpdateBackupRecordResponse{}, errors.New("failed"))
	mockClient.EXPECT().FindTiupTaskByID(gomock.Any(), gomock.Any()).Return(&dbpb.FindTiupTaskByIDResponse{TiupTask: &dbpb.TiupTask{
		Status: dbpb.TiupTaskStatus_Finished,
	}}, nil)
	client.DBClient = mockClient

	task := &TaskEntity{}
	flowCtx := NewFlowContext(ctx.TODO())
	flowCtx.SetData(contextClusterKey, &ClusterAggregation{
		LastBackupRecord: &BackupRecord{
			Id:   123,
			Size: 1000,
		},
		Cluster: &Cluster{
			Id: "test-tidb",
		},
	})
	flowCtx.SetData("backupTaskId", uint64(123))
	ret := updateBackupRecord(task, flowCtx)

	assert.Equal(t, false, ret)
}

func Test_backupCluster_case1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTiup := mocksecondparty.NewMockMicroSrv(ctrl)
	mockTiup.EXPECT().MicroSrvBackUp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(uint64(123), nil)
	secondparty.SecondParty = mockTiup

	task := &TaskEntity{
		Id: 123,
	}
	flowCtx := NewFlowContext(ctx.TODO())
	flowCtx.SetData(contextClusterKey, &ClusterAggregation{
		LastBackupRecord: &BackupRecord{
			Id:          123,
			StorageType: StorageTypeS3,
		},
		Cluster: &Cluster{
			Id:          "test-tidb123",
			ClusterName: "test-tidb",
		},
		CurrentTopologyConfigRecord: &TopologyConfigRecord{
			ConfigModel: &spec.Specification{
				TiDBServers: []*spec.TiDBSpec{
					{
						Host: "127.0.0.1",
						Port: 4000,
					},
				},
			},
		},
	})
	ret := backupCluster(task, flowCtx)
	assert.Equal(t, true, ret)
}

func Test_backupCluster_case2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTiup := mocksecondparty.NewMockMicroSrv(ctrl)
	mockTiup.EXPECT().MicroSrvBackUp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(uint64(123), nil)
	secondparty.SecondParty = mockTiup

	task := &TaskEntity{
		Id: 123,
	}
	flowCtx := NewFlowContext(ctx.TODO())
	flowCtx.SetData(contextClusterKey, &ClusterAggregation{
		LastBackupRecord: &BackupRecord{
			Id:          123,
			StorageType: StorageTypeS3,
		},
		Cluster: &Cluster{
			Id:          "test-tidb123",
			ClusterName: "test-tidb",
		},
		CurrentTopologyConfigRecord: &TopologyConfigRecord{
			ConfigModel: &spec.Specification{
				TiDBServers: []*spec.TiDBSpec{
					{
						Host: "127.0.0.1",
						Port: 4000,
					},
				},
			},
		},
	})
	ret := backupCluster(task, flowCtx)
	assert.Equal(t, true, ret)
}

func Test_backupCluster_case3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTiup := mocksecondparty.NewMockMicroSrv(ctrl)
	mockTiup.EXPECT().MicroSrvBackUp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(uint64(123), errors.New("failed"))
	secondparty.SecondParty = mockTiup

	task := &TaskEntity{
		Id: 123,
	}
	flowCtx := NewFlowContext(ctx.TODO())
	flowCtx.SetData(contextClusterKey, &ClusterAggregation{
		LastBackupRecord: &BackupRecord{
			Id:          123,
			StorageType: StorageTypeS3,
		},
		Cluster: &Cluster{
			Id:          "test-tidb123",
			ClusterName: "test-tidb",
		},
		CurrentTopologyConfigRecord: &TopologyConfigRecord{
			ConfigModel: &spec.Specification{
				TiDBServers: []*spec.TiDBSpec{
					{
						Host: "127.0.0.1",
						Port: 4000,
					},
				},
			},
		},
	})

	ret := backupCluster(task, flowCtx)
	assert.Equal(t, false, ret)
}

func Test_recoverFromSrcCluster_case1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTiup := mocksecondparty.NewMockMicroSrv(ctrl)
	mockTiup.EXPECT().MicroSrvRestore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(uint64(123), nil)
	secondparty.SecondParty = mockTiup

	mockDB := mockdb.NewMockTiEMDBService(ctrl)

	mockDB.EXPECT().QueryBackupRecords(gomock.Any(), gomock.Any()).Return(&dbpb.DBQueryBackupRecordResponse{
		Status: &dbpb.DBClusterResponseStatus{
			Code: service.ClusterSuccessResponseStatus.GetCode(),
		},
		BackupRecords: &dbpb.DBDBBackupRecordDisplayDTO{
			BackupRecord: &dbpb.DBBackupRecordDTO{
				StorageType: string(StorageTypeS3),
				FilePath:    "/tmp/test",
			},
		},
	}, nil)
	client.DBClient = mockDB

	task := &TaskEntity{
		Id: 123,
	}
	flowCtx := NewFlowContext(ctx.TODO())
	flowCtx.SetData(contextClusterKey, &ClusterAggregation{
		LastBackupRecord: &BackupRecord{
			Id:          123,
			StorageType: StorageTypeS3,
		},
		Cluster: &Cluster{
			Id:          "test-tidb123",
			ClusterName: "test-tidb",
			RecoverInfo: RecoverInfo{
				SourceClusterId: "src-tidb",
				BackupRecordId:  123,
			},
		},
		CurrentTopologyConfigRecord: &TopologyConfigRecord{
			ConfigModel: &spec.Specification{
				TiDBServers: []*spec.TiDBSpec{
					{
						Host: "127.0.0.1",
						Port: 4000,
					},
				},
			},
		},
	})
	flowCtx.SetData("startTaskId", uint64(123))
	ret := recoverFromSrcCluster(task, flowCtx)

	assert.Equal(t, true, ret)
}

func Test_recoverFromSrcCluster_case2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTiup := mocksecondparty.NewMockMicroSrv(ctrl)
	mockTiup.EXPECT().MicroSrvRestore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(uint64(123), nil)
	secondparty.SecondParty = mockTiup

	mockDB := mockdb.NewMockTiEMDBService(ctrl)
	mockDB.EXPECT().QueryBackupRecords(gomock.Any(), gomock.Any()).Return(&dbpb.DBQueryBackupRecordResponse{
		Status: &dbpb.DBClusterResponseStatus{
			Code: service.ClusterSuccessResponseStatus.GetCode(),
		},
		BackupRecords: &dbpb.DBDBBackupRecordDisplayDTO{
			BackupRecord: &dbpb.DBBackupRecordDTO{
				StorageType: string(StorageTypeS3),
				FilePath:    "/tmp/test",
			},
		},
	}, nil)
	client.DBClient = mockDB

	task := &TaskEntity{
		Id: 123,
	}
	flowCtx := NewFlowContext(ctx.TODO())
	flowCtx.SetData(contextClusterKey, &ClusterAggregation{
		LastBackupRecord: &BackupRecord{
			Id:          123,
			StorageType: StorageTypeS3,
		},
		Cluster: &Cluster{
			Id:          "test-tidb123",
			ClusterName: "test-tidb",
			RecoverInfo: RecoverInfo{
				SourceClusterId: "src-tidb",
				BackupRecordId:  123,
			},
		},
		CurrentTopologyConfigRecord: &TopologyConfigRecord{
			ConfigModel: &spec.Specification{
				TiDBServers: []*spec.TiDBSpec{
					{
						Host: "127.0.0.1",
						Port: 4000,
					},
				},
			},
		},
	})
	flowCtx.SetData("startTaskId", uint64(123))
	ret := recoverFromSrcCluster(task, flowCtx)

	assert.Equal(t, true, ret)
}

func Test_recoverFromSrcCluster_case3(t *testing.T) {
	task := &TaskEntity{
		Id: 123,
	}
	flowCtx := NewFlowContext(ctx.TODO())
	flowCtx.SetData(contextClusterKey, &ClusterAggregation{
		LastBackupRecord: &BackupRecord{
			Id:          123,
			StorageType: StorageTypeS3,
		},
		Cluster: &Cluster{
			Id:          "test-tidb123",
			ClusterName: "test-tidb",
			RecoverInfo: RecoverInfo{
				SourceClusterId: "",
				BackupRecordId:  123,
			},
		},
		CurrentTopologyConfigRecord: &TopologyConfigRecord{
			ConfigModel: &spec.Specification{
				TiDBServers: []*spec.TiDBSpec{
					{
						Host: "127.0.0.1",
						Port: 4000,
					},
				},
			},
		},
	})
	flowCtx.SetData("startTaskId", uint64(123))
	ret := recoverFromSrcCluster(task, flowCtx)

	assert.Equal(t, true, ret)
}

func Test_recoverFromSrcCluster_case4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockdb.NewMockTiEMDBService(ctrl)

	mockDB.EXPECT().QueryBackupRecords(gomock.Any(), gomock.Any()).Return(&dbpb.DBQueryBackupRecordResponse{
		Status: &dbpb.DBClusterResponseStatus{
			Code: service.ClusterSuccessResponseStatus.GetCode(),
		},
		BackupRecords: &dbpb.DBDBBackupRecordDisplayDTO{
			BackupRecord: &dbpb.DBBackupRecordDTO{
				StorageType: string(StorageTypeS3),
				FilePath:    "/tmp/test",
			},
		},
	}, errors.New("failed"))
	client.DBClient = mockDB

	task := &TaskEntity{
		Id: 123,
	}
	flowCtx := NewFlowContext(ctx.TODO())
	flowCtx.SetData(contextClusterKey, &ClusterAggregation{
		LastBackupRecord: &BackupRecord{
			Id:          123,
			StorageType: StorageTypeS3,
		},
		Cluster: &Cluster{
			Id:          "test-tidb123",
			ClusterName: "test-tidb",
			RecoverInfo: RecoverInfo{
				SourceClusterId: "src-tidb",
				BackupRecordId:  123,
			},
		},
		CurrentTopologyConfigRecord: &TopologyConfigRecord{
			ConfigModel: &spec.Specification{
				TiDBServers: []*spec.TiDBSpec{
					{
						Host: "127.0.0.1",
						Port: 4000,
					},
				},
			},
		},
	})
	flowCtx.SetData("startTaskId", uint64(123))
	ret := recoverFromSrcCluster(task, flowCtx)

	assert.Equal(t, false, ret)
}

func Test_recoverFromSrcCluster_case5(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockdb.NewMockTiEMDBService(ctrl)

	mockDB.EXPECT().QueryBackupRecords(gomock.Any(), gomock.Any()).Return(&dbpb.DBQueryBackupRecordResponse{
		Status: &dbpb.DBClusterResponseStatus{
			Code: service.ClusterSuccessResponseStatus.GetCode(),
		},
		BackupRecords: &dbpb.DBDBBackupRecordDisplayDTO{
			BackupRecord: &dbpb.DBBackupRecordDTO{
				StorageType: "oss",
				FilePath:    "/tmp/test",
			},
		},
	}, nil)
	client.DBClient = mockDB

	task := &TaskEntity{
		Id: 123,
	}
	flowCtx := NewFlowContext(ctx.TODO())
	flowCtx.SetData(contextClusterKey, &ClusterAggregation{
		LastBackupRecord: &BackupRecord{
			Id:          123,
			StorageType: StorageTypeS3,
		},
		Cluster: &Cluster{
			Id:          "test-tidb123",
			ClusterName: "test-tidb",
			RecoverInfo: RecoverInfo{
				SourceClusterId: "src-tidb",
				BackupRecordId:  123,
			},
		},
		CurrentTopologyConfigRecord: &TopologyConfigRecord{
			ConfigModel: &spec.Specification{
				TiDBServers: []*spec.TiDBSpec{
					{
						Host: "127.0.0.1",
						Port: 4000,
					},
				},
			},
		},
	})
	flowCtx.SetData("startTaskId", uint64(123))
	ret := recoverFromSrcCluster(task, flowCtx)

	assert.Equal(t, false, ret)
}

func Test_recoverFromSrcCluster_case6(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTiup := mocksecondparty.NewMockMicroSrv(ctrl)
	mockTiup.EXPECT().MicroSrvRestore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(uint64(123), errors.New("failed"))
	secondparty.SecondParty = mockTiup

	mockDB := mockdb.NewMockTiEMDBService(ctrl)

	mockDB.EXPECT().QueryBackupRecords(gomock.Any(), gomock.Any()).Return(&dbpb.DBQueryBackupRecordResponse{
		Status: &dbpb.DBClusterResponseStatus{
			Code: service.ClusterSuccessResponseStatus.GetCode(),
		},
		BackupRecords: &dbpb.DBDBBackupRecordDisplayDTO{
			BackupRecord: &dbpb.DBBackupRecordDTO{
				StorageType: string(StorageTypeS3),
				FilePath:    "/tmp/test",
			},
		},
	}, nil)
	client.DBClient = mockDB

	task := &TaskEntity{
		Id: 123,
	}
	flowCtx := NewFlowContext(ctx.TODO())
	flowCtx.SetData(contextClusterKey, &ClusterAggregation{
		LastBackupRecord: &BackupRecord{
			Id:          123,
			StorageType: StorageTypeS3,
		},
		Cluster: &Cluster{
			Id:          "test-tidb123",
			ClusterName: "test-tidb",
			RecoverInfo: RecoverInfo{
				SourceClusterId: "src-tidb",
				BackupRecordId:  123,
			},
		},
		CurrentTopologyConfigRecord: &TopologyConfigRecord{
			ConfigModel: &spec.Specification{
				TiDBServers: []*spec.TiDBSpec{
					{
						Host: "127.0.0.1",
						Port: 4000,
					},
				},
			},
		},
	})
	flowCtx.SetData("startTaskId", uint64(123))
	ret := recoverFromSrcCluster(task, flowCtx)

	assert.Equal(t, false, ret)
}
