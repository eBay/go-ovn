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

// LogicalSwitchMayExist allow logical switch exists and not fail creation
func LogicalSwitchMayExist(b bool) LogicalSwitchOpt {
	return func(o OVNRow) error {
		o["may_exist"] = b
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

	var ls *LogicalSwitch
	if err := imp.odbi.getRow(tableLogicalSwitch, row, &ls); err == nil {
		if v, ok := optRow["may_exist"]; ok && v.(bool) {
			return nil, nil
		}
	}

	if external_ids, ok := optRow["external_ids"]; ok {
		oMap, err := libovsdb.NewOvsMap(external_ids)
		if err != nil {
			return nil, err
		}
		row["external_ids"] = oMap
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
	optLSRow := newRow()

	if len(opts) == 0 {
		return nil, ErrorOption
	}

	// parse options
	for _, opt := range opts {
		if err := opt(optLSRow); err != nil {
			return nil, err
		}
	}

	lsRow := newRow()
	if uuid, ok := optLSRow["uuid"]; ok {
		lsRow["uuid"] = uuid
	} else if name, ok := optLSRow["name"]; ok {
		lsRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var ls *LogicalSwitch
	if err := imp.odbi.getRow(tableLogicalSwitch, lsRow, &ls); err != nil {
		return nil, err
	}

	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(ls.UUID))
	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableLogicalSwitch,
		Where: []interface{}{condition},
	}
	operations := []libovsdb.Operation{deleteOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lsImp) Get(opts ...LogicalSwitchOpt) (*LogicalSwitch, error) {
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
	if uuid, ok := optRow["uuid"]; ok {
		row["uuid"] = uuid
	} else if name, ok := optRow["name"]; ok {
		row["name"] = name
	} else {
		return nil, ErrorOption
	}

	var ls *LogicalSwitch
	if err := imp.odbi.getRow(tableLogicalSwitch, row, &ls); err != nil {
		return nil, err
	}

	return ls, nil
}

func (imp *lsImp) List() ([]*LogicalSwitch, error) {
	var lsList []*LogicalSwitch

	if err := imp.odbi.getRows(tableLogicalSwitch, nil, &lsList); err != nil {
		return nil, err
	}

	return lsList, nil
}

func (imp *lsImp) LBAdd(ls LogicalSwitchOpt, lb LoadBalancerOpt) (*OvnCommand, error) {
	optLSRow := newRow()
	if err := ls(optLSRow); err != nil {
		return nil, err
	}

	optLBRow := newRow()
	if err := lb(optLBRow); err != nil {
		return nil, err
	}

	lsRow := newRow()
	if uuid, ok := optLSRow["uuid"]; ok {
		lsRow["uuid"] = uuid
	} else if name, ok := optLSRow["name"]; ok {
		lsRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var lsItem *LogicalSwitch
	if err := imp.odbi.getRow(tableLogicalSwitch, lsRow, &lsItem); err != nil {
		return nil, err
	}

	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lsItem.UUID))

	lbRow := newRow()
	if uuid, ok := optLBRow["uuid"]; ok {
		lbRow["uuid"] = uuid
	} else if name, ok := optLBRow["name"]; ok {
		lbRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var lbItem *LoadBalancer
	if err := imp.odbi.getRow(tableLoadBalancer, lbRow, &lbItem); err != nil {
		return nil, err
	}

	mutateUUID := []libovsdb.UUID{stringToGoUUID(lbItem.UUID)}
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
	operations := []libovsdb.Operation{mutateOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lsImp) LBDel(ls LogicalSwitchOpt, lb LoadBalancerOpt) (*OvnCommand, error) {
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

	lsRow := newRow()
	if uuid, ok := optLSRow["uuid"]; ok {
		lsRow["uuid"] = uuid
	} else if name, ok := optLSRow["name"]; ok {
		lsRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var lsItem *LogicalSwitch
	if err := imp.odbi.getRow(tableLogicalSwitch, lsRow, &lsItem); err != nil {
		return nil, err
	}

	var lbUUIDs []string
	if lb == nil {
		lbUUIDs = append(lbUUIDs, lsItem.LoadBalancer...)
	} else {
		lbRow := newRow()
		if uuid, ok := optLBRow["uuid"]; ok {
			lbRow["uuid"] = uuid
		} else if name, ok := optLBRow["name"]; ok {
			lbRow["name"] = name
		} else {
			return nil, ErrorOption
		}
		var lbItem *LoadBalancer
		if err := imp.odbi.getRow(tableLoadBalancer, lbRow, &lbItem); err != nil {
			return nil, err
		}
		lbUUIDs = append(lbUUIDs, lbItem.UUID)
	}

	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lsItem.UUID))

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
	operations := []libovsdb.Operation{mutateOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lsImp) LBList(opt LogicalSwitchOpt) ([]*LoadBalancer, error) {
	var err error

	optLSRow := newRow()
	if err = opt(optLSRow); err != nil {
		return nil, err
	}

	lsRow := newRow()
	if uuid, ok := optLSRow["uuid"]; ok {
		lsRow["uuid"] = uuid
	} else if name, ok := optLSRow["name"]; ok {
		lsRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var ls *LogicalSwitch
	if err = imp.odbi.getRow(tableLogicalSwitch, lsRow, &ls); err != nil {
		return nil, err
	}

	lbList := make([]*LoadBalancer, len(ls.LoadBalancer))

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
	optLSRow := newRow()

	if len(opts) == 0 {
		return nil, ErrorOption
	}

	// parse options
	for _, opt := range opts {
		if err := opt(optLSRow); err != nil {
			return nil, err
		}
	}

	// set only external_ids now
	external_ids, ok := optLSRow["external_ids"]
	if !ok || len(external_ids.(map[string]string)) == 0 {
		return nil, fmt.Errorf("external_ids is nil or empty")
	}

	lsRow := newRow()
	if uuid, ok := optLSRow["uuid"]; ok {
		lsRow["uuid"] = uuid
	} else if name, ok := optLSRow["name"]; ok {
		lsRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var ls *LogicalSwitch
	if err := imp.odbi.getRow(tableLogicalSwitch, lsRow, &ls); err != nil {
		return nil, err
	}

	mutateSet, err := libovsdb.NewOvsMap(external_ids)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("external_ids", opInsert, mutateSet)
	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(ls.UUID))
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{mutateOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lsImp) DelExternalIDs(opts ...LogicalSwitchOpt) (*OvnCommand, error) {
	optLSRow := newRow()

	if len(opts) == 0 {
		return nil, ErrorOption
	}

	// parse options
	for _, opt := range opts {
		if err := opt(optLSRow); err != nil {
			return nil, err
		}
	}

	// set only external_ids now
	external_ids, ok := optLSRow["external_ids"]
	if !ok || len(external_ids.(map[string]string)) == 0 {
		return nil, fmt.Errorf("external_ids is nil or empty")
	}

	lsRow := newRow()
	if uuid, ok := optLSRow["uuid"]; ok {
		lsRow["uuid"] = uuid
	} else if name, ok := optLSRow["name"]; ok {
		lsRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var ls *LogicalSwitch
	if err := imp.odbi.getRow(tableLogicalSwitch, lsRow, &ls); err != nil {
		return nil, err
	}

	mutateSet, err := libovsdb.NewOvsMap(external_ids)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("external_ids", opDelete, mutateSet)
	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(ls.UUID))
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{mutateOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}
