/**
 * Copyright (c) 2017-2019 eBay Inc.
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
	"log"
	"reflect"

	"github.com/ebay/libovsdb"
)

var (
	// ErrorOption used when invalid args specified
	ErrorOption = errors.New("invalid option specified")
	// ErrorSchema used when something wrong in ovnnb
	ErrorSchema = errors.New("table schema error")
	// ErrorNotFound used when object not found in ovnnb
	ErrorNotFound = errors.New("object not found")
	// ErrorMultiple used when multiple object found, but needs only one
	ErrorMultiple = errors.New("multiple objects exists")
	// ErrorExist used when object already exists in ovnnb
	ErrorExist = errors.New("object exist")
)

// OVNRow ovnnb row
type OVNRow map[string]interface{}

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
				iface, err := rowUpdateToStruct(table, uuid, row.New)
				if err != nil {
					log.Printf("%v", err)
					continue
				}
				odbi.cache[table][uuid] = iface
				if odbi.signalCB != nil {
					switch table {
					case tableLogicalRouter:
						odbi.signalCB.OnLogicalRouterCreate(iface.(*LogicalRouter))
					case tableLogicalRouterPort:
						odbi.signalCB.OnLogicalRouterPortCreate(iface.(*LogicalRouterPort))
					case tableLogicalRouterStaticRoute:
						odbi.signalCB.OnLogicalRouterStaticRouteCreate(iface.(*LogicalRouterStaticRoute))
					case tableLogicalSwitch:
						odbi.signalCB.OnLogicalSwitchCreate(iface.(*LogicalSwitch))
					case tableLogicalSwitchPort:
						odbi.signalCB.OnLogicalSwitchPortCreate(iface.(*LogicalSwitchPort))
					case tableACL:
						odbi.signalCB.OnACLCreate(iface.(*ACL))
					case tableDHCPOptions:
						odbi.signalCB.OnDHCPOptionsCreate(iface.(*DHCPOptions))
					case tableQoS:
						odbi.signalCB.OnQoSCreate(iface.(*QoS))
					case tableLoadBalancer:
						odbi.signalCB.OnLoadBalancerCreate(iface.(*LoadBalancer))
					}
				}
			} else {
				defer delete(odbi.cache[table], uuid)
				if row.Old != nil && odbi.signalCB != nil {
					iface, err := rowUpdateToStruct(table, uuid, row.Old)
					if err != nil {
						log.Printf("%v", err)
						continue
					}
					defer func(table, uuid string) {
						switch table {
						case tableLogicalRouter:
							odbi.signalCB.OnLogicalRouterDelete(iface.(*LogicalRouter))
						case tableLogicalRouterPort:
							odbi.signalCB.OnLogicalRouterPortDelete(iface.(*LogicalRouterPort))
						case tableLogicalRouterStaticRoute:
							odbi.signalCB.OnLogicalRouterStaticRouteDelete(iface.(*LogicalRouterStaticRoute))
						case tableLogicalSwitch:
							odbi.signalCB.OnLogicalSwitchDelete(iface.(*LogicalSwitch))
						case tableLogicalSwitchPort:
							odbi.signalCB.OnLogicalSwitchPortDelete(iface.(*LogicalSwitchPort))
						case tableACL:
							odbi.signalCB.OnACLDelete(iface.(*ACL))
						case tableDHCPOptions:
							odbi.signalCB.OnDHCPOptionsDelete(iface.(*DHCPOptions))
						case tableQoS:
							odbi.signalCB.OnQoSDelete(iface.(*QoS))
						case tableLoadBalancer:
							odbi.signalCB.OnLoadBalancerDelete(iface.(*LoadBalancer))
						}
					}(table, uuid)
				}
			}
		}
	}
}

func stringToGoUUID(uuid string) libovsdb.UUID {
	return libovsdb.UUID{GoUUID: uuid}
}

func (odbi *ovndb) getRowByUUID(table string, uuid string, iface interface{}) error {
	row := newRow()
	row["uuid"] = uuid
	return odbi.getRow(table, row, iface)
}

func (odbi *ovndb) getRowByName(table string, name string, iface interface{}) error {
	row := newRow()
	row["name"] = name
	return odbi.getRow(table, row, iface)
}

func (odbi *ovndb) getRow(table string, lrow OVNRow, iface interface{}) error {
	var ifaces []interface{}

	err := odbi.getRows(table, lrow, &ifaces)
	if err != nil {
		return err
	}

	if len(ifaces) > 1 {
		return ErrorMultiple
	}

	valuePtr := reflect.ValueOf(iface)
	value := valuePtr.Elem()
	value.Set(reflect.ValueOf(ifaces[0]))

	return nil
}

func (odbi *ovndb) getRows(table string, lrow OVNRow, ifaces interface{}) error {
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()

	cacheTable, ok := odbi.cache[table]
	if !ok {
		return ErrorNotFound
	}

	valuePtr := reflect.ValueOf(ifaces)
	value := valuePtr.Elem()

	// list lookup
	if lrow == nil || len(lrow) == 0 {
		for _, iface := range cacheTable {
			value.Set(reflect.Append(value, reflect.ValueOf(iface)))
			//			ifaces = append(ifaces, iface)
		}
		return nil
	}

	// direct lookup by uuid if it specified
	if uuid, ok := lrow["uuid"]; ok {
		if iface, ok := cacheTable[uuid.(string)]; ok {
			value.Set(reflect.Append(value, reflect.ValueOf(iface)))
			//			ifaces = append(ifaces, iface)
			return nil
		}
		return ErrorNotFound
	}

	// lookup by other fields
	for _, iface := range cacheTable {
		crow, err := libovsdb.StructToMap(iface, "ovn")
		if err != nil {
			return err
		}
		if cmpRows(crow, lrow) {
			value.Set(reflect.Append(value, reflect.ValueOf(iface)))
			//ifaces = append(ifaces, iface)
		}
	}

	if value.Len() == 0 {
		return ErrorNotFound
	}

	return nil
}
