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
	// LBAdd add LoadBalancer to switch
	LBAdd(LogicalSwitchOpt, LoadBalancerOpt) (*OvnCommand, error)
}

type lsImp struct {
	odbi *ovndb
}

type LogicalSwitchOpt func(OVNRow) error

func LogicalSwitchName(n string) LogicalSwitchOpt {
	return func(o OVNRow) error {
		o["name"] = n
		return nil
	}
}

func LogicalSwitchUUID(n string) LogicalSwitchOpt {
	return func(o OVNRow) error {
		o["_uuid"] = n
		return nil
	}
}

func (imp *lsImp) Add(opts ...LogicalSwitchOpt) (*OvnCommand, error) {
	row := newRow()

	// parse options
	for _, opt := range opts {
		if err := opt(row); err != nil {
			return nil, err
		}
	}

	if _, ok := row["_uuid"]; ok {
		return nil, ErrorOption
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
	row := newRow()

	if len(opts) == 0 {
		return nil, ErrorOption
	}

	// parse options
	for _, opt := range opts {
		if err := opt(row); err != nil {
			return nil, err
		}
	}

	var condition []interface{}
	if uuid, ok := row["_uuid"]; ok {
		condition = libovsdb.NewCondition("_uuid", "==", uuid)
	} else if name, ok := row["name"]; ok {
		condition = libovsdb.NewCondition("name", "==", name)
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
	row := newRow()

	if len(opts) == 0 {
		return nil, ErrorOption
	}

	// parse options
	for _, opt := range opts {
		if err := opt(row); err != nil {
			return nil, err
		}
	}

	tableRow, err := imp.odbi.getRow(tableLogicalSwitch, row)
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

	lsRow := newRow()
	if err := ls(lsRow); err != nil {
		return nil, err
	}

	lbRow := newRow()
	if err := lb(lbRow); err != nil {
		return nil, err
	}

	var lsCondition []interface{}
	if uuid, ok := lsRow["_uuid"]; ok {
		lsCondition = libovsdb.NewCondition("_uuid", "==", uuid)
	} else if name, ok := lsRow["name"]; ok {
		lsCondition = libovsdb.NewCondition("name", "==", name)
	} else {
		return nil, ErrorOption
	}

	var mutateUUID []libovsdb.UUID

	if uuid, ok := lbRow["_uuid"]; ok {
		mutateUUID = []libovsdb.UUID{stringToGoUUID(uuid.(string))}
	} else {
		iface, err := imp.odbi.getRow(tableLoadBalancer, lbRow)
		if err != nil {
			return nil, err
		}
		mutateUUID = []libovsdb.UUID{stringToGoUUID(iface.(*libovsdb.LoadBalancer).UUID)}
	}

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

func (odbi *ovndb) lsExtIdsAddImp(ls string, external_ids map[string]string) (*OvnCommand, error) {
	var operations []libovsdb.Operation
	row := make(OVNRow)
	row["name"] = ls
	lsuuid := odbi.getRowUUID(tableLogicalSwitch, row)
	if len(lsuuid) == 0 {
		return nil, ErrorNotFound
	}
	if len(external_ids) == 0 {
		return nil, fmt.Errorf("external_ids is nil or empty")
	}
	mutateSet, err := libovsdb.NewOvsMap(external_ids)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("external_ids", opInsert, mutateSet)
	condition := libovsdb.NewCondition("name", "==", ls)
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations = append(operations, mutateOp)
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovndb) lsExtIdsDelImp(ls string, external_ids map[string]string) (*OvnCommand, error) {
	var operations []libovsdb.Operation
	row := make(OVNRow)
	row["name"] = ls
	lsuuid := odbi.getRowUUID(tableLogicalSwitch, row)
	if len(lsuuid) == 0 {
		return nil, ErrorNotFound
	}
	if len(external_ids) == 0 {
		return nil, fmt.Errorf("external_ids is nil or empty")
	}
	mutateSet, err := libovsdb.NewOvsMap(external_ids)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("external_ids", opDelete, mutateSet)
	condition := libovsdb.NewCondition("name", "==", ls)
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations = append(operations, mutateOp)
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}
*/
