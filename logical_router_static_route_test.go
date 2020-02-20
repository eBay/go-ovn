package goovn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	LR2      = "lr2"
	IPPREFIX = "10.0.0.1/24"
	NEXTHOP  = "10.3.0.1"
)

func TestLogicalRouterStaticRoute(t *testing.T) {
	ovndbapi := getOVNClient(dbNB)
	cmd, err := ovndbapi.LRAdd(LR2, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Add router %s  Done", LR2)

	// verify router create
	lrs, err := ovndbapi.LRGet(LR2)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrs) != 1 {
		t.Fatalf("lr not created %v", lrs)
	}

	//lr string, ip_prefix string, nexthop string, output_port []string, policy []string, external_ids map[string]string
	cmd, err = ovndbapi.LRSRAdd(LR2, IPPREFIX, NEXTHOP, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Adding static route to lr2 failed with err %v", err)
	}
	t.Logf("Adding static route %s to LRouter %s Done", IPPREFIX, LR2)
	// verify static route addition to lr2
	lrsr, err := ovndbapi.LRSRList(LR2)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrsr) == 0 {
		t.Fatalf("Static Route %s not created in %s", IPPREFIX, LR2)
	}
	assert.Equal(t, true, lrsr[0].IPPrefix == IPPREFIX, "Added static route to lr2")

	// Delete the static route from lr2
	cmd, err = ovndbapi.LRSRDel(LR2, IPPREFIX)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Deleting static route from lr2 failed with err %v", err)
	}
	t.Logf("Deleting static route %s from LRouter %s Done", IPPREFIX, LR2)

	// verify static route delete from lr2
	lrsr, err = ovndbapi.LRSRList(LR2)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrsr) != 0 {
		t.Fatalf("Static Route %s not deleted in %s", IPPREFIX, LR2)
	}
	assert.Equal(t, true, len(lrsr) == 0, "Deleted static route from lr2")

	// Delete the router
	cmd, err = ovndbapi.LRDel(LR2)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Delete router %s  Done", LR2)

	// verify router delete
	lrs, err = ovndbapi.LRGet(LR2)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrs) != 0 {
		t.Fatalf("lr not deleted %v", lrs)
	}

}
