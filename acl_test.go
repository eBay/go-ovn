/**
 * Copyright (c) 2017 eBay Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 **/

package goovn

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	aclTestLS    string = "testeb3109af-eb44-4b0b-b7a4-dcf439f21a33"
	aclTestMatch        = "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.61"
)

func TestACLAdd(t *testing.T) {
	lsUUID := newUUID(t)

	aclTestLS = "test" + lsUUID

	cmd, err := ovndbapi.LogicalSwitch.Add(LogicalSwitchName(aclTestLS))
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.ACL.Add(
		ACLEntityName(aclTestLS),
		ACLDirection("to-lport"),
		ACLMatch(aclTestMatch),
		ACLAction("drop"),
		ACLPriority(1001),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	ls, err := ovndbapi.LogicalSwitch.Get(LogicalSwitchName(aclTestLS))
	if err != nil {
		t.Fatal(err)
	}
	if len(ls.ACLs) != 1 {
		t.Fatalf("acl not added to %s", aclTestLS)
	}

	acls, err := ovndbapi.ACL.List(ACLEntityName(aclTestLS))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(acls) == 1 && acls[0].Match == aclTestMatch &&
		acls[0].Action == "drop" && acls[0].Priority == 1001, "test[%s] %s", "add acl", acls[0])
}

func TestACLAddDup(t *testing.T) {
	cmd, err := ovndbapi.ACL.Add(
		ACLEntityName(aclTestLS),
		ACLDirection("to-lport"),
		ACLMatch(aclTestMatch),
		ACLAction("drop"),
		ACLPriority(1001),
	)
	if err == nil {
		t.Fatal("the same acl must be added only one")
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

}

func TestACLDel(t *testing.T) {
	cmd, err := ovndbapi.ACL.Del(
		ACLEntityName(aclTestLS),
		ACLDirection("to-lport"),
		ACLMatch(aclTestMatch),
		ACLAction("drop"),
		ACLPriority(1001),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	ls, err := ovndbapi.LogicalSwitch.Get(LogicalSwitchName(aclTestLS))
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("EEE %#+v\n", ls)
	acls, err := ovndbapi.ACL.List(ACLEntityName(aclTestLS))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(acls) == 0, "fail remove %v", acls[0])

	/*
		cmd, err = ovndbapi.LogicalSwitch.Del(LogicalSwitchName(aclTestLS))
		if err != nil {
			t.Fatal(err)
		}
		err = ovndbapi.Execute(cmd)
		if err != nil {
			t.Fatal(err)
		}
	*/
}
