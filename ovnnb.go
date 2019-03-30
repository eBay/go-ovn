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

func (odb *OVNDB) LSAdd(lsw string) (*OvnCommand, error) {
	return odb.imp.lsAddImp(lsw)
}

func (odb *OVNDB) LSDel(lsw string) (*OvnCommand, error) {
	return odb.imp.lsDelImp(lsw)
}

func (odb *OVNDB) LSList() ([]*LogicalSwitch, error) {
	return odb.imp.lsListImp()
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

func (odb *OVNDB) LSLBAdd(lswitch string, lb string) (*OvnCommand, error) {
	return odb.imp.lslbAddImp(lswitch, lb)
}

func (odb *OVNDB) LSLBDel(lswitch string, lb string) (*OvnCommand, error) {
	return odb.imp.lslbDelImp(lswitch, lb)
}

func (odb *OVNDB) LSLBList(lswitch string) ([]*LoadBalancer, error) {
	return odb.imp.lslblistImp(lswitch)
}

func (odb *OVNDB) LRAdd(name string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.lrAddImp(name, external_ids)
}

func (odb *OVNDB) LRDel(name string) (*OvnCommand, error) {
	return odb.imp.lrDelImp(name)
}

func (odb *OVNDB) LRList() ([]*LogicalRouter, error) {
	return odb.imp.lrListImp()
}

func (odb *OVNDB) LRPAdd(lr string, lrp string, mac string, network []string, peer string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.lrpAddImp(lr, lrp, mac, network, peer, external_ids)
}

func (odb *OVNDB) LRPDel(lr string, lrp string) (*OvnCommand, error) {
	return odb.imp.lrpDelImp(lr, lrp)
}

func (odb *OVNDB) LRLBAdd(lr string, lb string) (*OvnCommand, error) {
	return odb.imp.lrlbAddImp(lr, lb)
}

func (odb *OVNDB) LRLBDel(lr string, lb string) (*OvnCommand, error) {
	return odb.imp.lrlbDelImp(lr, lb)
}

func (odb *OVNDB) LRLBList(lr string) ([]*LoadBalancer, error) {
	return odb.imp.lrlblistImp(lr)
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

func (odb *OVNDB) QoSAdd(ls string, direction string, priority int, match string, action map[string]int, bandwidth map[string]int, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.addQoSImp(ls, direction, priority, match, action, bandwidth, external_ids)
}

func (odb *OVNDB) QoSDel(ls string, direction string, priority int, match string) (*OvnCommand, error) {
	return odb.imp.delQoSImp(ls, direction, priority, match)
}

func (odb *OVNDB) QoSList(ls string) ([]*QoS, error) {
	return odb.imp.listQoSImp(ls)
}

func (odb *OVNDB) Execute(cmds ...*OvnCommand) error {
	return odb.imp.Execute(cmds...)
}

func (odb *OVNDB) GetLogicalSwitchByName(ls string) (*LogicalSwitch, error) {
	return odb.imp.GetLogicalSwitchByName(ls)
}

func (odb *OVNDB) LSPList(lsw string) ([]*LogicalSwitchPort, error) {
	return odb.imp.LSPList(lsw)
}

func (odb *OVNDB) GetLogicalRouterPortsByRouter(lr string) ([]*LogicalRouterPort, error) {
	return odb.imp.GetLogicalRouterPortsByRouter(lr)
}

func (odb *OVNDB) GetACLsBySwitch(lsw string) ([]*ACL, error) {
	return odb.imp.GetACLsBySwitch(lsw)
}

func (odb *OVNDB) GetAddressSets() ([]*AddressSet, error) {
	return odb.imp.GetAddressSets()
}

func (odb *OVNDB) GetASByName(name string) (*AddressSet, error) {
	return odb.imp.GetASByName(name)
}

func (odb *OVNDB) GetLogicalRouter(name string) ([]*LogicalRouter, error) {
	return odb.imp.GetLogicalRouter(name)
}

func (odb *OVNDB) GetLB(name string) ([]*LoadBalancer, error) {
	return odb.imp.GetLB(name)
}

func (odb *OVNDB) DHCPOptionsAdd(cidr string, options map[string]string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.dhcpOptionsAddImp(cidr, options, external_ids)
}

func (odb *OVNDB) DHCPOptionsSet(cidr string, options map[string]string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.dhcpOptionsSetImp(cidr, options, external_ids)
}

func (odb *OVNDB) DHCPOptionsDel(uuid string) (*OvnCommand, error) {
	return odb.imp.dhcpOptionsDelImp(uuid)
}

func (odb *OVNDB) DHCPOptionsList() ([]*DHCPOptions, error) {
	return odb.imp.dhcpOptionsListImp()
}

func (odb *OVNDB) SetCallBack(callback OVNSignal) {
	odb.imp.callback = callback
}
