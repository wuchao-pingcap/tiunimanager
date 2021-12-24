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

package gormreadwrite

import (
	mm "github.com/pingcap-inc/tiem/models/resource/management"
	resourcePool "github.com/pingcap-inc/tiem/models/resource/resourcepool"
	"os"
	"testing"

	"github.com/pingcap-inc/tiem/library/framework"
	"github.com/pingcap-inc/tiem/library/util/uuidutil"
	"github.com/pingcap-inc/tiem/models/resource"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var GormRW resource.ReaderWriter
var MetaDB *gorm.DB

func TestMain(m *testing.M) {
	testFilePath := "testdata/" + uuidutil.ShortId()
	os.MkdirAll(testFilePath, 0755)
	dbFile := testFilePath + "/test.db"

	defer func() {
		os.RemoveAll(testFilePath)
		os.Remove(testFilePath)
	}()

	framework.InitBaseFrameworkForUt(framework.MetaDBService,
		func(d *framework.BaseFramework) error {
			db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
			if err != nil || db.Error != nil {
				return err
			}
			MetaDB = db
			GormRW = NewGormResourceReadWrite(db)
			MetaDB.AutoMigrate(new(resourcePool.Host))
			MetaDB.AutoMigrate(new(resourcePool.Disk))
			MetaDB.AutoMigrate(new(resourcePool.Label))
			MetaDB.AutoMigrate(new(mm.UsedCompute))
			MetaDB.AutoMigrate(new(mm.UsedPort))
			MetaDB.AutoMigrate(new(mm.UsedDisk))
			return nil
		},
	)
	os.Exit(m.Run())
}
