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
	"sync"

	"github.com/ebay/libovsdb"
)

const (
	opInsert string = "insert"
	opMutate string = "mutate"
	opDelete string = "delete"
	opSelect string = "select"
	opUpdate string = "update"
)

const (
	NBDB string = "OVN_Northbound"
)

const (
	tableNBGlobal                 string = "NB_Global"
	tableLogicalSwitch            string = "Logical_Switch"
	tableLogicalSwitchPort        string = "Logical_Switch_Port"
	tableAddressSet               string = "Address_Set"
	tablePortGroup                string = "Port_Group"
	tableLoadBalancer             string = "Load_Balancer"
	tableACL                      string = "ACL"
	tableLogicalRouter            string = "Logical_Router"
	tableQoS                      string = "QoS"
	tableMeter                    string = "Meter"
	tableMeterBand                string = "Meter_Band"
	tableLogicalRouterPort        string = "Logical_Router_Port"
	tableLogicalRouterStaticRoute string = "Logical_Router_Static_Route"
	tableNAT                      string = "NAT"
	tableDHCPOptions              string = "DHCP_Options"
	tableConnection               string = "Connection"
	tableDNS                      string = "DNS"
	tableSSL                      string = "SSL"
	tableGatewayChassis           string = "Gateway_Chassis"
)

// OVN supporter protocols
const (
	UNIX string = "unix"
	TCP  string = "tcp"
	SSL  string = "ssl"
)

type ovnDBClient struct {
	socket   string
	server   string
	port     int
	protocol string
	dbclient *libovsdb.OvsdbClient
}

type ovnDBImp struct {
	client     *ovnDBClient
	cache      map[string]map[string]libovsdb.Row
	cachemutex sync.RWMutex
	tranmutex  sync.Mutex
	callback   OVNSignal
}

type OVNDB struct {
	imp *ovnDBImp
}

var once sync.Once
var ovnDBApi OVNDBApi

func GetInstance(socketfile string, proto string, server string, port int, callback OVNSignal) (OVNDBApi, error) {
	var err error

	once.Do(func() {
		var dbapi *OVNDB

		switch proto {
		case UNIX:
			dbapi, err = newNBBySocket(socketfile, callback)
		case TCP:
			dbapi, err = newNBByServer(server, port, callback, TCP)
		case SSL:
			dbapi, err = newNBByServer(server, port, callback, SSL)
		default:
			err = fmt.Errorf("the protocol [%s] is not supported", proto)
		}
		if err != nil {
			return
		}
		ovnDBApi = dbapi
	})

	return ovnDBApi, err
}

func SetCallBack(c OVNDBApi, callback OVNSignal) {
	if c != nil {
		c.SetCallBack(callback)
	}
}
