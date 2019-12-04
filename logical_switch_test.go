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

const (
	LS3             = "LS3"
	NEUTRON_NETWORK = "neutron:network"
	DUMMY           = "dummy"
	FOO             = "foo"
	BAR             = "bar"
	LS5             = "LS5"
	LR5             = "LR5"
	LSP5            = "ls5-lr5"
	LRP5            = "lr5-ls5"
)

func TestLSwitchExtIds(t *testing.T) {
	// create Switch
	t.Logf("Adding  %s to OVN", LS3)
	cmd, err := ovndbapi.LSAdd(LS3)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	ls, err := ovndbapi.LSGet(LS3)
	if err != nil {
		t.Fatal(err)
	}
	if ls[0].Name != LS3 {
		t.Fatalf("ls not created %v", LS3)
	}
	// Add external_id to LS3
	cmd, err = ovndbapi.LSExtIdsAdd(LS3, map[string]string{NEUTRON_NETWORK: DUMMY, FOO: BAR})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	// Get LS3 and get external_id NEUTRON_NETWORK
	lswitch, err := ovndbapi.LSGet(LS3)
	if err != nil {
		t.Fatal(err)
	}
	externalIDs := lswitch[0].ExternalID
	for key, val := range externalIDs {
		if key == NEUTRON_NETWORK {
			assert.Equal(t, true, val.(string) == DUMMY, "Got external ID dummy")
			t.Logf("Successfully validated external_id key NEUTRON_NETWORK to LS3")
		}
	}
	// Add empty external_ids to LS3
	cmd, err = ovndbapi.LSExtIdsAdd(LS3, nil)
	if err != nil {
		assert.Errorf(t, err, "Cannot update lswitch with empty ext_id")
		t.Logf("Adding empty external_id for LS3 validation is ok")
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	//delete external_id from LS3
	cmd, err = ovndbapi.LSExtIdsDel(LS3, map[string]string{"neutron:network": "dummy"})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	// Get LS3 and get external_id
	lswitch, err = ovndbapi.LSGet(LS3)
	if err != nil {
		t.Fatal(err)
	}
	externalIDs = lswitch[0].ExternalID
	for key, val := range externalIDs {
		if key == FOO {
			assert.Equal(t, true, val.(string) == BAR, "Externel id with value dummy deleted")
			t.Logf("Deleted external_id key NEUTRON_NETWORK from LS3")
		}
	}
	// Delete empty external_ids from LS3
	cmd, err = ovndbapi.LSExtIdsDel(LS3, nil)
	if err != nil {
		assert.Errorf(t, err, "Cannot update lswitch with empty ext_id")
		t.Logf("Deleting empty external_id from LS3 validation is ok")
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	// Finally delete Switch
	t.Logf("Deleting LS3")
	cmd, err = ovndbapi.LSDel(LS3)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}

}

func TestLinkSwitchToRouter(t *testing.T) {
	// create Switch
	t.Logf("Adding %s to OVN", LS5)
	cmd, err := ovndbapi.LSAdd(LS5)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
		return
	}
	cmd, err = ovndbapi.LRAdd(LR5, nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
		return
	}
	lrpMac := "12:34:56:78:90:ab"
	cmd, err = ovndbapi.LinkSwitchToRouter(LS5, LSP5, LR5, LRP5, lrpMac, []string{"10.10.10.0/24"}, nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
		return
	}
	ports, err := ovndbapi.LRPList(LR5)
	if err != nil {
		t.Fatal(err)
		return
	}
	found := false
	for _, port := range ports {
		if port.Name == LRP5 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("logical router port %s wasn't created", LRP5)
	}
}
