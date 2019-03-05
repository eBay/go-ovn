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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	OVS_RUNDIR   = "/var/run/openvswitch"
	OVNNB_SOCKET = "ovnnb_db.sock"
	LSW          = "TEST_LSW"
	LSP          = "TEST_LSP"
	LSP_SECOND   = "TEST_LSP_SECOND "
	ADDR         = "36:46:56:76:86:96 127.0.0.1"
	MATCH        = "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.61"
	MATCH_SECOND = "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.62"
)

var ovndbapi OVNDBApi

func init() {
	var ovs_rundir = os.Getenv("OVS_RUNDIR")
	if ovs_rundir == "" {
		ovs_rundir = OVS_RUNDIR
	}
	var socket = ovs_rundir + "/" + OVNNB_SOCKET
	ovndbapi = GetInstance(socket, UNIX, "", 0, nil)
}

func TestACLs(t *testing.T) {
	var c []*OvnCommand = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.LSWAdd(LSW))
	c = append(c, ovndbapi.LSPAdd(LSW, LSP))

	c = append(c, ovndbapi.LSPSetAddress(LSP, ADDR))
	c = append(c, ovndbapi.LSPSetPortSecurity(LSP, ADDR))
	c = append(c, ovndbapi.ACLAdd(LSW, "to-lport", MATCH, "drop", 1001, nil, true, ""))
	ovndbapi.Execute(c...)

	lsps := ovndbapi.GetLogicPortsBySwitch(LSW)
	assert.Equal(t, true, len(lsps) == 1 && lsps[0].Name == LSP, "test[%s]: %v", "added port", lsps)
	assert.Equal(t, true, len(lsps) == 1 && lsps[0].Addresses[0] == ADDR, "test[%s]", "setted port address")
	assert.Equal(t, true, len(lsps) == 1 && lsps[0].PortSecurity[0] == ADDR, "test[%s]", "setted port port security")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.LSPAdd(LSW, LSP_SECOND))
	ovndbapi.Execute(c...)
	lsps = ovndbapi.GetLogicPortsBySwitch(LSW)
	assert.Equal(t, true, len(lsps) == 2, "test[%s]: %+v", "added 2 ports", lsps)

	acls := ovndbapi.GetACLsBySwitch(LSW)
	assert.Equal(t, true, len(acls) == 1 && acls[0].Match == MATCH &&
		acls[0].Action == "drop" && acls[0].Priority == 1001 && acls[0].Log == true, "test[%s] %s", "add acl", acls[0])

	assert.Equal(t, true, nil == ovndbapi.ACLAdd(LSW, "to-lport", MATCH, "drop", 1001, nil, true, ""),
		"test[%s]", "add same acl twice, should only one added.")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ACLAdd(LSW, "to-lport", MATCH_SECOND, "drop", 1001, map[string]string{"A": "a", "B": "b"}, false, ""))
	ovndbapi.Execute(c...)
	acls = ovndbapi.GetACLsBySwitch(LSW)
	assert.Equal(t, true, len(acls) == 2, "test[%s]", "add second acl")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ACLAdd(LSW, "to-lport", MATCH_SECOND, "drop", 1001, map[string]string{"A": "b", "B": "b"}, false, ""))
	ovndbapi.Execute(c...)
	acls = ovndbapi.GetACLsBySwitch(LSW)
	assert.Equal(t, true, len(acls) == 3, "test[%s]", "add second acl")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ACLDel(LSW, "to-lport", MATCH, 1001, map[string]string{}))
	ovndbapi.Execute(c...)
	acls = ovndbapi.GetACLsBySwitch(LSW)
	assert.Equal(t, true, len(acls) == 2, "test[%s]", "acl remove")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ACLDel(LSW, "to-lport", MATCH_SECOND, 1001, map[string]string{"A": "a"}))
	ovndbapi.Execute(c...)
	acls = ovndbapi.GetACLsBySwitch(LSW)
	assert.Equal(t, true, len(acls) == 1, "test[%s]", "acl remove")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ACLDel(LSW, "to-lport", MATCH_SECOND, 1001, map[string]string{"A": "b"}))
	ovndbapi.Execute(c...)
	acls = ovndbapi.GetACLsBySwitch(LSW)
	assert.Equal(t, true, len(acls) == 2, "test[%s]", "acl remove")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.LSPDel(LSP))
	ovndbapi.Execute(c...)
	lsps = ovndbapi.GetLogicPortsBySwitch(LSW)
	assert.Equal(t, true, len(lsps) == 1, "test[%s]", "one port remove")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.LSPDel(LSP_SECOND))
	ovndbapi.Execute(c...)
	lsps = ovndbapi.GetLogicPortsBySwitch(LSW)
	assert.Equal(t, true, len(lsps) == 0, "test[%s]", "one port remove")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.LSWDel(LSW))
	ovndbapi.Execute(c...)

}

func findAS(name string) bool {
	as := ovndbapi.GetAddressSets()
	for _, a := range as {
		if a.Name == name {
			return true
		}
	}
	return false
}

func addressSetCmp(asname string, targetvalue []string) bool {
	as := ovndbapi.GetAddressSets()
	for _, a := range as {
		if a.Name == asname {
			if len(a.Addresses) == len(targetvalue) {
				addressSetMap := map[string]bool{}
				for _, i := range a.Addresses {
					addressSetMap[i] = true
				}
				for _, t := range targetvalue {
					if _, ok := addressSetMap[t]; !ok {
						return false
					}
				}
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func TestAddressSet(t *testing.T) {
	addressList := []string{"127.0.0.1"}
	var c []*OvnCommand = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ASAdd("AS1", addressList, map[string]string{}))
	ovndbapi.Execute(c...)
	as := ovndbapi.GetAddressSets()
	assert.Equal(t, true, addressSetCmp("AS1", addressList), "test[%s] and value[%v]", "address set 1 added.", as[0].Addresses)

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ASAdd("AS2", addressList, map[string]string{}))
	ovndbapi.Execute(c...)
	as = ovndbapi.GetAddressSets()
	assert.Equal(t, true, addressSetCmp("AS2", addressList), "test[%s] and value[%v]", "address set 2 added.", as[1].Addresses)

	addressList = []string{"127.0.0.4", "127.0.0.5", "127.0.0.6"}
	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ASUpdate("AS2", addressList, map[string]string{}))
	ovndbapi.Execute(c...)
	as = ovndbapi.GetAddressSets()
	assert.Equal(t, true, addressSetCmp("AS2", addressList), "test[%s] and value[%v]", "address set added with different list.", as[0].Addresses)

	addressList = []string{"127.0.0.4", "127.0.0.5"}
	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ASUpdate("AS2", addressList, map[string]string{}))
	ovndbapi.Execute(c...)
	as = ovndbapi.GetAddressSets()
	assert.Equal(t, true, addressSetCmp("AS2", addressList), "test[%s] and value[%v]", "address set updated.", as[0].Addresses)

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ASDel("AS1"))
	ovndbapi.Execute(c...)
	assert.Equal(t, false, findAS("AS1"), "test AS remove")

	c = make([]*OvnCommand, 0)
	c = append(c, ovndbapi.ASDel("AS2"))
	ovndbapi.Execute(c...)
	assert.Equal(t, false, findAS("AS2"), "test AS remove")
}
