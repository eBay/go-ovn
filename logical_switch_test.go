/**
 * Copyright (c) 2019 eBay Inc.
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

var (
	lsTestLS string
	lsTestLB string
)

func TestLogicalSwitchAdd(t *testing.T) {
	lsUUID := newUUID(t)

	lsTestLS = "test" + lsUUID
	cmd, err := ovndbapi.LogicalSwitch.Add(LogicalSwitchName(lsTestLS))
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogicalSwitchList(t *testing.T) {
	lsList, err := ovndbapi.LogicalSwitch.List()
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, ls := range lsList {
		if ls.Name == lsTestLS {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("logical switch add fail")
	}
}

func TestLogicalSwitchGet(t *testing.T) {
	ls, err := ovndbapi.LogicalSwitch.Get(LogicalSwitchName(lsTestLS))
	if err != nil {
		t.Fatal(err)
	}
	if ls == nil || ls.Name != lsTestLS {
		t.Fatalf("logical switch get fail: %s not found: %v", lsTestLS, ls)
	}
}

func TestLogicalSwitchSetExternalIDs(t *testing.T) {
	cmd, err := ovndbapi.LogicalSwitch.SetExternalIDs(
		LogicalSwitchName(lsTestLS),
		LogicalSwitchExternalIDs(map[string]string{"test": "true", "foo": "bar"}))
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	ls, err := ovndbapi.LogicalSwitch.Get(LogicalSwitchName(lsTestLS))
	if err != nil {
		t.Fatal(err)
	}
	for key, val := range ls.ExternalIDs {
		if key == "test" {
			assert.Equal(t, true, val == "true", "external_ids test")
		}
	}
	cmd, err = ovndbapi.LogicalSwitch.SetExternalIDs(LogicalSwitchName(lsTestLS))
	if err != nil {
		assert.Errorf(t, err, "cannot update %s with empty external_ids", lsTestLS)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

}

func TestLogicalSwitchLBAdd(t *testing.T) {
	lbUUID := newUUID(t)

	lsTestLB = "test" + lbUUID
	// alternative can be specified via map[string]string{"192.168.0.19:80":"10.0.0.11:80,10.0.0.12:80"}
	ocmd, err := ovndbapi.LoadBalancer.Add(
		LoadBalancerName(lsTestLB),
		LoadBalancerVIP("192.192.192.192:80"),
		LoadBalancerProtocol("tcp"),
		LoadBalancerIP([]string{"10.0.0.11:80", "10.0.0.12:80"}))
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatal(err)
	}

	// alternative can be specified via map[string]string{"192.168.0.19:80":"10.0.0.11:80,10.0.0.12:80"}
	ocmd, err = ovndbapi.LogicalSwitch.LBAdd(
		LogicalSwitchName(lsTestLS),
		LoadBalancerName(lsTestLB),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogicalSwitchLBList(t *testing.T) {
	lbList, err := ovndbapi.LogicalSwitch.LBList(
		LogicalSwitchName(lsTestLS),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lbList) == 0 {
		t.Fatalf("lb list empty")
	}

	var found bool
	for _, lb := range lbList {
		if lb.Name == lsTestLB {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("logical switch load balancer list fail")
	}
}

func TestLogicalSwitchLBDel(t *testing.T) {
	// alternative can be specified via map[string]string{"192.168.0.
	ocmd, err := ovndbapi.LogicalSwitch.LBDel(
		LogicalSwitchName(lsTestLS),
		LoadBalancerName(lsTestLB),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatal(err)
	}

	ls, err := ovndbapi.LogicalSwitch.Get(LogicalSwitchName(lsTestLS))
	if err != nil {
		t.Fatal(err)
	}
	if len(ls.LoadBalancer) != 0 {
		t.Fatal("load balancer not deleted from logical switch")
	}
}

func TestLogicalSwitchDelExternalIDs(t *testing.T) {
	cmd, err := ovndbapi.LogicalSwitch.DelExternalIDs(
		LogicalSwitchName(lsTestLS),
		LogicalSwitchExternalIDs(map[string]string{"test": "true"}))
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	ls, err := ovndbapi.LogicalSwitch.Get(LogicalSwitchName(lsTestLS))
	if err != nil {
		t.Fatal(err)
	}
	for key, val := range ls.ExternalIDs {
		if key == "foo" {
			assert.Equal(t, true, val == "bar", "remove external_ids test:true")
		}
	}
	cmd, err = ovndbapi.LogicalSwitch.DelExternalIDs(LogicalSwitchName(lsTestLS))
	if err != nil {
		assert.Errorf(t, err, "unable to delete empty external_ids")
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogicalSwitchDel(t *testing.T) {
	cmd, err := ovndbapi.LogicalSwitch.Del(LogicalSwitchName(lsTestLS))
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	lsList, err := ovndbapi.LogicalSwitch.List()
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, ls := range lsList {
		if ls.Name == lsTestLS {
			found = true
			break
		}
	}

	if found {
		t.Fatalf("logical switch %s not deleted", lsTestLS)
	}
}
