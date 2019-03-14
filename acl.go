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

type ACL struct {
	UUID      string
	Action    string
	Direction string
	Match     string
	Priority  int
	Log       bool
	Name      string
	Severity  string
	//Meter      string
	ExternalID map[interface{}]interface{}
}

func (odbi *ovnDBImp) getACLUUIDByRow(lsw, table string, row OVNRow) string {
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()

	for _, drows := range odbi.cache[tableLogicalSwitch] {
		if rlsw, ok := drows.Fields["name"].(string); ok && rlsw == lsw {
			acls := drows.Fields["acls"]
			if acls != nil {
				switch acls.(type) {
				case libovsdb.OvsSet:
					if as, ok := acls.(libovsdb.OvsSet); ok {
						for _, a := range as.GoSet {
							if va, ok := a.(libovsdb.UUID); ok {
								for field, value := range row {
									switch field {
									case "action":
										if odbi.cache[tableACL][va.GoUUID].Fields["action"].(string) != value {
											goto unmatched
										}
									case "direction":
										if odbi.cache[tableACL][va.GoUUID].Fields["direction"].(string) != value {
											goto unmatched
										}
									case "match":
										if odbi.cache[tableACL][va.GoUUID].Fields["match"].(string) != value {
											goto unmatched
										}
									case "priority":
										if odbi.cache[tableACL][va.GoUUID].Fields["priority"].(int) != value {
											goto unmatched
										}
									case "log":
										if odbi.cache[tableACL][va.GoUUID].Fields["log"].(bool) != value {
											goto unmatched
										}
									case "external_ids":
										if value != nil && !odbi.oMapContians(odbi.cache[tableACL][va.GoUUID].Fields["external_ids"].(libovsdb.OvsMap).GoMap, value.(*libovsdb.OvsMap).GoMap) {
											goto unmatched
										}
									}
								}
								return va.GoUUID
							}
						unmatched:
						}
						return ""
					}
				case libovsdb.UUID:
					if va, ok := acls.(libovsdb.UUID); ok {
						for field, value := range row {
							switch field {
							case "action":
								if odbi.cache[tableACL][va.GoUUID].Fields["action"].(string) != value {
									goto out
								}
							case "direction":
								if odbi.cache[tableACL][va.GoUUID].Fields["direction"].(string) != value {
									goto out
								}
							case "match":
								if odbi.cache[tableACL][va.GoUUID].Fields["match"].(string) != value {
									goto out
								}
							case "priority":
								if odbi.cache[tableACL][va.GoUUID].Fields["priority"].(int) != value {
									goto out
								}
							case "log":
								if odbi.cache[tableACL][va.GoUUID].Fields["log"].(bool) != value {
									goto out
								}
							case "external_ids":
								if value != nil && !odbi.oMapContians(odbi.cache[tableACL][va.GoUUID].Fields["external_ids"].(libovsdb.OvsMap).GoMap, value.(*libovsdb.OvsMap).GoMap) {
									goto out
								}
							}
						}
						return va.GoUUID
					out:
					}
				}
			}
		}
	}
	return ""
}

func (odbi *ovnDBImp) aclAddImp(lsw, direct, match, action string, priority int, external_ids map[string]string, logflag bool, meter string) (*OvnCommand, error) {
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

	row["action"] = action
	row["log"] = logflag
	if logflag && len(meter) > 0 {
		row["meter"] = meter
	}

	insertOp, err := odbi.insertRowOp(tableACL, row)
	if err != nil {
		return nil, err
	}

	mutateUUID := []libovsdb.UUID{{insertOp.UUIDName}}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("acls", opInsert, mutateSet)
	condition := libovsdb.NewCondition("name", "==", lsw)
	mutateOp := odbi.mutateRowOp(tableLogicalSwitch, mutation, condition)
	operations := []libovsdb.Operation{insertOp, mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) aclDelImp(lsw, direct, match string, priority int, external_ids map[string]string) (*OvnCommand, error) {
	row := make(OVNRow)
	if direct != "" {
		row["direction"] = direct
	}
	if match != "" {
		row["match"] = match
	}
	//in ovn pirority is greater than/equal 0,
	//if input the pirority < 0, lots of acls will be deleted if matches direct and match condition judgement.
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

	uuid := odbi.getACLUUIDByRow(lsw, tableACL, row)
	if len(uuid) == 0 {
		return nil, ErrorNotFound
	}

	condition := libovsdb.NewCondition("_uuid", "==", libovsdb.UUID{uuid})
	deleteOp := odbi.deleteRowOp(tableACL, condition)
	mutation := libovsdb.NewMutation("acls", opDelete, libovsdb.UUID{uuid})
	condition = libovsdb.NewCondition("name", "==", lsw)
	mutateOp := odbi.mutateRowOp(tableLogicalSwitch, mutation, condition)
	operations := []libovsdb.Operation{mutateOp, deleteOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) RowToACL(uuid string) *ACL {
	odbi.cachemutex.RLock()
	row := odbi.cache[tableACL][uuid]
	odbi.cachemutex.RUnlock()
	return &ACL{
		UUID:       uuid,
		Action:     row.Fields["action"].(string),
		Direction:  row.Fields["direction"].(string),
		Match:      row.Fields["match"].(string),
		Priority:   row.Fields["priority"].(int),
		Log:        row.Fields["log"].(bool),
		ExternalID: row.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}
}

// Get all acl by lswitch
func (odbi *ovnDBImp) GetACLsBySwitch(lsw string) ([]*ACL, error) {
	//TODO: should be improvement here, when have lots of acls.
	acllist := make([]*ACL, 0, 0)
	odbi.cachemutex.RLock()
	rows, ok := odbi.cache[tableLogicalSwitch]
	odbi.cachemutex.RUnlock()

	if !ok {
		return nil, ErrorNotFound
	}

	for _, drows := range rows {
		if rlsw, ok := drows.Fields["name"].(string); ok && rlsw == lsw {
			acls := drows.Fields["acls"]
			if acls != nil {
				switch acls.(type) {
				case libovsdb.OvsSet:
					if as, ok := acls.(libovsdb.OvsSet); ok {
						for _, a := range as.GoSet {
							if va, ok := a.(libovsdb.UUID); ok {
								ta := odbi.RowToACL(va.GoUUID)
								acllist = append(acllist, ta)
							}
						}
					}
				case libovsdb.UUID:
					if va, ok := acls.(libovsdb.UUID); ok {
						ta := odbi.RowToACL(va.GoUUID)
						acllist = append(acllist, ta)
					}
				}
			}
			break
		}
	}
	return acllist, nil
}
