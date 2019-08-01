package goovn

import (
	"testing"
)

var (
	lrTestLR string
	lrTestLB string
)

func TestLogicalRouterAdd(t *testing.T) {
	lrUUID := newUUID(t)

	lrTestLR = "test" + lrUUID
	cmd, err := ovndbapi.LogicalRouter.Add(LogicalRouterName(lrTestLR))
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogicalRouterGet(t *testing.T) {
	lr, err := ovndbapi.LogicalRouter.Get(LogicalRouterName(lrTestLR))
	if err != nil {
		t.Fatal(err)
	}
	if lr.Name != lrTestLR {
		t.Fatal("test lr not found")
	}
}

func TestLogicalRouterList(t *testing.T) {
	lrs, err := ovndbapi.LogicalRouter.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(lrs) == 0 {
		t.Fatal("no logical router found")
	}
	var found bool
	for _, lr := range lrs {
		if lr.Name == lrTestLR {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("no logical router found")
	}
}

func TestLogicalRouterLBAdd(t *testing.T) {
	lbUUID := newUUID(t)
	lrTestLB = "test" + lbUUID

	cmd, err := ovndbapi.LoadBalancer.Add(
		LoadBalancerName(lrTestLB),
		LoadBalancerVIP("192.168.0.20:80"),
		LoadBalancerProtocol("tcp"),
		LoadBalancerIP([]string{"10.10.10.21:80", "10.10.10.22:80"}),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.LogicalRouter.LBAdd(
		LogicalRouterName(lrTestLR),
		LoadBalancerName(lrTestLB),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lr, err := ovndbapi.LogicalRouter.Get(LogicalRouterName(lrTestLR))
	if err != nil {
		t.Fatal(err)
	}
	if len(lr.LoadBalancer) != 1 {
		t.Fatalf("load balancer not added to logical router: %v", lr)
	}
}

func TestLogicalRouterLBDel(t *testing.T) {
	cmd, err := ovndbapi.LogicalRouter.LBDel(
		LogicalRouterName(lrTestLR),
		LoadBalancerName(lrTestLB),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lr, err := ovndbapi.LogicalRouter.Get(LogicalRouterName(lrTestLR))
	if err != nil {
		t.Fatal(err)
	}
	if len(lr.LoadBalancer) != 0 {
		t.Fatalf("load balancer not deleted from logical router: %v", lr)
	}
}

func TestLogicalRouterDel(t *testing.T) {
	cmd, err := ovndbapi.LogicalRouter.Del(LogicalRouterName(lrTestLR))
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	lr, err := ovndbapi.LogicalRouter.Get(LogicalRouterName(lrTestLR))
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}

	if lr != nil {
		t.Fatal("test lr not found")
	}
}

/*
	cmd, err = ovndbapi.LRPAdd(LR, LRP, "54:54:54:54:54:54", []string{"192.168.0.1/24"}, "lrp2", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.LRSRAdd(LR, IPPREFIX, NEXTHOP, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lr, err := ovndbapi.LRGet(LR)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range lr {
		if len(v.LoadBalancer) != 1 {
			t.Fatal("Get Loadblancer Fail")
		}
		if len(v.Ports) != 1 {
			t.Fatal("Get Ports Fail")
		}
		if len(v.StaticRoutes) != 1 {
			t.Fatal("Get StaticRouter Fail")
		}
	}

	cmd, err = ovndbapi.LBDel(LB2)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.LRDel(LR)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

}
*/
