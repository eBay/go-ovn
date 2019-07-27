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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestACLs(t *testing.T) {
	var cmds []*OvnCommand
	var cmd *OvnCommand
	var err error

	cmds = make([]*OvnCommand, 0)
	cmd, err = ovndbapi.LSAdd(LSW)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)

	// execute to create lsw and lsp
	err = ovndbapi.Execute(cmds...)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		cmd, err = ovndbapi.LSDel(LSW)
		if err != nil {
			t.Fatal(err)
		}
		err = ovndbapi.Execute(cmd)
		if err != nil {
			t.Fatal(err)
		}
	}()

	lsws, err := ovndbapi.LSList()
	if err != nil {
		t.Fatal(err)
	}
	if len(lsws) != 1 {
		t.Fatalf("ls not created %d", len(lsws))
	}

	// nil cmds for next batch
	cmds = make([]*OvnCommand, 0)
	cmd, err = ovndbapi.ACLAdd(LSW, &ACL{Direction: "to-lport", Match: MATCH, Action: "drop", Priority: 1001})
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)

	err = ovndbapi.Execute(cmds...)
	if err != nil {
		t.Fatal(err)
	}

	acls, err := ovndbapi.ACLList(LSW)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(acls) == 1 && acls[0].Match == MATCH &&
		acls[0].Action == "drop" && acls[0].Priority == 1001, "test[%s] %s", "add acl", acls[0])

	cmd, err = ovndbapi.ACLAdd(LSW, &ACL{Direction: "to-lport", Match: MATCH, Action: "drop", Priority: 1001})
	// ACLAdd must return error
	assert.Equal(t, false, nil == err, "test[%s]", "add same acl twice, should only one added.")
	// cmd is nil, so this is noop
	err = ovndbapi.Execute(cmd)
	assert.Equal(t, true, nil == err, "test[%s]", "add same acl twice, should only one added.")

	cmd, err = ovndbapi.ACLAdd(LSW, &ACL{Direction: "to-lport", Match: MATCH_SECOND, Action: "drop", Priority: 1001, ExternalID: map[string]string{"A": "a", "B": "b"}})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	acls, err = ovndbapi.ACLList(LSW)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(acls) == 2, "test[%s]", "add second acl")

	cmd, err = ovndbapi.ACLAdd(LSW, &ACL{Direction: "to-lport", Match: MATCH_SECOND, Action: "drop", Priority: 1001, ExternalID: map[string]string{"A": "b", "B": "b"}})
	// ACLAdd must return error ovn-nbctl acl-add not allow to add the same acl with different external_ids
	assert.Equal(t, false, nil == err, "test[%s]", "add same acl twice, should only one added.")
	// cmd is noop
	err = ovndbapi.Execute(cmd)
	assert.Equal(t, true, nil == err, "test[%s]", "add same acl twice, should only one added.")

	acls, err = ovndbapi.ACLList(LSW)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(acls) == 2, "test[%s]", "add second acl")

	cmd, err = ovndbapi.ACLDel(LSW, &ACL{Direction: "to-lport", Match: MATCH, Priority: 1001})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	acls, err = ovndbapi.ACLList(LSW)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(acls) == 1, "test[%s]", "acl remove")

	cmd, err = ovndbapi.ACLDel(LSW, &ACL{Direction: "to-lport", Match: MATCH_SECOND, Priority: 1001, ExternalID: map[string]string{"A": "a"}})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	acls, err = ovndbapi.ACLList(LSW)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, len(acls) == 0, "test[%s]", "acl remove")

}
