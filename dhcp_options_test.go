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

func TestDHCPOptions(t *testing.T) {
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

	cmd, err = ovndbapi.LSPSetPortSecurity(LSP, ADDR)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)

	cmd, err = ovndbapi.DHCPOptionsAdd(
		"192.168.0.0/24",
		map[string]string{
			"server_id":  "192.168.1.1",
			"server_mac": "54:54:54:54:54:54",
			"lease_time": "6000",
		},
		nil)
	cmds = append(cmds, cmd)

	// execute to create lsw and lsp
	err = ovndbapi.Execute(cmds...)
	if err != nil {
		t.Fatal(err)
	}

	lsws, err := ovndbapi.LSList()
	if err != nil {
		t.Fatal(err)
	}

	if len(lsws) != 1 {
		t.Fatalf("ls not created %d", len(lsws))
	}

	dhcp_opts, err := ovndbapi.DHCPOptionsList()
	if err != nil {
		t.Fatal(err)
	}

	if len(dhcp_opts) != 1 {
		t.Fatalf("dhcp options not created %v", dhcp_opts)
	}

	cmd, err = ovndbapi.LSPSetDHCPv4Options(LSP, dhcp_opts[0].UUID)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	lsps, err := ovndbapi.LSPList(LSW)
	if err != nil {
		t.Fatal(err)
	}
	if len(lsps) != 1 {
		t.Fatalf("lsp not created %d", len(lsps))
	}

	assert.Equal(t, true, len(lsps) == 1 && lsps[0].Name == LSP, "test[%s]: %v", "added port", lsps)
	assert.Equal(t, true, len(lsps) == 1 && lsps[0].DHCPv4Options != "", "test[%s]", "setted dhcpv4_options")

	cmd, err = ovndbapi.DHCPOptionsDel(dhcp_opts[0].UUID)
	if err != nil {
		t.Fatal(err)
	}

	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	dhcp_opts, err = ovndbapi.DHCPOptionsList()
	if err != nil {
		t.Fatal(err)
	}
	if len(dhcp_opts) != 0 {
		t.Fatalf("dhcp options not deleted %#+v", dhcp_opts[0])
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

	assert.Equal(t, true, len(lsps) == 0, "test[%s]", "one port remove")

	cmd, err = ovndbapi.LSDel(LSW)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

}
