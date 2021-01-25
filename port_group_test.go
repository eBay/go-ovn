/**
 * Copyright (c) 2020 Red Hat Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless assertd by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 **/

package goovn

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"testing"
)

type pgConfig struct {
	pgName      string
	ports       []string
	externalIds map[string]string
}

type pgTest struct {
name        string
testFunc    func(group string, ports []string, external_ids map[string]string) (*OvnCommand, error)
startConfig pgConfig
testConfig  pgConfig
}

func compareExternalIds(want map[string]string, got map[interface{}]interface{}) bool {
	if len(want) != len(got) {
		return false
	}
	for key, w := range want {
		if w != got[key] {
			return false
		}
	}
	return true
}

func TestPortGroupAPI(t *testing.T) {
	ovndbapi := getOVNClient(DBNB)
	assert := assert.New(t)
	var cmd *OvnCommand
	var err error
	var lsp1UUID, lsp2UUID, lsp3UUID, lsp4UUID string

	// create Switch with four ports
	createSwitch := func(t *testing.T) {
		t.Helper()
		var cmds []*OvnCommand
		// Create Switch
		cmd, err := ovndbapi.LSAdd(PG_TEST_LS1)
		assert.Nil(err)
		cmds = append(cmds, cmd)
		// Add ports
		cmd, err = ovndbapi.LSPAdd(PG_TEST_LS1, PG_TEST_LSP1)
		assert.Nil(err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPAdd(PG_TEST_LS1, PG_TEST_LSP2)
		assert.Nil(err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPAdd(PG_TEST_LS1, PG_TEST_LSP3)
		assert.Nil(err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPAdd(PG_TEST_LS1, PG_TEST_LSP4)
		assert.Nil(err)
		cmds = append(cmds, cmd)

		result, err := ovndbapi.ExecuteR(cmds...)
		assert.Nil(err)

		assert.Equal(5, len(result))
		lsp1UUID = result[1]
		lsp2UUID = result[2]
		lsp3UUID = result[3]
		lsp4UUID = result[4]
	}

	// Delete Switch
	deleteSwitch := func(t *testing.T) {
		t.Helper()

		cmd, err = ovndbapi.LSDel(PG_TEST_LS1)
		assert.Nil(err)

		err = ovndbapi.Execute(cmd)
		assert.Nil(err)
	}

	// Create a switch w/four ports to be used for logical port tests
	createSwitch(t)

	portGroupAddTests := []pgTest {
		{
			name: "add an empty port group",
			testFunc: ovndbapi.PortGroupAdd,
			startConfig: pgConfig{
				pgName:"",
				ports: nil,
				externalIds: nil,
			},
			testConfig: pgConfig{
				pgName: PG_TEST_PG1,
				ports: nil,
				externalIds: nil,
			},
		},
		{
			name: "add a port group with two ports",
			testFunc: ovndbapi.PortGroupAdd,
			startConfig: pgConfig{
				pgName:"",
				ports:nil,
				externalIds:nil,
			},
			testConfig: pgConfig{
				pgName: PG_TEST_PG1,
				ports: []string{lsp1UUID, lsp2UUID},
				externalIds: nil,
			},
		},
		{
			name: "add a port group with four ports",
			testFunc: ovndbapi.PortGroupAdd,
			startConfig: pgConfig{
				pgName:"",
				ports:nil,
				externalIds:nil,
			},
			testConfig: pgConfig{
				pgName: PG_TEST_PG1,
				ports: []string{lsp1UUID, lsp2UUID, lsp3UUID, lsp4UUID},
				externalIds: nil,
			},
		},
		{
			name: "add two port groups with the same ports",
			testFunc: ovndbapi.PortGroupAdd,
			startConfig: pgConfig{
				pgName: PG_TEST_PG1,
				ports: []string{lsp1UUID, lsp2UUID},
				externalIds: nil,
			},
			testConfig: pgConfig{
				pgName: PG_TEST_PG2,
				ports: []string{lsp1UUID, lsp2UUID},
				externalIds: nil,
			},
		},
		{
			name: "add a port group with external ids",
			testFunc: ovndbapi.PortGroupAdd,
			startConfig: pgConfig{
				pgName:"",
				ports:nil,
				externalIds:nil,
			},
			testConfig: pgConfig{
				pgName: PG_TEST_PG1,
				ports: nil,
				externalIds: map[string]string{PG_TEST_KEY_1: PG_TEST_ID_1, PG_TEST_KEY_2: PG_TEST_ID_2},
			},
		},
		{
			name: "add a port group with ports and external ids",
			testFunc: ovndbapi.PortGroupAdd,
			startConfig: pgConfig{
				pgName:"",
				ports:nil,
				externalIds:nil,
			},
			testConfig: pgConfig{
				pgName: PG_TEST_PG1,
				ports: []string{lsp1UUID, lsp2UUID},
				externalIds: map[string]string{PG_TEST_KEY_1: PG_TEST_ID_1, PG_TEST_KEY_2: PG_TEST_ID_2},
			},
		},
		{
			name: "add a port group with empty name",
			testFunc: ovndbapi.PortGroupAdd,
			startConfig: pgConfig{
				pgName:"",
				ports:nil,
				externalIds:nil,
			},
			testConfig: pgConfig{
				pgName: "",
				ports:nil,
				externalIds:nil,
			},
		},
		{
			name: "add a port group when another exists",
			testFunc: ovndbapi.PortGroupAdd,
			startConfig: pgConfig{
				pgName: PG_TEST_PG1,
				ports: []string{lsp1UUID, lsp2UUID},
				externalIds: map[string]string{PG_TEST_KEY_1: PG_TEST_ID_1, PG_TEST_KEY_2: PG_TEST_ID_2},
			},
			testConfig: pgConfig{
				pgName: PG_TEST_PG2,
				ports: []string{lsp3UUID, lsp4UUID},
				externalIds: map[string]string{PG_TEST_KEY_1: PG_TEST_ID_1, PG_TEST_KEY_2: PG_TEST_ID_2},
			},
		},
		{
			name: "set ports and external ids on existing empty port group",
			testFunc: ovndbapi.PortGroupUpdate,
			startConfig: pgConfig{
				pgName:PG_TEST_PG1,
				ports:nil,
				externalIds:nil,
			},
			testConfig: pgConfig{
				pgName: PG_TEST_PG1,
				ports: []string{lsp1UUID, lsp2UUID},
				externalIds: map[string]string{PG_TEST_KEY_1: PG_TEST_ID_1, PG_TEST_KEY_2: PG_TEST_ID_2},
			},
		},
		{
			name: "set ports and external ids on existing port group with exiting config",
			testFunc: ovndbapi.PortGroupUpdate,
			startConfig: pgConfig{
				pgName:PG_TEST_PG1,
				ports: []string{lsp1UUID, lsp2UUID},
				externalIds: map[string]string{PG_TEST_KEY_1: PG_TEST_ID_1},
			},
			testConfig: pgConfig{
				pgName:PG_TEST_PG1,
				ports: []string{lsp3UUID, lsp4UUID},
				externalIds: map[string]string{PG_TEST_KEY_2: PG_TEST_ID_2, PG_TEST_KEY_3: PG_TEST_ID_3},
			},
		},
	}

	for _, tc := range portGroupAddTests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.startConfig.pgName != "" {
				// Add start config
				cmd, err = ovndbapi.PortGroupAdd(tc.startConfig.pgName, tc.startConfig.ports, tc.startConfig.externalIds)
				assert.Nil(err)
				err = ovndbapi.Execute(cmd)
				assert.Nil(err)
			}

			// Add or Set the port group
			cmd, err = tc.testFunc(tc.testConfig.pgName, tc.testConfig.ports, tc.testConfig.externalIds)
			assert.Nil(err)
			err = ovndbapi.Execute(cmd)
			assert.Nil(err)

			// Validate port group
			pg, err := ovndbapi.PortGroupGet(tc.testConfig.pgName)
			assert.Nil(err)
			assert.NotNil(pg)
			assert.Equal(pg.Name, tc.testConfig.pgName)
			if tc.testConfig.ports != nil {
				sort.Strings(tc.testConfig.ports)
				sort.Strings(pg.Ports)
				assert.True(reflect.DeepEqual(tc.testConfig.ports, pg.Ports))
			}
			if tc.testConfig.externalIds != nil {
				assert.True(compareExternalIds(tc.testConfig.externalIds, pg.ExternalID))
			}

			// Delete the port group
			cmd, err = ovndbapi.PortGroupDel(tc.testConfig.pgName)
			assert.Nil(err)
			err = ovndbapi.Execute(cmd)
			assert.Nil(err)

			// Confirm that it's deleted
			_, err = ovndbapi.PortGroupGet(tc.testConfig.pgName)
			assert.NotNil(err)

			if (tc.startConfig.pgName != "") &&  (tc.startConfig.pgName != tc.testConfig.pgName) {
				// Delete the start config port group.
				cmd, err = ovndbapi.PortGroupDel(tc.startConfig.pgName)
				assert.Nil(err)
				err = ovndbapi.Execute(cmd)
				assert.Nil(err)
			}
		})
	}

	// The following are negative/boundary cases that are not as easy to generalize

	t.Run("add duplicate port group", func(t *testing.T) {
		// Add first port group
		cmd, err = ovndbapi.PortGroupAdd(PG_TEST_PG1, []string{lsp1UUID, lsp2UUID}, nil)
		assert.Nil(err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(err)

		cmd, err = ovndbapi.PortGroupAdd(PG_TEST_PG1, []string{lsp1UUID, lsp2UUID}, nil)
		assert.NotNil(err)

		// Delete the first port group
		cmd, err = ovndbapi.PortGroupDel(PG_TEST_PG1)
		assert.Nil(err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(err)
	})

	t.Run("add port group with non-existent ports", func(t *testing.T) {
		badUUID, err := uuid.NewRandom()
		assert.Nil(err)
		ports := []string{lsp1UUID, badUUID.String()}

		// Add the port group
		cmd, err = ovndbapi.PortGroupAdd(PG_TEST_PG1, ports, nil)
		assert.Nil(err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(err)
		// The way this currently works, ovsdb doesn't care whether the ports exist,
		// and will add the port group regardless. Should it?

		// Validate port group
		pg, err := ovndbapi.PortGroupGet(PG_TEST_PG1)
		assert.Nil(err)
		assert.NotNil(pg)
		assert.Equal(pg.Name, PG_TEST_PG1)
		sort.Strings(ports)
		sort.Strings(pg.Ports)
		// Don't expect the badUUID to be in the list
		assert.False(reflect.DeepEqual(ports, pg.Ports))

		// Delete the port group
		cmd, err = ovndbapi.PortGroupDel(PG_TEST_PG1)
		assert.Nil(err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(err)
	})

	t.Run("add/delete ports to/from port group", func(t *testing.T) {
		ports := []string{lsp1UUID}

		// Add the port group
		cmd, err = ovndbapi.PortGroupAdd(PG_TEST_PG1, ports, nil)
		assert.Nil(err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(err)

		cmd, err = ovndbapi.PortGroupAddPort(PG_TEST_PG1, lsp2UUID)
		assert.Nil(err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(err)

		// Validate port group
		ports = append(ports, lsp2UUID)
		pg, err := ovndbapi.PortGroupGet(PG_TEST_PG1)
		assert.Nil(err)
		assert.NotNil(pg)
		sort.Strings(ports)
		sort.Strings(pg.Ports)
		assert.True(reflect.DeepEqual(ports, pg.Ports))

		// Add duplicate port
		cmd, err = ovndbapi.PortGroupAddPort(PG_TEST_PG1, lsp2UUID)
		assert.Equal(ErrorExist, err)

		// Add port to non-existent port group
		cmd, err = ovndbapi.PortGroupAddPort(PG_TEST_PG2, lsp1UUID)
		assert.Equal(ErrorNotFound, err)

		// Remove lsp2 from port group
		cmd, err = ovndbapi.PortGroupRemovePort(PG_TEST_PG1, lsp2UUID)
		assert.Nil(err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(err)
		pg, err = ovndbapi.PortGroupGet(PG_TEST_PG1)
		assert.Nil(err)
		ports = []string{lsp1UUID}
		assert.True(reflect.DeepEqual(ports, pg.Ports))

		// Remove lsp1 from port group
		cmd, err = ovndbapi.PortGroupRemovePort(PG_TEST_PG1, lsp1UUID)
		assert.Nil(err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(err)
		pg, err = ovndbapi.PortGroupGet(PG_TEST_PG1)
		assert.Nil(err)
		assert.True(len(pg.Ports) == 0)

		// Delete the port group
		cmd, err = ovndbapi.PortGroupDel(PG_TEST_PG1)
		assert.Nil(err)
		err = ovndbapi.Execute(cmd)
		assert.Nil(err)
	})

	t.Run("set port group that doesn't exist", func(t *testing.T) {
		cmd, err = ovndbapi.PortGroupUpdate(PG_TEST_PG1, []string{lsp1UUID, lsp2UUID}, nil)
		assert.NotNil(err)
	})

	deleteSwitch(t)
}
