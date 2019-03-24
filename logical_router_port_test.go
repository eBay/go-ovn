package goovn

import "testing"

func TestLogicalRouterPort(t *testing.T) {
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

	lrs, err := ovndbapi.GetLogicalRouters()
	if err != nil {
		t.Fatal(err)
	}
	if len(lrs) != 1 {
		t.Fatalf("lr not created %v", lrs)
	}

	// lr string, lrp string, mac string, network []string, peer string
	cmd, err = ovndbapi.LRPAdd(LR, LRP, "54:54:54:54:54:54", []string{"192.168.0.1/24"}, "lrp2", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lrps, err := ovndbapi.GetLogicalRouterPortsByRouter(LR)
	if err != nil {
		t.Fatal(err)
	}

	if len(lrps) != 1 {
		t.Fatalf("lrp not created %v", lrps)
	}

	cmd, err = ovndbapi.LRPDel(LR, LRP)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lrps, err = ovndbapi.GetLogicalRouterPortsByRouter(LR)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrps) != 0 {
		t.Fatalf("lrp not created %v", lrps)
	}

	cmd, err = ovndbapi.LRDel(LR)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lrs, err = ovndbapi.GetLogicalRouters()
	if err != nil {
		t.Fatal(err)
	}
	if len(lrs) != 0 {
		t.Fatalf("lr not deleted %v", lrs)
	}

}
