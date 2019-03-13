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
	// Create a logical switch named SWITCH
	LSWAdd(lsw string) (*OvnCommand, error)
	//delete SWITCH and all its ports
	LSWDel(lsw string) (*OvnCommand, error)
	// Print the names of all logical switches
	LSWList() (*OvnCommand, error)
	// Add logical port PORT on SWITCH
	LSPAdd(lsw, lsp string) (*OvnCommand, error)
	// Delete PORT from its attached switch
	LSPDel(lsp string) (*OvnCommand, error)
	// Set addressset per lport
	LSPSetAddress(lsp string, addresses ...string) (*OvnCommand, error)
	// Set port security per lport
	LSPSetPortSecurity(lsp string, security ...string) (*OvnCommand, error)
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
	// Add LB
	LBAdd(name string, vipPort string, protocol string, addrs []string) (*OvnCommand, error)
	// Delete LB with given name
	LBDel(name string) (*OvnCommand, error)
	// Update existing LB
	LBUpdate(name string, vipPort string, protocol string, addrs []string) (*OvnCommand, error)
	// Set options in lswtich
	LSSetOpt(lsp string, options map[string]string) (*OvnCommand, error)
	// Exec command, support mul-commands in one transaction.
	Execute(cmds ...*OvnCommand) error

	// Get all logical switches
	GetLogicSwitches() []*LogicalSwitch
	// Get all lport by lswitch
	GetLogicPortsBySwitch(lsw string) ([]*LogicalPort, error)
	// Get all acl by lswitch
	GetACLsBySwitch(lsw string) []*ACL

	GetAddressSets() []*AddressSet
	GetASByName(name string) *AddressSet
	// Get LB with given name
	GetLB(name string) []*LoadBalancer

	SetCallBack(callback OVNSignal)
}

type OVNSignal interface {
	OnLogicalSwitchCreate(ls *LogicalSwitch)
	OnLogicalSwitchDelete(ls *LogicalSwitch)

	OnLogicalPortCreate(lp *LogicalPort)
	OnLogicalPortDelete(lp *LogicalPort)

	OnACLCreate(acl *ACL)
	OnACLDelete(acl *ACL)
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

type LogicalSwitch struct {
	UUID       string
	Name       string
	ExternalID map[interface{}]interface{}
}

type LoadBalancer struct {
	UUID       string
	Name       string
	vips       map[interface{}]interface{}
	protocol   string
	ExternalID map[interface{}]interface{}
}

type LogicalPort struct {
	UUID         string
	Name         string
	Addresses    []string
	PortSecurity []string
}

type ACL struct {
	UUID       string
	Action     string
	Direction  string
	Match      string
	Priority   int
	Log        bool
	ExternalID map[interface{}]interface{}
}

type AddressSet struct {
	UUID       string
	Name       string
	Addresses  []string
	ExternalID map[interface{}]interface{}
}
