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

	"github.com/socketplane/libovsdb"
)

func (odbi *ovnDBImp) lspAddImp(lsw, lsp string) (*OvnCommand, error) {
	namedUUID, err := newUUID()
	if err != nil {
		return nil, err
	}
	lsprow := make(OVNRow)
	lsprow["name"] = lsp

	if uuid := odbi.getRowUUID(tableLogicalSwitchPort, lsprow); len(uuid) > 0 {
		return nil, ErrorExist
	}

	insertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableLogicalSwitchPort,
		Row:      lsprow,
		UUIDName: namedUUID,
	}

	mutateUUID := []libovsdb.UUID{{namedUUID}}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}

	mutation := libovsdb.NewMutation("ports", opInsert, mutateSet)
	condition := libovsdb.NewCondition("name", "==", lsw)

	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{insertOp, mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lspDelImp(lsp string) (*OvnCommand, error) {
	lsprow := make(OVNRow)
	lsprow["name"] = lsp

	lspUUID := odbi.getRowUUID(tableLogicalSwitchPort, lsprow)
	if len(lspUUID) == 0 {
		return nil, ErrorNotFound
	}

	mutateUUID := []libovsdb.UUID{{lspUUID}}
	condition := libovsdb.NewCondition("name", "==", lsp)
	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableLogicalSwitchPort,
		Where: []interface{}{condition},
	}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}
	mutation := libovsdb.NewMutation("ports", opDelete, mutateSet)
	ucondition, err := odbi.getRowUUIDContainsUUID(tableLogicalSwitch, "ports", lspUUID)
	if err != nil {
		return nil, err
	}

	mucondition := libovsdb.NewCondition("_uuid", "==", libovsdb.UUID{ucondition})
	// simple mutate operation
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitch,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{mucondition},
	}
	operations := []libovsdb.Operation{deleteOp, mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lspSetAddressImp(lsp string, addr ...string) (*OvnCommand, error) {
	row := make(OVNRow)
	addresses, err := libovsdb.NewOvsSet(addr)
	if err != nil {
		return nil, err
	}
	row["addresses"] = addresses
	condition := libovsdb.NewCondition("name", "==", lsp)
	updateOp := libovsdb.Operation{
		Op:    opUpdate,
		Table: tableLogicalSwitchPort,
		Row:   row,
		Where: []interface{}{condition},
	}
	operations := []libovsdb.Operation{updateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lspSetPortSecurityImp(lsp string, security ...string) (*OvnCommand, error) {
	row := make(OVNRow)
	port_security, err := libovsdb.NewOvsSet(security)
	if err != nil {
		return nil, err
	}
	row["port_security"] = port_security
	condition := libovsdb.NewCondition("name", "==", lsp)
	updateOp := libovsdb.Operation{
		Op:    opUpdate,
		Table: tableLogicalSwitchPort,
		Row:   row,
		Where: []interface{}{condition},
	}
	operations := []libovsdb.Operation{updateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) LSSetOpt(lsp string, options map[string]string) (*OvnCommand, error) {
	mutatemap, _ := libovsdb.NewOvsMap(options)
	mutation := libovsdb.NewMutation("options", opInsert, mutatemap)
	condition := libovsdb.NewCondition("name", "==", lsp)

	// simple mutate operation
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitchPort,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) RowToLogicalPort(uuid string) *LogicalPort {
	lp := &LogicalPort{
		UUID: uuid,
		Name: odbi.cache[tableLogicalSwitchPort][uuid].Fields["name"].(string),
	}
	addr := odbi.cache[tableLogicalSwitchPort][uuid].Fields["addresses"]
	switch addr.(type) {
	case string:
		lp.Addresses = []string{addr.(string)}
	case libovsdb.OvsSet:
		lp.Addresses = odbi.ConvertGoSetToStringArray(addr.(libovsdb.OvsSet))
	default:
		//	glog.V(OVNLOGLEVEL).Info("Unsupport type found in lport address.")
	}
	portsecurity := odbi.cache[tableLogicalSwitchPort][uuid].Fields["port_security"]
	switch portsecurity.(type) {
	case string:
		lp.PortSecurity = []string{portsecurity.(string)}
	case libovsdb.OvsSet:
		lp.PortSecurity = odbi.ConvertGoSetToStringArray(portsecurity.(libovsdb.OvsSet))
	default:
		//glog.V(OVNLOGLEVEL).Info("Unsupport type found in lport port security.")
	}
	return lp
}

// Get all lport by lswitch
func (odbi *ovnDBImp) GetLogicPortsBySwitch(lsw string) ([]*LogicalPort, error) {
	var lplist = []*LogicalPort{}
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	for _, drows := range odbi.cache[tableLogicalSwitch] {
		if rlsw, ok := drows.Fields["name"].(string); ok && rlsw == lsw {
			ports := drows.Fields["ports"]
			if ports != nil {
				switch ports.(type) {
				case libovsdb.OvsSet:
					if ps, ok := ports.(libovsdb.OvsSet); ok {
						for _, p := range ps.GoSet {
							if vp, ok := p.(libovsdb.UUID); ok {
								tp := odbi.RowToLogicalPort(vp.GoUUID)
								lplist = append(lplist, tp)
							}
						}
					} else {
						return nil, fmt.Errorf("type libovsdb.OvsSet casting failed")
					}
				case libovsdb.UUID:
					if vp, ok := ports.(libovsdb.UUID); ok {
						tp := odbi.RowToLogicalPort(vp.GoUUID)
						lplist = append(lplist, tp)
					} else {
						return nil, fmt.Errorf("type libovsdb.UUID casting failed")
					}
				default:
					return nil, fmt.Errorf("Unsupport type found in ovsdb rows")
				}
			}
			break
		}
	}
	return lplist, nil
}
