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
	ovndbapi := getOVNClient(DBNB)
	var cmds []*OvnCommand
	var cmd *OvnCommand
	var err error

	cmds = make([]*OvnCommand, 0)
	cmd, err = ovndbapi.LSAdd(LSW)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)

	cmd, err = ovndbapi.LSPAdd(LSW, LSP)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)

	cmd, err = ovndbapi.LSPSetAddress(LSP, ADDR)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)

	cmd, err = ovndbapi.LSPSetPortSecurity(LSP, ADDR)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)
	cmd, err = ovndbapi.MeterAdd("meter1", "drop", 101, "kbps", nil, 300)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)

	// execute to create lsw and lsp
	err = ovndbapi.Execute(cmds...)
	if err != nil {
		t.Fatal(err)
	}

	// nil cmds for next batch
	cmds = make([]*OvnCommand, 0)
	cmd, err = ovndbapi.ACLAdd(LSW, "to-lport", MATCH, "drop", 1001, nil, true, "meter1", "alert")
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)

	err = ovndbapi.Execute(cmds...)
	if err != nil {
		t.Fatal(err)
	}

	lsws, err := ovndbapi.LSGet(LSW)
	if err != nil {
		t.Fatal(err)
	}

	if len(lsws) == 0 {
		t.Fatalf("ls not created %d", len(lsws))
	}

	lsps, err := ovndbapi.LSPList(LSW)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, len(lsps) == 1 && lsps[0].Name == LSP, "test[%s]: %v", "added port", lsps)
	assert.Equal(t, true, len(lsps) == 1 && lsps[0].Addresses[0] == ADDR, "test[%s]", "setted port address")
	assert.Equal(t, true, len(lsps) == 1 && lsps[0].PortSecurity[0] == ADDR, "test[%s]", "setted port port security")

	cmd, err = ovndbapi.LSPAdd(LSW, LSP_SECOND)
	if err != nil {
		t.Fatal(err)
	}

	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lsps, err = ovndbapi.LSPList(LSW)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, len(lsps) == 2, "test[%s]: %+v", "added 2 ports", lsps)

	acls, err := ovndbapi.ACLList(LSW)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(acls) == 1 && acls[0].Match == MATCH &&
		acls[0].Action == "drop" && acls[0].Priority == 1001 && acls[0].Log == true && acls[0].Meter[0] == "meter1" && acls[0].Severity == "alert", "test[%s] %s", "add acl", acls[0])

	cmd, err = ovndbapi.ACLAdd(LSW, "to-lport", MATCH, "drop", 1001, nil, true, "", "")
	// ACLAdd must return error
	assert.Equal(t, true, nil != err, "test[%s]", "add same acl twice, should only one added.")
	// cmd is nil, so this is noop
	err = ovndbapi.Execute(cmd)
	assert.Equal(t, true, nil == err, "test[%s]", "add same acl twice, should only one added.")

	cmd, err = ovndbapi.ACLAdd(LSW, "to-lport", MATCH_SECOND, "drop", 1001, map[string]string{"A": "a", "B": "b"}, false, "", "")
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

	cmd, err = ovndbapi.ACLAdd(LSW, "to-lport", MATCH_SECOND, "drop", 1001, map[string]string{"A": "b", "B": "b"}, false, "", "")
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

	cmd, err = ovndbapi.LSPDel(LSP)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lsps, err = ovndbapi.LSPList(LSW)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, len(lsps) == 1, "test[%s]", "one port remove")

	cmd, err = ovndbapi.LSPDel(LSP_SECOND)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lsps, err = ovndbapi.LSPList(LSW)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, len(lsps) == 0, "test[%s]", "one port remove")

	cmd, err = ovndbapi.LSDel(LSW)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	cmd, err = ovndbapi.MeterDel("meter1")
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
}
