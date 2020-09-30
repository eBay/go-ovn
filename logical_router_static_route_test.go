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

var nextHop2 = "10.3.0.2"

func TestLogicalRouterStaticRoute(t *testing.T) {
	ovndbapi := getOVNClient(DBNB)
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
	t.Logf("Adding static route %s via %s to LRouter %s Done", IPPREFIX, NEXTHOP, LR2)
	// verify static route addition to lr2
	lrsr, err := ovndbapi.LRSRList(LR2)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrsr) == 0 {
		t.Fatalf("Static Route %s not created in %s", IPPREFIX, LR2)
	}
	assert.Equal(t, true, lrsr[0].IPPrefix == IPPREFIX, "Added static route to lr2")
	// add static route IPPREFIX via nextHop2
	cmd, err = ovndbapi.LRSRAdd(LR2, IPPREFIX, nextHop2, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Adding static route to lr2 failed with err %v", err)
	}
	t.Logf("Adding static route %s via %s to LRouter %s Done", IPPREFIX, nextHop2, LR2)
	// verify static route addition to lr2 via nexthop2
	lrsr, err = ovndbapi.LRSRList(LR2)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrsr) < 2 {
		t.Fatalf("Static Route %s via %s not created in %s", IPPREFIX, nextHop2, LR2)
	}
	found := false
	var secondSr *LogicalRouterStaticRoute
	for _, sr := range lrsr {
		if sr.Nexthop == nextHop2 && sr.IPPrefix == IPPREFIX {
			found = true
			secondSr = sr
		}
	}
	assert.Equal(t, true, found, "Added second static route to lr2")
	// delete static route via nextHop2
	cmd, err = ovndbapi.LRSRDelByUUID(LR2, secondSr.UUID)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Deleting static route from lr2 via %s failed with err %v", nextHop2, err)
	}
	t.Logf("Deleted static route %s via %s from LRouter %s", IPPREFIX, nextHop2, LR2)
	// verify static route via nexthop2 delete from lr2
	lrsr, err = ovndbapi.LRSRList(LR2)
	if err != nil {
		t.Fatal(err)
	}
	found = false
	for _, sr := range lrsr {
		if sr.Nexthop == nextHop2 && sr.IPPrefix == IPPREFIX {
			found = true
		}
	}
	assert.Equal(t, false, found, "Deleted second static route from lr2")

	// Delete the static route from lr2
	cmd, err = ovndbapi.LRSRDel(LR2, IPPREFIX, nil, nil, nil)
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

	// Add static route with policy and output_port.
	outputPort := "lsp1"
	policy := "src-ip"
	cmd, err = ovndbapi.LRSRAdd(LR2, IPPREFIX, NEXTHOP, []string{outputPort}, []string{policy}, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Adding static route to lr2 failed with err %v", err)
	}
	lrsr, err = ovndbapi.LRSRList(LR2)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrsr) < 1 {
		t.Fatalf("Static Route %s using %s with policy %s not created in %s", IPPREFIX, outputPort, policy, LR2)
	}
	assert.Equal(t, outputPort, lrsr[0].OutputPort[0])
	assert.Equal(t, policy, lrsr[0].Policy[0])

	// Delete the static route with outputPort specified
	cmd, err = ovndbapi.LRSRDel(LR2, IPPREFIX, nil, &outputPort, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Deleting static route from lr2 failed with err %v", err)
	}
	t.Logf("Deleting static route %s from LRouter %s Done", IPPREFIX, LR2)

	lrsr, err = ovndbapi.LRSRList(LR2)
	if err != nil {
		t.Fatal(err)
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

	// verify static route list for non-existing routers
	lrsr, err = ovndbapi.LRSRList(FAKENOROUTER)
	if err != nil {
		assert.EqualError(t, ErrorNotFound, err.Error())
	}
}
