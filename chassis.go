/**
 * Copyright (c) 2020 eBay Inc.
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

// Chassis table OVN SB
type Chassis struct {
	UUID                string
	Encaps              []string
	ExternalID          map[interface{}]interface{}
	Hostname            string
	Name                string
	NbCfg               int
	TransportZones      []string
	VtepLogicalSwitches []string
}

func (odbi *ovndb) chassisAddImp(name string, hostname string, etype string, ip string) (*OvnCommand, error) {
	// / Prepare for encap record
	enCapUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}
	row := make(OVNRow)

	if len(name) > 0 {
		row["chassis_name"] = name
	}

	if len(etype) > 0 {
		row["type"] = etype
	}

	if len(ip) > 0 {
		row["ip"] = ip
	}
	if uuid := odbi.getRowUUID(tableEncap, row); len(uuid) > 0 {
		return nil, ErrorExist
	}

	insertEncapOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableEncap,
		Row:      row,
		UUIDName: enCapUUID,
	}
	// Prepare for chassis record
	ChassisUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}

	rowChassis := make(OVNRow)

	var encap_id []string
	var namedUUID = "named-uuid"
	encap_id = append(encap_id, namedUUID)
	encap_id = append(encap_id, enCapUUID)

	rowChassis["encaps"] = encap_id
	if len(name) > 0 {
		rowChassis["name"] = name
	}
	if len(hostname) > 0 {
		rowChassis["hostname"] = hostname
	}
	if uuid := odbi.getRowUUID(tableChassis, row); len(uuid) > 0 {
		return nil, ErrorExist
	}

	insertChassisOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableChassis,
		Row:      rowChassis,
		UUIDName: ChassisUUID,
	}
	operations := []libovsdb.Operation{insertEncapOp, insertChassisOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil

}

func (odbi *ovndb) chassisDelImp(name string) (*OvnCommand, error) {
	var operations []libovsdb.Operation

	condition := libovsdb.NewCondition("name", "==", name)
	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableChassis,
		Where: []interface{}{condition},
	}
	operations = append(operations, deleteOp)
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovndb) chassisGetImp(hostname string) ([]*Chassis, error) {
	var listChassis []*Chassis

	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheChassis, ok := odbi.cache[tableChassis]

	if !ok {
		return nil, ErrorSchema
	}

	for uuid, drows := range cacheChassis {
		if chName, ok := drows.Fields["hostname"].(string); ok && chName == hostname {
			ch, err := odbi.rowToChassis(uuid)
			if err != nil {
				return nil, err
			}
			listChassis = append(listChassis, ch)
		}
	}
	return listChassis, nil
}

func (odbi *ovndb) rowToChassis(uuid string) (*Chassis, error) {
	ch := &Chassis{
		UUID:       uuid,
		Name:       odbi.cache[tableChassis][uuid].Fields["name"].(string),
		Hostname:   odbi.cache[tableChassis][uuid].Fields["hostname"].(string),
		ExternalID: odbi.cache[tableChassis][uuid].Fields["external_ids"].(libovsdb.OvsMap).GoMap,
		NbCfg:      odbi.cache[tableChassis][uuid].Fields["nb_cfg"].(int),
	}
	return ch, nil
}
