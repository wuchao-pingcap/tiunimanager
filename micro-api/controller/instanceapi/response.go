package instanceapi

type InstanceInfo struct {
	InstanceId 		string 	`json:"instanceId"`
	InstanceName 	string 	`json:"instanceName"`
	InstanceStatus 	int 	`json:"instanceStatus"`
	InstanceVersion string 	`json:"instanceVersion"`
}