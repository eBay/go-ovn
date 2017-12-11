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

	"github.com/ebay/go-ovn/goovn"
	"os"
)

const (
	OVS_RUNDIR = "/var/run/openvswitch"
	OVNNB_SOCKET = "ovnnb_db.sock"
)

var ovndbapi goovn.OVNDBApi

func init() {
	var ovs_rundir = os.Getenv("OVS_RUNDIR")
	if ovs_rundir == "" {
		ovs_rundir = OVS_RUNDIR
	}
	var socket = ovs_rundir + "/" + OVNNB_SOCKET
	ovndbapi = goovn.GetInstance(socket, goovn.UNIX, "", 0, nil)
}

func LSWAdd() {
	ocmd := ovndbapi.LSWAdd("ls1")
	ovndbapi.Execute(ocmd)
	fmt.Printf("return: %v", ocmd.Results)
}

func LSPAdd() {
	ocmd := ovndbapi.LSPAdd("ls1", "test")
	ovndbapi.Execute(ocmd)
}

func LSPSetAddress() {
	ocmd := ovndbapi.LSPSetAddress("test", "mac ip")
	ovndbapi.Execute(ocmd)
}

func LSPDel() {
	ocmd := ovndbapi.LSPDel("test")
	ovndbapi.Execute(ocmd)
}

func LSWDel() {
	ocmd := ovndbapi.LSWDel("ls1")
	ovndbapi.Execute(ocmd)
}

func ACLAdd() {
	ocmd := ovndbapi.ACLAdd("ls1", "to-lport", "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.61", "drop", 1001, nil, false)
	ovndbapi.Execute(ocmd)
}

func ACLDel() {
	ocmd := ovndbapi.ACLDel("ls1", "to-lport", "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.61", 1001)
	ovndbapi.Execute(ocmd)
}


func LISTLS() {
	ocmd := ovndbapi.LSWList()
	ovndbapi.Execute(ocmd)
	fmt.Printf("return: %v", ocmd.Results)
}


func main() {

	LSWAdd()
	LSPAdd()
	LSPSetAddress()
	ACLAdd()
	LISTLS()
	ACLDel()
	LSPDel()
	LSWDel()

}
