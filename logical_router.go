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

type LogicalRouter struct {
	UUID    string
	Name    string
	Enabled bool

	Ports        []string
	StaticRoutes []string
	NAT          []string
	LoadBalancer []string

	Options    map[interface{}]interface{}
	ExternalID map[interface{}]interface{}
}

func (odbi *ovnDBImp) lrAddImp(name string, external_ids map[string]string) (*OvnCommand, error) {
	namedUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}

	row := make(OVNRow)
	row["name"] = name

	if external_ids != nil {
		oMap, err := libovsdb.NewOvsMap(external_ids)
		if err != nil {
			return nil, err
		}
		row["external_ids"] = oMap
	}

	if uuid := odbi.getRowUUID(tableLogicalRouter, row); len(uuid) > 0 {
		return nil, ErrorExist
	}

	insertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableLogicalRouter,
		Row:      row,
		UUIDName: namedUUID,
	}

	operations := []libovsdb.Operation{insertOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lrDelImp(name string) (*OvnCommand, error) {
	condition := libovsdb.NewCondition("name", "==", name)
	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableLogicalRouter,
		Where: []interface{}{condition},
	}
	operations := []libovsdb.Operation{deleteOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) GetLogicalRouter(name string) ([]*LogicalRouter, error) {
	var lrList []*LogicalRouter

	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheLogicalRouter, ok := odbi.cache[tableLogicalRouter]
	if !ok {
		return nil, ErrorNotFound
	}

	for uuid, drows := range cacheLogicalRouter {
		if lrName, ok := drows.Fields["name"].(string); ok && lrName == name {
			lr := odbi.rowToLogicalRouter(uuid)
			lrList = append(lrList, lr)
		}
	}
	return lrList, nil
}

func (odbi *ovnDBImp) rowToLogicalRouter(uuid string) *LogicalRouter {
	lr := &LogicalRouter{
		UUID:       uuid,
		Name:       odbi.cache[tableLogicalRouter][uuid].Fields["name"].(string),
		Options:    odbi.cache[tableLogicalRouter][uuid].Fields["options"].(libovsdb.OvsMap).GoMap,
		ExternalID: odbi.cache[tableLogicalRouter][uuid].Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}

	if enabled, ok := odbi.cache[tableLogicalRouter][uuid].Fields["enabled"]; ok {
		switch enabled.(type) {
		case bool:
			lr.Enabled = enabled.(bool)
		case libovsdb.OvsSet:
			if enabled.(libovsdb.OvsSet).GoSet == nil {
				lr.Enabled = true
			}
		}
	}

	ports := odbi.cache[tableLogicalRouter][uuid].Fields["ports"]
	switch ports.(type) {
	case string:
		lr.Ports = []string{ports.(string)}
	case libovsdb.OvsSet:
		lr.Ports = odbi.ConvertGoSetToStringArray(ports.(libovsdb.OvsSet))
	}

	return lr
}

// Get all logical switches
func (odbi *ovnDBImp) GetLogicalRouters() ([]*LogicalRouter, error) {
	var listLR []*LogicalRouter

	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheLogicalRouter, ok := odbi.cache[tableLogicalRouter]
	if !ok {
		return nil, ErrorNotFound
	}

	for uuid, _ := range cacheLogicalRouter {
		listLR = append(listLR, odbi.rowToLogicalRouter(uuid))
	}

	return listLR, nil
}
