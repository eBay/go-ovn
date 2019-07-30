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
	// Get load balancer with LoadBalancerName or LoadBalancerUUID
	Get(...LoadBalancerOpt) (*libovsdb.LoadBalancer, error)
	// Set load balancer with LoadBalancerName or LoadBalancerUUID and other opts
	Set(...LoadBalancerOpt) (*OvnCommand, error)
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
		o["uuid"] = n
		return nil
	}
}

func LoadBalancerVIP(vip string) LoadBalancerOpt {
	return func(o OVNRow) error {
		o["vip"] = vip
		return nil
	}
}

func LoadBalancerIP(ip []string) LoadBalancerOpt {
	return func(o OVNRow) error {
		o["ip"] = ip
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
	if opts == nil || len(opts) < 3 {
		return nil, ErrorOption
	}

	optRow := newRow()

	// parse options
	for _, opt := range opts {
		if err := opt(optRow); err != nil {
			return nil, err
		}
	}

	if _, ok := optRow["uuid"]; ok {
		return nil, ErrorOption
	}

	var operations []libovsdb.Operation
	namedUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}

	row := newRow()
	if vip, ok := optRow["vip"]; ok {
		vips := make(map[string]string)
		vips[vip.(string)] = strings.Join(optRow["ip"].([]string), ",")
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
	var operations []libovsdb.Operation

	if opts == nil || len(opts) == 0 {
		return nil, ErrorOption
	}

	optRow := newRow()

	// parse options
	for _, opt := range opts {
		if err := opt(optRow); err != nil {
			return nil, err
		}
	}

	row := newRow()
	if vip, ok := optRow["vip"]; ok {
		vips := make(map[string]string)
		vips[vip.(string)] = strings.Join(optRow["ip"].([]string), ",")
		row["vips"] = vips
	} else if vips, ok := optRow["vips"]; ok {
		row["vips"] = vips
	}

	if protocol, ok := optRow["protocol"]; ok {
		row["protocol"] = protocol
	}
	if uuid, ok := optRow["uuid"]; ok {
		row["uuid"] = uuid
	}
	if name, ok := optRow["name"]; ok {
		row["name"] = name
	}

	var lb *libovsdb.LoadBalancer
	if err := imp.odbi.getRow(tableLoadBalancer, row, &lb); err != nil {
		return nil, err
	}

	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lb.UUID))

	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableLoadBalancer,
		Where: []interface{}{condition},
	}

	mutateUUID := []libovsdb.UUID{stringToGoUUID(lb.UUID)}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("load_balancer", opDelete, mutateSet)

	var lsList []*libovsdb.LogicalSwitch
	err = imp.odbi.getRows(tableLogicalSwitch, map[string]interface{}{"load_balancer": []string{lb.UUID}}, &lsList)
	if err != nil && err != ErrorNotFound {
		return nil, err
	} else if err == nil {
		// mutate all matching logical switches for the corresponding load_balancer
		for _, ls := range lsList {
			mucondition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(ls.UUID))
			mutateOp := libovsdb.Operation{
				Op:        opMutate,
				Table:     tableLogicalSwitch,
				Mutations: []interface{}{mutation},
				Where:     []interface{}{mucondition},
			}
			operations = append(operations, mutateOp)
		}
	}
	operations = append(operations, deleteOp)
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lbImp) Get(opts ...LoadBalancerOpt) (*libovsdb.LoadBalancer, error) {
	optRow := newRow()

	if len(opts) == 0 {
		return nil, ErrorOption
	}

	// parse options
	for _, opt := range opts {
		if err := opt(optRow); err != nil {
			return nil, err
		}
	}

	row := newRow()
	if vip, ok := optRow["vip"]; ok {
		vips := make(map[string]string)
		vips[vip.(string)] = strings.Join(optRow["ip"].([]string), ",")
		row["vips"] = vips
	} else if vips, ok := optRow["vips"]; ok {
		row["vips"] = vips
	}

	if protocol, ok := optRow["protocol"]; ok {
		row["protocol"] = protocol
	}
	if uuid, ok := optRow["uuid"]; ok {
		row["uuid"] = uuid
	}
	if name, ok := optRow["name"]; ok {
		row["name"] = name
	}

	var lb *libovsdb.LoadBalancer
	if err := imp.odbi.getRow(tableLoadBalancer, row, &lb); err != nil {
		return nil, err
	}

	return lb, nil
}

func (imp *lbImp) List() ([]*libovsdb.LoadBalancer, error) {
	return nil, nil
}

func (imp *lbImp) Set(opts ...LoadBalancerOpt) (*OvnCommand, error) {
	optRow := newRow()

	if len(opts) == 0 {
		return nil, ErrorOption
	}

	// parse options
	for _, opt := range opts {
		if err := opt(optRow); err != nil {
			return nil, err
		}
	}

	row := newRow()
	if vip, ok := optRow["vip"]; ok {
		vips := make(map[string]string)
		vips[vip.(string)] = strings.Join(optRow["ip"].([]string), ",")
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

	if protocol, ok := optRow["protocol"]; ok {
		row["protocol"] = protocol
	}

	var lbUUID string
	if uuid, ok := optRow["uuid"]; ok {
		lbUUID = uuid.(string)
	} else {
		var lb *libovsdb.LoadBalancer
		if err := imp.odbi.getRowByName(tableLoadBalancer, optRow["name"].(string), &lb); err != nil {
			return nil, err
		}
		lbUUID = lb.UUID
	}

	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lbUUID))

	insertOp := libovsdb.Operation{
		Op:    opUpdate,
		Table: tableLoadBalancer,
		Row:   row,
		Where: []interface{}{condition},
	}
	operations := []libovsdb.Operation{insertOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}
