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

// ACL ovnnb item
type ACL struct {
	UUID       string
	Action     string
	Direction  string
	Match      string
	Priority   int
	Log        bool
	ExternalID map[interface{}]interface{}
}

func (odbi *ovndb) getACLUUIDByRow(lsw, table string, row OVNRow) (string, error) {
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheLogicalSwitch, ok := odbi.cache[tableLogicalSwitch]
	if !ok {
		return "", ErrorSchema
	}

	for _, drows := range cacheLogicalSwitch {
		if rlsw, ok := drows.Fields["name"].(string); ok && rlsw == lsw {
			acls := drows.Fields["acls"]
			if acls != nil {
				switch acls.(type) {
				case libovsdb.OvsSet:
					if as, ok := acls.(libovsdb.OvsSet); ok {
						for _, a := range as.GoSet {
							if va, ok := a.(libovsdb.UUID); ok {
								cacheACL, ok := odbi.cache[tableACL][va.GoUUID]
								if !ok {
									return "", ErrorSchema
								}
								for field, value := range row {
									switch field {
									case "action":
										if cacheACL.Fields["action"].(string) != value {
											goto unmatched
										}
									case "direction":
										if cacheACL.Fields["direction"].(string) != value {
											goto unmatched
										}
									case "match":
										if cacheACL.Fields["match"].(string) != value {
											goto unmatched
										}
									case "priority":
										if cacheACL.Fields["priority"].(int) != value {
											goto unmatched
										}
									case "log":
										if cacheACL.Fields["log"].(bool) != value {
											goto unmatched
										}
									case "external_ids":
										if value != nil && !odbi.oMapContians(cacheACL.Fields["external_ids"].(libovsdb.OvsMap).GoMap, value.(*libovsdb.OvsMap).GoMap) {
											goto unmatched
										}
									}
								}
								return va.GoUUID, nil
							}
						unmatched:
						}
						return "", ErrorNotFound
					}
				case libovsdb.UUID:
					if va, ok := acls.(libovsdb.UUID); ok {
						cacheACL, ok := odbi.cache[tableACL][va.GoUUID]
						if !ok {
							return "", ErrorSchema
						}

						for field, value := range row {
							switch field {
							case "action":
								if cacheACL.Fields["action"].(string) != value {
									goto out
								}
							case "direction":
								if cacheACL.Fields["direction"].(string) != value {
									goto out
								}
							case "match":
								if cacheACL.Fields["match"].(string) != value {
									goto out
								}
							case "priority":
								if cacheACL.Fields["priority"].(int) != value {
									goto out
								}
							case "log":
								if cacheACL.Fields["log"].(bool) != value {
									goto out
								}
							case "external_ids":
								if value != nil && !odbi.oMapContians(cacheACL.Fields["external_ids"].(libovsdb.OvsMap).GoMap, value.(*libovsdb.OvsMap).GoMap) {
									goto out
								}
							}
						}
						return va.GoUUID, nil
					out:
					}
				}
			}
		}
	}
	return "", ErrorNotFound
}

func (odbi *ovndb) aclAddImp(lsw, direct, match, action string, priority int, external_ids map[string]string, logflag bool, meter string) (*OvnCommand, error) {
	namedUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}
	row := make(OVNRow)
	row["direction"] = direct
	row["match"] = match
	row["priority"] = priority

	if external_ids != nil {
		oMap, err := libovsdb.NewOvsMap(external_ids)
		if err != nil {
			return nil, err
		}
		row["external_ids"] = oMap
	}

	_, err = odbi.getACLUUIDByRow(lsw, tableACL, row)
	switch err {
	case ErrorNotFound:
		break
	case nil:
		return nil, ErrorExist
	default:
		return nil, err
	}

	row["action"] = action
	row["log"] = logflag
	if logflag && len(meter) > 0 {
		row["meter"] = meter
	}
	insertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableACL,
		Row:      row,
		UUIDName: namedUUID,
	}

	mutateUUID := []libovsdb.UUID{{namedUUID}}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("acls", opInsert, mutateSet)
	condition := libovsdb.NewCondition("name", "==", lsw)

	// simple mutate operation
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{insertOp, mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

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

	uuidcondition := libovsdb.NewCondition("_uuid", "==", libovsdb.UUID{aclUUID})
	wherecondition = append(wherecondition, uuidcondition)
	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableACL,
		Where: wherecondition,
	}

	mutation := libovsdb.NewMutation("acls", opDelete, libovsdb.UUID{aclUUID})
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

// Get all acl by lswitch
func (odbi *ovndb) aclListImp(lsw string) ([]*ACL, error) {
	var listACL []*ACL

	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheLogicalSwitch, ok := odbi.cache[tableLogicalSwitch]
	if !ok {
		return nil, ErrorNotFound
	}

	for _, drows := range cacheLogicalSwitch {
		if rlsw, ok := drows.Fields["name"].(string); ok && rlsw == lsw {
			acls := drows.Fields["acls"]
			if acls != nil {
				switch acls.(type) {
				case libovsdb.OvsSet:
					if as, ok := acls.(libovsdb.OvsSet); ok {
						for _, a := range as.GoSet {
							if va, ok := a.(libovsdb.UUID); ok {
								ta := odbi.rowToACL(va.GoUUID)
								listACL = append(listACL, ta)
							}
						}
					}
				case libovsdb.UUID:
					if va, ok := acls.(libovsdb.UUID); ok {
						ta := odbi.rowToACL(va.GoUUID)
						listACL = append(listACL, ta)
					}
				}
			}
			break
		}
	}
	return listACL, nil
}
