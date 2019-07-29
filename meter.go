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
	"math"
	"strings"

	"github.com/ebay/libovsdb"
)

type Meter struct {
	UUID        string
	Name        string                      `json:"name"`
	Unit        string                      `json:"unit"`
	Bands       []string                    `json:"bands"`
	ExternalIds map[interface{}]interface{} `json:"external_ids"`
}

type MeterBand struct {
	UUID        string
	Action      string                      `json:"action"`
	Rate        int                         `json:"rate"`
	BurstSize   int                         `json:"burst_size"`
	ExternalIds map[interface{}]interface{} `json:"external_ids"`
}

/*
The unit specifies the unit for the rate argument; valid  values
are  kbps and pktps for kilobits per second and packets per sec‐
ond, respectively. The burst option configures the maximum burst
allowed for the band in kilobits or packets depending on whether
the unit chosen was kbps or pktps, respectively. If a  burst  is
not supplied, the switch is free to select some reasonable value
depending on its configuration.

ovn-nbctl only supports adding a meter with a single  band,  but
the other commands support meters with multiple bands.

*/

func (odbi *ovndb) rowToMeter(uuid string) *Meter {
	cacheMeter, ok := odbi.cache[tableMeter][uuid]
	if !ok {
		return nil
	}
	meter := &Meter{
		UUID:        uuid,
		Name:        cacheMeter.Fields["name"].(string),
		Unit:        cacheMeter.Fields["unit"].(string),
		Bands:       []string{cacheMeter.Fields["bands"].(libovsdb.UUID).GoUUID},
		ExternalIds: cacheMeter.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}
	return meter
}

func (odbi *ovndb) rowToMeterBand(uuid string) (*MeterBand, error) {
	cacheMeterBand, ok := odbi.cache[tableMeterBand][uuid]
	if !ok {
		return nil, ErrorNotFound
	}
	meterBand := &MeterBand{
		UUID:        uuid,
		Action:      cacheMeterBand.Fields["action"].(string),
		Rate:        cacheMeterBand.Fields["rate"].(int),
		BurstSize:   cacheMeterBand.Fields["burst_size"].(int),
		ExternalIds: cacheMeterBand.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}
	return meterBand, nil
}

/*
ovn-nbctl only supports adding a meter with a single  band,  but
the other commands support meters with multiple bands.

meter-add name action rate unit [burst]
Adds the specified meter. name must be a unique name to identify
this  meter.  The  action  argument specifies what should happen
when this meter is exceeded.
*/
func (odbi *ovndb) meterAddImp(name, action string, rate int, unit string, external_ids map[string]string, burst ...int) (*OvnCommand, error) {

	//Names  that  start  with "__" (two underscores) are reserved for
	//internal use by OVN, so ovn-nbctl does not allow adding them.
	if strings.HasPrefix(name, "__") {
		return nil, ErrorOption
	}

	// The only supported action is drop.
	// If add with wrong option , libovsdb won't accept but doesn't have any feedback. So add some judgment is necessary.
	if action != "drop" {
		return nil, ErrorOption
	}

	MeterUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}
	MeterBandUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}
	//meter row
	mRow := make(OVNRow)

	mRow["name"] = name
	if uuid := odbi.getRowUUID(tableMeter, mRow); len(uuid) > 0 {
		return nil, ErrorExist
	}

	mRow["bands"] = libovsdb.UUID{GoUUID: MeterBandUUID}

	switch unit {
	case "kbps", "pktps":
		mRow["unit"] = unit
	default:
		return nil, ErrorOption
	}

	//Meter Band row
	mbRow := make(OVNRow)

	mbRow["action"] = action

	//rate must be in the range 1...4294967295
	if rate < 0 || rate > math.MaxInt32 {
		return nil, ErrorSchema
	}
	mbRow["rate"] = rate

	if len(burst) != 0 {
		//burst must be in the range 0...4294967295
		if burst[0] < 0 || burst[0] > math.MaxInt32 {
			return nil, ErrorSchema
		}
		mbRow["burst_size"] = burst[0]
	}

	if external_ids != nil {
		oMap, err := libovsdb.NewOvsMap(external_ids)
		if err != nil {
			return nil, err
		}
		mRow["external_ids"] = oMap
		//mbRow["external_ids"] = oMap
	}

	mbInsterOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableMeterBand,
		Row:      mbRow,
		UUIDName: MeterBandUUID,
	}

	mInsertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableMeter,
		Row:      mRow,
		UUIDName: MeterUUID,
	}
	operations := []libovsdb.Operation{mbInsterOp, mInsertOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

/*
meter-del [name]
Deletes meters. By default, all meters are deleted. If  name  is
supplied, only the meter with that name will be deleted.
*/
func (odbi *ovndb) meterDelImp(name ...string) (*OvnCommand, error) {
	var operations []libovsdb.Operation
	var err error

	switch len(name) {
	case 0:
		for uuid := range odbi.cache[tableMeter] {
			name := odbi.cache[tableMeter][uuid].Fields["name"].(string)
			operations, err = odbi.singleMeterDel(name, operations)
			if err != nil {
				return nil, err
			}
		}
	case 1:
		operations, err = odbi.singleMeterDel(name[0], operations)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrorOption
	}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil

}

//meter-list
//Lists all meters.
//but not like ovn-nbctl , it can't show meter bands information
func (odbi *ovndb) meterListImp() ([]*Meter, error) {
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()
	var ListMeter []*Meter
	cacheMeter, ok := odbi.cache[tableMeter]
	if !ok {
		return nil, ErrorNotFound
	}
	for uuid := range cacheMeter {
		ListMeter = append(ListMeter, odbi.rowToMeter(uuid))
	}
	return ListMeter, nil
}

//Because meterList can't show meter bands , add this method as a solution
func (odbi *ovndb) meterBandsListImp() ([]*MeterBand, error) {
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()
	var ListMeterBands []*MeterBand
	cacheMeterBands, ok := odbi.cache[tableMeterBand]
	if !ok {
		return nil, ErrorNotFound
	}
	for uuid := range cacheMeterBands {
		meterBand, err := odbi.rowToMeterBand(uuid)
		if err != nil {
			return nil, ErrorNotFound
		}
		ListMeterBands = append(ListMeterBands, meterBand)
	}
	return ListMeterBands, nil
}

func (odbi *ovndb) singleMeterDel(name string, operations []libovsdb.Operation) ([]libovsdb.Operation, error) {
	meterName := name
	row := make(OVNRow)
	row["name"] = meterName
	meterUUID := odbi.getRowUUID(tableMeter, row)
	if len(meterUUID) == 0 {
		return nil, ErrorNotFound
	}
	bands := odbi.cache[tableMeter][meterUUID].Fields["bands"].(libovsdb.UUID)
	mCondition := libovsdb.NewCondition("name", "==", meterName)
	mDeleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableMeter,
		Where: []interface{}{mCondition},
	}

	bCondition := libovsdb.NewCondition("_uuid", "==", bands)
	bDeleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableMeterBand,
		Where: []interface{}{bCondition},
	}
	operations = append(operations, bDeleteOp)
	operations = append(operations, mDeleteOp)
	return operations, nil
}
