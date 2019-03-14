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

type AddressSet struct {
	UUID       string
	Name       string
	Addresses  []string
	ExternalID map[interface{}]interface{}
}

func (odbi *ovnDBImp) ASUpdate(name string, addrs []string, external_ids map[string]string) (*OvnCommand, error) {
	row := make(OVNRow)
	row["name"] = name
	addresses, err := libovsdb.NewOvsSet(addrs)
	if err != nil {
		return nil, err
	}

	row["addresses"] = addresses
	if external_ids != nil {
		oMap, err := libovsdb.NewOvsMap(external_ids)
		if err != nil {
			return nil, err
		}
		row["external_ids"] = oMap
	}
	condition := libovsdb.NewCondition("name", "==", name)
	updateOp := odbi.updateRowOp(tableAddressSet, row, condition)
	operations := []libovsdb.Operation{updateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) ASAdd(name string, addrs []string, external_ids map[string]string) (*OvnCommand, error) {
	row := make(OVNRow)
	row["name"] = name
	//should support the -is-exist flag here.

	if external_ids != nil {
		oMap, err := libovsdb.NewOvsMap(external_ids)
		if err != nil {
			return nil, err
		}
		row["external_ids"] = oMap
	}
	addresses, err := libovsdb.NewOvsSet(addrs)
	if err != nil {
		return nil, err
	}
	row["addresses"] = addresses
	insertOp, err := odbi.insertRowOp(tableAddressSet, row)
	if err != nil {
		return nil, err
	}
	operations := []libovsdb.Operation{insertOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) GetASByName(name string) (*AddressSet, error) {
	addresssets, err := odbi.GetAddressSets()
	if err != nil {
		return nil, err
	}
	for _, s := range addresssets {
		if s.Name == name {
			return s, nil
		}
	}
	return nil, ErrorNotFound
}

func (odbi *ovnDBImp) ASDel(name string) (*OvnCommand, error) {
	condition := libovsdb.NewCondition("name", "==", name)
	deleteOp := odbi.deleteRowOp(tableAddressSet, condition)
	operations := []libovsdb.Operation{deleteOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) RowToAddressSet(uuid string) *AddressSet {
	odbi.cachemutex.RLock()
	row := odbi.cache[tableAddressSet][uuid]
	odbi.cachemutex.RUnlock()
	return &AddressSet{
		UUID:       uuid,
		Name:       row.Fields["name"].(string),
		ExternalID: row.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}
}

// Get all addressset
func (odbi *ovnDBImp) GetAddressSets() ([]*AddressSet, error) {
	adlist := make([]*AddressSet, 0, 0)
	odbi.cachemutex.RLock()
	rows, ok := odbi.cache[tableAddressSet]
	odbi.cachemutex.RUnlock()

	if !ok {
		return nil, ErrorNotFound
	}

	for uuid, drows := range rows {
		ta := odbi.RowToAddressSet(uuid)
		addresses := []string{}
		as := drows.Fields["addresses"]
		switch as.(type) {
		case libovsdb.OvsSet:
			//TODO: is it possible return interface type directly instead of GoSet
			if goset, ok := drows.Fields["addresses"].(libovsdb.OvsSet); ok {
				for _, i := range goset.GoSet {
					addresses = append(addresses, i.(string))
				}
			}
		case string:
			if v, ok := drows.Fields["addresses"].(string); ok {
				addresses = append(addresses, v)
			}
		}
		ta.Addresses = addresses
		adlist = append(adlist, ta)
	}
	return adlist, nil
}
