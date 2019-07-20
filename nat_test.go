package goovn

import (
	"testing"
)

func TestNAT(t *testing.T) {
	var cmd *OvnCommand
	var err error
	defer func() {
		cmd, err = ovndbapi.LRDel(LR)
		if err != nil {
			t.Fatal(err)
		}
		err = ovndbapi.Execute(cmd)
		if err != nil {
			t.Fatal(err)
		}
	}()

	cmd, err = ovndbapi.LRAdd(LR, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	cmd, err = ovndbapi.LRNATAdd(LR, "snat", "10.127.0.129", "", "172.16.255.128/25", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	natlist, err := ovndbapi.LRNATList(LR)
	if err != nil {
		t.Fatal(err)
	}
	if len(natlist) != 1 {
		t.Fatal("nat not add yet!")
	}
	cmd, err = ovndbapi.LRNATDel(LR, "snat", "172.16.255.128/25")
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	natlist, err = ovndbapi.LRNATList(LR)
	if len(natlist) != 0 {
		t.Fatal("nat not Delete!")
	}
}
