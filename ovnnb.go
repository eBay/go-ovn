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
	"errors"

	"github.com/ebay/libovsdb"
)

func newNBClient(socketfile string, proto string, server string, port int) (*ovnDBClient, error) {
	client := &ovnDBClient{
		socket:   socketfile,
		server:   server,
		port:     port,
		protocol: UNIX,
	}

	switch proto {
	case UNIX:
		clt, err := libovsdb.ConnectWithUnixSocket(socketfile)
		if err != nil {
			return nil, err
		}
		client.dbclient = clt
		return client, nil
	case TCP:
		clt, err := libovsdb.Connect(server, port, proto)

		if err != nil {
			return nil, err
		}
		client.dbclient = clt
		return client, nil
	case SSL:
		// for connection using SSL, make sure to set CLIENT_CERT_CA_CERT
		// and CLIENT_PRIVKEY in the env variable. CLIENT_CERT_CA_CERT is a
		// combination of client cert and ca cert appended in the same file.
		clt, err := libovsdb.Connect(server, port, proto)

		if err != nil {
			return nil, err
		}
		client.dbclient = clt
		return client, nil
	}
	return nil, errors.New("OVN DB initial failed: (unsupported protocol)")
}

func newNBBySocket(socketfile string, callback OVNSignal) (*OVNDB, error) {
	odb, err := newNBClient(socketfile, UNIX, "", 0)
	if err != nil {
		return nil, err
	}

	imp, err := newNBImp(odb, callback)
	if err != nil {
		return nil, err
	}

	return &OVNDB{imp}, nil
}

func newNBByServer(server string, port int, callback OVNSignal, protocol string) (*OVNDB, error) {
	odb, err := newNBClient("", protocol, server, port)
	if err != nil {
		return nil, err
	}

	imp, err := newNBImp(odb, callback)
	if err != nil {
		return nil, err
	}

	return &OVNDB{imp}, nil
}

func (odb *OVNDB) LSWAdd(lsw string) (*OvnCommand, error) {
	return odb.imp.lswAddImp(lsw)
}

func (odb *OVNDB) LSWDel(lsw string) (*OvnCommand, error) {
	return odb.imp.lswDelImp(lsw)
}

func (odb *OVNDB) LSWList() (*OvnCommand, error) {
	return odb.imp.lswListImp()
}

func (odb *OVNDB) LSPAdd(lsw string, lsp string) (*OvnCommand, error) {
	return odb.imp.lspAddImp(lsw, lsp)
}

func (odb *OVNDB) LSPDel(lsp string) (*OvnCommand, error) {
	return odb.imp.lspDelImp(lsp)
}

func (odb *OVNDB) LSPSetAddress(lsp string, addresses ...string) (*OvnCommand, error) {
	return odb.imp.lspSetAddressImp(lsp, addresses...)
}

func (odb *OVNDB) LSPSetPortSecurity(lsp string, security ...string) (*OvnCommand, error) {
	return odb.imp.lspSetPortSecurityImp(lsp, security...)
}

func (odb *OVNDB) LRAdd(name string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.lrAddImp(name, external_ids)
}

func (odb *OVNDB) LRDel(name string) (*OvnCommand, error) {
	return odb.imp.lrDelImp(name)
}

func (odb *OVNDB) LRPAdd(lr string, lrp string, mac string, network []string, peer string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.lrpAddImp(lr, lrp, mac, network, peer, external_ids)
}

func (odb *OVNDB) LRPDel(lr string, lrp string) (*OvnCommand, error) {
	return odb.imp.lrpDelImp(lr, lrp)
}

func (odb *OVNDB) LBAdd(name string, vipPort string, protocol string, addrs []string) (*OvnCommand, error) {
	return odb.imp.lbAddImp(name, vipPort, protocol, addrs)
}

func (odb *OVNDB) LBUpdate(name string, vipPort string, protocol string, addrs []string) (*OvnCommand, error) {
	return odb.imp.lbUpdateImp(name, vipPort, protocol, addrs)
}

func (odb *OVNDB) LBDel(name string) (*OvnCommand, error) {
	return odb.imp.lbDelImp(name)
}

func (odb *OVNDB) ACLAdd(lsw, direct, match, action string, priority int, external_ids map[string]string, logflag bool, meter string) (*OvnCommand, error) {
	return odb.imp.aclAddImp(lsw, direct, match, action, priority, external_ids, logflag, meter)
}

func (odb *OVNDB) ACLDel(lsw, direct, match string, priority int, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.aclDelImp(lsw, direct, match, priority, external_ids)
}

func (odb *OVNDB) ASAdd(name string, addrs []string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.ASAdd(name, addrs, external_ids)
}

func (odb *OVNDB) ASDel(name string) (*OvnCommand, error) {
	return odb.imp.ASDel(name)
}

func (odb *OVNDB) ASUpdate(name string, addrs []string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.ASUpdate(name, addrs, external_ids)
}

func (odb *OVNDB) LSPSetDHCPv4Options(lsp string, options string) (*OvnCommand, error) {
	return odb.imp.LSPSetDHCPv4Options(lsp, options)
}

func (odb *OVNDB) LSPGetDHCPv4Options(lsp string) (*DHCPOptions, error) {
	return odb.imp.LSPGetDHCPv4Options(lsp)
}

func (odb *OVNDB) LSPSetDHCPv6Options(lsp string, options string) (*OvnCommand, error) {
	return odb.imp.LSPSetDHCPv6Options(lsp, options)
}

func (odb *OVNDB) LSPGetDHCPv6Options(lsp string) (*DHCPOptions, error) {
	return odb.imp.LSPGetDHCPv6Options(lsp)
}

func (odb *OVNDB) LSPSetOpt(lsp string, options map[string]string) (*OvnCommand, error) {
	return odb.imp.LSPSetOpt(lsp, options)
}

func (odb *OVNDB) Execute(cmds ...*OvnCommand) error {
	return odb.imp.Execute(cmds...)
}

func (odb *OVNDB) GetLogicSwitches() []*LogicalSwitch {
	return odb.imp.GetLogicSwitches()
}

func (odb *OVNDB) GetLogicPortsBySwitch(lsw string) ([]*LogicalSwitchPort, error) {
	return odb.imp.GetLogicPortsBySwitch(lsw)
}

func (odb *OVNDB) GetLogicalRouterPortsByRouter(lr string) ([]*LogicalRouterPort, error) {
	return odb.imp.GetLogicalRouterPortsByRouter(lr)
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

func (odb *OVNDB) GetLogicalRouters() []*LogicalRouter {
	return odb.imp.GetLogicalRouters()
}

func (odb *OVNDB) GetLB(name string) []*LoadBalancer {
	return odb.imp.GetLB(name)
}

func (odb *OVNDB) AddDHCPOptions(cidr string, options map[string]string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.addDHCPOptionsImp(cidr, options, external_ids)
}

func (odb *OVNDB) SetDHCPOptions(cidr string, options map[string]string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.setDHCPOptionsImp(cidr, options, external_ids)
}

func (odb *OVNDB) DelDHCPOptions(uuid string) (*OvnCommand, error) {
	return odb.imp.delDHCPOptionsImp(uuid)
}

func (odb *OVNDB) GetDHCPOptions() []*DHCPOptions {
	return odb.imp.getDHCPOptionsImp()
}

func (odb *OVNDB) SetCallBack(callback OVNSignal) {
	odb.imp.callback = callback
}
