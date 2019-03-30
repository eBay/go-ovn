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
	"github.com/ebay/libovsdb"
)

type OvnCommand struct {
	Operations []libovsdb.Operation
	Exe        Execution
	Results    [][]map[string]interface{}
}

type Execution interface {
	//Excute multi-commands
	Execute(cmds ...*OvnCommand) error
}

// North bound api set
type OVNDBApi interface {
	// Create ls named SWITCH
	LSAdd(ls string) (*OvnCommand, error)
	// Del ls and all its ports
	LSDel(ls string) (*OvnCommand, error)
	// Get all logical switches
	LSList() ([]*LogicalSwitch, error)

	// Add logical port PORT on SWITCH
	LSPAdd(ls string, lsp string) (*OvnCommand, error)
	// Delete PORT from its attached switch
	LSPDel(lsp string) (*OvnCommand, error)
	// Set addressset per lport
	LSPSetAddress(lsp string, addresses ...string) (*OvnCommand, error)
	// Set port security per lport
	LSPSetPortSecurity(lsp string, security ...string) (*OvnCommand, error)
	// Get all lport by lswitch
	LSPList(ls string) ([]*LogicalSwitchPort, error)

	// Add LB to LSW
	LSLBAdd(lswitch string, lb string) (*OvnCommand, error)
	// Delete LB from LSW
	LSLBDel(lswitch string, lb string) (*OvnCommand, error)
	// List Load balancers for a LSW
	LSLBList(lswitch string) ([]*LoadBalancer, error)

	// Add ACL
	ACLAdd(lsw, direct, match, action string, priority int, external_ids map[string]string, logflag bool, meter string) (*OvnCommand, error)
	// Delete acl
	ACLDel(lsw, direct, match string, priority int, external_ids map[string]string) (*OvnCommand, error)

	// Update address set
	ASUpdate(name string, addrs []string, external_ids map[string]string) (*OvnCommand, error)
	// Add addressset
	ASAdd(name string, addrs []string, external_ids map[string]string) (*OvnCommand, error)
	// Delete addressset
	ASDel(name string) (*OvnCommand, error)

	// Add LR with given name
	LRAdd(name string, external_ids map[string]string) (*OvnCommand, error)
	// Delete LR with given name
	LRDel(name string) (*OvnCommand, error)
	// Get LRs
	LRList() ([]*LogicalRouter, error)

	// Add LRP with given name on given lr
	LRPAdd(lr string, lrp string, mac string, network []string, peer string, external_ids map[string]string) (*OvnCommand, error)
	// Delete LRP with given name on given lr
	LRPDel(lr string, lrp string) (*OvnCommand, error)

	// Add LB to LR
	LRLBAdd(lr string, lb string) (*OvnCommand, error)
	// Delete LB from LR
	LRLBDel(lr string, lb string) (*OvnCommand, error)
	// List Load balancers for a LR
	LRLBList(lr string) ([]*LoadBalancer, error)

	// Add LB
	LBAdd(name string, vipPort string, protocol string, addrs []string) (*OvnCommand, error)
	// Delete LB with given name
	LBDel(name string) (*OvnCommand, error)
	// Update existing LB
	LBUpdate(name string, vipPort string, protocol string, addrs []string) (*OvnCommand, error)

	// Set dhcp4_options uuid on lsp
	LSPSetDHCPv4Options(lsp string, options string) (*OvnCommand, error)
	// Get dhcp4_options from lsp
	LSPGetDHCPv4Options(lsp string) (*DHCPOptions, error)
	// Set dhcp6_options uuid on lsp
	LSPSetDHCPv6Options(lsp string, options string) (*OvnCommand, error)
	// Get dhcp6_options from lsp
	LSPGetDHCPv6Options(lsp string) (*DHCPOptions, error)

	// Set options in LSP
	LSPSetOpt(lsp string, options map[string]string) (*OvnCommand, error)

	// Add dhcp options for cidr and provided external_ids
	DHCPOptionsAdd(cidr string, options map[string]string, external_ids map[string]string) (*OvnCommand, error)
	// Set dhcp options for specific cidr and provided external_ids
	DHCPOptionsSet(cidr string, options map[string]string, external_ids map[string]string) (*OvnCommand, error)
	// Del dhcp options via provided external_ids
	DHCPOptionsDel(uuid string) (*OvnCommand, error)
	// List dhcp options
	DHCPOptionsList() ([]*DHCPOptions, error)

	// Add qos rule
	QoSAdd(ls string, direction string, priority int, match string, action map[string]int, bandwidth map[string]int, external_ids map[string]string) (*OvnCommand, error)
	// Del qos rule, to delete wildcard specify priority -1 and string options as ""
	QoSDel(ls string, direction string, priority int, match string) (*OvnCommand, error)
	// Get qos rules by logical switch
	QoSList(ls string) ([]*QoS, error)

	// Get logical switch by name
	GetLogicalSwitchByName(ls string) (*LogicalSwitch, error)
	// Get all lrp by lr
	GetLogicalRouterPortsByRouter(lr string) ([]*LogicalRouterPort, error)

	// Get all acl by lswitch
	GetACLsBySwitch(lsw string) ([]*ACL, error)

	GetAddressSets() ([]*AddressSet, error)
	GetASByName(name string) (*AddressSet, error)
	// Get LB with given name
	GetLB(name string) ([]*LoadBalancer, error)
	// Get LR with given name
	GetLogicalRouter(name string) ([]*LogicalRouter, error)

	// Exec command, support mul-commands in one transaction.
	Execute(cmds ...*OvnCommand) error
	SetCallBack(callback OVNSignal)
}

type OVNSignal interface {
	OnLogicalSwitchCreate(ls *LogicalSwitch)
	OnLogicalSwitchDelete(ls *LogicalSwitch)

	OnLogicalPortCreate(lp *LogicalSwitchPort)
	OnLogicalPortDelete(lp *LogicalSwitchPort)

	OnLogicalRouterCreate(lr *LogicalRouter)
	OnLogicalRouterDelete(lr *LogicalRouter)

	OnLogicalRouterPortCreate(lrp *LogicalRouterPort)
	OnLogicalRouterPortDelete(lrp *LogicalRouterPort)

	OnACLCreate(acl *ACL)
	OnACLDelete(acl *ACL)

	OnDHCPOptionsCreate(dhcp *DHCPOptions)
	OnDHCPOptionsDelete(dhcp *DHCPOptions)

	OnQoSCreate(qos *QoS)
	OnQoSDelete(qos *QoS)

	OnLoadBalancerCreate(ls *LoadBalancer)
	OnLoadBalancerDelete(ls *LoadBalancer)
}

// Notifier
type OVNNotifier interface {
	Update(context interface{}, tableUpdates libovsdb.TableUpdates)
	Locked([]interface{})
	Stolen([]interface{})
	Echo([]interface{})
	Disconnected(client *libovsdb.OvsdbClient)
}

func (ocmd *OvnCommand) Execute() error {
	return ocmd.Exe.Execute()
}
