package goovn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	LR6               = "lr6"
	PRIORITY          = 10
	LR_POLICY_MATCH_2 = "ip4.src == 3.3.3.0/24"
	ACTION_DROP       = "drop"
	ACTION_ALLOW      = "allow"
	ACTION_REROUTE    = "reroute"
)

var (
	LR_POLICY_MATCH_1 = "ip4.src == 1.1.1.0/24"
)

func TestLogicalRouterPolicy(t *testing.T) {
	ovndbapi := getOVNClient(DBNB)
	cmd, err := ovndbapi.LRAdd(LR6, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Add router %s Done", LR6)

	// verify router create
	lrs, err := ovndbapi.LRGet(LR6)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrs) != 1 {
		t.Fatalf("lr not created %v", lrs)
	}

	cmd, err = ovndbapi.LRPolicyAdd(LR6, PRIORITY, LR_POLICY_MATCH_1, ACTION_DROP, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Adding policy to LR6 failed with err %v", err)
	}

	// verify policy addition to LR6
	lrpolicy, err := ovndbapi.LRPolicyList(LR6)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrpolicy) == 0 {
		t.Fatalf("policy %s not created in %s", LR_POLICY_MATCH_1, LR6)
	}
	assert.Equal(t, true, lrpolicy[0].Match == LR_POLICY_MATCH_1, "Added policy to LR6")

	// add another policy
	cmd, err = ovndbapi.LRPolicyAdd(LR6, PRIORITY, LR_POLICY_MATCH_2, ACTION_ALLOW, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Adding policy to LR6 failed with err %v", err)
	}

	// verify policy addition to LR6
	lrpolicy, err = ovndbapi.LRPolicyList(LR6)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrpolicy) < 2 {
		t.Fatalf("policy %s not created in %s", LR_POLICY_MATCH_2, LR6)
	}
	found := false
	var secondPolicy *LogicalRouterPolicy
	for _, p := range lrpolicy {
		if p.Match == LR_POLICY_MATCH_2 {
			found = true
			secondPolicy = p
		}
	}
	assert.Equal(t, true, found, "Added second policy to LR6")

	// delete policy
	cmd, err = ovndbapi.LRPolicyDelByUUID(LR6, secondPolicy.UUID)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Deleting policy from LR6 via %s failed with err %v", LR_POLICY_MATCH_2, err)
	}
	t.Logf("Deleted policy %s from LRouter %s", LR_POLICY_MATCH_2, LR6)

	// verify policy via delete from LR6
	lrpolicy, err = ovndbapi.LRPolicyList(LR6)
	if err != nil {
		t.Fatal(err)
	}
	found = false
	for _, p := range lrpolicy {
		if p.Match == LR_POLICY_MATCH_2 {
			found = true
		}
	}
	assert.Equal(t, false, found, "Deleted second policy from LR6")

	// Delete the policy first policy from LR6
	cmd, err = ovndbapi.LRPolicyDel(LR6, PRIORITY, &LR_POLICY_MATCH_1)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Deleting policy from LR6 failed with err %v", err)
	}
	t.Logf("Deleting policy %s from LRouter %s Done", LR_POLICY_MATCH_1, LR6)

	// verify policy delete from LR6
	lrpolicy, err = ovndbapi.LRPolicyList(LR6)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrpolicy) != 0 {
		t.Fatalf("policy %s not deleted in %s", LR_POLICY_MATCH_1, LR6)
	}
	assert.Equal(t, true, len(lrpolicy) == 0, "Deleted policy from LR6")

	// Add re-route policy
	nexthop := "4.4.4.4"
	cmd, err = ovndbapi.LRPolicyAdd(LR6, PRIORITY, LR_POLICY_MATCH_1, ACTION_REROUTE, &nexthop, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Adding policy to LR6 failed with err %v", err)
	}
	lrpolicy, err = ovndbapi.LRPolicyList(LR6)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrpolicy) < 1 {
		t.Fatalf("policy %s with nexthop %s not created in %s", LR_POLICY_MATCH_1, nexthop, LR6)
	}
	assert.Equal(t, nexthop, *lrpolicy[0].Nexthop)
	assert.Equal(t, LR_POLICY_MATCH_1, lrpolicy[0].Match)

	// Delete the policy by priority
	cmd, err = ovndbapi.LRPolicyDel(LR6, PRIORITY, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Deleting policy from LR6 failed with err %v", err)
	}
	t.Logf("Deleting policy %s from LRouter %s Done", LR_POLICY_MATCH_1, LR6)

	lrpolicy, err = ovndbapi.LRPolicyList(LR6)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(lrpolicy) == 0, "Deleted policy from LR6")

	// once more
	cmd, err = ovndbapi.LRPolicyAdd(LR6, PRIORITY, LR_POLICY_MATCH_1, ACTION_DROP, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Adding policy to LR6 failed with err %v", err)
	}

	// verify policy addition to LR6
	lrpolicy, err = ovndbapi.LRPolicyList(LR6)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrpolicy) == 0 {
		t.Fatalf("policy %s not created in %s", LR_POLICY_MATCH_1, LR6)
	}
	assert.Equal(t, true, lrpolicy[0].Match == LR_POLICY_MATCH_1, "Added policy to LR6")

	// delete all policies on the LR
	cmd, err = ovndbapi.LRPolicyDelAll(LR6)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("Deleting policy from LR6 failed with err %v", err)
	}
	t.Logf("Deleting policy %s from LRouter %s Done", LR_POLICY_MATCH_1, LR6)

	lrpolicy, err = ovndbapi.LRPolicyList(LR6)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(lrpolicy) == 0, "Deleted policy from LR6")

	// Delete the router
	cmd, err = ovndbapi.LRDel(LR6)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Delete router %s  Done", LR6)

	// verify router delete
	lrs, err = ovndbapi.LRGet(LR6)
	if err != nil {
		t.Fatal(err)
	}
	if len(lrs) != 0 {
		t.Fatalf("lr not deleted %v", lrs)
	}

	// verify policy list for non-existing routers
	_, err = ovndbapi.LRPolicyList(FAKENOROUTER)
	if err != nil {
		assert.EqualError(t, ErrorNotFound, err.Error())
	}
}
