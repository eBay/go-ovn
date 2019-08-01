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

/*
func (odbi *ovndb) lrlbAddImp(lr string, lb string) (*OvnCommand, error) {
	var operations []libovsdb.Operation
	row := make(OVNRow)
	row["name"] = lb
	lbuuid := odbi.getRowUUID(tableLoadBalancer, row)
	if len(lbuuid) == 0 {
		return nil, ErrorNotFound
	}
	mutateUUID := []libovsdb.UUID{stringToGoUUID(lbuuid)}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	mutation := libovsdb.NewMutation("load_balancer", opInsert, mutateSet)
	if err != nil {
		return nil, err
	}
	row = make(OVNRow)
	row["name"] = lr
	lruuid := odbi.getRowUUID(tableLogicalRouter, row)
	if len(lruuid) == 0 {
		return nil, ErrorNotFound
	}
	condition := libovsdb.NewCondition("name", "==", lr)
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalRouter,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations = append(operations, mutateOp)
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovndb) lrlbDelImp(lr string, lb string) (*OvnCommand, error) {
	var operations []libovsdb.Operation
	row := make(OVNRow)
	row["name"] = lb
	lbuuid := odbi.getRowUUID(tableLoadBalancer, row)
	if len(lbuuid) == 0 {
		return nil, ErrorNotFound
	}
	mutateUUID := []libovsdb.UUID{stringToGoUUID(lbuuid)}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}
	row = make(OVNRow)
	row["name"] = lr
	lruuid := odbi.getRowUUID(tableLogicalRouter, row)
	if len(lruuid) == 0 {
		return nil, ErrorNotFound
	}
	mutation := libovsdb.NewMutation("load_balancer", opDelete, mutateSet)
	// mutate  lswitch for the corresponding load_balancer
	mucondition := libovsdb.NewCondition("name", "==", lr)
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalRouter,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{mucondition},
	}
	operations = append(operations, mutateOp)
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovndb) lrlbListImp(lr string) ([]*LoadBalancer, error) {
	var listLB []*LoadBalancer
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheLogicalRouter, ok := odbi.cache[tableLogicalRouter]
	if !ok {
		return nil, ErrorSchema
	}

	for _, drows := range cacheLogicalRouter {
		if router, ok := drows.Fields["name"].(string); ok && router == lr {
			lbs := drows.Fields["load_balancer"]
			if lbs != nil {
				switch lbs.(type) {
				case libovsdb.OvsSet:
					if lb, ok := lbs.(libovsdb.OvsSet); ok {
						for _, l := range lb.GoSet {
							if lb, ok := l.(libovsdb.UUID); ok {
								lb, err := odbi.rowToLB(lb.GoUUID)
								if err != nil {
									return nil, err
								}
								listLB = append(listLB, lb)
							}
						}
					} else {
						return nil, fmt.Errorf("type libovsdb.OvsSet casting failed")
					}
				case libovsdb.UUID:
					if lb, ok := lbs.(libovsdb.UUID); ok {
						lb, err := odbi.rowToLB(lb.GoUUID)
						if err != nil {
							return nil, err
						}
						listLB = append(listLB, lb)
					} else {
						return nil, fmt.Errorf("type libovsdb.UUID casting failed")
					}
				default:
					return nil, fmt.Errorf("Unsupport type found in ovsdb rows")
				}
			}
			break
		}
	}
	return listLB, nil
}
*/
