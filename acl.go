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

	getRow := newRow()
	if uuid, ok := optRow["entity_uuid"]; ok {
		getRow["uuid"] = uuid
	} else if name, ok := optRow["entity_name"]; ok {
		getRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var ls *LogicalSwitch
	if err := imp.odbi.getRow(tableLogicalSwitch, getRow, &ls); err != nil && err != ErrorNotFound {
		return nil, err
	}
	var pg *PortGroup
	if err := imp.odbi.getRow(tablePortGroup, getRow, &pg); err != nil && err != ErrorNotFound {
		return nil, err
	}

	var acls []*ACL
	if err := imp.odbi.getRows(tableACL, row, &acls); err != nil && err != ErrorNotFound {
		return nil, err
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

	if ls != nil {
		for _, uuid := range ls.ACLs {
			for _, acl := range acls {
				if uuid == acl.UUID {
					return nil, ErrorExist
				}
			}
		}
	}
	if pg != nil {
		for _, uuid := range pg.ACLs {
			for _, acl := range acls {
				if uuid == acl.UUID {
					return nil, ErrorExist
				}
			}
		}
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

func (imp *aclImp) Del(opts ...ACLOpt) (*OvnCommand, error) {
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
	if external_ids, ok := optRow["external_ids"]; ok {
		row["external_ids"] = external_ids
	}

	lsRow := newRow()
	if uuid, ok := optRow["entity_uuid"]; ok {
		lsRow["uuid"] = uuid
	} else if name, ok := optRow["entity_name"]; ok {
		lsRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var ls *LogicalSwitch
	if err := imp.odbi.getRow(tableLogicalSwitch, lsRow, &ls); err != nil {
		return nil, err
	}

	var acls []*ACL
	if err := imp.odbi.getRows(tableACL, row, &acls); err != nil {
		return nil, err
	}

	var aclUUID string

	// suboptimal search
forLoop:
	for _, lsACL := range ls.ACLs {
		for _, acl := range acls {
			if acl.UUID == lsACL {
				aclUUID = acl.UUID
				break forLoop
			}
		}
	}

	aclCondition := libovsdb.NewCondition("_uuid", "==", stringToGoUUID(aclUUID))
	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableACL,
		Where: aclCondition,
	}

	mutateUUID := []libovsdb.UUID{stringToGoUUID(aclUUID)}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}

	mutation := libovsdb.NewMutation("acls", opDelete, mutateSet)
	lsCondition := libovsdb.NewCondition("_uuid", "==", ls.UUID)

	// Simple mutate operation
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{lsCondition},
	}
	operations := []libovsdb.Operation{mutateOp, deleteOp}
	return &OvnCommand{operations, imp.odbi, make([][]map[string]interface{}, len(operations))}, nil
}

// List acls for logical switch or port group
func (imp *aclImp) List(opt ACLOpt) ([]*ACL, error) {
	var err error

	optGetRow := newRow()
	if err := opt(optGetRow); err != nil {
		return nil, err
	}

	getRow := newRow()
	var aclUUIDs []string
	if uuid, ok := optGetRow["entity_uuid"]; ok {
		getRow["uuid"] = uuid
	} else if name, ok := optGetRow["entity_name"]; ok {
		getRow["name"] = name
	} else {
		return nil, ErrorOption
	}

	var ls *LogicalSwitch
	if err := imp.odbi.getRow(tableLogicalSwitch, getRow, &ls); err != nil && err != ErrorNotFound {
		return nil, err
	}
	var pg *PortGroup
	if err := imp.odbi.getRow(tablePortGroup, getRow, &pg); err != nil && err != ErrorNotFound {
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

	aclList := make([]*ACL, len(aclUUIDs))

	for i := 0; i < len(aclUUIDs); i++ {
		if err = imp.odbi.getRowByUUID(tableACL, aclUUIDs[i], &aclList[i]); err != nil {
			return nil, err
		}
	}

	return aclList, nil
}
