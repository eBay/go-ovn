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

var (
	aclLS = "acl_ls"
)

func TestACLAdd(t *testing.T) {
	var cmd *OvnCommand
	var err error
	_ = cmd
	/*
		cmd, err := ovndbapi.LogicalSwitch.Add(LogicalSwitchName(aclLS))
			if err != nil {
				t.Fatal(err)
			}
			err = ovndbapi.Execute(cmd)
			if err != nil {
				t.Fatal(err)
			}

			cmd, err = ovndbapi.ACL.Add(
				ACLEntityName(aclLS),
				ACLDirection("to-lport"),
				ACLMatch(MATCH),
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
	*/
	ls, err := ovndbapi.LogicalSwitch.Get(LogicalSwitchName(aclLS))
	if err != nil {
		t.Fatal(err)
	}
	_ = ls
	acls, err := ovndbapi.ACL.List(ACLEntityName(aclLS))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(acls) == 1 && acls[0].Match == MATCH &&
		acls[0].Action == "drop" && acls[0].Priority == 1001, "test[%s] %s", "add acl", acls[0])

	/*
		cmd, err = ovndbapi.ACLAdd(LSW, "to-lport", MATCH, "drop", 1001, nil, true, "")
		// ACLAdd must return error
		assert.Equal(t, true, nil != err, "test[%s]", "add same acl twice, should only one added.")
		// cmd is nil, so this is noop
		err = ovndbapi.Execute(cmd)
		assert.Equal(t, true, nil == err, "test[%s]", "add same acl twice, should only one added.")

		cmd, err = ovndbapi.ACLAdd(LSW, "to-lport", MATCH_SECOND, "drop", 1001, map[string]string{"A": "a", "B": "b"}, false, "")
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

		cmd, err = ovndbapi.ACLAdd(LSW, "to-lport", MATCH_SECOND, "drop", 1001, map[string]string{"A": "b", "B": "b"}, false, "")
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
		assert.Equal(t, true, len(acls) == 3, "test[%s]", "add second acl")

		cmd, err = ovndbapi.ACLDel(LSW, "to-lport", MATCH, 1001, map[string]string{})
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
		assert.Equal(t, true, len(acls) == 2, "test[%s]", "acl remove")

		cmd, err = ovndbapi.ACLDel(LSW, "to-lport", MATCH_SECOND, 1001, map[string]string{"A": "a"})
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

		cmd, err = ovndbapi.ACLDel(LSW, "to-lport", MATCH_SECOND, 1001, map[string]string{"A": "b"})
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

		cmd, err = ovndbapi.LSDel(LSW)
		if err != nil {
			t.Fatal(err)
		}
		err = ovndbapi.Execute(cmd)
		if err != nil {
			t.Fatal(err)
		}
	*/
}
