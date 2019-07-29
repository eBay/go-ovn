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
)

const LB1 = "lb1"

func TestLoadBalancerAdd(t *testing.T) {
	t.Logf("Add LoadBalancer")
	// alternative can be specified via map[string]string{"192.168.0.19:80":"10.0.0.11:80,10.0.0.12:80"}
	ocmd, err := ovndbapi.LoadBalancer.Add(
		LoadBalancerName(LB1),
		LoadBalancerVIP("192.168.0.19:80"),
		LoadBalancerProtocol("tcp"),
		LoadBalancerIP([]string{"10.0.0.11:80", "10.0.0.12:80"}))
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Adding LB OVN failed with err %v", err)
	}

	t.Logf("Add LoadBalancer new VIP")
	ocmd, err = ovndbapi.LoadBalancer.Add(
		LoadBalancerName(LB1),
		LoadBalancerVIP("192.168.0.10:80"),
		LoadBalancerProtocol("tcp"),
		LoadBalancerIP([]string{"10.10.10.127:8080", "10.10.10.120:8080"}))
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Updating LB OVN failed with err %v", err)
	}

	t.Logf("Get LoadBalancer")
	lb, err := ovndbapi.LoadBalancer.Get(LoadBalancerName(LB1))
	if err != nil {
		t.Fatal(err)
	}
	if lb.Name != LB1 {
		t.Fatalf("no load balancer: %v\n", lb)
	}

	t.Logf("Del LoadBalancer")
	ocmd, err = ovndbapi.LoadBalancer.Del(LoadBalancerName(LB1))
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}

	lb, err = ovndbapi.LoadBalancer.Get(LoadBalancerName(LB1))
	if err == nil {
		t.Fatalf("load balancer not deleted: %v", lb)
	}
}
