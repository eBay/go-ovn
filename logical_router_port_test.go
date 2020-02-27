package goovn

import "testing"

const LR4 = "lr4"

func TestLogicalRouterPort(t *testing.T) {
	ovndbapi := getOVNClient(DBNB)
	var cmds []*OvnCommand
	var cmd *OvnCommand
	var err error

	cmds = make([]*OvnCommand, 0)
	cmd, err = ovndbapi.LRAdd(LR4, nil)
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

	// lr string, lrp string, mac string, network []string, peer string
	cmd, err = ovndbapi.LRPAdd(LR4, LRP, "54:54:54:54:54:54", []string{"192.168.0.1/24"}, "lrp2", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lrps, err := ovndbapi.LRPList(LR4)
	if err != nil {
		t.Fatal(err)
	}

	if len(lrps) != 1 {
		t.Fatalf("lrp not created %v", lrps)
	}

	cmd, err = ovndbapi.LRPDel(LR4, LRP)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lrps, err = ovndbapi.LRPList(LR4)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrps) != 0 {
		t.Fatalf("lrp not created %v", lrps)
	}

	cmd, err = ovndbapi.LRDel(LR4)
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
			if lr.Name == LR4 {
				t.Fatalf("lr not deleted %v", LR4)
				break
			}
		}
	} else {
		t.Logf("Successfully deleted router %s", LR4)
	}

}
