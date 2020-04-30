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
	PORT_TEST_EXT_ID_MAC_2   = "01:02:03:05:05:06"
	PORT_TEST_EXT_ID_IP_KEY  = "ip_addr"
	PORT_TEST_EXT_ID_IP      = "169.254.1.1"
	PORT_TEST_OPT_1_KEY      = "foo1"
	PORT_TEST_OPT_1_VAL      = "bar1"
	PORT_TEST_OPT_1_VAL_2    = "baz1"
	PORT_TEST_OPT_2_KEY      = "foo2"
	PORT_TEST_OPT_2_VAL      = "bar2"
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

	// test external id APIs
	extIds, err := ovndbapi.LSPGetExternalIds(PORT_TEST_LSP1)
	assert.Equal(t, extIds != nil, true)

	t.Logf("Setting external ids with two keys")
	extIds[PORT_TEST_EXT_ID_MAC_KEY] = PORT_TEST_EXT_ID_MAC
	extIds[PORT_TEST_EXT_ID_IP_KEY] = PORT_TEST_EXT_ID_IP

	cmd, err = ovndbapi.LSPSetExternalIds(PORT_TEST_LSP1, extIds)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Validating the external ids are set correctly")
	extIdsRet, err := ovndbapi.LSPGetExternalIds(PORT_TEST_LSP1)
	if err != nil {
		t.Fatal(err)
	}

	extIdMac, ok := extIdsRet[PORT_TEST_EXT_ID_MAC_KEY]

	assert.Equal(t, ok, true)
	assert.Equal(t, extIdMac, PORT_TEST_EXT_ID_MAC)

	extIdIP, ok := extIdsRet[PORT_TEST_EXT_ID_IP_KEY]
	assert.Equal(t, ok, true)
	assert.Equal(t, extIdIP, PORT_TEST_EXT_ID_IP)

	t.Logf("Validated the external ids are set correctly")
	// update one of the existing keys, remove one key
	// and make sure it's a complete update.
	t.Logf("Validating that keys get clobbered by Set API correctly")
	extIds[PORT_TEST_EXT_ID_MAC_KEY] = PORT_TEST_EXT_ID_MAC_2
	delete(extIds, PORT_TEST_EXT_ID_IP_KEY)
	t.Logf("Remove one of the external ids and verify the contents")
	cmd, err = ovndbapi.LSPSetExternalIds(PORT_TEST_LSP1, extIds)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	extIdsRet, err = ovndbapi.LSPGetExternalIds(PORT_TEST_LSP1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Validating that only the MAC key is present in the external ids")
	extIdMac, ok = extIdsRet[PORT_TEST_EXT_ID_MAC_KEY]

	assert.Equal(t, ok, true)
	assert.Equal(t, extIdMac, PORT_TEST_EXT_ID_MAC_2)

	t.Logf("Validated that the MAC key's value has updated in the external ids")

	// make sure IP key is not present in external ids
	extIdIP, ok = extIdsRet[PORT_TEST_EXT_ID_IP_KEY]
	assert.Equal(t, !ok, true)

	t.Logf("Validated that the IP key is not present in the external ids")

	// test options API
	options, err := ovndbapi.LSPGetOptions(PORT_TEST_LSP1)
	assert.Equal(t, options != nil, true)

	t.Logf("Validating the options are set correctly")

	options[PORT_TEST_OPT_1_KEY] = PORT_TEST_OPT_1_VAL
	options[PORT_TEST_OPT_2_KEY] = PORT_TEST_OPT_2_VAL

	cmd, err = ovndbapi.LSPSetOptions(PORT_TEST_LSP1, options)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	optionsRet, err := ovndbapi.LSPGetOptions(PORT_TEST_LSP1)
	if err != nil {
		t.Fatal(err)
	}

	option1, ok := optionsRet[PORT_TEST_OPT_1_KEY]

	assert.Equal(t, ok, true)
	assert.Equal(t, option1, PORT_TEST_OPT_1_VAL)

	option2, ok := optionsRet[PORT_TEST_OPT_2_KEY]
	assert.Equal(t, ok, true)
	assert.Equal(t, option2, PORT_TEST_OPT_2_VAL)

	t.Logf("Validated that multiple options are set correctly")

	// update one of the existing keys, remove one key
	// and make sure it's a complete update.
	t.Logf("Validating that keys get clobbered by Set API correctly")

	options[PORT_TEST_OPT_1_KEY] = PORT_TEST_OPT_1_VAL_2
	delete(options, PORT_TEST_OPT_2_KEY)
	t.Logf("Remove one of the options and verify the contents")

	cmd, err = ovndbapi.LSPSetOptions(PORT_TEST_LSP1, options)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	optionsRet, err = ovndbapi.LSPGetOptions(PORT_TEST_LSP1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Validating that only the OPT_1 key is present in the options")

	option1, ok = optionsRet[PORT_TEST_OPT_1_KEY]

	assert.Equal(t, ok, true)
	assert.Equal(t, option1, PORT_TEST_OPT_1_VAL_2)
	t.Logf("Validated that the OPT_1 key's value has updated in the options")

	// make sure IP key is not present in external ids
	option2, ok = optionsRet[PORT_TEST_OPT_2_KEY]
	assert.Equal(t, !ok, true)
	t.Logf("Validated that the OPT_2 key is not present in the options")

	//validate that setting fields on an empty LSP string gives a nil cmd and an error
	cmd, err = ovndbapi.LSPSetDynamicAddresses(PORT_TEST_LSP3, PORT_TEST_LSP1DYNADDR1)
	assert.Equal(t, cmd == nil, true)
	assert.Equal(t, err == nil, false)

	cmd, err = ovndbapi.LSPSetExternalIds(PORT_TEST_LSP3, extIds)
	assert.Equal(t, cmd == nil, true)
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
