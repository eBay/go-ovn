package goovn

import (
	"fmt"
	"log"
	"reflect"

	"github.com/ebay/libovsdb"
)

func (odbi *ovndb) populateCacheORM(updates libovsdb.TableUpdates) {
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	api, err := odbi.client.ORM(odbi.db)
	if err != nil {
		// TODO: Propagate errors back to main thread
		log.Printf("Failed to get ORM API")
		return
	}

	for table := range odbi.dbModel.types {
		tableUpdate, ok := updates.Updates[table]
		if !ok {
			continue
		}
		if _, ok := odbi.ormCache[table]; !ok {
			odbi.ormCache[table] = make(ORMTableCache)
		}

		for uuid, row := range tableUpdate.Rows {
			// TODO: this is a workaround for the problem of
			// missing json number conversion in libovsdb
			odbi.float64_to_int(row.New)
			if !reflect.DeepEqual(row.New, emptyRow) {
				model, err := odbi.dbModel.newModel(table)
				if err != nil {
					// TODO: Propagate errors back to main thread
					log.Printf("Error creating Model from table %s %s\n", table, err.Error())
				}
				err = api.GetRowData(table, &row.New, model)
				if err != nil {
					// TODO: Propagate errors back to main thread
					log.Printf("Error getting row data %s\n", err.Error())
				}
				odbi.setUUID(model, uuid)

				// Store model in cache
				if reflect.DeepEqual(model, odbi.ormCache[table][uuid]) {
					continue
				}
				odbi.ormCache[table][uuid] = model
				if odbi.ormSignalCB != nil {
					odbi.ormSignalCB.OnCreated(model)
				}
			} else {
				defer delete(odbi.ormCache[table], uuid)
				if odbi.ormSignalCB != nil {
					defer odbi.ormSignalCB.OnDeleted(odbi.ormCache[table][uuid].(Model))
				}

			}
		}
	}
}

// List is a generic function capable of returning (through a provided pointer)
// a list of instances of any row in the cache. It only works on ORM mode.
// 'result' must be a pointer to an slice of the ORM structs that shall be retrived
// The items are appended to the given (pointer to) slice until its capability is reached.
// If the slice is null, all of the table cache will be returned
func (odbi *ovndb) List(result interface{}) error {
	if odbi.mode != ORM {
		return fmt.Errorf("List() is only available in ORM mode")
	}

	resultPtr := reflect.ValueOf(result)
	if resultPtr.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("List() result must be a pointer")
	}

	resultVal := reflect.Indirect(resultPtr)
	if resultVal.Type().Kind() != reflect.Slice {
		return fmt.Errorf("List() result must be a pointer to slice")
	}

	// DBModel stores pointer to structs, slice should have structs, so calling PtrTo
	table := odbi.findTable(reflect.PtrTo(resultVal.Type().Elem()))
	if table == "" {
		return fmt.Errorf("Schema error: finding table for types %s. Table content %+v", reflect.PtrTo(resultVal.Type().Elem()), odbi.dbModel.types)
	}

	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	tableCache, ok := odbi.ormCache[table]
	if !ok {
		return ErrorNotFound
	}

	// If given a null slice, fill it in the cache table completely, if not, just up to
	// its capability
	if resultVal.IsNil() {
		resultVal.Set(reflect.MakeSlice(resultVal.Type(), 0, len(tableCache)))
	}
	i := resultVal.Len()
	for _, elem := range tableCache {
		if i >= resultVal.Cap() {
			break
		}
		resultVal.Set(reflect.Append(resultVal, reflect.Indirect(reflect.ValueOf(elem))))
		i++
	}
	return nil
}

// findTable returns the TableName associated with a reflect.Type or ""
func (odbi *ovndb) findTable(mType reflect.Type) TableName {
	for table, tType := range odbi.dbModel.types {
		if tType == mType {
			return table
		}
	}
	return ""
}

func (odbi *ovndb) setUUID(model Model, uuid string) {
	uField := reflect.Indirect(reflect.ValueOf(model)).FieldByName("UUID")
	uField.Set(reflect.ValueOf(uuid))
}
