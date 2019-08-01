package goovn

import "testing"

var (
	lrTestLR string
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
	lrs, err := ovndbapi.LRList()
	if err != nil {
		t.Fatal(err)
	}
	if len(lrs) != 1 {
		t.Fatalf("lr not created %v", lrs)
	}

	cmd, err = ovndbapi.LRDel(LR)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lrs, err = ovndbapi.LRList()
	if err != nil {
		t.Fatal(err)
	}
	if len(lrs) != 0 {
		t.Fatalf("lr not deleted %v", lrs)
	}

	cmd, err = ovndbapi.LRAdd(LR, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.LBAdd(LB2, "192.168.0.20:80", "tcp", []string{"10.0.0.21:80", "10.0.0.22:80"})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.LRLBAdd(LR, LB2)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

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
