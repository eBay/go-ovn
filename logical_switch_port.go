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

type LogicalSwitchPort struct {
	UUID         string
	Name         string
	Addresses    []string
	PortSecurity []string
}

func (odbi *ovnDBImp) lspAddImp(lsw, lsp string) (*OvnCommand, error) {
	row := make(OVNRow)
	row["name"] = lsp

	insertOp, err := odbi.insertRowOp(tableLogicalSwitchPort, row)
	if err != nil {
		return nil, err
	}

	mutateUUID := []libovsdb.UUID{{insertOp.UUIDName}}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}

	mutation := libovsdb.NewMutation("ports", opInsert, mutateSet)
	condition := libovsdb.NewCondition("name", "==", lsw)
	mutateOp := odbi.mutateRowOp(tableLogicalSwitch, mutation, condition)
	operations := []libovsdb.Operation{insertOp, mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lspDelImp(lsp string) (*OvnCommand, error) {
	row := make(OVNRow)
	row["name"] = lsp

	lspUUID := odbi.getRowUUID(tableLogicalSwitchPort, row)
	if len(lspUUID) == 0 {
		return nil, ErrorNotFound
	}

	condition := libovsdb.NewCondition("name", "==", lsp)
	deleteOp := odbi.deleteRowOp(tableLogicalSwitchPort, condition)

	mutateUUID := []libovsdb.UUID{{lspUUID}}
	mutateSet, err := libovsdb.NewOvsSet(mutateUUID)
	if err != nil {
		return nil, err
	}

	mutation := libovsdb.NewMutation("ports", opDelete, mutateSet)
	portUUID := odbi.getRowUUIDContainsUUID(tableLogicalSwitch, "ports", lspUUID)
	if len(portUUID) == 0 {
		return nil, err
	}

	condition = libovsdb.NewCondition("_uuid", "==", libovsdb.UUID{portUUID})
	mutateOp := odbi.mutateRowOp(tableLogicalSwitch, mutation, condition)
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
	updateOp := odbi.updateRowOp(tableLogicalSwitchPort, row, condition)
	operations := []libovsdb.Operation{updateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lspSetPortSecurityImp(lsp string, security ...string) (*OvnCommand, error) {
	row := make(OVNRow)
	portsec, err := libovsdb.NewOvsSet(security)
	if err != nil {
		return nil, err
	}
	row["port_security"] = portsec
	condition := libovsdb.NewCondition("name", "==", lsp)
	updateOp := odbi.updateRowOp(tableLogicalSwitchPort, row, condition)
	operations := []libovsdb.Operation{updateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) LSPSetOpt(lsp string, options map[string]string) (*OvnCommand, error) {
	mutatemap, _ := libovsdb.NewOvsMap(options)
	mutation := libovsdb.NewMutation("options", opInsert, mutatemap)
	condition := libovsdb.NewCondition("name", "==", lsp)
	mutateOp := odbi.mutateRowOp(tableLogicalSwitchPort, mutation, condition)
	operations := []libovsdb.Operation{mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) RowToLogicalPort(uuid string) *LogicalSwitchPort {
	odbi.cachemutex.RLock()
	row := odbi.cache[tableLogicalSwitchPort][uuid]
	odbi.cachemutex.RUnlock()

	lp := &LogicalSwitchPort{
		UUID: uuid,
		Name: row.Fields["name"].(string),
	}
	addr := row.Fields["addresses"]
	portsec := row.Fields["port_security"]

	switch addr.(type) {
	case string:
		lp.Addresses = []string{addr.(string)}
	case libovsdb.OvsSet:
		lp.Addresses = odbi.ConvertGoSetToStringArray(addr.(libovsdb.OvsSet))
	default:
		//	glog.V(OVNLOGLEVEL).Info("Unsupport type found in lport address.")
	}
	switch portsec.(type) {
	case string:
		lp.PortSecurity = []string{portsec.(string)}
	case libovsdb.OvsSet:
		lp.PortSecurity = odbi.ConvertGoSetToStringArray(portsec.(libovsdb.OvsSet))
	default:
		//glog.V(OVNLOGLEVEL).Info("Unsupport type found in lport port security.")
	}
	return lp
}

// Get all lport by lswitch
func (odbi *ovnDBImp) GetLogicalPortsBySwitch(lsw string) ([]*LogicalSwitchPort, error) {
	var lplist = []*LogicalSwitchPort{}
	odbi.cachemutex.RLock()
	rows, ok := odbi.cache[tableLogicalSwitch]
	odbi.cachemutex.RUnlock()
	if !ok {
		return nil, ErrorNotFound
	}

	for _, drows := range rows {
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
