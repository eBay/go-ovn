package goovn

import (
	"testing"
)

const LR3 = "lr3"

func TestNAT(t *testing.T) {
	var cmd *OvnCommand
	var err error
	defer func() {
		cmd, err = ovndbapi.LRDel(LR3)
		if err != nil {
			t.Fatal(err)
		}
		err = ovndbapi.Execute(cmd)
		if err != nil {
			t.Fatal(err)
		}
	}()

	cmd, err = ovndbapi.LRAdd(LR3, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.LRNATAdd(LR3, "snat", "10.127.0.129", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	cmd, err = ovndbapi.LRNATAdd(LR3, "snat", "10.127.0.128", "172.16.255.127/25", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	cmd, err = ovndbapi.LRNATAdd(LR3, "dnat_and_snat", "10.127.0.128", "172.16.255.127/25", nil, "br-int", "55.55.55.55.55.55")
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.LRNATAdd(LR3, "dnat", "10.127.0.127", "172.16.255.128/24", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	natlist, err := ovndbapi.LRNATList(LR3)
	if err != nil {
		t.Fatal(err)
	}
	if len(natlist) != 4 {
		t.Fatal("nat not add yet!")
	}

	cmd, err = ovndbapi.LRNATDel(LR3, "snat")
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	natlist, err = ovndbapi.LRNATList(LR3)

	if len(natlist) != 2 {
		t.Fatal("snat not Delete!")
	}

	cmd, err = ovndbapi.LRNATDel(LR3, "dnat")
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	natlist, err = ovndbapi.LRNATList(LR3)
	if err != nil {
		t.Fatal(err)
	}
	if len(natlist) != 1 {
		t.Fatal("dnat not Delete!")
	}

	cmd, err = ovndbapi.LRNATAdd(LR3, "snat", "10.127.0.128", "172.16.255.128/24", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.LRNATAdd(LR3, "dnat", "10.127.0.127", "172.16.255.128/24", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	natlist, err = ovndbapi.LRNATList(LR3)
	if err != nil {
		t.Fatal(err)
	}
	if len(natlist) != 3 {
		t.Fatal("nat not add yet!")
	}

	cmd, err = ovndbapi.LRNATDel(LR3, "")
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	natlist, err = ovndbapi.LRNATList(LR3)
	if err != nil {
		t.Fatal(err)
	}
	if len(natlist) != 0 {
		t.Fatal("nat not delete yet!")
	}
}
