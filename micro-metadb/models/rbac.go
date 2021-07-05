package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model

	TenantId 		uint		`gorm:"size:255"`
	Name 			string		`gorm:"size:255"`
	Salt 			string		`gorm:"size:255"`
	FinalHash 		string		`gorm:"size:255"`
	Status 			int8		`gorm:"size:255"`
}

type Role struct {
	gorm.Model

	TenantId 		uint		`gorm:"size:255"`
	Name    		string		`gorm:"size:255"`
	Desc    		string		`gorm:"size:255"`
	Status 			int8		`gorm:"size:255"`
}

type PermissionBinding struct {
	gorm.Model

	TenantId 		uint		`gorm:"size:255"`
	RoleId 			uint		`gorm:"size:255"`
	PermissionId	uint		`gorm:"size:255"`
	Status 			int8		`gorm:"size:255"`
}

type RoleBinding struct {
	gorm.Model

	TenantId 		uint		`gorm:"size:255"`
	RoleId 			uint		`gorm:"size:255"`
	AccountId	 	uint		`gorm:"size:255"`
	Status 			int8		`gorm:"size:255"`
}

type Permission struct {
	gorm.Model

	TenantId 		uint		`gorm:"size:255"`
	Code   			string		`gorm:"size:255"`
	Name  	 		string		`gorm:"size:255"`
	Type   			int8		`gorm:"size:255"`
	Desc   			string		`gorm:"size:255"`
	Status 			int8		`gorm:"size:255"`
}

func AddAccount(tenantId uint, name string, salt string, finalHash string, status int8) (result Account, err error) {
	result.TenantId = tenantId
	result.Name = name
	result.Salt = salt
	result.FinalHash = finalHash
	result.Status = status

	MetaDB.Create(&result)
	return
}

func FindAccount(name string) (result Account, err error) {
	MetaDB.Where(&Account{Name: name}).First(&result)
	return
}

func AddRole(tenantId uint, name string, desc string, status int8) (result Role, err error) {
	result.TenantId = tenantId
	result.Name = name
	result.Desc = desc
	result.Status = status

	MetaDB.Create(&result)
	return
}

func FetchRole(tenantId uint, name string) (result Role, err error) {
	MetaDB.Where(&Role{TenantId: tenantId, Name: name}).First(&result)
	return
}

func AddPermission(tenantId uint, code, name, desc string, permissionType, status int8) (result Permission, err error) {
	result.TenantId = tenantId
	result.Code = code
	result.Name = name
	result.Desc = desc
	result.Type = permissionType
	result.Status = status

	MetaDB.Create(&result)
	return
}

func FetchPermission(tenantId uint, code string) (result Permission, err error) {
	MetaDB.Where(&Permission{TenantId: tenantId, Code: code}).First(&result)
	return
}

func FetchAllRolesByAccount(tenantId uint, accountId uint) (result []Role, err error) {
	//service.DB.Where("account_id = ?", accountId).Limit(50).Find(&result)

	var roleBinds []RoleBinding
	MetaDB.Where("tenant_id = ? and account_id = ? and status = 0", tenantId, accountId).Limit(50).Find(&roleBinds)

	var roleIds []uint
	for _, v := range roleBinds {
		roleIds = append(roleIds, v.RoleId)
	}

	result, err = FetchRolesByIds(roleIds)
	return
}

func FetchRolesByIds(roleIds []uint) (result []Role, err error){
	MetaDB.Where("id in ?", roleIds).Find(&result)
	return
}

func FetchAllRolesByPermission(tenantId uint, permissionId uint) (result []Role, err error) {
	var permissionBindings []PermissionBinding
	MetaDB.Where("tenant_id = ? and permission_id = ? and status = 0", tenantId, permissionId).Limit(50).Find(&permissionBindings)

	var roleIds []uint
	for _, v := range permissionBindings {
		roleIds = append(roleIds, v.RoleId)
	}

	result, err = FetchRolesByIds(roleIds)
	return
}

func AddPermissionBindings(bindings []PermissionBinding) error {
	MetaDB.Create(&bindings)
	return nil
}

func AddRoleBindings(bindings []RoleBinding) error{
	MetaDB.Create(&bindings)
	return nil
}