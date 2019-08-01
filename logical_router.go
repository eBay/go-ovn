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
	"log"

	"github.com/ebay/libovsdb"
)

// LogicalRouterOpt option to logical router commands
type LogicalRouterOpt func(OVNRow) error

// LogicalRouterName pass name for logical router
func LogicalRouterName(n string) LogicalRouterOpt {
	return func(o OVNRow) error {
		o["name"] = n
		return nil
	}
}

// LogicalSwitchUUID pass uuid for logical switch
func LogicalRouterUUID(n string) LogicalRouterOpt {
	return func(o OVNRow) error {
		o["uuid"] = n
		return nil
	}
}

// LogicalRouterMayExist allow logical router exists and not fail creation
func LogicalRouterMayExist(b bool) LogicalRouterOpt {
	return func(o OVNRow) error {
		o["may_exist"] = b
		return nil
	}
}

type lrImp struct {
	odbi *ovndb
}

func (imp *lrImp) Add(opts ...LogicalRouterOpt) (*OvnCommand, error) {
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
	} else {
		return nil, ErrorOption
	}

	var lr *LogicalRouter
	if err := imp.odbi.getRow(tableLogicalRouter, row, &lr); err == nil {
		if v, ok := optRow["may_exist"]; ok && v.(bool) {
			return nil, nil
		}
	}

	namedUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}

	if external_ids, ok := optRow["external_ids"]; ok {
		oMap, err := libovsdb.NewOvsMap(external_ids)
		if err != nil {
			return nil, err
		}
		row["external_ids"] = oMap
	}

	insertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableLogicalRouter,
		Row:      row,
		UUIDName: namedUUID,
	}

	operations := []libovsdb.Operation{insertOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lrImp) Del(opts ...LogicalRouterOpt) (*OvnCommand, error) {
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

	var lr *LogicalRouter
	if err := imp.odbi.getRow(tableLogicalRouter, row, &lr); err != nil {
		return nil, err
	}

	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lr.UUID))
	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableLogicalRouter,
		Where: []interface{}{condition},
	}
	operations := []libovsdb.Operation{deleteOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lrImp) Get(opts ...LogicalRouterOpt) (*LogicalRouter, error) {
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

	var lr *LogicalRouter
	if err := imp.odbi.getRow(tableLogicalRouter, row, &lr); err != nil {
		return nil, err
	}

	return lr, nil
}

// Get all logical switches
func (imp *lrImp) List() ([]*LogicalRouter, error) {
	var lrList []*LogicalRouter

	if err := imp.odbi.getRows(tableLogicalRouter, nil, &lrList); err != nil {
		return nil, err
	}

	return lrList, nil
}

func (imp *lrImp) LBAdd(lr LogicalRouterOpt, lb LoadBalancerOpt) (*OvnCommand, error) {
	optLRRow := newRow()
	if err := lr(optLRRow); err != nil {
		return nil, err
	}

	optLBRow := newRow()
	if err := lb(optLBRow); err != nil {
		return nil, err
	}

	lrRow := newRow()
	if uuid, ok := optLRRow["uuid"]; ok {
		lrRow["uuid"] = uuid
	} else if name, ok := optLRRow["name"]; ok {
		lrRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var lrItem *LogicalRouter
	if err := imp.odbi.getRow(tableLogicalRouter, lrRow, &lrItem); err != nil {
		return nil, err
	}

	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lrItem.UUID))

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
		log.Printf("ZZZ\n")
		return nil, err
	}

	mutateUUID := []libovsdb.UUID{stringToGoUUID(lbItem.UUID)}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	mutation := libovsdb.NewMutation("load_balancer", opInsert, mutateSet)
	if err != nil {
		return nil, err
	}

	condition = libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lrItem.UUID))
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalRouter,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{mutateOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lrImp) LBDel(lr LogicalRouterOpt, lb LoadBalancerOpt) (*OvnCommand, error) {
	optLRRow := newRow()
	if err := lr(optLRRow); err != nil {
		return nil, err
	}

	optLBRow := newRow()
	if lb != nil {
		if err := lb(optLBRow); err != nil {
			return nil, err
		}
	}

	lrRow := newRow()
	if uuid, ok := optLRRow["uuid"]; ok {
		lrRow["uuid"] = uuid
	} else if name, ok := optLRRow["name"]; ok {
		lrRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var lrItem *LogicalRouter
	if err := imp.odbi.getRow(tableLogicalRouter, lrRow, &lrItem); err != nil {
		return nil, err
	}

	var lbUUIDs []string
	if lb == nil {
		lbUUIDs = append(lbUUIDs, lrItem.LoadBalancer...)
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

	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lrItem.UUID))

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
		Table:     tableLogicalRouter,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{mutateOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (imp *lrImp) LBList(opt LogicalRouterOpt) ([]*LoadBalancer, error) {
	var err error

	optLRRow := newRow()
	if err = opt(optLRRow); err != nil {
		return nil, err
	}

	lrRow := newRow()
	if uuid, ok := optLRRow["uuid"]; ok {
		lrRow["uuid"] = uuid
	} else if name, ok := optLRRow["name"]; ok {
		lrRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var lr *LogicalRouter
	if err = imp.odbi.getRow(tableLogicalSwitch, lrRow, &lr); err != nil {
		return nil, err
	}

	lbList := make([]*LoadBalancer, len(lr.LoadBalancer))

	for i := 0; i < len(lr.LoadBalancer); i++ {
		if err = imp.odbi.getRowByUUID(tableLoadBalancer, lr.LoadBalancer[i], &lbList[i]); err != nil {
			return nil, err
		}
	}

	if len(lbList) == 0 {
		return nil, ErrorNotFound
	}

	return lbList, nil
}
