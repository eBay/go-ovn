package libovndb

import (
	"errors"
	"github.com/golang/glog"

	"github.com/socketplane/libovsdb"
)

func newOvnDbClient(socketfile string, protocol string, server string, port int) (*ovnDBClient, error) {
	client := &ovnDBClient{
		socket:   socketfile,
		server:   server,
		port:     port,
		protocol: UNIX,
	}

	if protocol == UNIX {
		clt, err := libovsdb.ConnectWithUnixSocket(socketfile)
		if err != nil {
			glog.Fatalf("OVN DB initial failed: (%v) with socket file %s.", err, socketfile)
			return nil, err
		}
		client.dbclient = clt
		return client, nil

	} else if protocol == TCP {
		clt, err := libovsdb.Connect(server, port)
		if err != nil {
			glog.Fatalf("OVN DB initial failed: (%v) on %s:%d", err, server, port)
			return nil, err
		}
		client.dbclient = clt
		return client, nil
	}
	return nil, errors.New("OVN DB initial failed: (unsupported protocol)")
}

func newNBCtlBySocket(socketfile string, callback OVNSignal) (*OVNDB, error) {
	odb, err := newOvnDbClient(socketfile, UNIX, "", 0)
	if err == nil {
		return &OVNDB{newNBCtlImp(odb, callback)}, nil
	} else {
		return nil, err
	}
}

func newNBCtlByServer(server string, port int, callback OVNSignal) (*OVNDB, error) {
	odb, err := newOvnDbClient("", TCP, server, port)
	if err != nil {
		return &OVNDB{newNBCtlImp(odb, callback)}, nil
	} else {
		return nil, err
	}
}

func (odb *OVNDB) LSWAdd(lsw string) *OvnCommand {
	return odb.imp.lswAddImp(lsw)
}

func (odb *OVNDB) LSWDel(lsw string) *OvnCommand {
	return odb.imp.lswDelImp(lsw)
}

func (odb *OVNDB) LSWList() *OvnCommand {
	return odb.imp.lswListImp()
}

func (odb *OVNDB) LSPAdd(lsw string, lsp string) *OvnCommand {
	return odb.imp.lspAddImp(lsw, lsp)
}

func (odb *OVNDB) LSPDel(lsp string) *OvnCommand {
	return odb.imp.lspDelImp(lsp)
}


func (odb *OVNDB) LSPSetAddress(lsp string, addresses ...string) *OvnCommand {
	return odb.imp.lspSetAddressImp(lsp, addresses...)
}

func (odb *OVNDB) LSPSetPortSecurity(lsp string, security ...string) *OvnCommand {
	return odb.imp.lspSetPortSecurityImp(lsp, security...)
}

func (odb *OVNDB) ACLAdd(lsw, direct, match, action string, priority int, external_ids map[string]string, logflag bool) *OvnCommand {
	return odb.imp.aclAddImp(lsw, direct, match, action, priority, external_ids, logflag)
}

func (odb *OVNDB) ACLDel(lsw, direct, match string, priority int) *OvnCommand {
	return odb.imp.aclDelImp(lsw, direct, match, priority)
}

func (odb *OVNDB) ASAdd(name string, addrs []string) *OvnCommand {
	return odb.imp.ASAdd(name, addrs)
}

func (odb *OVNDB) ASDel(name string) *OvnCommand {
	return odb.imp.ASDel(name)
}

func (odb *OVNDB) ASUpdate(name string, addrs []string) *OvnCommand {
	return odb.imp.ASUpdate(name, addrs)
}

func (odb *OVNDB) LSSetOpt(lsp string, options map[string]string) *OvnCommand {
	return odb.imp.LSSetOpt(lsp, options)
}

func (odb *OVNDB) Execute(cmds ...*OvnCommand) error {
	return odb.imp.Execute(cmds...)
}

func (odb *OVNDB) GetLogicPortsBySwitch(lsw string) []*LogcalPort {
	return odb.imp.GetLogicPortsBySwitch(lsw)
}

func (odb *OVNDB) GetACLsBySwitch(lsw string) []*ACL {
	return odb.imp.GetACLsBySwitch(lsw)
}

func (odb *OVNDB) GetAddressSets() []*AddressSet {
	return odb.imp.GetAddressSets()
}

func (odb *OVNDB) GetASByName(name string) *AddressSet {
	return odb.imp.GetASByName(name)
}

func (odb *OVNDB) SetCallBack(callback OVNSignal) {
	odb.imp.callback = callback
}
