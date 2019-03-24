package goovn

import "testing"

func TestLogicalRouter(t *testing.T) {
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
