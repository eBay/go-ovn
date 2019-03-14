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

func (odbi *ovnDBImp) lbUpdateImp(name string, vipPort string, protocol string, addrs []string) (*OvnCommand, error) {
	//row to update
	lb := make(OVNRow)

	// prepare vips map
	vipMap := make(map[string]string)
	members := strings.Join(addrs, ",")
	vipMap[vipPort] = members

	oMap, err := libovsdb.NewOvsMap(vipMap)
	if err != nil {
		return nil, err
	}

	lb["vips"] = oMap
	lb["protocol"] = protocol

	condition := libovsdb.NewCondition("name", "==", name)

	insertOp := libovsdb.Operation{
		Op:    opUpdate,
		Table: tableLoadBalancer,
		Row:   lb,
		Where: []interface{}{condition},
	}
	operations := []libovsdb.Operation{insertOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lbAddImp(name string, vipPort string, protocol string, addrs []string) (*OvnCommand, error) {
	namedUUID, err := newUUID()
	if err != nil {
		return nil, err
	}
	//row to insert
	lb := make(OVNRow)
	lb["name"] = name

	if uuid := odbi.getRowUUID(tableLoadBalancer, lb); len(uuid) > 0 {
		return nil, ErrorExist
	}

	// prepare vips map
	vipMap := make(map[string]string)
	members := strings.Join(addrs, ",")
	vipMap[vipPort] = members

	oMap, err := libovsdb.NewOvsMap(vipMap)
	if err != nil {
		return nil, err
	}
	lb["vips"] = oMap
	lb["protocol"] = protocol

	insertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableLoadBalancer,
		Row:      lb,
		UUIDName: namedUUID,
	}

	mutateUUID := []libovsdb.UUID{{namedUUID}}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}

	mutation := libovsdb.NewMutation("load_balancer", opInsert, mutateSet)
	// TODO: Add filter for LS name
	condition := libovsdb.NewCondition("name", "!=", "")

	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{insertOp, mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lbDelImp(name string) (*OvnCommand, error) {
	condition := libovsdb.NewCondition("name", "==", name)
	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableLoadBalancer,
		Where: []interface{}{condition},
	}
	operations := []libovsdb.Operation{deleteOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) GetLB(name string) []*LoadBalancer {
	var lbList []*LoadBalancer
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()

	for uuid, drows := range odbi.cache[tableLoadBalancer] {
		if lbName, ok := drows.Fields["name"].(string); ok && lbName == name {
			lb := odbi.RowToLB(uuid)
			lbList = append(lbList, lb)
		}
	}
	return lbList
}

func (odbi *ovnDBImp) RowToLB(uuid string) *LoadBalancer {
	return &LoadBalancer{
		UUID:       uuid,
		protocol:   odbi.cache[tableLoadBalancer][uuid].Fields["protocol"].(string),
		Name:       odbi.cache[tableLoadBalancer][uuid].Fields["name"].(string),
		vips:       odbi.cache[tableLoadBalancer][uuid].Fields["vips"].(libovsdb.OvsMap).GoMap,
		ExternalID: odbi.cache[tableLoadBalancer][uuid].Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}
}
