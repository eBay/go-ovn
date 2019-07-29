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
	// LBAdd add LoadBalancer to switch
	LBAdd(LogicalSwitchOpt, LoadBalancerOpt) (*OvnCommand, error)
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

	var condition []interface{}
	if uuid, ok := optRow["uuid"]; ok {
		condition = libovsdb.NewCondition("_uuid", "==", stringToGoUUID(uuid.(string)))
	} else if name, ok := optRow["name"]; ok {
		condition = libovsdb.NewCondition("name", "==", name.(string))
	} else {
		return nil, ErrorOption
	}

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

	tableRow, err := imp.odbi.getRow(tableLogicalSwitch, optRow)
	if err != nil {
		return nil, err
	}

	return tableRow.(*libovsdb.LogicalSwitch), nil
}

func (imp *lsImp) List() ([]*libovsdb.LogicalSwitch, error) {
	rows, err := imp.odbi.getRows(tableLogicalSwitch, nil)
	if err != nil {
		return nil, err
	}
	lsList := make([]*libovsdb.LogicalSwitch, len(rows))

	for i, row := range rows {
		lsList[i] = row.(*libovsdb.LogicalSwitch)
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

	var lsCondition []interface{}
	var lsUUID string
	if uuid, ok := optLSRow["uuid"]; ok {
		lsCondition = libovsdb.NewCondition("_uuid", "==", stringToGoUUID(uuid.(string)))
	} else {
		iface, err := imp.odbi.getRowByName(tableLogicalSwitch, optLSRow)
		if err != nil {
			return nil, err
		}
		lsUUID = iface.(*libovsdb.LogicalSwitch).UUID
	}

	lsCondition = libovsdb.NewCondition("_uuid", "==", stringToGoUUID(lsUUID))

	var mutateUUID []libovsdb.UUID
	var lbUUID string

	if uuid, ok := optLBRow["uuid"]; ok {
		lbUUID = uuid.(string)
	} else {
		iface, err := imp.odbi.getRowByName(tableLoadBalancer, optLBRow)
		if err != nil {
			return nil, err
		}
		lbUUID = iface.(*libovsdb.LoadBalancer).UUID
	}

	mutateUUID = []libovsdb.UUID{stringToGoUUID(lbUUID)}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	mutation := libovsdb.NewMutation("load_balancer", opInsert, mutateSet)
	if err != nil {
		return nil, err
	}

	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{lsCondition},
	}
	operations = append(operations, mutateOp)
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

/*
func (odbi *ovndb) lslbDelImp(lswitch string, lb string) (*OvnCommand, error) {
	var operations []libovsdb.Operation
	row := make(OVNRow)
	row["name"] = lb
	lbuuid := odbi.getRowUUID(tableLoadBalancer, row)
	if len(lbuuid) == 0 {
		return nil, ErrorNotFound
	}
	row = make(OVNRow)
	row["name"] = lswitch
	lsuuid := odbi.getRowUUID(tableLogicalSwitch, row)
	if len(lsuuid) == 0 {
		return nil, ErrorNotFound
	}
	mutateUUID := []libovsdb.UUID{stringToGoUUID(lbuuid)}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("load_balancer", opDelete, mutateSet)
	// mutate  lswitch for the corresponding load_balancer
	mucondition := libovsdb.NewCondition("name", "==", lswitch)
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{mucondition},
	}
	operations = append(operations, mutateOp)
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovndb) lslbListImp(lswitch string) ([]*LoadBalancer, error) {
	var listLB []*LoadBalancer
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheLogicalSwitch, ok := odbi.cache[tableLogicalSwitch]
	if !ok {
		return nil, ErrorSchema
	}

	for _, drows := range cacheLogicalSwitch {
		if rlsw, ok := drows.Fields["name"].(string); ok && rlsw == lswitch {
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
		iface, err := imp.odbi.getRowByName(tableLogicalSwitch, optRow)
		if err != nil {
			return nil, err
		}
		lsUUID = iface.(*libovsdb.LogicalSwitch).UUID
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
		iface, err := imp.odbi.getRow(tableLogicalSwitch, optRow)
		if err != nil {
			return nil, err
		}
		lsUUID = iface.(*libovsdb.LogicalSwitch).UUID
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
