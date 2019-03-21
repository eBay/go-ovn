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

type LogicalSwitch struct {
	UUID       string
	Name       string
	ExternalID map[interface{}]interface{}
}

func (odbi *ovnDBImp) lswListImp() (*OvnCommand, error) {
	condition := libovsdb.NewCondition("name", "!=", "")
	selectOp := libovsdb.Operation{
		Op:    opSelect,
		Table: tableLogicalSwitch,
		Where: []interface{}{condition},
	}

	operations := []libovsdb.Operation{selectOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lswAddImp(lsw string) (*OvnCommand, error) {
	namedUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}

	//row to insert
	lswitch := make(OVNRow)
	lswitch["name"] = lsw

	if uuid := odbi.getRowUUID(tableLogicalSwitch, lswitch); len(uuid) > 0 {
		return nil, ErrorExist
	}

	insertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableLogicalSwitch,
		Row:      lswitch,
		UUIDName: namedUUID,
	}
	operations := []libovsdb.Operation{insertOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lswDelImp(lsw string) (*OvnCommand, error) {
	condition := libovsdb.NewCondition("name", "==", lsw)
	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableLogicalSwitch,
		Where: []interface{}{condition},
	}
	operations := []libovsdb.Operation{deleteOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) rowToLogicalSwitch(uuid string) *LogicalSwitch {
	cacheLogicalSwitch, ok := odbi.cache[tableLogicalSwitch][uuid]
	if !ok {
		return nil
	}

	ls := &LogicalSwitch{
		UUID:       uuid,
		Name:       cacheLogicalSwitch.Fields["name"].(string),
		ExternalID: cacheLogicalSwitch.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}
	return ls
}

// Get all logical switches
func (odbi *ovnDBImp) GetLogicalSwitches() ([]*LogicalSwitch, error) {
	var listLS []*LogicalSwitch

	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheLogicalSwitch, ok := odbi.cache[tableLogicalSwitch]
	if !ok {
		return nil, ErrorNotFound
	}

	for uuid, _ := range cacheLogicalSwitch {
		listLS = append(listLS, odbi.rowToLogicalSwitch(uuid))
	}

	return listLS, nil
}
