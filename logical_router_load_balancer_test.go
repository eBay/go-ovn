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
	LB2 = "lb2"
	LR1 = "lr1"
)

func TestLRLoadBalancer(t *testing.T) {
	ovndbapi := getOVNClient(DBNB)
	// Add LR
	cmd, err := ovndbapi.LRAdd(LR1, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	// verify router create
	lrs, err := ovndbapi.LRGet(LR1)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrs) != 1 {
		t.Fatalf("lr not created %v", lrs)
	}
	// Add LB
	cmd, err = ovndbapi.LBAdd(LB2, "192.168.0.20:80", "tcp", []string{"10.0.0.21:80", "10.0.0.22:80"})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Adding LB OVN failed with err %v", err)
	}
	// Add LB to router
	t.Logf("Adding LB to LRouter %s", LR1)
	cmd, err = ovndbapi.LRLBAdd(LR1, LB2)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Adding LB lb2 to LRouter failed with err %v", err)
	}
	t.Logf("Adding LB lb2 to LRouter %s Done", LR1)
	// verify LB addition
	lbs, err := ovndbapi.LRLBList(LR1)
	if err != nil {
		t.Fatal(err)
	}
	if len(lbs) == 0 {
		t.Fatalf("lbs not created in %s", LR1)
	}
	assert.Equal(t, true, lbs[0].Name == LB2, "Added lb to lr")
	// Delete LB from router
	t.Logf("Delete LB from LRouter %s", LR1)
	cmd, err = ovndbapi.LRLBDel(LR1, LB2)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Deleting LB lb2 from LRouter failed with err %v", err)
	}
	t.Logf("Deleting LB lb2 to LRouter %s Done", LR1)
	// verify lb delete from lr
	lbs, err = ovndbapi.LRLBList(LR1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(lbs) == 0, "Deleted lb from lr")
	//Delete LB
	cmd, err = ovndbapi.LBDel(LB2)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}
	// Verify deletion
	lb, err := ovndbapi.LBGet(LB2)
	if err != nil {
		t.Fatal(err)
	}
	if len(lb) != 0 {
		t.Fatalf("error: lb deletion not done, total:%v", len(lb))
	}
	t.Logf("LB lb2 deleted")
	// Delete router
	cmd, err = ovndbapi.LRDel(LR1)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Router %s deleted", LR1)

	// verify lb list for non-existing routers
	_, err = ovndbapi.LRLBList(FAKENOROUTER)
	if err != nil {
		assert.EqualError(t, ErrorNotFound, err.Error())
	}
}
