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
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
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

var ovndbapi Client

func TestMain(m *testing.M) {
	var api Client
	var err error

	cfg := &Config{}
	var ovs_rundir = os.Getenv("OVS_RUNDIR")
	if ovs_rundir == "" {
		ovs_rundir = OVS_RUNDIR
	}
	var ovn_nb_db = os.Getenv("OVN_NB_DB")
	if ovn_nb_db == "" {
		cfg.Addr = "unix:" + ovs_rundir + "/" + OVNNB_SOCKET
		api, err = NewClient(cfg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		strs := strings.Split(ovn_nb_db, ":")
		if len(strs) < 2 || len(strs) > 3 {
			log.Fatal("Unexpected format of $OVN_NB_DB")
		}
		if len(strs) == 2 {
			cfg.Addr = "unix:" + ovs_rundir + "/" + strs[1]
			api, err = NewClient(cfg)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			port, _ := strconv.Atoi(strs[2])
			cfg.Addr = fmt.Sprintf("%s:%s:%d", strs[0], strs[1], port)
			api, err = NewClient(cfg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	ovndbapi = api
	os.Exit(m.Run())
}
