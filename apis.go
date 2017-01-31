package libovndb

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
	// Add ACL
	ACLAdd(lsw, direct, match, action string, priority int, external_ids *map[string]string, logflag bool) *OvnCommand
	// Delete acl
	ACLDel(lsw, direct, match string, priority int) *OvnCommand
	//add addresset
	ASAdd(name string, addrs []string) *OvnCommand
	// Delete addressset
	ASDel(name string) *OvnCommand
	// Set options in lswtich
	LSSetOpt(lsp string, options map[string]string) *OvnCommand
	// Exec command, support mul-commands in one transaction.
	Execute(cmds ...*OvnCommand) error

	// Get all lport by lswitch
	GetLogicPortsBySwitch(lsw string) []*LogcalPort
	// Get all acl by lswitch
	GetACLsBySwitch(lsw string) []*ACL

	GetAddressSets() []*AddressSet
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

type LogcalPort struct {
	UUID	  string
	Name      string
	Addresses []string
}

type ACL struct {
	UUID	  string
	Action    string
	Direction string
	Match     string
	Priority  int
}

type AddressSet struct {
	UUID	  string
	Name      string
	Addresses []string
}
