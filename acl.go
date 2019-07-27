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
	//	"github.com/google/go-cmp/cmp"
)

// ACL ovnnb item
type ACL struct {
	UUID       string            `ovn:"uuid"`
	Action     string            `ovn:"action"`
	Direction  string            `ovn:"direction"`
	Match      string            `ovn:"match"`
	Priority   int               `ovn:"priority"`
	Log        bool              `ovn:"log"`
	Name       string            `ovn:"name"`
	Severity   string            `ovn:"severity"`
	Meter      string            `ovn:"meter"`
	ExternalID map[string]string `ovn:"external_ids"`
}

func (odbi *ovndb) aclExist(item *ACL) bool {
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	chk := item
	chk.ExternalID = nil
	if _, err := odbi.getRows(tableACL, chk); err == nil {
		return true
	}

	return false
}

func (odbi *ovndb) aclUUID(item *ACL) (string, error) {
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	acls, err := odbi.getRows(tableACL, item)
	if err != nil {
		return "", err
	} else if len(acls) != 1 {
		return "", ErrorOption
	}

	if acl, ok := acls[0].(*ACL); ok {
		return acl.UUID, nil
	}

	return "", ErrorOption
}

func (odbi *ovndb) aclAddImp(ls string, item *ACL) (*OvnCommand, error) {
	if odbi.aclExist(item) {
		return nil, ErrorExist
	}

	namedUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}
	row := rowMarshal(item)

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
	condition := libovsdb.NewCondition("name", "==", ls)

	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{insertOp, mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovndb) aclDelImp(ls string, item *ACL) (*OvnCommand, error) {
	if !odbi.aclExist(item) {
		return nil, ErrorNotFound
	}

	wherecondition := []interface{}{}

	aclUUID, err := odbi.aclUUID(item)
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
	condition := libovsdb.NewCondition("name", "==", ls)

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
	cacheACL, ok := odbi.cache2[tableACL][uuid]
	if !ok {
		return nil
	}

	return cacheACL.(*ACL)
}

func (odbi *ovndb) aclListImp(ls string) ([]*ACL, error) {
	var listACL []*ACL

	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheLogicalSwitch, ok := odbi.cache[tableLogicalSwitch]
	if !ok {
		return nil, ErrorNotFound
	}

	for _, drows := range cacheLogicalSwitch {
		if lsw, ok := drows.Fields["name"].(string); ok && lsw == ls {
			acls := drows.Fields["acls"]
			if acls == nil {
				break
			}
			switch acls.(type) {
			case libovsdb.OvsSet:
				if as, ok := acls.(libovsdb.OvsSet); ok {
					for _, a := range as.GoSet {
						if va, ok := a.(libovsdb.UUID); ok {
							listACL = append(listACL, odbi.cache2[tableACL][va.GoUUID].(*ACL))
						}
					}
				}
			case libovsdb.UUID:
				if va, ok := acls.(libovsdb.UUID); ok {
					listACL = append(listACL, odbi.cache2[tableACL][va.GoUUID].(*ACL))
				}
			}
		}
		break
	}

	return listACL, nil
}
