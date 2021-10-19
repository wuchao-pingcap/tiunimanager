package ports

import (
	"github.com/pingcap-inc/tiem/micro-cluster/service/user/domain"
)

type TenantRepository interface {
	AddTenant(*domain.Tenant) error

	LoadTenantByName(name string)  (domain.Tenant, error)

	LoadTenantById(id string)  (domain.Tenant, error)
}

type RbacRepository interface {

	AddAccount(a *domain.Account) error

	LoadAccountByName(name string) (domain.Account, error)

	LoadAccountAggregation(name string) (domain.AccountAggregation, error)

	LoadAccountById(id string) (domain.Account, error)

	AddRole(r *domain.Role) error

	LoadRole(tenantId string, name string) (domain.Role, error)

	AddPermission(r *domain.Permission) error

	LoadPermissionAggregation(tenantId string, code string) (domain.PermissionAggregation, error)

	LoadPermission(tenantId string, code string) (domain.Permission, error)

	LoadAllRolesByAccount(account *domain.Account) ([]domain.Role, error)

	LoadAllRolesByPermission(permission *domain.Permission) ([]domain.Role, error)

	AddPermissionBindings(bindings []domain.PermissionBinding) error

	AddRoleBindings(bindings []domain.RoleBinding) error
}

type TokenHandler interface {

	// Provide 提供一个有效的token
	Provide (tiEMToken *domain.TiEMToken) (string, error)

	// Modify 修改token
	Modify (tiEMToken *domain.TiEMToken) error

	// GetToken 获取一个token
	GetToken(tokenString string) (domain.TiEMToken, error)
}
