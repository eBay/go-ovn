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
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	OVS_RUNDIR   = "/var/run/openvswitch"
	OVNNB_SOCKET = "nb1.ovsdb"
	LR           = "TEST_LR"
	LRP          = "TEST_LRP"
	LSW          = "TEST_LSW"
	LSP          = "TEST_LSP"
	LSP_SECOND   = "TEST_LSP_SECOND "
	ADDR         = "36:46:56:76:86:96 127.0.0.1"
	MATCH        = "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.61"
	MATCH_SECOND = "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.62"
)

var ovndbapi OVNDBApi

func TestMain(m *testing.M) {
	var api OVNDBApi
	var err error

	var ovs_rundir = os.Getenv("OVS_RUNDIR")
	if ovs_rundir == "" {
		ovs_rundir = OVS_RUNDIR
	}
	var ovn_nb_db = os.Getenv("OVN_NB_DB")
	if ovn_nb_db == "" {
		var socket = ovs_rundir + "/" + OVNNB_SOCKET
		api, err = GetInstance(socket, UNIX, "", 0, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		strs := strings.Split(ovn_nb_db, ":")
		if len(strs) < 2 || len(strs) > 3 {
			log.Fatal("Unexpected format of $OVN_NB_DB")
		}
		if len(strs) == 2 {
			var socket = ovs_rundir + "/" + strs[1]
			api, err = GetInstance(socket, UNIX, "", 0, nil)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			port, _ := strconv.Atoi(strs[2])
			api, err = GetInstance("", strs[0], strs[1], port, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	ovndbapi = api
	os.Exit(m.Run())
}

func TestACLs(t *testing.T) {
	var cmds []*OvnCommand
	var cmd *OvnCommand
	var err error

	cmds = make([]*OvnCommand, 0)
	cmd, err = ovndbapi.LSWAdd(LSW)
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

	// execute to create lsw and lsp
	err = ovndbapi.Execute(cmds...)
	if err != nil {
		t.Fatal(err)
	}

	// nil cmds for next batch
	cmds = make([]*OvnCommand, 0)
	cmd, err = ovndbapi.ACLAdd(LSW, "to-lport", MATCH, "drop", 1001, nil, true, "")
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)

	err = ovndbapi.Execute(cmds...)
	if err != nil {
		t.Fatal(err)
	}

	lsws, err := ovndbapi.GetLogicalSwitches()
	if err != nil {
		t.Fatal(err)
	}
	if len(lsws) != 1 {
		t.Fatalf("ls not created %d", len(lsws))
	}

	lsps, err := ovndbapi.GetLogicalSwitchPortsBySwitch(LSW)
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

	lsps, err = ovndbapi.GetLogicalSwitchPortsBySwitch(LSW)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, len(lsps) == 2, "test[%s]: %+v", "added 2 ports", lsps)

	acls, err := ovndbapi.GetACLsBySwitch(LSW)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(acls) == 1 && acls[0].Match == MATCH &&
		acls[0].Action == "drop" && acls[0].Priority == 1001 && acls[0].Log == true, "test[%s] %s", "add acl", acls[0])

	cmd, err = ovndbapi.ACLAdd(LSW, "to-lport", MATCH, "drop", 1001, nil, true, "")
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

	acls, err = ovndbapi.GetACLsBySwitch(LSW)
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

	acls, err = ovndbapi.GetACLsBySwitch(LSW)
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

	acls, err = ovndbapi.GetACLsBySwitch(LSW)
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

	acls, err = ovndbapi.GetACLsBySwitch(LSW)
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

	acls, err = ovndbapi.GetACLsBySwitch(LSW)
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

	lsps, err = ovndbapi.GetLogicalSwitchPortsBySwitch(LSW)
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

	lsps, err = ovndbapi.GetLogicalSwitchPortsBySwitch(LSW)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, len(lsps) == 0, "test[%s]", "one port remove")

	cmd, err = ovndbapi.LSWDel(LSW)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

}

func findAS(name string) bool {
	as, err := ovndbapi.GetAddressSets()
	if err != nil {
		return false
	}
	for _, a := range as {
		if a.Name == name {
			return true
		}
	}
	return false
}

func addressSetCmp(asname string, targetvalue []string) bool {
	as, err := ovndbapi.GetAddressSets()
	if err != nil {
		return false
	}
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
	var cmd *OvnCommand
	var err error

	/*
		// can not call like:
		// ovndbapi.ASAdd("AS1", addressList, map[string][]{})
		// it will not be successful when input empty map.
	*/
	cmd, err = ovndbapi.ASAdd("AS1", addressList, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	as, err := ovndbapi.GetAddressSets()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, addressSetCmp("AS1", addressList), "test[%s] and value[%v]", "address set 1 added.", as[0].Addresses)

	cmd, err = ovndbapi.ASAdd("AS2", addressList, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	as, err = ovndbapi.GetAddressSets()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, addressSetCmp("AS2", addressList), "test[%s] and value[%v]", "address set 2 added.", as[1].Addresses)

	addressList = []string{"127.0.0.4", "127.0.0.5", "127.0.0.6"}
	cmd, err = ovndbapi.ASUpdate("AS2", addressList, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	as, err = ovndbapi.GetAddressSets()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, addressSetCmp("AS2", addressList), "test[%s] and value[%v]", "address set added with different list.", as[0].Addresses)

	addressList = []string{"127.0.0.4", "127.0.0.5"}
	cmd, err = ovndbapi.ASUpdate("AS2", addressList, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	as, err = ovndbapi.GetAddressSets()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, addressSetCmp("AS2", addressList), "test[%s] and value[%v]", "address set updated.", as[0].Addresses)

	cmd, err = ovndbapi.ASDel("AS1")
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, false, findAS("AS1"), "test AS remove")

	cmd, err = ovndbapi.ASDel("AS2")
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, false, findAS("AS2"), "test AS remove")
}

func TestLoadBalancer(t *testing.T) {
	t.Logf("Adding LB to OVN")
	ocmd, err := ovndbapi.LBAdd("lb1", "192.168.0.19:80", "tcp", []string{"10.0.0.11:80", "10.0.0.12:80"})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Adding LB OVN failed with err %v", err)
	}
	t.Logf("Adding LB to OVN Done")

	t.Logf("Updating LB to OVN")
	ocmd, err = ovndbapi.LBUpdate("lb1", "192.168.0.10:80", "tcp", []string{"10.10.10.127:8080", "10.10.10.120:8080"})
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Updating LB OVN failed with err %v", err)
	}
	t.Logf("Updating LB to OVN done")

	t.Logf("Gettting LB by name")
	lb, err := ovndbapi.GetLB("lb1")
	if err != nil {
		t.Fatal(err)
	}
	if len(lb) != 1 {
		t.Fatalf("err getting lbs, total:%v", len(lb))
	}
	t.Logf("Lb found:%+v", lb[0])

	t.Logf("Deleting LB")
	ocmd, err = ovndbapi.LBDel("lb1")
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}

	// Verify deletion
	lb, err = ovndbapi.GetLB("lb1")
	if err != nil {
		t.Fatal(err)
	}
	if len(lb) != 0 {
		t.Fatalf("error: lb deletion not done, total:%v", len(lb))
	}
	t.Logf("LB deletion done")
}
