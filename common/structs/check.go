/******************************************************************************
 * Copyright (c)  2022 PingCAP, Inc.                                          *
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

package structs

import "github.com/pingcap-inc/tiem/common/constants"

type CheckString struct {
	Valid         bool   `json:"valid"`
	RealValue     string `json:"realValue"`
	ExpectedValue string `json:"expectedValue"`
}

type CheckInt32 struct {
	Valid         bool  `json:"valid"`
	RealValue     int32 `json:"realValue"`
	ExpectedValue int32 `json:"expectedValue"`
}

type CheckBool struct {
	Valid         bool `json:"valid"`
	RealValue     bool `json:"realValue"`
	ExpectedValue bool `json:"expectedValue"`
}

type CheckAny struct {
	Valid         bool        `json:"valid"`
	RealValue     interface{} `json:"realValue"`
	ExpectedValue interface{} `json:"expectedValue"`
}

type CheckStatus struct {
	Health  bool   `json:"health"`
	Message string `json:"message"`
}

type TenantCheck struct {
	ClusterCount int32          `json:"clusterCount"`
	CPURatio     float32        `json:"cpuRatio"`
	MemoryRatio  float32        `json:"memoryRatio"`
	StorageRatio float32        `json:"storageRatio"`
	Clusters     []ClusterCheck `json:"clusters"`
}

type ClusterCheck struct {
	ID                string                             `json:"clusterID"`
	MaintenanceStatus constants.ClusterMaintenanceStatus `json:"maintenanceStatus"`
	ConnectionCount   int32                              `json:"connectionCount"`
	CPU               int32                              `json:"cpu"`
	Memory            int32                              `json:"memory"`
	Storage           int32                              `json:"storage"`
	Copies            CheckInt32                         `json:"copies"`
	TLS               CheckBool                          `json:"tls"`
	Versions          map[string]CheckString             `json:"versions"`
	AccountStatus     CheckStatus                        `json:"accountStatus"`
	Topology          CheckString                        `json:"topology"`
	Parameters        map[string]CheckAny                `json:"parameters"`
	RegionStatus      CheckStatus                        `json:"regionStatus"`
	Instances         []InstanceCheck                    `json:"instances"`
	BackupStrategy    CheckString                        `json:"backupStrategy"`
	BackupRecordValid map[string]bool                    `json:"backupRecordValid"`
}

type InstanceCheck struct {
	ID     string      `json:"instanceID"`
	Status CheckStatus `json:"status"`
}

type CheckKey struct {
	Source, Target string
}

type CheckRangeInt32 struct {
	Valid         bool    `json:"valid"`
	RealValue     int32   `json:"realValue"`
	ExpectedRange []int32 `json:"expectedRange"`	// Left closed right closed interval
}
type HostsCheck struct {
	NTP                 map[CheckKey]CheckRangeInt32 `json:"ntp"`
	TimeZoneConsistency map[CheckKey]bool            `json:"timeZoneConsistency"`
	Ping                map[CheckKey]bool            `json:"ping"`
	Delay               map[CheckKey]CheckRangeInt32 `json:"delay"`
	Hosts               map[string]HostCheck         `json:"hosts"`
}

type CheckSwitch struct {
	Enable bool `json:"enable"`
}

type CheckError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type HostCheck struct {
	OSVersion    CheckString  `json:"version"`
	SELinux      CheckSwitch  `json:"selinux"`
	Firewall     CheckSwitch  `json:"firewall"`
	Swap         CheckSwitch  `json:"swap"`
	Memory       CheckInt32   `json:"memory"`
	CPU          CheckInt32   `json:"cpu"`
	Disk         CheckInt32   `json:"disk"`
	StorageRatio float32      `json:"storageRatio"`
	Errors       []CheckError `json:"errors"`
}
