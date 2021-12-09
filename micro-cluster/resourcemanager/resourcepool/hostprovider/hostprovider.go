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

package hostprovider

import (
	"context"

	"github.com/pingcap-inc/tiem/common/structs"
)

type HostProvider interface {
	ImportHosts(ctx context.Context, hosts []structs.HostInfo) (hostIds []string, err error)
	DeleteHosts(ctx context.Context, hostIds []string) (err error)
	QueryHosts(ctx context.Context, filter *structs.HostFilter, page *structs.PageRequest) (hosts []structs.HostInfo, err error)
	UpdateHostStatus(ctx context.Context, hostId []string, status string) (err error)
	UpdateHostReserved(ctx context.Context, hostId []string, reserved bool) (err error)

	GetHierarchy(ctx context.Context, filter *structs.HostFilter, level int, depth int) (root *structs.HierarchyTreeNode, err error)
	GetStocks(ctx context.Context, location structs.Location, hostFilter structs.HostFilter, diskFilter structs.DiskFilter) (*structs.Stocks, error)
}
