package goovn

import (
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

func (odbi *ovndb) setUUID(model Model, uuid string) {
	uField := reflect.Indirect(reflect.ValueOf(model)).FieldByName("UUID")
	uField.Set(reflect.ValueOf(uuid))
}
