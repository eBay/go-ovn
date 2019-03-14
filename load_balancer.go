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
	"strings"

	"github.com/socketplane/libovsdb"
)

type LoadBalancer struct {
	UUID       string
	Name       string
	vips       map[interface{}]interface{}
	protocol   string
	ExternalID map[interface{}]interface{}
}

func (odbi *ovnDBImp) lbUpdateImp(name string, vipPort string, proto string, addrs []string) (*OvnCommand, error) {
	row := make(OVNRow)

	// prepare vips map
	vipMap := make(map[string]string)
	vipMap[vipPort] = strings.Join(addrs, ",")

	oMap, err := libovsdb.NewOvsMap(vipMap)
	if err != nil {
		return nil, err
	}
	row["vips"] = oMap
	row["protocol"] = proto

	condition := libovsdb.NewCondition("name", "==", name)
	updateOp := odbi.updateRowOp(tableLoadBalancer, row, condition)
	operations := []libovsdb.Operation{updateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lbAddImp(name string, vipPort string, proto string, addrs []string) (*OvnCommand, error) {
	row := make(OVNRow)
	row["name"] = name

	vipMap := make(map[string]string)
	vipMap[vipPort] = strings.Join(addrs, ",")

	oMap, err := libovsdb.NewOvsMap(vipMap)
	if err != nil {
		return nil, err
	}
	row["vips"] = oMap
	row["protocol"] = proto

	insertOp, err := odbi.insertRowOp(tableLoadBalancer, row)
	if err != nil {
		return nil, err
	}

	mutateUUID := []libovsdb.UUID{{insertOp.UUIDName}}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}

	mutation := libovsdb.NewMutation("load_balancer", opInsert, mutateSet)
	// TODO: Add filter for LS name
	condition := libovsdb.NewCondition("name", "!=", "")
	mutateOp := odbi.mutateRowOp(tableLogicalSwitch, mutation, condition)
	operations := []libovsdb.Operation{insertOp, mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lbDelImp(name string) (*OvnCommand, error) {
	condition := libovsdb.NewCondition("name", "==", name)
	deleteOp := odbi.deleteRowOp(tableLoadBalancer, condition)
	operations := []libovsdb.Operation{deleteOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) GetLB(name string) ([]*LoadBalancer, error) {
	var lbList []*LoadBalancer
	odbi.cachemutex.RLock()
	rows, ok := odbi.cache[tableLoadBalancer]
	odbi.cachemutex.RUnlock()
	if !ok {
		return nil, ErrorNotFound
	}

	for uuid, drows := range rows {
		if lbName, ok := drows.Fields["name"].(string); ok && lbName == name {
			lb := odbi.RowToLB(uuid)
			lbList = append(lbList, lb)
		}
	}

	return lbList, nil
}

func (odbi *ovnDBImp) RowToLB(uuid string) *LoadBalancer {
	odbi.cachemutex.RLock()
	row := odbi.cache[tableLoadBalancer][uuid]
	odbi.cachemutex.RUnlock()
	return &LoadBalancer{
		UUID:       uuid,
		protocol:   row.Fields["protocol"].(string),
		Name:       row.Fields["name"].(string),
		vips:       row.Fields["vips"].(libovsdb.OvsMap).GoMap,
		ExternalID: row.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}
}
