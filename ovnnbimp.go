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
	"errors"
	"fmt"

	"github.com/ebay/libovsdb"
)

var (
	// ErrorOption used when invalid args specified
	ErrorOption = errors.New("invalid option specified")
	// ErrorSchema used when something wrong in ovnnb
	ErrorSchema = errors.New("table schema error")
	// ErrorNotFound used when object not found in ovnnb
	ErrorNotFound = errors.New("object not found")
	// ErrorExist used when object already exists in ovnnb
	ErrorExist = errors.New("object exist")
)

// OVNRow ovnnb row
type OVNRow map[string]interface{}

/*
func (odbi *ovndb) getRowUUIDs(table string, row OVNRow) []string {
	var uuids []string
	var wildcard bool

	if reflect.DeepEqual(row, make(OVNRow)) {
		wildcard = true
	}

	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheTable, ok := odbi.cache[table]
	if !ok {
		return nil
	}

	for uuid, drows := range cacheTable {
		if wildcard {
			uuids = append(uuids, uuid)
			continue
		}

		found := false
		for field, value := range row {
			if v, ok := drows.Fields[field]; ok {
				if v == value {
					found = true
				} else {
					found = false
					break
				}
			}
		}
		if found {
			uuids = append(uuids, uuid)
		}
	}

	return uuids
}
*/
/*
func (odbi *ovndb) getRowUUID(table string, row OVNRow) string {
	uuids := odbi.getRowUUIDs(table, row)
	if len(uuids) > 0 {
		return uuids[0]
	}
	return ""
}
*/
//test if map s contains t
//This function is not both s and t are nil at same time
func (odbi *ovndb) oMapContians(s, t map[interface{}]interface{}) bool {
	if s == nil || t == nil {
		return false
	}

	for tk, tv := range t {
		if sv, ok := s[tk]; !ok {
			return false
		} else if tv != sv {
			return false
		}
	}
	return true
}

/*
func (odbi *ovndb) getRowUUIDContainsUUID(table, field, uuid string) (string, error) {
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheTable, ok := odbi.cache[table]
	if !ok {
		return "", ErrorSchema
	}

	for id, drows := range cacheTable {
		v := fmt.Sprintf("%s", drows.Fields[field])
		if strings.Contains(v, uuid) {
			return id, nil
		}
	}
	return "", ErrorNotFound
}

func (odbi *ovndb) getRowsMatchingUUID(table, field, uuid string) ([]string, error) {
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	var uuids []string
	for id, drows := range odbi.cache[table] {
		v := fmt.Sprintf("%s", drows.Fields[field])
		if strings.Contains(v, uuid) {
			uuids = append(uuids, id)
		}
	}
	if len(uuids) == 0 {
		return uuids, ErrorNotFound
	}
	return uuids, nil
}
*/

func (odbi *ovndb) transact(ops ...libovsdb.Operation) ([]libovsdb.OperationResult, error) {
	// Only support one trans at same time now.
	odbi.tranmutex.Lock()
	defer odbi.tranmutex.Unlock()
	reply, err := odbi.client.Transact(dbNB, ops...)

	if err != nil {
		return reply, err
	}

	if len(reply) < len(ops) {
		for i, o := range reply {
			if o.Error != "" && i < len(ops) {
				return nil, fmt.Errorf("Transaction Failed due to an error : %v details: %v in %v", o.Error, o.Details, ops[i])
			}
		}
		return reply, fmt.Errorf("Number of Replies should be atleast equal to number of operations")
	}
	return reply, nil
}

func (odbi *ovndb) execute(cmds ...*OvnCommand) error {
	if cmds == nil {
		return nil
	}
	var ops []libovsdb.Operation
	for _, cmd := range cmds {
		if cmd != nil {
			ops = append(ops, cmd.Operations...)
		}
	}
	_, err := odbi.transact(ops...)
	if err != nil {
		return err
	}
	return nil
}

func (odbi *ovndb) float64_to_int(row libovsdb.Row) {
	for field, value := range row.Fields {
		if v, ok := value.(float64); ok {
			n := int(v)
			if float64(n) == v {
				row.Fields[field] = n
			}
		}
	}
}

func (odbi *ovndb) populateCache(updates libovsdb.TableUpdates) {
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()

	for _, table := range tablesOrder {
		tableUpdate, ok := updates.Updates[table]
		if !ok {
			continue
		}

		if _, ok := odbi.cache[table]; !ok {
			odbi.cache[table] = make(map[string]interface{})
		}

		for uuid, row := range tableUpdate.Rows {
			if row.New != nil {
				odbi.cache[table][uuid] = row.New
				if odbi.signalCB != nil {
					switch table {
					case tableLogicalRouter:
						odbi.signalCB.OnLogicalRouterCreate(row.New.(*libovsdb.LogicalRouter))
					case tableLogicalRouterPort:
						odbi.signalCB.OnLogicalRouterPortCreate(row.New.(*libovsdb.LogicalRouterPort))
					case tableLogicalRouterStaticRoute:
						odbi.signalCB.OnLogicalRouterStaticRouteCreate(row.New.(*libovsdb.LogicalRouterStaticRoute))
					case tableLogicalSwitch:
						odbi.signalCB.OnLogicalSwitchCreate(row.New.(*libovsdb.LogicalSwitch))
					case tableLogicalSwitchPort:
						odbi.signalCB.OnLogicalSwitchPortCreate(row.New.(*libovsdb.LogicalSwitchPort))
					case tableACL:
						odbi.signalCB.OnACLCreate(row.New.(*libovsdb.ACL))
					case tableDHCPOptions:
						odbi.signalCB.OnDHCPOptionsCreate(row.New.(*libovsdb.DHCPOptions))
					case tableQoS:
						odbi.signalCB.OnQoSCreate(row.New.(*libovsdb.QoS))
					case tableLoadBalancer:
						odbi.signalCB.OnLoadBalancerCreate(row.New.(*libovsdb.LoadBalancer))
					}
				}
			} else {
				defer delete(odbi.cache[table], uuid)

				if odbi.signalCB != nil {
					defer func(table, uuid string) {
						switch table {
						case tableLogicalRouter:
							odbi.signalCB.OnLogicalRouterDelete(row.Old.(*libovsdb.LogicalRouter))
						case tableLogicalRouterPort:
							odbi.signalCB.OnLogicalRouterPortDelete(row.Old.(*libovsdb.LogicalRouterPort))
						case tableLogicalRouterStaticRoute:
							odbi.signalCB.OnLogicalRouterStaticRouteDelete(row.Old.(*libovsdb.LogicalRouterStaticRoute))
						case tableLogicalSwitch:
							odbi.signalCB.OnLogicalSwitchDelete(row.Old.(*libovsdb.LogicalSwitch))
						case tableLogicalSwitchPort:
							odbi.signalCB.OnLogicalSwitchPortDelete(row.Old.(*libovsdb.LogicalSwitchPort))
						case tableACL:
							odbi.signalCB.OnACLDelete(row.Old.(*libovsdb.ACL))
						case tableDHCPOptions:
							odbi.signalCB.OnDHCPOptionsDelete(row.Old.(*libovsdb.DHCPOptions))
						case tableQoS:
							odbi.signalCB.OnQoSDelete(row.Old.(*libovsdb.QoS))
						case tableLoadBalancer:
							odbi.signalCB.OnLoadBalancerDelete(row.Old.(*libovsdb.LoadBalancer))
						}
					}(table, uuid)
				}
			}
		}
	}
}

func (odbi *ovndb) ConvertGoSetToStringArray(oset libovsdb.OvsSet) []string {
	var ret = []string{}
	for _, s := range oset.GoSet {
		value, ok := s.(string)
		if ok {
			ret = append(ret, value)
		}
	}
	return ret
}

func stringToGoUUID(uuid string) libovsdb.UUID {
	return libovsdb.UUID{GoUUID: uuid}
}

func (odbi *ovndb) getRowByName(table string, lrow OVNRow) (interface{}, error) {
	if name, ok := lrow["name"]; ok {
		row := newRow()
		row["name"] = name
		return odbi.getRow(table, row)
	}
	return nil, ErrorNotFound
}

func (odbi *ovndb) getRow(table string, lrow OVNRow) (interface{}, error) {
	ifaces, err := odbi.getRows(table, lrow)
	if err != nil {
		return nil, err
	}

	if len(ifaces) > 1 {
		return nil, ErrorOption
	}

	return ifaces[0], nil
}

func (odbi *ovndb) getRows(table string, lrow OVNRow) ([]interface{}, error) {
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()

	cacheTable, ok := odbi.cache[table]
	if !ok {
		return nil, ErrorNotFound
	}

	var rows []interface{}

	// list lookup
	if lrow == nil || len(lrow) == 0 {
		for _, iface := range cacheTable {
			//			log.Printf("AAA %#+v\n", cacheTable)
			rows = append(rows, iface)
		}
		return rows, nil
	}

	// direct lookup by uuid if it specified
	if uuid, ok := lrow["_uuid"]; ok {
		if iface, ok := cacheTable[uuid.(string)]; ok {
			rows = append(rows, iface)
			return rows, nil
		}
		return nil, ErrorNotFound
	}

	// lookup by other fields
	for _, iface := range cacheTable {
		crow, err := structToMap(iface)
		if err != nil {
			return nil, err
		}
		if cmpRows(crow, lrow) {
			rows = append(rows, iface)
		}
	}

	if len(rows) == 0 {
		return nil, ErrorNotFound
	}

	return rows, nil
}
