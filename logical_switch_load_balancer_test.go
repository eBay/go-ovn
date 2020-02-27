/**
 * Copyright (c) 2017 eBay Inc.
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
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	LB3  = "lb3"
	LSW1 = "LSW1"
)

func TestLSLoadBalancer(t *testing.T) {
	ovndbapi := getOVNClient(DBNB)
	// create Switch
	t.Logf("Adding  %s to OVN", LSW1)
	cmd, err := ovndbapi.LSAdd(LSW1)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	ls, err := ovndbapi.LSGet(LSW1)
	if err != nil {
		t.Fatal(err)
	}
	if ls[0].Name != LSW1 {
		t.Fatalf("ls not created %v", LSW1)
	}
	// Create LB LB3
	t.Logf("Adding LB %s to OVN", LB3)
	cmd, err = ovndbapi.LBAdd(LB3, "192.168.0.21:80", "tcp", []string{"10.0.0.21:80", "10.0.0.22:80"})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Adding LB OVN failed with err %v", err)
	}
	// Add lb  to switch
	t.Logf("Adding LB LB3 to LSW1itches %s", LSW1)
	cmd, err = ovndbapi.LSLBAdd(LSW1, LB3)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Adding LB LB3 to LSW1itch failed with err %v", err)
	}
	t.Logf("Adding LB LB3 to LSW1itch %s Done", LSW1)
	// verify if lb addition
	lbs, err := ovndbapi.LSLBList(LSW1)
	if err != nil {
		t.Fatal(err)
	}
	if len(lbs) == 0 {
		t.Fatalf("lbs not created in %s", LSW1)
	}
	assert.Equal(t, true, lbs[0].Name == LB3, "Added lb to lswitch")
	// delete lb from switch
	t.Logf("Deleting LB LB3 to LSW1itches %s", LSW1)
	cmd, err = ovndbapi.LSLBDel(LSW1, LB3)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Deleting LB LB3 to LSW1itch failed with err %v", err)
	}
	t.Logf("Deleting LB LB3 to LSW1itch %s Done", LSW1)
	// verify lb deletion from switch
	lbs, err = ovndbapi.LSLBList(LSW1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(lbs) == 0, "Deleted lb from lswitch")
	// Delete LB
	t.Logf("Deleting LB")
	cmd, err = ovndbapi.LBDel(LB3)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}
	// Finally delete Switch
	t.Logf("Deleting LS")
	cmd, err = ovndbapi.LSDel(LSW1)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}

}
