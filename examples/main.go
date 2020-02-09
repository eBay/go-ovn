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

package main

import (
	"fmt"
	"os"

	goovn "github.com/ebay/go-ovn"
)

const (
	ovsRundir   = "/var/run/openvswitch"
	ovnnbSocket = "ovnnb_db.sock"
	matchFirst  = "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.61"
	matchSecond = "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.62"
)

var ovndbapi goovn.Client

func init() {
	var err error
	var ovs_rundir = os.Getenv("OVS_RUNDIR")
	if ovs_rundir == "" {
		ovs_rundir = ovsRundir
	}
	ovndbapi, err = goovn.NewClient(&goovn.Config{Addr: "unix:" + ovs_rundir + "/" + ovnnbSocket},
	"OVN_Northbound")
	if err != nil {
		panic(err)
	}
}

func main() {

	ocmd, _ := ovndbapi.LSAdd("ls1")
	ovndbapi.Execute(ocmd)
	ocmd, _ = ovndbapi.LSPAdd("ls1", "test")
	ovndbapi.Execute(ocmd)
	ocmd, _ = ovndbapi.LSPSetAddress("test", "12:34:56:78:90 10.10.10.1")
	ovndbapi.Execute(ocmd)

	lports, _ := ovndbapi.LSPList("ls1")
	for _, lp := range lports {
		fmt.Printf("%v\n", *lp)
	}

	ocmd, _ = ovndbapi.ACLAdd("ls1", "to-lport", matchFirst, "drop", 1001, nil, true, "", "")
	ovndbapi.Execute(ocmd)

	ocmd, _ = ovndbapi.ACLAdd("ls1", "to-lport", matchSecond, "drop", 1001, map[string]string{"A": "a", "B": "b"}, false, "", "")
	ovndbapi.Execute(ocmd)

	ocmd, _ = ovndbapi.ACLAdd("ls1", "to-lport", matchSecond, "drop", 1001, map[string]string{"A": "b", "B": "b"}, false, "", "")
	ovndbapi.Execute(ocmd)

	acls, _ := ovndbapi.ACLList("ls1")
	for _, acl := range acls {
		fmt.Printf("%v\n", *acl)
	}
	fmt.Println()

	ocmd, _ = ovndbapi.ACLDel("ls1", "to-lport", matchFirst, 1001, map[string]string{})
	ovndbapi.Execute(ocmd)
	acls, _ = ovndbapi.ACLList("ls1")
	for _, acl := range acls {
		fmt.Printf("%v\n", *acl)
	}

	fmt.Println()
	ocmd, _ = ovndbapi.ACLDel("ls1", "to-lport", matchFirst, 1001, map[string]string{"A": "a"})
	ovndbapi.Execute(ocmd)
	acls, _ = ovndbapi.ACLList("ls1")
	for _, acl := range acls {
		fmt.Printf("%v\n", *acl)
	}
	fmt.Println()
	ocmd, _ = ovndbapi.ACLDel("ls1", "to-lport", matchSecond, 1001, map[string]string{"A": "b"})
	ovndbapi.Execute(ocmd)
	acls, _ = ovndbapi.ACLList("ls1")
	for _, acl := range acls {
		fmt.Printf("%v\n", *acl)
	}
	ocmd, _ = ovndbapi.LSPDel("test")
	ovndbapi.Execute(ocmd)
	ocmd, _ = ovndbapi.LSDel("ls1")
	ovndbapi.Execute(ocmd)
}
