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
	"fmt"
	"github.com/ebay/libovsdb"
)

type NAT struct {
	UUID        string
	Type        string
	ExternalIP  string
	ExternalMAC string
	LogicalIP   string
	LogicalPort string
	ExternalID  map[interface{}]interface{}
}

func (odbi *ovndb) rowToNat(uuid string) *NAT{
	cacheNAT , ok := odbi.cache[tableNAT][uuid]
	if !ok{
		fmt.Println("error :",uuid)
		return nil
	}
	nat := &NAT{
		UUID:	uuid,
		Type:   cacheNAT.Fields["type"].(string),
		ExternalIP: cacheNAT.Fields["external_ip"].(string),
		LogicalIP:   cacheNAT.Fields["logical_ip"].(string),
		ExternalID:	cacheNAT.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}
	if mac ,ok := cacheNAT.Fields["external_mac"];ok{
		switch mac.(type){
		case libovsdb.UUID:
			nat.ExternalMAC = mac.(libovsdb.UUID).GoUUID
		case string:
			nat.ExternalMAC = mac.(string)
		default:
		}
	}
	if lip ,ok:= cacheNAT.Fields["logical_port"];ok{
		switch lip.(type) {
		case libovsdb.UUID:
			nat.LogicalIP = lip.(libovsdb.UUID).GoUUID
		case string:
			nat.LogicalIP = lip.(string)
		default:
		}

	}

	return nat
}



func (odbi *ovndb) natAddImp(lr string,Type string, externalIp string,externalMac string,logicalIp string,logicalPort string,external_ids map[string]string)(*OvnCommand, error){
	nameUUID,err := newRowUUID()
	if err != nil {
		return nil, err
	}
	row :=make(OVNRow)

	row["external_ip"] = externalIp

	row["logical_ip"] = logicalIp

	switch Type {
	case "snat":
		row["type"] = Type
	case "dnat":
		row["type"] = Type
	case "dnat_and_snat":
		row["type"] = Type
	default:
		return nil,ErrorSchema
	}

	if uuid := odbi.getRowUUID(tableNAT, row); len(uuid) > 0 {
		return nil, ErrorExist
	}
	// The logical_port and  external_mac  are  only  accepted
	//when  router  is  a  distributed  router  (rather than a gateway
	//router) and type is dnat_and_snat.
	if externalMac != ""{
		if row["type"] != "dnat_and_snat"{
			return nil,ErrorSchema
		}
		row["external_mac"] = externalMac
	}

	row["logical_port"] = logicalPort

	if external_ids != nil {
		oMap, err := libovsdb.NewOvsMap(external_ids)
		if err != nil {
			return nil, err
		}
		row["external_ids"] = oMap
	}
	insertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableNAT,
		Row:      row,
		UUIDName: nameUUID,
	}
	mutateUUID := []libovsdb.UUID{{nameUUID}}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("nat", opInsert, mutateSet)
	condition := libovsdb.NewCondition("name", "==", lr)
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalRouter,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{insertOp, mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}


//Deletes  NATs  from  router. If only router is supplied, all the
//NATs from the logical router are deleted. If type is also speci‐
//fied, then all the NATs that match the type will be deleted from
//the logical router. If all the fields are given, then  a  single
//NAT  rule that matches all the fields will be deleted. When type
//	is snat, the ip should be  logical_ip.  When  type  is  dnat  or
//dnat_and_snat, the ip shoud be external_ip.
func (odbi *ovndb) natDelImp(lr string,Type string,ip ...string)(*OvnCommand, error){
	var operations []libovsdb.Operation
	row := make(OVNRow)
	var lrNatUUID string
	if Type != ""{
		switch Type {
		case "snat":
			row["type"] = Type
			if len(ip) != 0{
				row["logical_ip"] = ip[0]
			}
		case "dnat":
			row["type"] = Type
			if len(ip) != 0{
			row["external_ip"] = ip[0]
		}
		case "dnat_and_snat":
			row["type"] = Type
			if len(ip) != 0{
				row["external_ip"] = ip[0]
			}
		default:
			return nil,ErrorSchema
		}
		lrNatUUID = odbi.getRowUUID(tableNAT,row)
		if len(lrNatUUID) == 0 {
			return nil, ErrorNotFound
		}
	}
	var natlist []string
	LR ,_:= odbi.LRGet(lr)
	if len(LR) == 0 {
		return nil, ErrorNotFound
	}
	for _,v := range LR  {
		natlist = v.NAT
	}

	var mutateUUID []libovsdb.UUID
	if lrNatUUID != ""{
		for _, i := range natlist{
			if lrNatUUID == i{
				mutateUUID = append(mutateUUID,libovsdb.UUID{GoUUID: lrNatUUID})
			}
		}
	}else{
		for _, i := range natlist{
			mutateUUID = append(mutateUUID,libovsdb.UUID{GoUUID: i})
			}
	}

	mutateSet ,err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}

	lrNatUUID = odbi.getRowUUID(tableNAT,row)
	if len(lrNatUUID) == 0 {
		return nil, ErrorNotFound
	}
	row = make(OVNRow)
	row["name"] = lr
	mutation := libovsdb.NewMutation("nat", opDelete, mutateSet)
	mucondition := libovsdb.NewCondition("name", "==", lr)
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalRouter,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{mucondition},
	}
	operations = append(operations, mutateOp)
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil

}

func (odbi *ovndb) natListImp(lr string) ([]*NAT,error){
	var NATList []*NAT
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	LR ,_:= odbi.LRGet(lr)
	if len(LR) == 0 {
		return nil, ErrorNotFound
	}
	for _,v :=range LR[0].NAT{
		NATList = append(NATList,odbi.rowToNat(v))
		return NATList,nil
	}

	return NATList,nil
}

