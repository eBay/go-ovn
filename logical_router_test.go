package goovn

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const LB4 = "lb4"

func TestLogicalRouter(t *testing.T) {
	ovndbapi := getOVNClient(DBNB)
	var cmds []*OvnCommand
	var cmd *OvnCommand
	var err error

	cmds = make([]*OvnCommand, 0)
	cmd, err = ovndbapi.LRAdd(LR, nil)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)
	err = ovndbapi.Execute(cmds...)
	if err != nil {
		t.Fatal(err)
	}

	lrs, err := ovndbapi.LRList()
	if err != nil {
		t.Fatal(err)
	}
	if len(lrs) == 0 {
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
	if len(lrs) > 0 {
		for _, lr := range lrs {
			if lr.Name == LR {
				t.Fatalf("lr not deleted %v", LR)
				break
			}
		}
	} else {
		t.Logf("Successfully deleted router %s", LR)
	}

	cmd, err = ovndbapi.LRAdd(LR, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.LBAdd(LB4, "192.168.0.20:80", "tcp", []string{"10.0.0.21:80", "10.0.0.22:80"})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.LRLBAdd(LR, LB4)
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

	cmd, err = ovndbapi.LRSRAdd(LR, IPPREFIX, NEXTHOP, "", "", nil)
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
	// Delete LB from router
	t.Logf("Delete LB from LRouter %s", LR)
	cmd, err = ovndbapi.LRLBDel(LR, LB4)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Deleting LB lb2 from LRouter failed with err %v", err)
	}
	t.Logf("Deleting LB lb2 to LRouter %s Done", LR)
	// verify lb delete from lr
	lbs, err := ovndbapi.LRLBList(LR)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(lbs) == 0, "Deleted lb from lr")
	//Delete LB
	cmd, err = ovndbapi.LBDel(LB4)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}
	// Verify deletion
	lb, err := ovndbapi.LBGet(LB4)
	if err != nil {
		t.Fatal(err)
	}
	if len(lb) != 0 {
		t.Fatalf("error: lb deletion not done, total:%v", len(lb))
	}
	t.Logf("LB lb4 deleted")

	cmd, err = ovndbapi.LRDel(LR)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

}
