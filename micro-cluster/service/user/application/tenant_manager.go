
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

package application

import (
	"fmt"
	"github.com/pingcap-inc/tiem/micro-cluster/service/user/domain"
	"github.com/pingcap-inc/tiem/micro-cluster/service/user/ports"
)

type TenantManager struct {
	tenantRepo ports.TenantRepository
}

func NewTenantManager(tenantRepo ports.TenantRepository) *TenantManager {
	return &TenantManager{tenantRepo : tenantRepo}
}

// CreateTenant 创建租户
func (p *TenantManager) CreateTenant(name string) (*domain.Tenant, error) {
	existed, e := p.FindTenant(name)

	if e == nil && existed != nil {
		return existed, fmt.Errorf("tenant already exist")
	}

	tenant := domain.Tenant{Name: name, Type: domain.InstanceWorkspace, Status: domain.Valid}

	return &tenant, nil
}

// FindTenant 查找租户
func (p *TenantManager) FindTenant(name string) (*domain.Tenant, error) {
	tenant, err := p.tenantRepo.LoadTenantByName(name)
	return &tenant, err
}

func (p *TenantManager) FindTenantById(id string) (*domain.Tenant, error) {
	tenant, err := p.tenantRepo.LoadTenantById(id)
	return &tenant, err
}