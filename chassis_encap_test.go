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
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestEncaps(t *testing.T) {
	db = getOVNDB()
	if db == dbNB {
		t.Skip("Skip running chassis test against nb db")
	}
	t.Logf("Adding Chassis to OVN SB DB")
	ocmd, err := ovndbapi.ChassisAdd(CHASSIS_NAME, CHASSIS_HOSTNAME, ENCAP_TYPES, IP, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Adding Chassis to OVN failed with err %v", err)
	}
	t.Logf("Adding Chassis to OVN Done")

	t.Logf("Gettting Chassis by hostname")
	chassis, err := ovndbapi.ChassisGet(CHASSIS_HOSTNAME)
	if err != nil {
		t.Fatal(err)
	}
	if len(chassis) != 1 {
		t.Fatalf("err getting chassis, total:%v", len(chassis))
	}
	chName := chassis[0].Name
	ChEnCaps := chassis[0].Encaps
	t.Logf("Chassis found:%+v", chName)

	// Verify if encaps are created
	t.Logf("Gettting Encaps by chassis name")
	encaps, err := ovndbapi.EncapList(CHASSIS_NAME)
	if err != nil {
		t.Fatal(err)
	}
	if len(encaps) != 3 {
		t.Fatalf("err getting encaps, total:%v", len(chassis))
	}

	assert.Equal(t, true, len(encaps) == 3, "Encaps exist")

	// verify if chassis have the encaps too
	var encapsUid []string
	for _, v := range encaps {
		encapsUid = append(encapsUid, v.UUID)
	}
	assert.Equal(t, true, reflect.DeepEqual(encapsUid, ChEnCaps), "Same encaps exists in chassis too")

	t.Logf("Deleting Chassis: %v", chName)
	ocmd, err = ovndbapi.ChassisDel(chName)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}

	// Verify deletion
	chassis, err = ovndbapi.ChassisGet(CHASSIS_HOSTNAME)
	if err != nil {
		t.Fatal(err)
	}
	if len(chassis) != 0 {
		t.Fatalf("error: Chassis deletion not done, total:%v", len(chassis))
	}
	t.Logf("Chassis %s deleted", chName)
}
