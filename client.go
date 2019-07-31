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
	"sync"

	"github.com/ebay/libovsdb"
)

// Client ovnnb client
type Client struct {
	ovndb         *ovndb
	LogicalSwitch *lsImp
	LoadBalancer  *lbImp
	ACL           *aclImp
}

type ovndb struct {
	client       *libovsdb.OvsdbClient
	cache        map[string]map[string]interface{}
	cachemutex   sync.RWMutex
	tranmutex    sync.Mutex
	signalCB     OVNSignal
	disconnectCB OVNDisconnectedCallback
}

func (c *Client) Close() error {
	c.ovndb.client.Disconnect()
	return nil
}

func (c *Client) Execute(cmds ...*OvnCommand) error {
	return c.ovndb.execute(cmds...)
}

func NewClient(cfg *Config) (*Client, error) {
	odbi := &ovndb{
		cache:        make(map[string]map[string]interface{}),
		signalCB:     cfg.SignalCB,
		disconnectCB: cfg.DisconnectCB,
	}

	ovncli, err := libovsdb.Connect(cfg.Addr, cfg.TLSConfig)
	if err != nil {
		return nil, err
	}
	odbi.client = ovncli

	initial, err := odbi.client.MonitorAll(dbNB, "")
	if err != nil {
		return nil, err
	}

	odbi.populateCache(*initial)
	notifier := ovnNotifier{odbi}
	odbi.client.Register(notifier)

	cli := &Client{
		ovndb:         odbi,
		LogicalSwitch: &lsImp{odbi: odbi},
		LoadBalancer:  &lbImp{odbi: odbi},
		ACL:           &aclImp{odbi: odbi},
	}

	return cli, nil
}

func (c *ovndb) Execute(cmds ...*OvnCommand) error {
	return c.execute(cmds...)
}

func newRow() OVNRow {
	return make(OVNRow)
}
