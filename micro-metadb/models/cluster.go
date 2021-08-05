package models

import (
	"errors"
)

type ClusterDO struct {
	Entity
	ClusterName 			string
	DbPassword 				string
	ClusterType 			string
	ClusterVersion 			string
	Tls 					bool
	Tags           			string
	OwnerId 				string		`gorm:"not null;type:varchar(36);default:null"`
	CurrentTiupConfigId     uint
	CurrentDemandId 		uint
	CurrentFlowId			uint
}

func (d ClusterDO) TableName() string {
	return "clusters"
}

type DemandRecordDO struct {
	Record
	ClusterId 			string		`gorm:"not null;type:varchar(36);default:null"`
	Content 			string		`gorm:"type:text"`
}

func (d DemandRecordDO) TableName() string {
	return "demand_records"
}

type TiUPConfigDO struct {
	Record
	ClusterId			string		`gorm:"not null;type:varchar(36);default:null"`
	Content 			string		`gorm:"type:text"`
}

func (d TiUPConfigDO) TableName() string {
	return "tiup_configs"
}

func UpdateClusterStatus(clusterId string, status int8) (cluster *ClusterDO, err error) {
	if clusterId == ""{
		return nil, errors.New("cluster id is empty")
	}
	cluster = &ClusterDO{}
	err = MetaDB.Model(cluster).Where("id = ?", clusterId).Update("status", status).Error
	return
}

func UpdateClusterDemand(clusterId string, content string, tenantId string) (cluster *ClusterDO, demand *DemandRecordDO, err error) {
	demand = &DemandRecordDO{
		ClusterId: clusterId,
		Content: content,
		Record: Record{
			TenantId: tenantId,
		},
	}

	err = MetaDB.Create(demand).Error
	if err != nil {
		return
	}

	cluster = &ClusterDO{}
	err = MetaDB.Model(cluster).Where("id = ?", clusterId).Update("current_demand_id", demand.ID).Find(cluster).Error
	return
}

func UpdateClusterFlowId(clusterId string, flowId uint) (cluster *ClusterDO, err error) {
	if clusterId == ""{
		return nil, errors.New("cluster id is empty")
	}
	cluster = &ClusterDO{}

	err = MetaDB.Model(cluster).Where("id = ?", clusterId).Update("current_flow_id", flowId).Find(cluster).Error

	return
}

func UpdateTiUPConfig(clusterId string, content string, tenantId string) (cluster *ClusterDO, err error) {
	cluster = &ClusterDO{}
	record := &TiUPConfigDO{
		ClusterId: clusterId,
		Content: content,
		Record: Record{
			TenantId: tenantId,
		},
	}

	err = MetaDB.Create(record).Error
	if err != nil {
		return
	}

	err = MetaDB.Model(cluster).Where("id = ?", clusterId).Update("current_tiup_config_id", record.ID).Find(cluster).Error

	return
}

func DeleteCluster(clusterId string) (cluster *ClusterDO, err error) {
	if clusterId == ""{
		 return nil, errors.New("empty cluster id")
	}
	cluster = &ClusterDO{}
	err = MetaDB.Find(cluster, "id = ?", clusterId).Error

	if err != nil {
		return
	}

	err = MetaDB.Delete(cluster).Error
	return
}

func FetchCluster(clusterId string) (result *ClusterFetchResult, err error) {
	result = &ClusterFetchResult{
		Cluster: &ClusterDO{},
		TiUPConfig: &TiUPConfigDO{},
		Flow: &FlowDO{},
	}

	err = MetaDB.Find(result.Cluster, clusterId).Error
	if err != nil {
		return
	}

	cluster := result.Cluster
	if cluster.CurrentDemandId > 0 {

		err = MetaDB.Find(result.DemandRecord, cluster.CurrentDemandId).Error
		if err != nil {
			return
		}
	}

	if cluster.CurrentTiupConfigId > 0 {
		err = MetaDB.Find(result.TiUPConfig, cluster.CurrentTiupConfigId).Error
		if err != nil {
			return
		}
	}

	if cluster.CurrentFlowId > 0 {
		err = MetaDB.Find(result.Flow, cluster.CurrentFlowId).Error
	}
	return
}

type ClusterFetchResult struct {
	Cluster *ClusterDO
	Flow *FlowDO
	DemandRecord *DemandRecordDO
	TiUPConfig *TiUPConfigDO
}

func ListClusterDetails(clusterId string,
	clusterName string,
	clusterType string,
	clusterStatus string,
	clusterTag string,
	offset int, length int) (result []*ClusterFetchResult, total int64, err error){

	clusters, total, err := ListClusters(clusterId, clusterName, clusterType, clusterStatus, clusterTag, offset, length)

	flowIds := make([]uint, len(clusters), len(clusters))
	demandIds := make([]uint, len(clusters), len(clusters))
	tiupConfigIds := make([]uint, len(clusters), len(clusters))

	result = make([]*ClusterFetchResult, len(clusters), len(clusters))
	clusterMap := make(map[string]*ClusterFetchResult)

	for i,c := range clusters {
		flowIds[i] = c.CurrentFlowId
		demandIds[i] = c.CurrentDemandId
		tiupConfigIds[i] = c.CurrentTiupConfigId
		result[i] = &ClusterFetchResult{
			Cluster: c,
		}
		clusterMap[c.ID] = result[i]
	}

	flows := make([]*FlowDO, len(clusters), len(clusters))
	err = MetaDB.Find(&flows, flowIds).Error
	for _,v := range flows {
		clusterMap[v.BizId].Flow = v
	}

	demands := make([]*DemandRecordDO, len(clusters), len(clusters))
	err = MetaDB.Find(&demands, demandIds).Error
	for _,v := range demands {
		clusterMap[v.ClusterId].DemandRecord = v
	}

	tiupConfigs := make([]*TiUPConfigDO, len(clusters), len(clusters))
	err = MetaDB.Find(&tiupConfigs, tiupConfigIds).Error
	for _,v := range tiupConfigs {
		clusterMap[v.ClusterId].TiUPConfig = v
	}

	return
}

func ListClusters(clusterId string,
	clusterName string,
	clusterType string,
	clusterStatus string,
	clusterTag string,
	offset int, length int) (clusters []*ClusterDO, total int64, err error){

	clusters = make([]*ClusterDO, length, length)

	db := MetaDB.Table("clusters")

	if clusterId != ""{
		db = db.Where("id = ?", clusterId)
	}

	if clusterName != ""{
		db = db.Where("cluster_name like '%" + clusterName + "%'")
	}

	if clusterType != ""{
		db = db.Where("cluster_type = ?", clusterType)
	}

	if clusterStatus != ""{
		db = db.Where("status = ?", clusterStatus)
	}

	if clusterTag != ""{
		db = db.Where("tags like '%," + clusterTag + ",%'")
	}

	err = db.Count(&total).Offset(offset).Limit(length).Find(&clusters).Error

	return
}

func CreateCluster(
		ClusterName 			string,
		DbPassword 				string,
		ClusterType 			string,
		ClusterVersion 			string,
		Tls 					bool,
		Tags           			string,
		OwnerId 				string,
		TenantId    			string,
	) (cluster *ClusterDO, err error){
	cluster = &ClusterDO{}
	cluster.ClusterName = ClusterName
	cluster.DbPassword = DbPassword
	cluster.ClusterType = ClusterType
	cluster.ClusterVersion = ClusterVersion
	cluster.Tls = Tls
	cluster.Tags = Tags
	cluster.OwnerId = OwnerId
	cluster.TenantId = TenantId

	err = MetaDB.Create(cluster).Error
	if err != nil {
		return
	}

	return
}

type BackupRecordDO struct {
	Record
	ClusterId   string		`gorm:"not null;type:varchar(36);default:null"`
	BackupRange int8
	BackupType  int8
	OperatorId  string		`gorm:"not null;type:varchar(36);default:null"`

	FilePath 		string
	FlowId			uint
}

func (d BackupRecordDO) TableName() string {
	return "backup_records"
}

type RecoverRecordDO struct {
	Record
	ClusterId 		string		`gorm:"not null;type:varchar(36);default:null"`

	OperatorId 		string		`gorm:"not null;type:varchar(36);default:null"`
	BackupRecordId  uint
	FlowId			uint
}

func (d RecoverRecordDO) TableName() string {
	return "recover_records"
}

type ParametersRecordDO struct {
	Record
	ClusterId 		string		`gorm:"not null;type:varchar(36);default:null"`

	OperatorId 		string		`gorm:"not null;type:varchar(36);default:null"`
	Content 		string		`gorm:"type:text"`
	FlowId 			uint
}

func (d ParametersRecordDO) TableName() string {
	return "parameters_records"
}

func SaveParameters(tenantId, clusterId, operatorId string, flowId uint, content string) (do *ParametersRecordDO, err error) {
	do = &ParametersRecordDO {
		Record: Record{
			TenantId: tenantId,
		},
		OperatorId: operatorId,
		ClusterId: clusterId,
		Content: content,
		FlowId: flowId,
	}

	err = MetaDB.Create(do).Error
	return
}

func GetCurrentParameters(clusterId string) (do *ParametersRecordDO, err error) {
	do = &ParametersRecordDO{}
	err = MetaDB.Where("cluster_id = ?", clusterId).Last(do).Error
	return
}

func DeleteBackupRecord(id uint) (record *BackupRecordDO, err error) {
	record = &BackupRecordDO{}
	err = MetaDB.Find(record, "id = ?", id).Error

	if err != nil {
		return
	}

	err = MetaDB.Delete(record).Error
	return
}

func SaveBackupRecord(tenantId, clusterId, operatorId string,
	backupRange, backupType int8, flowId uint,
	filePath string) (do *BackupRecordDO, err error){
	do = &BackupRecordDO{
		Record: Record{
			TenantId: tenantId,
		},
		ClusterId: clusterId,
		OperatorId: operatorId,
		BackupRange: backupRange,
		BackupType: backupType,
		FlowId: flowId,
		FilePath: filePath,
	}

	err = MetaDB.Create(do).Error
	return
}


type BackupRecordFetchResult struct {
	BackupRecordDO *BackupRecordDO
	Flow *FlowDO
}

func ListBackupRecords(clusterId string,
	offset, length int) (dos []*BackupRecordFetchResult, total int64, err error) {

	records := make([]*BackupRecordDO, length, length)
	err = MetaDB.Table("backup_records").
		Where("cluster_id = ?", clusterId).
		Count(&total).Order("id desc").Offset(offset).Limit(length).
		Find(&records).
		Error

	if err != nil {return}
	// query flows
	flowIds := make([]uint, len(records), len(records))

	dos = make([]*BackupRecordFetchResult, len(records), len(records))

	for i,r := range records {
		flowIds[i] = r.FlowId
		dos[i] = &BackupRecordFetchResult{
			BackupRecordDO: r,
		}
	}

	flows := make([]*FlowDO, len(records), len(records))
	err = MetaDB.Find(&flows, flowIds).Error
	if err != nil {return}

	flowMap := make(map[uint]*FlowDO)

	for _,v := range flows {
		flowMap[v.ID] = v
	}

	for i,v := range records {
		dos[i].BackupRecordDO = v
		dos[i].Flow = flowMap[v.FlowId]
	}

	return
}

func SaveRecoverRecord(tenantId, clusterId, operatorId string,
	backupRecordId uint,
	flowId uint) (do *RecoverRecordDO, err error) {
	do = &RecoverRecordDO{
		Record: Record{
			TenantId: tenantId,
		},
		ClusterId: clusterId,
		OperatorId: operatorId,
		FlowId: flowId,
		BackupRecordId: backupRecordId,
	}

	err = MetaDB.Create(do).Error
	return
}

