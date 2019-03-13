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
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/socketplane/libovsdb"
)

var (
	ErrorNotFound = errors.New("object not found")
	ErrorExist    = errors.New("object exist")
)

type OVNRow map[string]interface{}

func newNBImp(client *ovnDBClient, callback OVNSignal) (*ovnDBImp, error) {
	nbimp := &ovnDBImp{
		client: client,
		cache:  make(map[string]map[string]libovsdb.Row),
	}
	initial, err := nbimp.client.dbclient.MonitorAll(NBDB, "")
	if err != nil {
		return nil, err
	}
	nbimp.populateCache(*initial)
	notifier := ovnNotifier{nbimp}
	nbimp.client.dbclient.Register(notifier)
	nbimp.callback = callback
	return nbimp, nil
}

func (odbi *ovnDBImp) getRowUUID(table string, row OVNRow) string {
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	for uuid, drows := range odbi.cache[table] {
		found := false
		for field, value := range row {
			if v, ok := drows.Fields[field]; ok {
				if v == value {
					found = true
				} else {
					found = false
					break
				}
			}
		}
		if found {
			return uuid
		}
	}
	return ""
}

//test if map s contains t
//This function is not both s and t are nil at same time
func (odbi *ovnDBImp) oMapContians(s, t map[interface{}]interface{}) bool {
	if s == nil || t == nil {
		return false
	}

	for tk, tv := range t {
		if sv, ok := s[tk]; !ok {
			return false
		} else if tv != sv {
			return false
		}
	}
	return true
}

func (odbi *ovnDBImp) getRowUUIDContainsUUID(table, field, uuid string) (string, error) {
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	for id, drows := range odbi.cache[table] {
		v := fmt.Sprintf("%s", drows.Fields[field])
		if strings.Contains(v, uuid) {
			return id, nil
		}
	}
	return "", ErrorNotFound
}

func (odbi *ovnDBImp) transact(ops ...libovsdb.Operation) ([]libovsdb.OperationResult, error) {
	// Only support one trans at same time now.
	odbi.tranmutex.Lock()
	defer odbi.tranmutex.Unlock()
	reply, err := odbi.client.dbclient.Transact(NBDB, ops...)

	if err != nil {
		return reply, err
	}

	if len(reply) < len(ops) {
		for i, o := range reply {
			if o.Error != "" && i < len(ops) {
				return nil, errors.New(fmt.Sprint("Transaction Failed due to an error :", o.Error, " details:", o.Details, " in ", ops[i]))
			}
		}
		return reply, errors.New(fmt.Sprint("Number of Replies should be atleast equal to number of operations"))
	}
	return reply, nil
}

func (odbi *ovnDBImp) Execute(cmds ...*OvnCommand) error {
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

func (odbi *ovnDBImp) float64_to_int(row libovsdb.Row) {
	for field, value := range row.Fields {
		if v, ok := value.(float64); ok {
			n := int(v)
			if float64(n) == v {
				row.Fields[field] = n
			}
		}
	}
}

func (odbi *ovnDBImp) populateCache(updates libovsdb.TableUpdates) {
	empty := libovsdb.Row{}
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	for table, tableUpdate := range updates.Updates {
		if _, ok := odbi.cache[table]; !ok {
			odbi.cache[table] = make(map[string]libovsdb.Row)
		}
		for uuid, row := range tableUpdate.Rows {
			// TODO: this is a workaround for the problem of
			// missing json number conversion in libovsdb
			odbi.float64_to_int(row.New)

			if !reflect.DeepEqual(row.New, empty) {
				odbi.cache[table][uuid] = row.New
				if odbi.callback != nil {
					switch table {
					case tableLogicalSwitch:
						ls := odbi.RowToLogicalSwitch(uuid)
						odbi.callback.OnLogicalSwitchCreate(ls)
					case tableLogicalSwitchPort:
						lp := odbi.RowToLogicalPort(uuid)
						odbi.callback.OnLogicalPortCreate(lp)
					case tableACL:
						acl := odbi.RowToACL(uuid)
						odbi.callback.OnACLCreate(acl)
					case tableDHCPOptions:
						log.Printf("row %#+v uuid %v\n", row, uuid)

						dhcp := odbi.RowToDHCPOptions(uuid)
						odbi.callback.OnDHCPOptionsCreate(dhcp)
					}
				}
			} else {
				if odbi.callback != nil {
					switch table {
					case tableLogicalSwitch:
						ls := odbi.RowToLogicalSwitch(uuid)
						odbi.callback.OnLogicalSwitchDelete(ls)
					case tableLogicalSwitchPort:
						lp := odbi.RowToLogicalPort(uuid)
						odbi.callback.OnLogicalPortDelete(lp)
					case tableACL:
						acl := odbi.RowToACL(uuid)
						odbi.callback.OnACLDelete(acl)
					case tableDHCPOptions:
						dhcp := odbi.RowToDHCPOptions(uuid)
						odbi.callback.OnDHCPOptionsDelete(dhcp)
					}
				}
				delete(odbi.cache[table], uuid)
			}
		}
	}
}

func (odbi *ovnDBImp) ConvertGoSetToStringArray(oset libovsdb.OvsSet) []string {
	var ret = []string{}
	for _, s := range oset.GoSet {
		value, ok := s.(string)
		if ok {
			ret = append(ret, value)
		}
	}
	return ret
}
