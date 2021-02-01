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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogicalSwitchACLs(t *testing.T) {
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

	// The following should fail because it is considered a duplicate of an existing ACL
	cmd, err = ovndbapi.ACLAdd(LSW, "to-lport", MATCH_SECOND, "drop", 1001, map[string]string{"A": "b", "B": "b"}, false, "", "")
	if err == nil {
		t.Fatal(err)
	}
	if cmd != nil {
		t.Fatal(err)
	}
	// cmd is nil, so this is noop
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	acls, err = ovndbapi.ACLList(LSW)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(acls) == 2, "test[%s]", "add second acl")

	// The following should fail because it is considered a duplicate of an existing ACL
	cmd, err = ovndbapi.ACLAdd(LSW, "to-lport", MATCH_SECOND, "drop", 1001, nil, false, "", "")
	if err == nil {
		t.Fatal(err)
	}

	// Different priority is a different ACL, so this should succeed
	cmd, err = ovndbapi.ACLAdd(LSW, "to-lport", MATCH_SECOND, "drop", 1002, map[string]string{"A": "a", "B": "b"}, false, "", "")
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

	// Different direction is a different ACL, so this should succeed
	cmd, err = ovndbapi.ACLAdd(LSW, "from-lport", MATCH_SECOND, "drop", 1001, map[string]string{"A": "a", "B": "b"}, false, "", "")
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
	assert.Equal(t, true, len(acls) == 4, "test[%s]", "add second acl")

	cmd, err = ovndbapi.ACLDel(LSW, "to-lport", MATCH_SECOND, 1002, nil)
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

	cmd, err = ovndbapi.ACLDel(LSW, "from-lport", MATCH_SECOND, 1001, nil)
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

	assert.Equal(t, true, len(acls) == 1, "test[%s]", "acl remove")

	//The following delete should fail because external ids are provided, but they don't exist in any ACL
	cmd, err = ovndbapi.ACLDel(LSW, "to-lport", MATCH_SECOND, 1001, map[string]string{"A": "b"})
	if err == nil {
		t.Fatal(err)
	}

	//The following delete should succeed because the provided external_ids provided are a subset of thoe in an existing ACL
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
	assert.Equal(t, true, len(acls) == 0, "test[%s]", "acl remove")

	// The following ACLDel should fail because all the ACLs have been deleted.
	cmd, err = ovndbapi.ACLDel(LSW, "to-lport", MATCH_SECOND, 1001, map[string]string{"A": "b"})
	if err == nil {
		t.Fatal(err)
	}

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
	// verify ACL list for non-existing switch
	_, err = ovndbapi.ACLList(FAKENOSWITCH)
	if err != nil {
		assert.EqualError(t, ErrorNotFound, err.Error())
	}
}

func compareMeterSlices(s1, s2 []string) bool {
	if (s1 == nil || s1[0] == "") && (s2 == nil || s2[0] == "") {
		return true
	}
	if len(s1) != len(s2) {
		return false
	}
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}

// Returns true if an acl is in aclList
func containsACL(aclList []*ACL, acl *ACL) bool {
	for _, a := range aclList {
		if a.UUID == acl.UUID &&
			a.Name == acl.Name &&
			a.Action == acl.Action &&
			a.Direction == acl.Direction &&
			a.Match == acl.Match &&
			a.Priority == acl.Priority &&
			a.Log == acl.Log &&
			compareMeterSlices(a.Meter, acl.Meter) &&
			a.Severity == acl.Severity &&
			compareExternalIds(iMapToSMap(a.ExternalID), acl.ExternalID) {
			return true
		}
	}
	return false
}

// converts and interface{} map to a string map
func iMapToSMap(iMap map[interface{}]interface{}) map[string]string {
	if iMap == nil {
		return nil
	}
	sMap := make(map[string]string, len(iMap))
	for k, v := range iMap {
		sMap[fmt.Sprintf("%v", k)] = fmt.Sprintf("%v", v)
	}
	return sMap
}

func TestACLs(t *testing.T) {
	ovndbapi := getOVNClient(DBNB)
	var cmd *OvnCommand
	var cmds []*OvnCommand
	var err error
	var lsp1UUID string
	var lsp2UUID string

	t.Run("create switch, ports, port group, and meter for ACL testing", func(t *testing.T) {
		cmds = make([]*OvnCommand, 0)

		// Create switch and ports
		cmd, err := ovndbapi.LSAdd(PG_TEST_LS1)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		// Add ports
		cmd, err = ovndbapi.LSPAdd(PG_TEST_LS1, PG_TEST_LSP1)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPSetAddress(PG_TEST_LSP1, ADDR)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPSetPortSecurity(PG_TEST_LSP1, ADDR)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPAdd(PG_TEST_LS1, PG_TEST_LSP2)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPSetAddress(PG_TEST_LSP2, ADDR2)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPSetPortSecurity(PG_TEST_LSP2, ADDR2)
		cmds = append(cmds, cmd)
		assert.Nil(t, err)

		result, err := ovndbapi.ExecuteR(cmds...)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(result))

		if len(result) >= 3 {
			lsp1UUID = result[1]
			lsp2UUID = result[2]
		}

		// Create port group
		cmd, err = ovndbapi.PortGroupAdd(PG_TEST_PG1, []string{lsp1UUID, lsp2UUID}, nil)
		assert.Nil(t, err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(t, err)

		// Create some meters
		cmds = make([]*OvnCommand, 0)
		cmd, err = ovndbapi.MeterAdd(METER1, "drop", 101, "kbps", nil, 300)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.MeterAdd(METER2, "drop", 101, "kbps", nil, 300)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		// execute to create meters
		err = ovndbapi.Execute(cmds...)
		assert.Nil(t, err)
	})

	type aclTest struct {
		acl         ACL
		newName     string
		newMatch    string
		newMeter    string
		newSeverity string
	}

	type aclEntityTest struct {
		entityType EntityType
		entity     string
		aclTests   []aclTest
	}

	aclEntityTests := []aclEntityTest{
		{
			PORT_GROUP, PG_TEST_PG1,
			[]aclTest{
				{
					ACL{"", "", "drop", "from-lport", MATCH3, 1001, false, []string{""}, "", nil},
					ACL_NAME_2, MATCH_SECOND, METER1, SEVERITY_INFO,
				},
				{
					ACL{"", ACL_NAME_1, "drop", "to-lport", MATCH, 1001, true, []string{METER1}, SEVERITY_ALERT, nil},
					ACL_NAME_3, MATCH_SECOND, METER2, SEVERITY_INFO,
				},
				{
					ACL{"", ACL_NAME_2, "drop", "from-lport", MATCH, 1002, true, []string{METER2}, SEVERITY_INFO, map[interface{}]interface{}{"A": "a", "B": "b"}},
					ACL_NAME_4, MATCH3, METER1, SEVERITY_WARNING,
				},
				{
					ACL{"", ACL_NAME_3, "drop", "to-lport", MATCH, 1002, true, []string{METER1}, SEVERITY_ALERT, map[interface{}]interface{}{"A": "b", "B": "a"}},
					ACL_NAME_5, MATCH_SECOND, METER2, SEVERITY_INFO,
				},
			},
		},
		{
			LOGICAL_SWITCH, PG_TEST_LS1,
			[]aclTest{
				{
					ACL{"", "", "drop", "from-lport", MATCH3, 1001, false, []string{""}, "", nil},
					ACL_NAME_2, MATCH_SECOND, METER1, SEVERITY_INFO,
				},
				{
					ACL{"", ACL_NAME_1, "drop", "to-lport", MATCH, 1001, true, []string{METER1}, SEVERITY_ALERT, nil},
					ACL_NAME_3, MATCH_SECOND, METER2, SEVERITY_INFO,
				},
				{
					ACL{"", ACL_NAME_2, "drop", "from-lport", MATCH, 1002, true, []string{METER2}, SEVERITY_INFO, map[interface{}]interface{}{"A": "a", "B": "b"}},
					ACL_NAME_4, MATCH3, METER1, SEVERITY_WARNING,
				},
				{
					ACL{"", ACL_NAME_3, "drop", "to-lport", MATCH, 1002, true, []string{METER1}, SEVERITY_ALERT, map[interface{}]interface{}{"A": "b", "B": "a"}},
					ACL_NAME_5, MATCH_SECOND, METER2, SEVERITY_INFO,
				},
			},
		},
	}

	// Test ACLAddEntity()
	for i, e := range aclEntityTests {
		for j, _ := range e.aclTests {
			tcName := fmt.Sprintf("add ACL to %s %s #%d", e.entityType, e.entity, j)
			t.Run(tcName, func(t *testing.T) {
				acl := &aclEntityTests[i].aclTests[j].acl
				cmd, err = ovndbapi.ACLAddEntity(e.entityType, e.entity, acl.Name, acl.Direction, acl.Match, acl.Action, acl.Priority, iMapToSMap(acl.ExternalID), acl.Log, acl.Meter[0], acl.Severity)
				assert.Nil(t, err)
				result, err := ovndbapi.ExecuteR(cmd)
				assert.Nil(t, err)
				assert.Len(t, result, 1)
				if len(result) > 0 {
					acl.UUID = result[0]
				}
				acls, err := ovndbapi.ACLListEntity(e.entityType, e.entity)
				assert.Nil(t, err)
				assert.Equal(t, j+1, len(acls))
				assert.True(t, containsACL(acls, acl))
			})
		}
	}

	// Test adding duplicate ACLs
	for i, e := range aclEntityTests {
		for j, _ := range e.aclTests {
			tcName := fmt.Sprintf("add duplicate ACL to %s %s #%d", e.entityType, e.entity, j)
			t.Run(tcName, func(t *testing.T) {
				acl := &aclEntityTests[i].aclTests[j].acl
				cmd, err = ovndbapi.ACLAddEntity(e.entityType, e.entity, acl.Name, acl.Direction, acl.Match, acl.Action, acl.Priority, iMapToSMap(acl.ExternalID), acl.Log, acl.Meter[0], acl.Severity)
				assert.NotNil(t, err)
				result, err := ovndbapi.ExecuteR(cmd)
				assert.Nil(t, err)
				assert.Len(t, result, 0)
			})
		}
	}

	// Test ACLSetName()
	for i, e := range aclEntityTests {
		for j, tc := range e.aclTests {
			tcName := fmt.Sprintf("set ACL match for %s %s #%d", e.entityType, e.entity, j)
			t.Run(tcName, func(t *testing.T) {
				acl := &aclEntityTests[i].aclTests[j].acl
				cmd, err = ovndbapi.ACLSetName(acl.UUID, tc.newName)
				assert.Nil(t, err)
				err = ovndbapi.Execute(cmd)
				assert.Nil(t, err)
				acls, err := ovndbapi.ACLListEntity(e.entityType, e.entity)
				assert.Nil(t, err)
				acl.Name = tc.newName
				assert.True(t, containsACL(acls, acl))
			})
		}
	}

	// Test ACLSetMatch()
	for i, e := range aclEntityTests {
		for j, tc := range e.aclTests {
			tcName := fmt.Sprintf("set ACL match for %s %s #%d", e.entityType, e.entity, j)
			t.Run(tcName, func(t *testing.T) {
				acl := &aclEntityTests[i].aclTests[j].acl
				cmd, err = ovndbapi.ACLSetMatch(acl.UUID, tc.newMatch)
				assert.Nil(t, err)
				err = ovndbapi.Execute(cmd)
				assert.Nil(t, err)
				acls, err := ovndbapi.ACLListEntity(e.entityType, e.entity)
				assert.Nil(t, err)
				acl.Match = tc.newMatch
				assert.True(t, containsACL(acls, acl))
			})
		}
	}

	// Test ACLSetLogging()
	for i, e := range aclEntityTests {
		for j, tc := range e.aclTests {
			tcName := fmt.Sprintf("set ACL match for %s %s #%d", e.entityType, e.entity, j)
			t.Run(tcName, func(t *testing.T) {
				acl := &aclEntityTests[i].aclTests[j].acl
				var newLogflag = false
				if tc.newMeter != "" {
					newLogflag = true
				}
				cmd, err = ovndbapi.ACLSetLogging(acl.UUID, newLogflag, tc.newMeter, tc.newSeverity)
				assert.Nil(t, err)
				err = ovndbapi.Execute(cmd)
				assert.Nil(t, err)

				acl.Log = newLogflag
				if newLogflag {
					acl.Meter[0] = tc.newMeter
					switch tc.newSeverity {
					case "alert", "debug", "info", "notice", "warning":
						acl.Severity = tc.newSeverity
					case "":
						acl.Severity = "info"
					default:
						assert.Error(t, ErrorOption)
					}
				}

				acls, err := ovndbapi.ACLListEntity(e.entityType, e.entity)
				assert.Nil(t, err)
				assert.True(t, containsACL(acls, acl))
			})
		}
	}

	t.Run("set/delete non-existent ACL", func(t *testing.T) {
		cmd, err = ovndbapi.ACLSetLogging(NONEXISTENT_UUID, true, METER1, SEVERITY_INFO)
		assert.NotNil(t, err)
		cmd, err = ovndbapi.ACLSetMatch(NONEXISTENT_UUID, MATCH)
		assert.NotNil(t, err)
		cmd, err = ovndbapi.ACLSetName(NONEXISTENT_UUID, ACL_NAME_1)
		assert.NotNil(t, err)
		cmd, err = ovndbapi.ACLDelEntity(PORT_GROUP, PG_TEST_PG1, NONEXISTENT_UUID)
		assert.NotNil(t, err)
	})

	t.Run("delete ACL from non-existent entities", func(t *testing.T) {
		cmd, err = ovndbapi.ACLDelEntity(PORT_GROUP, PG_TEST_PG2, aclEntityTests[0].aclTests[0].acl.UUID)
		assert.NotNil(t, err)
		cmd, err = ovndbapi.ACLDelEntity(PORT_GROUP, LSW2, aclEntityTests[0].aclTests[0].acl.UUID)
		assert.NotNil(t, err)
	})

	// Test ACLDelEntity()
	for i, e := range aclEntityTests {
		for j, tc := range e.aclTests {
			tcName := fmt.Sprintf("delete ACL from %s %s #%d", e.entityType, e.entity, j)
			t.Run(tcName, func(t *testing.T) {
				cmd, err = ovndbapi.ACLDelEntity(e.entityType, e.entity, tc.acl.UUID)
				assert.Nil(t, err)
				result, err := ovndbapi.ExecuteR(cmd)
				assert.Nil(t, err)
				assert.Len(t, result, 0)
				acls, err := ovndbapi.ACLListEntity(e.entityType, e.entity)
				assert.Nil(t, err)
				assert.Equal(t, len(aclEntityTests[i].aclTests)-1-j, len(acls))
				assert.False(t, containsACL(acls, &tc.acl))
			})
		}
	}

	// Delete everything else
	t.Run("delete meter, switch, ports and port group used for ACL testing", func(t *testing.T) {
		cmd, err = ovndbapi.MeterDel(METER1)
		assert.Nil(t, err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(t, err)

		cmd, err = ovndbapi.MeterDel(METER2)
		assert.Nil(t, err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(t, err)

		cmd, err = ovndbapi.LSDel(PG_TEST_LS1)
		assert.Nil(t, err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(t, err)

		cmd, err = ovndbapi.PortGroupDel(PG_TEST_PG1)
		assert.Nil(t, err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(t, err)
	})
}
