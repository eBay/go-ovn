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
	"strings"

	"github.com/ebay/libovsdb"
)

type LoadBalancer interface {
	// Add load balancer with options
	Add(...LoadBalancerOpt) (*OvnCommand, error)
	// Del load balancer with LoadBalancerName and optional LoadBalancerVIP
	Del(...LoadBalancerOpt) (*OvnCommand, error)
	// Get logical switch with LogicalSwitchName or LogicalSwitchUUID
	Get(...LoadBalancerOpt) (*libovsdb.LoadBalancer, error)
	// List load balancers
	List() ([]*libovsdb.LoadBalancer, error)
}

type LoadBalancerOpt func(OVNRow) error

func LoadBalancerName(n string) LoadBalancerOpt {
	return func(o OVNRow) error {
		o["name"] = n
		return nil
	}
}

func LoadBalancerUUID(n string) LoadBalancerOpt {
	return func(o OVNRow) error {
		o["_uuid"] = n
		return nil
	}
}

func LoadBalancerVIP(vip string) LoadBalancerOpt {
	return func(o OVNRow) error {
		o["_vip"] = vip
		return nil
	}
}

func LoadBalancerIP(ip []string) LoadBalancerOpt {
	return func(o OVNRow) error {
		o["_ip"] = ip
		return nil
	}
}

func LoadBalancerVIPs(vips map[string]string) LoadBalancerOpt {
	return func(o OVNRow) error {
		if vips == nil || len(vips) == 0 {
			return ErrorOption
		}
		mp := make(map[string]string)
		for k, v := range vips {
			mp[k] = v
		}
		o["vips"] = mp
		return nil
	}
}

func LoadBalancerProtocol(p string) LoadBalancerOpt {
	return func(o OVNRow) error {
		o["protocol"] = p
		return nil
	}
}

type lbImp struct {
	odbi *ovndb
}

func (imp *lbImp) Add(opts ...LoadBalancerOpt) (*OvnCommand, error) {
	optRow := newRow()

	// parse options
	for _, opt := range opts {
		if err := opt(optRow); err != nil {
			return nil, err
		}
	}

	if _, ok := optRow["_uuid"]; ok {
		return nil, ErrorOption
	}

	var operations []libovsdb.Operation
	namedUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}

	row := newRow()
	if vip, ok := optRow["_vip"]; ok {
		vips := make(map[string]string)
		vips[vip.(string)] = strings.Join(optRow["_ip"].([]string), ",")
		oMap, err := libovsdb.NewOvsMap(vips)
		if err != nil {
			return nil, err
		}
		row["vips"] = oMap
	} else if vips, ok := optRow["vips"]; ok {
		oMap, err := libovsdb.NewOvsMap(vips.(map[string]string))
		if err != nil {
			return nil, err
		}
		row["vips"] = oMap
	}

	row["protocol"] = optRow["protocol"]
	row["name"] = optRow["name"]

	insertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableLoadBalancer,
		Row:      row,
		UUIDName: namedUUID,
	}
	operations = append(operations, insertOp)
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lbImp) Del(opts ...LoadBalancerOpt) (*OvnCommand, error) {
	return nil, nil
}

func (imp *lbImp) Get(opts ...LoadBalancerOpt) (*libovsdb.LoadBalancer, error) {
	return nil, nil
}

func (imp *lbImp) List() ([]*libovsdb.LoadBalancer, error) {
	return nil, nil
}
