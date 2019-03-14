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
	"github.com/socketplane/libovsdb"
)

type LogicalSwitch struct {
	UUID       string
	Name       string
	ExternalID map[interface{}]interface{}
}

func (odbi *ovnDBImp) lswListImp() (*OvnCommand, error) {
	condition := libovsdb.NewCondition("name", "!=", "")
	selectOp := odbi.selectRowOp(tableLogicalSwitch, condition)
	operations := []libovsdb.Operation{selectOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lswAddImp(name string) (*OvnCommand, error) {
	row := make(OVNRow)
	row["name"] = name
	insertOp, err := odbi.insertRowOp(tableLogicalSwitch, row)
	if err != nil {
		return nil, err
	}
	operations := []libovsdb.Operation{insertOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lswDelImp(name string) (*OvnCommand, error) {
	condition := libovsdb.NewCondition("name", "==", name)
	deleteOp := odbi.deleteRowOp(tableLogicalSwitch, condition)
	operations := []libovsdb.Operation{deleteOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) RowToLogicalSwitch(uuid string) *LogicalSwitch {
	odbi.cachemutex.RLock()
	row := odbi.cache[tableLogicalSwitch][uuid]
	odbi.cachemutex.RUnlock()

	return &LogicalSwitch{
		UUID:       uuid,
		Name:       row.Fields["name"].(string),
		ExternalID: row.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}
}

// Get all logical switches
func (odbi *ovnDBImp) GetLogicalSwitches() ([]*LogicalSwitch, error) {
	var lslist = []*LogicalSwitch{}
	odbi.cachemutex.RLock()
	rows, ok := odbi.cache[tableLogicalSwitch]
	odbi.cachemutex.RUnlock()
	if !ok {
		return nil, ErrorNotFound
	}

	for uuid, _ := range rows {
		lslist = append(lslist, odbi.RowToLogicalSwitch(uuid))
	}
	return lslist, nil
}
