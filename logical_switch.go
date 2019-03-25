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
	UUID         string
	Name         string
	Ports        []string
	LoadBalancer []string
	ACLs         []string
	QoSRules     []string
	DNSRecords   []string
	OtherConfig  map[interface{}]interface{}
	ExternalID   map[interface{}]interface{}
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
		UUID:        uuid,
		Name:        cacheLogicalSwitch.Fields["name"].(string),
		OtherConfig: cacheLogicalSwitch.Fields["other_config"].(libovsdb.OvsMap).GoMap,
		ExternalID:  cacheLogicalSwitch.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}
	if ports, ok := cacheLogicalSwitch.Fields["ports"]; ok {
		switch ports.(type) {
		case libovsdb.UUID:
			ls.Ports = []string{ports.(libovsdb.UUID).GoUUID}
		case libovsdb.OvsSet:
			ls.Ports = odbi.ConvertGoSetToStringArray(ports.(libovsdb.OvsSet))
		}
	}
	if lbs, ok := cacheLogicalSwitch.Fields["load_balancer"]; ok {
		switch lbs.(type) {
		case libovsdb.UUID:
			ls.LoadBalancer = []string{lbs.(libovsdb.UUID).GoUUID}
		case libovsdb.OvsSet:
			ls.LoadBalancer = odbi.ConvertGoSetToStringArray(lbs.(libovsdb.OvsSet))
		}
	}
	if acls, ok := cacheLogicalSwitch.Fields["acls"]; ok {
		switch acls.(type) {
		case libovsdb.UUID:
			ls.ACLs = []string{acls.(libovsdb.UUID).GoUUID}
		case libovsdb.OvsSet:
			ls.ACLs = odbi.ConvertGoSetToStringArray(acls.(libovsdb.OvsSet))
		}
	}
	if qosrules, ok := cacheLogicalSwitch.Fields["qos_rules"]; ok {
		switch qosrules.(type) {
		case libovsdb.UUID:
			ls.QoSRules = []string{qosrules.(libovsdb.UUID).GoUUID}
		case libovsdb.OvsSet:
			ls.QoSRules = odbi.ConvertGoSetToStringArray(qosrules.(libovsdb.OvsSet))
		}
	}
	if dnsrecords, ok := cacheLogicalSwitch.Fields["dns_records"]; ok {
		switch dnsrecords.(type) {
		case libovsdb.UUID:
			ls.DNSRecords = []string{dnsrecords.(libovsdb.UUID).GoUUID}
		case libovsdb.OvsSet:
			ls.DNSRecords = odbi.ConvertGoSetToStringArray(dnsrecords.(libovsdb.OvsSet))
		}
	}

	return ls
}

func (odbi *ovnDBImp) GetLogicalSwitchByName(ls string) (*LogicalSwitch, error) {
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheLogicalSwitch, ok := odbi.cache[tableLogicalSwitch]
	if !ok {
		return nil, ErrorNotFound
	}

	for uuid, drows := range cacheLogicalSwitch {
		if rlsw, ok := drows.Fields["name"].(string); ok && rlsw == ls {
			return odbi.rowToLogicalSwitch(uuid), nil
		}
	}

	return nil, ErrorNotFound
}

// Get all logical switches
func (odbi *ovnDBImp) GetLogicalSwitches() ([]*LogicalSwitch, error) {
	var listLS []*LogicalSwitch

	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheLogicalSwitch, ok := odbi.cache[tableLogicalSwitch]
	if !ok {
		return nil, ErrorSchema
	}

	for uuid, _ := range cacheLogicalSwitch {
		listLS = append(listLS, odbi.rowToLogicalSwitch(uuid))
	}

	return listLS, nil
}
