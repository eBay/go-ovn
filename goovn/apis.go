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
	LSWAdd(lsw string) *OvnCommand
	//delete SWITCH and all its ports
	LSWDel(lsw string) *OvnCommand
	// Print the names of all logical switches
	LSWList() *OvnCommand
	// Add logical port PORT on SWITCH
	LSPAdd(lsw, lsp string) *OvnCommand
	// Delete PORT from its attached switch
	LSPDel(lsp string) *OvnCommand
	// Set addressset per lport
	LSPSetAddress(lsp string, addresses ...string) *OvnCommand
	// Set port security per lport
	LSPSetPortSecurity(lsp string, security ...string) *OvnCommand
	// Add ACL
	ACLAdd(lsw, direct, match, action string, priority int, external_ids map[string]string, logflag bool, meter string) *OvnCommand
	// Delete acl
	ACLDel(lsw, direct, match string, priority int, external_ids map[string]string) *OvnCommand
	// Update address set
	ASUpdate(name string, addrs []string, external_ids map[string]string) *OvnCommand
	// Add addressset
	ASAdd(name string, addrs []string, external_ids map[string]string) *OvnCommand
	// Delete addressset
	ASDel(name string) *OvnCommand
	// Set options in lswtich
	LSSetOpt(lsp string, options map[string]string) *OvnCommand
	// Exec command, support mul-commands in one transaction.
	Execute(cmds ...*OvnCommand) error

	// Get all logical switches
	GetLogicSwitches() []*LogicalSwitch
	// Get all lport by lswitch
	GetLogicPortsBySwitch(lsw string) []*LogcalPort
	// Get all acl by lswitch
	GetACLsBySwitch(lsw string) []*ACL

	GetAddressSets() []*AddressSet
	GetASByName(name string) *AddressSet

	SetCallBack(callback OVNSignal)
}

type OVNSignal interface {
	OnLogicalSwitchCreate(ls *LogicalSwitch)
	OnLogicalSwitchDelete(ls *LogicalSwitch)

	OnLogicalPortCreate(lp *LogcalPort)
	OnLogicalPortDelete(lp *LogcalPort)

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

const (
	OVNLOGLEVEL = 4
)

type LogicalSwitch struct {
	UUID       string
	Name       string
	ExternalID map[interface{}]interface{}
}

type LogcalPort struct {
	UUID	  string
	Name      string
	Addresses []string
	PortSecurity []string
}

type ACL struct {
	UUID	  string
	Action    string
	Direction string
	Match     string
	Priority  int
	Log       bool
	ExternalID map[interface{}]interface{}
}

type AddressSet struct {
	UUID	  string
	Name      string
	Addresses []string
	ExternalID map[interface{}]interface{}
}
