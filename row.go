package goovn

import (
	"fmt"
	"strings"

	"github.com/socketplane/libovsdb"
)

func (odbi *ovnDBImp) mutateRowOp(table string, mutation []interface{}, condition []interface{}) libovsdb.Operation {
	return libovsdb.Operation{
		Op:        opMutate,
		Table:     table,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
}

func (odbi *ovnDBImp) selectRowOp(table string, condition []interface{}) libovsdb.Operation {
	return libovsdb.Operation{
		Op:    opSelect,
		Table: table,
		Where: []interface{}{condition},
	}
}

func (odbi *ovnDBImp) insertRowOp(table string, row OVNRow) (libovsdb.Operation, error) {
	namedUUID, err := newUUID()
	if err != nil {
		return libovsdb.Operation{}, err
	}
	if uuid := odbi.getRowUUID(table, row); len(uuid) > 0 {
		return libovsdb.Operation{}, ErrorExist
	}
	return libovsdb.Operation{
		Op:       opInsert,
		Table:    table,
		Row:      row,
		UUIDName: namedUUID,
	}, nil
}

func (odbi *ovnDBImp) deleteRowOp(table string, condition []interface{}) libovsdb.Operation {
	return libovsdb.Operation{
		Op:    opDelete,
		Table: table,
		Where: []interface{}{condition},
	}
}

func (odbi *ovnDBImp) updateRowOp(table string, row OVNRow, condition []interface{}) libovsdb.Operation {
	return libovsdb.Operation{
		Op:    opUpdate,
		Table: table,
		Row:   row,
		Where: []interface{}{condition},
	}
}

func (odbi *ovnDBImp) getRowUUID(table string, row OVNRow) string {
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()
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

func (odbi *ovnDBImp) getRowUUIDContainsUUID(table, field, uuid string) string {
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()
	for id, drows := range odbi.cache[table] {
		v := fmt.Sprintf("%s", drows.Fields[field])
		if strings.Contains(v, uuid) {
			return id
		}
	}
	return ""
}
