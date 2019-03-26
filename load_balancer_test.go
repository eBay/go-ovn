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
)

const LB1 = "lb1"

func TestLoadBalancer(t *testing.T) {
	t.Logf("Adding LB to OVN")
	ocmd, err := ovndbapi.LBAdd(LB1, "192.168.0.19:80", "tcp", []string{"10.0.0.11:80", "10.0.0.12:80"})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Adding LB OVN failed with err %v", err)
	}
	t.Logf("Adding LB to OVN Done")

	t.Logf("Updating LB to OVN")
	ocmd, err = ovndbapi.LBUpdate(LB1, "192.168.0.10:80", "tcp", []string{"10.10.10.127:8080", "10.10.10.120:8080"})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Updating LB OVN failed with err %v", err)
	}
	t.Logf("Updating LB to OVN done")

	t.Logf("Gettting LB by name")
	lb, err := ovndbapi.GetLB(LB1)
	if err != nil {
		t.Fatal(err)
	}
	if len(lb) != 1 {
		t.Fatalf("err getting lbs, total:%v", len(lb))
	}
	t.Logf("Lb found:%+v", lb[0])

	t.Logf("Deleting LB")
	ocmd, err = ovndbapi.LBDel(LB1)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}

	// Verify deletion
	lb, err = ovndbapi.GetLB(LB1)
	if err != nil {
		t.Fatal(err)
	}
	if len(lb) != 0 {
		t.Fatalf("error: lb deletion not done, total:%v", len(lb))
	}
	t.Logf("LB deletion done")
}
