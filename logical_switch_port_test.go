/**
 * Copyright (c) 2019 eBay Inc.
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

const (
	PORT_TEST_LS1            = "LogicalSwitch1"
	PORT_TEST_LSP1           = "LogicalSwitchPort1"
	PORT_TEST_LSP2           = "LogicalSwitchPort2"
	PORT_TEST_LSP3           = ""
	PORT_TEST_LSP1DYNADDR1   = "a.b.c.d"
	PORT_TEST_LSP2DYNADDR2   = ""
	PORT_TEST_EXT_ID_MAC_KEY = "mac_addr"
	PORT_TEST_EXT_ID_MAC     = "00:01:02:03:04:05"
)

func TestLogicalSwitchPortAPI(t *testing.T) {
	ovndbapi := getOVNClient(DBNB)
	// create Switch
	t.Logf("Adding  %s to OVN", PORT_TEST_LS1)
	cmd, err := ovndbapi.LSAdd(PORT_TEST_LS1)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	ls, err := ovndbapi.LSGet(PORT_TEST_LS1)
	if err != nil {
		t.Fatal(err)
	}

	if ls[0].Name != PORT_TEST_LS1 {
		t.Fatalf("ls not created %v", PORT_TEST_LS1)
	}

	// create logical switch port 1
	cmd, err = ovndbapi.LSPAdd(PORT_TEST_LS1, PORT_TEST_LSP1)
	if err != nil {
		t.Fatal(err)
	}

	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	// create logical switch port 2
	cmd, err = ovndbapi.LSPAdd(PORT_TEST_LS1, PORT_TEST_LSP2)
	if err != nil {
		t.Fatal(err)
	}

	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	// test LSPGet API
	lsp1, err := ovndbapi.LSPGet(PORT_TEST_LSP1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, lsp1.Name, PORT_TEST_LSP1)

	//check for set/get dynamic addresses with non-empty string
	cmd, err = ovndbapi.LSPSetDynamicAddresses(PORT_TEST_LSP1, PORT_TEST_LSP1DYNADDR1)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	// make sure the cache now shows the updated lsp object
	lsp1, err = ovndbapi.LSPGet(PORT_TEST_LSP1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, lsp1.DynamicAddresses, PORT_TEST_LSP1DYNADDR1)

	dynAddr1, err := ovndbapi.LSPGetDynamicAddresses(PORT_TEST_LSP1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, dynAddr1, PORT_TEST_LSP1DYNADDR1)

	//check for Set/Get Dynamic Address with empty string
	cmd, err = ovndbapi.LSPSetDynamicAddresses(PORT_TEST_LSP2, PORT_TEST_LSP2DYNADDR2)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lsp2, err := ovndbapi.LSPGet(PORT_TEST_LSP2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, lsp2.DynamicAddresses, PORT_TEST_LSP2DYNADDR2)

	extIds, err := ovndbapi.LSPGetExternalIds(PORT_TEST_LSP1)
	assert.Equal(t, extIds != nil, true)

	extIds[PORT_TEST_EXT_ID_MAC_KEY] = PORT_TEST_EXT_ID_MAC

	cmd, err = ovndbapi.LSPSetExternalIds(PORT_TEST_LSP1, extIds)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	extIdsRet, err := ovndbapi.LSPGetExternalIds(PORT_TEST_LSP1)
	if err != nil {
		t.Fatal(err)
	}

	extIdMac, ok := extIdsRet[PORT_TEST_EXT_ID_MAC_KEY]

	assert.Equal(t, ok, true)
	assert.Equal(t, extIdMac, PORT_TEST_EXT_ID_MAC)

	//validate that setting fields on an empty LSP string gives a nil cmd and an error
	cmd, err = ovndbapi.LSPSetDynamicAddresses(PORT_TEST_LSP3, PORT_TEST_LSP1DYNADDR1)
	assert.Equal(t, cmd, nil)
	assert.Equal(t, err == nil, false)

	cmd, err = ovndbapi.LSPSetExternalIds(PORT_TEST_LSP3, extIds)
	assert.Equal(t, cmd, nil)
	assert.Equal(t, err == nil, false)

	cmd, err = ovndbapi.LSDel(PORT_TEST_LS1)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

}
