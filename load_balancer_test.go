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

var (
	lbTestLB string
)

func TestLoadBalancerAdd(t *testing.T) {
	lbUUID := newUUID(t)

	lbTestLB = "test" + lbUUID
	// alternative can be specified via map[string]string{"192.168.0.19:80":"10.0.0.11:80,10.0.0.12:80"}
	ocmd, err := ovndbapi.LoadBalancer.Add(
		LoadBalancerName(lbTestLB),
		LoadBalancerVIP("192.168.0.19:80"),
		LoadBalancerProtocol("tcp"),
		LoadBalancerIP([]string{"10.0.0.11:80", "10.0.0.12:80"}))
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerGet(t *testing.T) {
	lb, err := ovndbapi.LoadBalancer.Get(LoadBalancerName(lbTestLB))
	if err != nil {
		t.Fatal(err)
	}
	if lb.Name != lbTestLB {
		t.Fatalf("failed to get load balancer: %v", lb)
	}
}

func TestLoadBalancerAddVIP(t *testing.T) {
	ocmd, err := ovndbapi.LoadBalancer.Set(
		LoadBalancerName(lbTestLB),
		LoadBalancerVIP("192.168.0.10:80"),
		LoadBalancerProtocol("tcp"),
		LoadBalancerIP([]string{"10.10.10.127:8080", "10.10.10.120:8080"}))
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("update vip in lb failed %v", err)
	}

	lb, err := ovndbapi.LoadBalancer.Get(
		LoadBalancerName(lbTestLB),
	)
	if err != nil {
		t.Fatal(err)
	}
	if lb.Name != lbTestLB {
		t.Fatalf("no load balancer %s found: %v", lbTestLB, lb)
	}
}

func TestLoadBalancerDel(t *testing.T) {
	ocmd, err := ovndbapi.LoadBalancer.Del(
		LoadBalancerName(lbTestLB),
		LoadBalancerVIP("192.168.0.19:80"),
		LoadBalancerProtocol("tcp"),
		LoadBalancerIP([]string{"10.0.0.11:80", "10.0.0.12:80"}),
	)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}

	ocmd, err = ovndbapi.LoadBalancer.Del(
		LoadBalancerName(lbTestLB),
		LoadBalancerVIP("192.168.0.10:80"),
		LoadBalancerProtocol("tcp"),
		LoadBalancerIP([]string{"10.10.10.127:8080", "10.10.10.120:8080"}),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatal(err)
	}

	lb, err := ovndbapi.LoadBalancer.Get(LoadBalancerName(lbTestLB))
	if err == nil {
		t.Fatalf("load balancer not deleted: %v", lb)
	}
}
