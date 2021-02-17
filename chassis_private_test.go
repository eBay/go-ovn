/**
 * Copyright (c) 2020 eBay Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 **/

package goovn

import (
	"fmt"
	"testing"
)

func TestChassisPrivate(t *testing.T) {
	ovndbapi := getOVNClient(DBSB)
	t.Logf("Adding row to OVN SB DB chassis_private table")
	ovn, ok := ovndbapi.(*ovndb)
	if !ok {
		t.Fatal(fmt.Errorf("Invalid type assertion"))
	}
	ocmd, err := ovn.chassisPrivateAdd(CHASSIS_NAME, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("err adding row to ovn-sb chassis_private table failed: %v", err)
	}
	t.Logf("Adding row to ovn-sb chassis_private table Done")

	t.Logf("Listing all the chassis in ovn-sb chassis_private table")
	chassisPrivate, err := ovndbapi.ChassisPrivateList()
	if err != nil {
		t.Fatal(err)
	}
	if len(chassisPrivate) != 1 {
		t.Fatalf("err listing chassis in ovn-sb chassis_private table, total:%v",
			len(chassisPrivate))
	}
	t.Logf("Single chassis found in ovn-sb chassis_private table: %+v", chassisPrivate[0].Name)

	t.Logf("Gettting chassis by name in chassis_private table")
	chassisPrivate, err = ovndbapi.ChassisPrivateGet(CHASSIS_NAME)
	if err != nil {
		t.Fatal(err)
	}
	if len(chassisPrivate) != 1 {
		t.Fatalf("err getting chassis in ovn-sb chassis_private table, total:%v",
			len(chassisPrivate))
	}
	chName := chassisPrivate[0].Name
	t.Logf("Chassis found in chassis-private table:%+v", chName)

	// deleting the row in chassis_private table
	t.Logf("Deleting row in OVN SB DB chassis_private table:%v", chName)
	ocmd, err = ovndbapi.ChassisPrivateDel(chName)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("err executing command %v: %v", ocmd, err)
	}
}
