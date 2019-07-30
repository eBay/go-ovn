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

	"github.com/ebay/libovsdb"
)

// LogicalSwitch control ovn logical switches
type LogicalSwitch interface {
	// Add logical switch with optional LogicalSwitchName
	Add(...LogicalSwitchOpt) (*OvnCommand, error)
	// Del logical switch with LogicalSwitchName or LogicalSwitchUUID
	Del(...LogicalSwitchOpt) (*OvnCommand, error)
	// Get logical switch with LogicalSwitchName or LogicalSwitchUUID
	Get(...LogicalSwitchOpt) (*libovsdb.LogicalSwitch, error)
	// List logical switches
	List() ([]*libovsdb.LogicalSwitch, error)
	// SetExternalIDs logical switch with LogicalSwitchName or LogicalSwitchUUID
	SetExternalIDs(...LogicalSwitchOpt) (*OvnCommand, error)
	// DelExternalIDs logical switch with LogicalSwitchName or LogicalSwitchUUID
	DelExternalIDs(...LogicalSwitchOpt) (*OvnCommand, error)
	// LBAdd add load balancer to switch
	LBAdd(LogicalSwitchOpt, LoadBalancerOpt) (*OvnCommand, error)
	// LBDel delete load balancer or all LoadBalancers from logical switch
	LBDel(LogicalSwitchOpt, LoadBalancerOpt) (*OvnCommand, error)
	// LBList list load balancers from logical switch
	LBList(LogicalSwitchOpt) ([]*libovsdb.LoadBalancer, error)
}

type lsImp struct {
	odbi *ovndb
}

// LogicalSwitchOpt option to logical switch commands
type LogicalSwitchOpt func(OVNRow) error

// LogicalSwitchName pass name for logical switch
func LogicalSwitchName(n string) LogicalSwitchOpt {
	return func(o OVNRow) error {
		o["name"] = n
		return nil
	}
}

// LogicalSwitchUUID pass uuid for logical switch
func LogicalSwitchUUID(n string) LogicalSwitchOpt {
	return func(o OVNRow) error {
		o["uuid"] = n
		return nil
	}
}

// LogicalSwitchExternalIDs pass external_ids for logical switch
func LogicalSwitchExternalIDs(m map[string]string) LogicalSwitchOpt {
	return func(o OVNRow) error {
		if m == nil || len(m) == 0 {
			return ErrorOption
		}

		mp := make(map[string]string)
		for k, v := range m {
			mp[k] = v
		}

		o["external_ids"] = mp
		return nil
	}
}

func (imp *lsImp) Add(opts ...LogicalSwitchOpt) (*OvnCommand, error) {
	optRow := newRow()

	// parse options
	for _, opt := range opts {
		if err := opt(optRow); err != nil {
			return nil, err
		}
	}

	row := newRow()
	if name, ok := optRow["name"]; ok {
		row["name"] = name
	}

	namedUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}

	insertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableLogicalSwitch,
		Row:      row,
		UUIDName: namedUUID,
	}

	operations := []libovsdb.Operation{insertOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lsImp) Del(opts ...LogicalSwitchOpt) (*OvnCommand, error) {
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

	var lsUUID string
	if uuid, ok := optRow["uuid"]; ok {
		lsUUID = uuid.(string)
	} else {
		var ls *libovsdb.LogicalSwitch
		if err := imp.odbi.getRowByName(tableLogicalSwitch, optRow["name"].(string), &ls); err != nil {
			return nil, err
		}
		lsUUID = ls.UUID
	}

	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lsUUID))
	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableLogicalSwitch,
		Where: []interface{}{condition},
	}
	operations := []libovsdb.Operation{deleteOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lsImp) Get(opts ...LogicalSwitchOpt) (*libovsdb.LogicalSwitch, error) {
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

	var ls *libovsdb.LogicalSwitch
	if err := imp.odbi.getRow(tableLogicalSwitch, optRow, &ls); err != nil {
		return nil, err
	}

	return ls, nil
}

func (imp *lsImp) List() ([]*libovsdb.LogicalSwitch, error) {
	var lsList []*libovsdb.LogicalSwitch

	if err := imp.odbi.getRows(tableLogicalSwitch, nil, &lsList); err != nil {
		return nil, err
	}

	return lsList, nil
}

func (imp *lsImp) LBAdd(ls LogicalSwitchOpt, lb LoadBalancerOpt) (*OvnCommand, error) {
	var operations []libovsdb.Operation

	optLSRow := newRow()
	if err := ls(optLSRow); err != nil {
		return nil, err
	}

	optLBRow := newRow()
	if err := lb(optLBRow); err != nil {
		return nil, err
	}

	var lsUUID string
	if uuid, ok := optLSRow["uuid"]; ok {
		lsUUID = uuid.(string)
	} else {
		var ls *libovsdb.LogicalSwitch
		if err := imp.odbi.getRowByName(tableLogicalSwitch, optLSRow["name"].(string), &ls); err != nil {
			return nil, err
		}
		lsUUID = ls.UUID
	}

	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lsUUID))

	var mutateUUID []libovsdb.UUID
	var lbUUID string

	if uuid, ok := optLBRow["uuid"]; ok {
		lbUUID = uuid.(string)
	} else {
		var lb *libovsdb.LoadBalancer
		if err := imp.odbi.getRowByName(tableLoadBalancer, optLBRow["name"].(string), &lb); err != nil {
			return nil, err
		}
		lbUUID = lb.UUID
	}

	mutateUUID = []libovsdb.UUID{stringToGoUUID(lbUUID)}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("load_balancer", opInsert, mutateSet)

	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations = append(operations, mutateOp)
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lsImp) LBDel(ls LogicalSwitchOpt, lb LoadBalancerOpt) (*OvnCommand, error) {
	var operations []libovsdb.Operation

	optLSRow := newRow()
	if err := ls(optLSRow); err != nil {
		return nil, err
	}

	optLBRow := newRow()
	if lb != nil {
		if err := lb(optLBRow); err != nil {
			return nil, err
		}
	}

	var lbUUIDs []string
	if lb != nil {
		if uuid, ok := optLBRow["uuid"]; ok {
			lbUUIDs = append(lbUUIDs, uuid.(string))
		} else {
			var lb *libovsdb.LoadBalancer
			if err := imp.odbi.getRowByName(tableLoadBalancer, optLBRow["name"].(string), &lb); err != nil {
				return nil, err
			}
			lbUUIDs = append(lbUUIDs, lb.UUID)
		}
	}

	var lsUUID string
	if uuid, ok := optLSRow["uuid"]; ok {
		lsUUID = uuid.(string)
	} else {
		var ls *libovsdb.LogicalSwitch
		if err := imp.odbi.getRowByName(tableLogicalSwitch, optLSRow["name"].(string), &ls); err != nil {
			return nil, err
		}
		lsUUID = ls.UUID
		if len(lbUUIDs) == 0 {
			lbUUIDs = append(lbUUIDs, ls.LoadBalancer...)
		}
	}
	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lsUUID))

	mutateUUID := make([]libovsdb.UUID, len(lbUUIDs))
	for i, lbUUID := range lbUUIDs {
		mutateUUID[i] = stringToGoUUID(lbUUID)
	}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("load_balancer", opDelete, mutateSet)
	// mutate  lswitch for the corresponding load_balancer
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations = append(operations, mutateOp)
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lsImp) LBList(opt LogicalSwitchOpt) ([]*libovsdb.LoadBalancer, error) {
	var err error

	optLSRow := newRow()
	if err = opt(optLSRow); err != nil {
		return nil, err
	}

	var ls *libovsdb.LogicalSwitch
	if err = imp.odbi.getRow(tableLogicalSwitch, optLSRow, &ls); err != nil {
		return nil, err
	}

	lbList := make([]*libovsdb.LoadBalancer, len(ls.LoadBalancer))

	for i := 0; i < len(ls.LoadBalancer); i++ {
		if err = imp.odbi.getRowByUUID(tableLoadBalancer, ls.LoadBalancer[i], &lbList[i]); err != nil {
			return nil, err
		}
	}

	if len(lbList) == 0 {
		return nil, ErrorNotFound
	}

	return lbList, nil
}

func (imp *lsImp) SetExternalIDs(opts ...LogicalSwitchOpt) (*OvnCommand, error) {
	var operations []libovsdb.Operation

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

	// set only external_ids now
	external_ids, ok := optRow["external_ids"]
	if !ok || len(external_ids.(map[string]string)) == 0 {
		return nil, fmt.Errorf("external_ids is nil or empty")
	}

	var lsUUID string
	if uuid, ok := optRow["uuid"]; ok {
		lsUUID = uuid.(string)
	} else {
		var ls *libovsdb.LogicalSwitch
		if err := imp.odbi.getRowByName(tableLogicalSwitch, optRow["name"].(string), &ls); err != nil {
			return nil, err
		}
		lsUUID = ls.UUID
	}

	mutateSet, err := libovsdb.NewOvsMap(external_ids)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("external_ids", opInsert, mutateSet)
	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lsUUID))
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations = append(operations, mutateOp)
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lsImp) DelExternalIDs(opts ...LogicalSwitchOpt) (*OvnCommand, error) {
	var operations []libovsdb.Operation

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

	// set only external_ids now
	external_ids, ok := optRow["external_ids"]
	if !ok || len(external_ids.(map[string]string)) == 0 {
		return nil, fmt.Errorf("external_ids is nil or empty")
	}

	var lsUUID string
	if uuid, ok := optRow["uuid"]; ok {
		lsUUID = uuid.(string)
	} else {
		var ls *libovsdb.LogicalSwitch
		if err := imp.odbi.getRow(tableLogicalSwitch, optRow, &ls); err != nil {
			return nil, err
		}
		lsUUID = ls.UUID
	}

	mutateSet, err := libovsdb.NewOvsMap(external_ids)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("external_ids", opDelete, mutateSet)
	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lsUUID))
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations = append(operations, mutateOp)
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}
