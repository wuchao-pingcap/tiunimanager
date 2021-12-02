/******************************************************************************
 * Copyright (c)  2021 PingCAP, Inc.                                          *
 * Licensed under the Apache License, Version 2.0 (the "License");            *
 * you may not use this file except in compliance with the License.           *
 * You may obtain a copy of the License at                                    *
 *                                                                            *
 * http://www.apache.org/licenses/LICENSE-2.0                                 *
 *                                                                            *
 *  Unless required by applicable law or agreed to in writing, software       *
 *  distributed under the License is distributed on an "AS IS" BASIS,         *
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  *
 *  See the License for the specific language governing permissions and       *
 *  limitations under the License.                                            *
 ******************************************************************************/

package importexport

import (
	"context"
	"github.com/pingcap-inc/tiem/library/common"
	"github.com/pingcap-inc/tiem/library/framework"
	"github.com/pingcap-inc/tiem/models"
	"time"
)

type ImportExportReadWrite struct {
}

func NewImportExportReadWrite() *ImportExportReadWrite {
	m := new(ImportExportReadWrite)
	return m
}

func (m *ImportExportReadWrite) CreateDataTransportRecord(ctx context.Context, record *DataTransportRecord) (err error) {
	return models.DB(ctx).Create(record).Error
}

func (m *ImportExportReadWrite) UpdateDataTransportRecord(ctx context.Context, recordId string, status string, endTime time.Time) (err error) {
	if "" == recordId {
		return framework.SimpleError(common.TIEM_PARAMETER_INVALID)
	}

	record := &DataTransportRecord{}
	err = models.DB(ctx).First(record, "id = ?", recordId).Error
	if err != nil {
		return framework.SimpleError(common.TIEM_TRANSPORT_RECORD_NOT_FOUND)
	}

	return models.DB(ctx).Model(record).
		Update("status", status).
		Update("end_time", endTime).Error
}

func (m *ImportExportReadWrite) GetDataTransportRecord(ctx context.Context, recorId string) (record *DataTransportRecord, err error) {
	if "" == recorId {
		return nil, framework.SimpleError(common.TIEM_PARAMETER_INVALID)
	}
	record = &DataTransportRecord{}
	err = models.DB(ctx).First(record, "id = ?", recorId).Error
	if err != nil {
		return nil, framework.SimpleError(common.TIEM_TRANSPORT_RECORD_NOT_FOUND)
	}
	return record, err
}

func (m *ImportExportReadWrite) QueryDataTransportRecords(ctx context.Context, recordId, clusterId string, reImport bool, startTime, endTime time.Time, page int, pageSize int) (records []*DataTransportRecord, total int64, err error) {
	records = make([]*DataTransportRecord, pageSize)
	query := models.DB(ctx).Model(DataTransportRecord{})
	if recordId != "" {
		query = query.Where("record_id = ?", recordId)
	}
	if clusterId != "" {
		query = query.Where("cluster_id = ?", clusterId)
	}
	if reImport {
		query = query.Where("re_import_support = ?", reImport)
	}
	if !startTime.IsZero() {
		query = query.Where("start_time >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("end_time <= ?", endTime)
	}
	err = query.Order("id desc").Count(&total).Offset(pageSize * (page - 1)).Limit(pageSize).Find(&records).Error
	return records, total, err
}

func (m *ImportExportReadWrite) DeleteDataTransportRecord(ctx context.Context, recordId string) (err error) {
	if "" == recordId {
		return framework.SimpleError(common.TIEM_PARAMETER_INVALID)
	}
	record := &DataTransportRecord{}
	return models.DB(ctx).First(record, "id = ?", recordId).Delete(record).Error
}
