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

func main() {
	ocmd := ovndbapi.LSWAdd("ls1")
	ovndbapi.Execute(ocmd)
	ocmd = ovndbapi.LSPAdd("ls1", "test")
	ovndbapi.Execute(ocmd)
	ocmd = ovndbapi.LSPSetAddress("test", "12:34:56:78:90 10.10.10.1")
	ovndbapi.Execute(ocmd)

	lports := ovndbapi.GetLogicPortsBySwitch("ls1")
	for _, lp := range(lports) {
		fmt.Printf("%v\n", *lp)
	}

	ocmd = ovndbapi.LSPDel("test")
	ovndbapi.Execute(ocmd)
	ocmd = ovndbapi.LSWDel("ls1")
	ovndbapi.Execute(ocmd)

}
