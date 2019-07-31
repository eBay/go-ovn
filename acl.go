/**
 * Copyright (c) 2019 eBay Inc.
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

// ACL implementation methods
type aclImp struct {
	odbi *ovndb
}

// ACLOpt option to acl commands
type ACLOpt func(OVNRow) error

// ACLEntityName pass entity name for acl
func ACLEntityName(e string) ACLOpt {
	return func(o OVNRow) error {
		if len(e) == 0 {
			return ErrorOption
		}
		o["entity_name"] = e
		return nil
	}
}

// ACLEntityUUID pass entity uuid for acl
func ACLEntityUUID(e string) ACLOpt {
	return func(o OVNRow) error {
		if len(e) == 0 {
			return ErrorOption
		}
		o["entity_uuid"] = e
		return nil
	}
}

// ACLDirection pass direction for acl
func ACLDirection(d string) ACLOpt {
	return func(o OVNRow) error {
		switch d {
		case "from-lport", "to-lport":
			o["direction"] = d
		default:
			return ErrorOption
		}
		return nil
	}
}

// ACLPiority pass priority for acl
func ACLPriority(p int) ACLOpt {
	return func(o OVNRow) error {
		if p < 0 || p > 32767 {
			return ErrorOption
		}
		o["priority"] = p
		return nil
	}
}

// ACLLog pass log for acl
func ACLLog(b bool) ACLOpt {
	return func(o OVNRow) error {
		o["log"] = b
		return nil
	}
}

// ACLLog pass severity for acl
func ACLSeverity(s string) ACLOpt {
	return func(o OVNRow) error {
		switch s {
		case "alert", "debug", "info", "notice", "warning":
			o["severity"] = s
			o["log"] = true
		default:
			return ErrorOption
		}
		return nil
	}
}

// ACLName pass name for acl log
func ACLName(n string) ACLOpt {
	return func(o OVNRow) error {
		if l := len(n); l < 1 || l > 63 {
			return ErrorOption
		}
		o["name"] = n
		return nil
	}
}

// ACLMeter pass meter for acl log
func ACLMeter(m string) ACLOpt {
	return func(o OVNRow) error {
		if len(m) == 0 {
			return ErrorOption
		}
		o["meter"] = m
		return nil
	}
}

// ACLMatch pass match for acl
func ACLMatch(m string) ACLOpt {
	return func(o OVNRow) error {
		if len(m) == 0 {
			return ErrorOption
		}
		o["match"] = m
		return nil
	}
}

// ACLAction pass action for acl
func ACLAction(a string) ACLOpt {
	return func(o OVNRow) error {
		switch a {
		case "allow-related", "allow", "drop", "reject":
			o["action"] = a
		default:
			return ErrorOption
		}
		return nil
	}
}

// ACLExternalIDs pass external_ids for acl
func ACLExternalIDs(m map[string]string) ACLOpt {
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

func (imp *aclImp) Add(opts ...ACLOpt) (*OvnCommand, error) {
	optRow := newRow()

	if len(opts) < 5 {
		return nil, ErrorOption
	}

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
	if direction, ok := optRow["direction"]; ok {
		row["direction"] = direction
	}
	if priority, ok := optRow["priority"]; ok {
		row["priority"] = priority
	}
	if match, ok := optRow["match"]; ok {
		row["match"] = match
	}
	if action, ok := optRow["action"]; ok {
		row["action"] = action
	}
	if log, ok := optRow["log"]; ok {
		row["log"] = log
	}
	if severity, ok := optRow["severity"]; ok {
		row["severity"] = severity
	}
	if meter, ok := optRow["meter"]; ok {
		row["meter"] = meter
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

	var entityName string
	var entityUUID string

	var ls *LogicalSwitch
	if uuid, ok := optRow["entity_uuid"]; ok {
		if err := imp.odbi.getRowByUUID(tableLogicalSwitch, uuid, &ls); err != nil && err != ErrorNotFound {
			return nil, err
		}
	} else if name, ok := optRow["entity_name"]; ok {
		if err := imp.odbi.getRowByName(tableLogicalSwitch, name, &ls); err != nil && err != ErrorNotFound {
			return nil, err
		}
	} else {
		return nil, ErrorOption
	}

	var pg *PortGroup
	if uuid, ok := optRow["entity_uuid"]; ok {
		if err := imp.odbi.getRowByUUID(tablePortGroup, uuid, &pg); err != nil && err != ErrorNotFound {
			return nil, err
		}
	} else if name, ok := optRow["entity_name"]; ok {
		if err := imp.odbi.getRowByName(tablePortGroup, name, &pg); err != nil && err != ErrorNotFound {
			return nil, err
		}
	} else {
		return nil, ErrorOption
	}

	if ls != nil && pg != nil {
		if etype, ok := optRow["entity_type"]; ok {
			switch etype.(string) {
			case "switch":
				entityName = tableLogicalSwitch
				entityUUID = ls.UUID
			case "port-group":
				entityName = tablePortGroup
				entityUUID = pg.UUID
			}
		} else {
			return nil, ErrorOption
		}
	} else if ls != nil {
		entityName = tableLogicalSwitch
		entityUUID = ls.UUID
	} else if pg != nil {
		entityName = tablePortGroup
		entityUUID = pg.UUID
	} else {
		return nil, ErrorNotFound
	}

	insertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableACL,
		Row:      row,
		UUIDName: namedUUID,
	}

	mutateUUID := []libovsdb.UUID{stringToGoUUID(namedUUID)}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("acls", opInsert, mutateSet)
	condition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(entityUUID))

	// simple mutate operation
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     entityName,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{insertOp, mutateOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

/*
func (odbi *ovndb) aclDelImp(lsw, direct, match string, priority int, external_ids map[string]string) (*OvnCommand, error) {
	row := make(OVNRow)

	wherecondition := []interface{}{}
	if direct != "" {
		row["direction"] = direct
	}
	if match != "" {
		row["match"] = match
	}
	//in ovn pirority is greater than/equal 0,
	//if input the priority < 0, lots of acls will be deleted if matches direct and match condition judgement.
	if priority >= 0 {
		row["priority"] = priority
	}

	if external_ids != nil {
		oMap, err := libovsdb.NewOvsMap(external_ids)
		if err != nil {
			return nil, err
		}
		row["external_ids"] = oMap
	}

	aclUUID, err := odbi.getACLUUIDByRow(lsw, tableACL, row)
	if err != nil {
		return nil, err
	}

	uuidcondition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(aclUUID))
	wherecondition = append(wherecondition, uuidcondition)
	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableACL,
		Where: wherecondition,
	}

	mutation := libovsdb.NewMutation("acls", opDelete, stringToGoUUID(aclUUID))
	condition := libovsdb.NewCondition("name", "==", lsw)

	// Simple mutate operation
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{mutateOp, deleteOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovndb) rowToACL(uuid string) *ACL {
	cacheACL, ok := odbi.cache[tableACL][uuid]
	if !ok {
		return nil
	}

	acl := &ACL{
		UUID:       uuid,
		Action:     cacheACL.Fields["action"].(string),
		Direction:  cacheACL.Fields["direction"].(string),
		Match:      cacheACL.Fields["match"].(string),
		Priority:   cacheACL.Fields["priority"].(int),
		Log:        cacheACL.Fields["log"].(bool),
		ExternalID: cacheACL.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}

	return acl
}
*/

// Get all acl by lswitch
func (imp *aclImp) List(opt ACLOpt) ([]*ACL, error) {
	var err error

	optLSRow := newRow()
	if err := opt(optLSRow); err != nil {
		return nil, err
	}

	var aclUUIDs []string
	if uuid, ok := optLSRow["entity_uuid"]; ok {
		var ls *LogicalSwitch
		if err := imp.odbi.getRowByUUID(tableLogicalSwitch, uuid, &ls); err != nil && err != ErrorNotFound {
			return nil, err
		}
		var pg *PortGroup
		if err := imp.odbi.getRowByUUID(tablePortGroup, uuid, &pg); err != nil && err != ErrorNotFound {
			return nil, err
		}
		if ls == nil && pg == nil {
			return nil, ErrorNotFound
		} else if ls != nil && pg != nil {
			return nil, ErrorMultiple
		} else if ls != nil && len(ls.ACLs) > 0 {
			aclUUIDs = append(aclUUIDs, ls.ACLs...)
		} else if pg != nil && len(pg.ACLs) > 0 {
			aclUUIDs = append(aclUUIDs, pg.ACLs...)
		}
	} else if name, ok := optLSRow["entity_name"]; ok {
		var ls *LogicalSwitch
		if err := imp.odbi.getRowByName(tableLogicalSwitch, name, &ls); err != nil && err != ErrorNotFound {
			return nil, err
		}
		var pg *PortGroup
		if err := imp.odbi.getRowByName(tablePortGroup, name, &pg); err != nil && err != ErrorNotFound {
			return nil, err
		}
		if ls == nil && pg == nil {
			return nil, ErrorNotFound
		} else if ls != nil && pg != nil {
			return nil, ErrorMultiple
		} else if ls != nil && len(ls.ACLs) > 0 {
			aclUUIDs = append(aclUUIDs, ls.ACLs...)
		} else if pg != nil && len(pg.ACLs) > 0 {
			aclUUIDs = append(aclUUIDs, pg.ACLs...)
		}
	} else {
		return nil, ErrorOption
	}

	aclList := make([]*ACL, len(aclUUIDs))

	for i := 0; i < len(aclUUIDs); i++ {
		if err = imp.odbi.getRowByUUID(tableACL, aclUUIDs[i], &aclList[i]); err != nil {
			return nil, err
		}
	}

	if len(aclList) == 0 {
		return nil, ErrorNotFound
	}

	return aclList, nil
}
