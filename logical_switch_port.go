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

type LogicalSwitchPort struct {
	UUID          string
	Name          string
	Type          string
	Options       map[interface{}]interface{}
	Addresses     []string
	PortSecurity  []string
	DHCPv4Options string
	DHCPv6Options string
	ExternalID    map[interface{}]interface{}
}

func (odbi *ovnDBImp) lspAddImp(lsw, lsp string) (*OvnCommand, error) {
	namedUUID, err := newRowUUID()
	if err != nil {
		return nil, err
	}
	row := make(OVNRow)
	row["name"] = lsp

	if uuid := odbi.getRowUUID(tableLogicalSwitchPort, row); len(uuid) > 0 {
		return nil, ErrorExist
	}

	insertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableLogicalSwitchPort,
		Row:      row,
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
	row := make(OVNRow)
	row["name"] = lsp

	lspUUID := odbi.getRowUUID(tableLogicalSwitchPort, row)
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

func (odbi *ovnDBImp) LSPSetDHCPv4Options(lsp string, uuid string) (*OvnCommand, error) {
	row := make(OVNRow)
	row["dhcpv4_options"] = libovsdb.UUID{uuid}
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

func (odbi *ovnDBImp) LSPGetDHCPv4Options(lsp string) (*DHCPOptions, error) {
	lp, err := odbi.GetLogicalSwitchPortByName(lsp)
	if err != nil {
		return nil, err
	}
	return odbi.rowToDHCPOptions(lp.DHCPv4Options), nil
}

func (odbi *ovnDBImp) LSPSetDHCPv6Options(lsp string, options string) (*OvnCommand, error) {
	mutation := libovsdb.NewMutation("dhcpv6_options", opInsert, options)
	condition := libovsdb.NewCondition("name", "==", lsp)
	mutateOp := libovsdb.Operation{
		Op:        opMutate,
		Table:     tableLogicalSwitchPort,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) LSPGetDHCPv6Options(lsp string) (*DHCPOptions, error) {
	lp, err := odbi.GetLogicalSwitchPortByName(lsp)
	if err != nil {
		return nil, err
	}
	return odbi.rowToDHCPOptions(lp.DHCPv6Options), nil
}

func (odbi *ovnDBImp) LSPSetOpt(lsp string, options map[string]string) (*OvnCommand, error) {
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

func (odbi *ovnDBImp) rowToLogicalPort(uuid string) *LogicalSwitchPort {
	lp := &LogicalSwitchPort{
		UUID:       uuid,
		Name:       odbi.cache[tableLogicalSwitchPort][uuid].Fields["name"].(string),
		Type:       odbi.cache[tableLogicalSwitchPort][uuid].Fields["type"].(string),
		ExternalID: odbi.cache[tableLogicalSwitchPort][uuid].Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}

	if dhcpv4, ok := odbi.cache[tableLogicalSwitchPort][uuid].Fields["dhcpv4_options"]; ok {
		switch dhcpv4.(type) {
		case libovsdb.UUID:
			lp.DHCPv4Options = dhcpv4.(libovsdb.UUID).GoUUID
		case libovsdb.OvsSet:
		default:
		}
	}
	if dhcpv6, ok := odbi.cache[tableLogicalSwitchPort][uuid].Fields["dhcpv6_options"]; ok {
		switch dhcpv6.(type) {
		case libovsdb.UUID:
			lp.DHCPv6Options = dhcpv6.(libovsdb.UUID).GoUUID
		case libovsdb.OvsSet:
		default:
		}
	}

	if addr, ok := odbi.cache[tableLogicalSwitchPort][uuid].Fields["addresses"]; ok {
		switch addr.(type) {
		case string:
			lp.Addresses = []string{addr.(string)}
		case libovsdb.OvsSet:
			lp.Addresses = odbi.ConvertGoSetToStringArray(addr.(libovsdb.OvsSet))
		default:
			//	glog.V(OVNLOGLEVEL).Info("Unsupport type found in lport address.")
		}
	}

	if portsecurity, ok := odbi.cache[tableLogicalSwitchPort][uuid].Fields["port_security"]; ok {
		switch portsecurity.(type) {
		case string:
			lp.PortSecurity = []string{portsecurity.(string)}
		case libovsdb.OvsSet:
			lp.PortSecurity = odbi.ConvertGoSetToStringArray(portsecurity.(libovsdb.OvsSet))
		default:
			//glog.V(OVNLOGLEVEL).Info("Unsupport type found in lport port security.")
		}
	}

	if options, ok := odbi.cache[tableLogicalSwitchPort][uuid].Fields["options"]; ok {
		lp.Options = options.(libovsdb.OvsMap).GoMap
	}

	return lp
}

// Get lsp by name
func (odbi *ovnDBImp) GetLogicalSwitchPortByName(lsp string) (*LogicalSwitchPort, error) {
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheLogicalSwitchPort, ok := odbi.cache[tableLogicalSwitchPort]
	if !ok {
		return nil, ErrorSchema
	}

	for uuid, drows := range cacheLogicalSwitchPort {
		if rlsp, ok := drows.Fields["name"].(string); ok && rlsp == lsp {
			return odbi.rowToLogicalPort(uuid), nil
		}
	}
	return nil, ErrorNotFound
}

// Get all lport by lswitch
func (odbi *ovnDBImp) GetLogicalSwitchPortsBySwitch(lsw string) ([]*LogicalSwitchPort, error) {
	var listLSP []*LogicalSwitchPort

	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheLogicalSwitch, ok := odbi.cache[tableLogicalSwitch]
	if !ok {
		return nil, ErrorSchema
	}

	for _, drows := range cacheLogicalSwitch {
		if rlsw, ok := drows.Fields["name"].(string); ok && rlsw == lsw {
			ports := drows.Fields["ports"]
			if ports != nil {
				switch ports.(type) {
				case libovsdb.OvsSet:
					if ps, ok := ports.(libovsdb.OvsSet); ok {
						for _, p := range ps.GoSet {
							if vp, ok := p.(libovsdb.UUID); ok {
								tp := odbi.rowToLogicalPort(vp.GoUUID)
								listLSP = append(listLSP, tp)
							}
						}
					} else {
						return nil, fmt.Errorf("type libovsdb.OvsSet casting failed")
					}
				case libovsdb.UUID:
					if vp, ok := ports.(libovsdb.UUID); ok {
						tp := odbi.rowToLogicalPort(vp.GoUUID)
						listLSP = append(listLSP, tp)
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
	return listLSP, nil
}
