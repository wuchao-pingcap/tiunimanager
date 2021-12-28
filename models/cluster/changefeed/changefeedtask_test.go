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

package changefeed

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/pingcap-inc/tiem/common/constants"
	"github.com/pingcap-inc/tiem/models/common"
	"github.com/stretchr/testify/assert"
)

func TestChangeFeedTask_Locked(t1 *testing.T) {
	type fields struct {
		Entity            common.Entity
		Name              string
		ClusterId         string
		DownstreamType    constants.DownstreamType
		StartTS           int64
		FilterRulesConfig string
		DownstreamConfig  string
		StatusLock        sql.NullTime
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"locked", fields{StatusLock: sql.NullTime{Time: time.Now(), Valid: true}}, true},
		{"invalid", fields{StatusLock: sql.NullTime{Time: time.Now(), Valid: false}}, false},
		{"expired", fields{StatusLock: sql.NullTime{Time: time.Now().Add(time.Minute * -2), Valid: true}}, false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := ChangeFeedTask{
				Entity:            tt.fields.Entity,
				Name:              tt.fields.Name,
				ClusterId:         tt.fields.ClusterId,
				Type:              tt.fields.DownstreamType,
				StartTS:           tt.fields.StartTS,
				FilterRulesConfig: tt.fields.FilterRulesConfig,
				DownstreamConfig:  tt.fields.DownstreamConfig,
				StatusLock:        tt.fields.StatusLock,
			}
			if got := t.Locked(); got != tt.want {
				t1.Errorf("Locked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormChangeFeedReadWrite_Create(t *testing.T) {
	existed, _ := testRW.Create(context.TODO(), &ChangeFeedTask{Entity: common.Entity{TenantId: "111"}})
	defer testRW.Delete(context.TODO(), existed.ID)

	type args struct {
		ctx  context.Context
		task *ChangeFeedTask
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"normal", args{context.TODO(), &ChangeFeedTask{Entity: common.Entity{TenantId: "111"}, ClusterId: "dafadsfefesdf"}}, false},
		{"without tenant", args{context.TODO(), &ChangeFeedTask{Entity: common.Entity{ID: existed.ID}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testRW.Create(tt.args.ctx, tt.args.task)
			defer testRW.Delete(context.TODO(), got.ID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				assert.NotEmpty(t, got.ID)
			}
		})
	}
}

func TestGormChangeFeedReadWrite_Delete(t *testing.T) {
	existed, _ := testRW.Create(context.TODO(), &ChangeFeedTask{Entity: common.Entity{TenantId: "111"}})
	deleted, _ := testRW.Create(context.TODO(), &ChangeFeedTask{Entity: common.Entity{TenantId: "111"}, ClusterId: "111"})
	testRW.DB(context.TODO()).Delete(deleted)

	type args struct {
		ctx    context.Context
		taskId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"normal", args{context.TODO(), existed.ID}, false},
		{"deleted", args{context.TODO(), deleted.ID}, true},
		{"not existed", args{context.TODO(), "dfadsafsa"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := testRW

			if err := m.Delete(tt.args.ctx, tt.args.taskId); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				_, e := m.Get(tt.args.ctx, tt.args.taskId)
				assert.Error(t, e)
			}

		})
	}
}

func TestGormChangeFeedReadWrite_LockStatus(t *testing.T) {
	locked, _ := testRW.Create(context.TODO(), &ChangeFeedTask{
		Entity: common.Entity{TenantId: "111"},
		StatusLock: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})

	unlocked, _ := testRW.Create(context.TODO(), &ChangeFeedTask{
		Entity: common.Entity{TenantId: "111"},
	})
	testRW.DB(context.TODO()).First(&ChangeFeedTask{}, "id = ?", unlocked.ID).Update("status_lock", sql.NullTime{
		Time:  time.Now(),
		Valid: false,
	})
	notExisted := "111"
	defer testRW.Delete(context.TODO(), locked.ID)
	defer testRW.Delete(context.TODO(), unlocked.ID)

	type args struct {
		ctx    context.Context
		taskId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"locked", args{context.TODO(), locked.ID}, true},
		{"unlocked", args{context.TODO(), unlocked.ID}, false},
		{"not existed", args{context.TODO(), notExisted}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := testRW
			if err := m.LockStatus(tt.args.ctx, tt.args.taskId); (err != nil) != tt.wantErr {
				t.Errorf("LockStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				updated, e := m.Get(tt.args.ctx, tt.args.taskId)
				assert.NoError(t, e)
				assert.True(t, updated.Locked())
			}

		})
	}
}

func TestGormChangeFeedReadWrite_QueryByClusterId(t *testing.T) {
	t1, _ := testRW.Create(context.TODO(), &ChangeFeedTask{Entity: common.Entity{TenantId: "111"}, ClusterId: "6666"})
	anotherCluster, _ := testRW.Create(context.TODO(), &ChangeFeedTask{Entity: common.Entity{TenantId: "111"}, ClusterId: "3121"})
	deleted, _ := testRW.Create(context.TODO(), &ChangeFeedTask{Entity: common.Entity{TenantId: "111"}, ClusterId: "6666"})
	t4, _ := testRW.Create(context.TODO(), &ChangeFeedTask{Entity: common.Entity{TenantId: "111"}, ClusterId: "6666", StartTS: int64(9999)})
	t5, _ := testRW.Create(context.TODO(), &ChangeFeedTask{Entity: common.Entity{TenantId: "111"}, ClusterId: "6666"})
	defer testRW.Delete(context.TODO(), t1.ID)
	defer testRW.Delete(context.TODO(), anotherCluster.ID)
	testRW.DB(context.TODO()).Delete(deleted)
	defer testRW.Delete(context.TODO(), t4.ID)
	defer testRW.Delete(context.TODO(), t5.ID)

	tasks, total, err := testRW.QueryByClusterId(context.TODO(), "6666", 0, 2)
	assert.NoError(t, err)
	assert.Equal(t, 3, int(total))
	assert.Equal(t, 9999, int(tasks[1].StartTS))

	_, _, err = testRW.QueryByClusterId(context.TODO(), "", 0, 2)
	assert.Error(t, err)
}

func TestGormChangeFeedReadWrite_UnlockStatus(t *testing.T) {
	newStatus := constants.ChangeFeedStatusStopped
	locked, _ := testRW.Create(context.TODO(), &ChangeFeedTask{
		Entity: common.Entity{TenantId: "111"},
		StatusLock: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	unlocked, _ := testRW.Create(context.TODO(), &ChangeFeedTask{
		Entity: common.Entity{TenantId: "111"},
	})
	testRW.DB(context.TODO()).First(&ChangeFeedTask{}, "id = ?", unlocked.ID).Update("status_lock", sql.NullTime{
		Time:  time.Now().Add(time.Minute * -3),
		Valid: true,
	})

	notExisted := "111"
	defer testRW.Delete(context.TODO(), locked.ID)
	defer testRW.Delete(context.TODO(), unlocked.ID)

	type args struct {
		ctx          context.Context
		taskId       string
		targetStatus constants.ChangeFeedStatus
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		finalStatus constants.ChangeFeedStatus
	}{
		{"locked", args{context.TODO(), locked.ID, newStatus}, false, newStatus},
		{"unlocked", args{context.TODO(), unlocked.ID, newStatus}, true, "0"},
		{"notExisted", args{context.TODO(), notExisted, newStatus}, true, "0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := testRW
			if err := m.UnlockStatus(tt.args.ctx, tt.args.taskId, tt.args.targetStatus); (err != nil) != tt.wantErr {
				t.Errorf("UnlockStatus() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.finalStatus != "0" {
				final, _ := m.Get(tt.args.ctx, tt.args.taskId)
				assert.Equal(t, string(tt.finalStatus), final.Status)
				assert.False(t, final.Locked())
			}
		})
	}
}

func TestGormChangeFeedReadWrite_UpdateConfig(t *testing.T) {
	existed, _ := testRW.Create(context.TODO(), &ChangeFeedTask{Entity: common.Entity{TenantId: "111"}})
	defer testRW.Delete(context.TODO(), existed.ID)

	newString := "new"
	newInt := 99
	existed.Downstream = &TiDBDownstream{
		Password: "updated",
	}
	existed.Type = constants.DownstreamTypeTiDB
	existed.FilterRules = []string{newString}
	existed.ClusterId = newString
	existed.StartTS = int64(newInt)
	existed.Entity.Status = "99"

	type args struct {
		ctx            context.Context
		updateTemplate *ChangeFeedTask
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"normal", args{context.TODO(), existed}, false},
		{"not existed", args{context.TODO(), &ChangeFeedTask{
			Entity: common.Entity{ID: "111"},
		}}, true},
		{"without id", args{context.TODO(), &ChangeFeedTask{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := testRW
			if err := m.UpdateConfig(tt.args.ctx, tt.args.updateTemplate); (err != nil) != tt.wantErr {
				t.Errorf("UpdateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				updated, _ := m.Get(tt.args.ctx, tt.args.updateTemplate.ID)
				assert.Equal(t, "updated", updated.Downstream.(*TiDBDownstream).Password)

				assert.Equal(t, "tidb", string(updated.Type))
				assert.Equal(t, []string{newString}, updated.FilterRules)
				assert.NotEqual(t, newString, updated.ClusterId)
				assert.NotEqual(t, int8(newInt), updated.Status)
			}
		})
	}
}

func TestConvertStatus(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus constants.ChangeFeedStatus
		wantErr    bool
	}{
		{"Initial", args{"Initial"}, constants.ChangeFeedStatusInitial, false},
		{"Normal", args{"Normal"}, constants.ChangeFeedStatusNormal, false},
		{"Stopped", args{"Stopped"}, constants.ChangeFeedStatusStopped, false},
		{"Finished", args{"Finished"}, constants.ChangeFeedStatusFinished, false},
		{"Error", args{"Error"}, constants.ChangeFeedStatusError, false},
		{"Failed", args{"Failed"}, constants.ChangeFeedStatusFailed, false},
		{"Unknown", args{"Unknown"}, constants.ChangeFeedStatusUnknown, true},
		{"whatever", args{"Unknown"}, constants.ChangeFeedStatusUnknown, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, err := constants.ConvertStatus(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("ConvertStatus() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestStatus_IsFinal(t *testing.T) {
	tests := []struct {
		name string
		s    constants.ChangeFeedStatus
		want bool
	}{
		{"Initial", constants.ChangeFeedStatusInitial, false},
		{"Normal", constants.ChangeFeedStatusNormal, false},
		{"Stopped", constants.ChangeFeedStatusStopped, false},
		{"Finished", constants.ChangeFeedStatusFinished, true},
		{"Error", constants.ChangeFeedStatusError, false},
		{"Failed", constants.ChangeFeedStatusFailed, true},
		{"Unknown", constants.ChangeFeedStatusUnknown, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsFinal(); got != tt.want {
				t.Errorf("IsFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}
