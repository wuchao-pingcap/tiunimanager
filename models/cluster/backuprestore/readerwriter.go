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

package backuprestore

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type DataTransportReaderWriter interface {
	SetDb(db *gorm.DB)
	Db(ctx context.Context) *gorm.DB

	CreateBackupRecord(ctx context.Context, record *BackupRecord) (err error)
	UpdateBackupRecord(ctx context.Context, recordId string, status string, size uint64, backupTso int64, endTime time.Time) (err error)
	GetBackupRecord(ctx context.Context, backupId string) (record *BackupRecord, err error)
	QueryBackupRecords(ctx context.Context, backupId string, startTime, endTime time.Time, page int, pageSize int) (records []*BackupRecord, total int64, err error)
	DeleteBackupRecord(ctx context.Context, backupId string) (err error)

	CreateBackupStrategy(ctx context.Context, strategy *BackupStrategy) (err error)
	UpdateBackupStrategy(ctx context.Context, strategy *BackupStrategy) (err error)
	GetBackupStrategy(ctx context.Context, clusterId string) (strategy *BackupStrategy, err error)
	QueryBackupStrategy(ctx context.Context, weekDay string, startHour uint32) (strategies []*BackupStrategy, err error)
	DeleteBackupStrategy(ctx context.Context, clusterId string) (err error)
}
