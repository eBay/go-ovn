/**
 * Copyright (c) 2019 eBay Inc.
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

// NBGlobal row
type NBGlobal struct {
	UUID        string            `ovn:"uuid"`
	NBCfg       int               `ovn:"nb_cfg"`
	SBCfg       int               `ovn:"sb_cfg"`
	HVCfg       int               `ovn:"hv_cfg"`
	ExternalIDs map[string]string `ovn:"external_ids"`
	Options     map[string]string `ovn:"options"`
	Connections []string          `ovn:"connections"`
	SSL         string            `ovn:"ssl"`
	IPSec       bool              `ovn:"ipsec"`
}

// LogicalSwitch row
type LogicalSwitch struct {
	UUID         string            `ovn:"uuid"`
	Name         string            `ovn:"name"`
	Ports        []string          `ovn:"ports"`
	LoadBalancer []string          `ovn:"load_balancer"`
	ACLs         []string          `ovn:"acls"`
	QoSRules     []string          `ovn:"qos_rules"`
	DNSRecords   []string          `ovn:"dns_records"`
	OtherConfig  map[string]string `ovn:"other_config"`
	ExternalIDs  map[string]string `ovn:"external_ids"`
}

// LogicalSwitchPort row
type LogicalSwitchPort struct {
	UUID             string            `ovn:"uuid"`
	Name             string            `ovn:"name"`
	Type             string            `ovn:"type"`
	Options          map[string]string `ovn:options"`
	ParentName       string            `ovn:"parent_name"`
	TagRequest       int               `ovn:"tag_request"`
	Tag              int               `ovn:"tag"`
	Up               bool              `ovn:"up"`
	Enabled          bool              `ovn:"enabled"`
	Addresses        []string          `ovn:"addresses"`
	DynamicAddresses string            `ovn:"dynamic_addresses"`
	PortSecurity     []string          `ovn:"port_security"`
	DHCPv4Options    string            `ovn:"dhcpv4_options"`
	DHCPv6Options    string            `ovn:"dhcpv6_options"`
	ExternalIDs      map[string]string `ovn:"external_ids"`
}

// AddressSet row
type AddressSet struct {
	UUID        string            `ovn:"uuid"`
	Name        string            `ovn:"name"`
	Addresses   []string          `ovn:"addresses"`
	ExternalIDs map[string]string `ovn:"external_ids"`
}

// PortGroup row
type PortGroup struct {
	UUID        string            `ovn:"uuid"`
	Name        string            `ovn:"name"`
	Ports       []string          `ovn:"ports"`
	ACLs        []string          `ovns:"acls"`
	ExternalIDs map[string]string `ovn:"external_ids"`
}

// LoadBalancer row
type LoadBalancer struct {
	UUID        string            `ovn:"uuid"`
	Name        string            `ovn:"name"`
	VIPs        map[string]string `ovn:"vips"`
	Protocol    string            `ovn:"protocol"`
	ExternalIDs map[string]string `ovn:"external_ids"`
}

// ACL row
type ACL struct {
	UUID        string            `ovn:"uuid"`
	Action      string            `ovn:"action"`
	Direction   string            `ovn:"direction"`
	Match       string            `ovn:"match"`
	Priority    int               `ovn:"priority"`
	Log         bool              `ovn:"log"`
	Name        string            `ovn:"name"`
	Severity    string            `ovn:"severity"`
	Meter       string            `ovn:"meter"`
	ExternalIDs map[string]string `ovn:"external_ids"`
}

// LogicalRouter row
type LogicalRouter struct {
	UUID         string            `ovn:"uuid"`
	Name         string            `ovn:"name"`
	Enabled      bool              `ovn:"enabled"`
	Ports        []string          `ovn:"ports"`
	StaticRoutes []string          `ovn:"static_routes"`
	NAT          []string          `ovn:"nat"`
	LoadBalancer []string          `ovn:"load_balancer"`
	Options      map[string]string `ovn:"options"`
	ExternalIDs  map[string]string `ovn:"external_ids"`
}

// QoS row
type QoS struct {
	UUID        string            `ovn:"uuid"`
	Priority    int               `ovn:"priority"`
	Direction   string            `ovn:"direction"`
	Match       string            `ovn:"match"`
	Action      map[string]int    `ovn:"action"`
	Bandwidth   map[string]int    `ovn:"bandwidth"`
	ExternalIDs map[string]string `ovn:"external_ids"`
}

// Meter row
type Meter struct {
	UUID        string            `ovn:"uuid"`
	Name        string            `ovn:"name"`
	Unit        string            `ovn:"unit"`
	Bands       []string          `ovn:"bands"`
	ExternalIDs map[string]string `ovn:"external_ids"`
}

// MeterBand row
type MeterBand struct {
	UUID        string            `ovn:"uuid"`
	Action      string            `ovn:"action"`
	Rate        int64             `ovn:"rate"`
	BurstSize   int64             `ovn:"burst_size"`
	ExternalIDs map[string]string `ovn:"external_ids"`
}

// LogicalRouterPort row
type LogicalRouterPort struct {
	UUID           string            `ovn:"uuid"`
	Name           string            `ovn:"name"`
	GatewayChassis []string          `ovn:"gateway_chassis"`
	Networks       []string          `ovn:"networks"`
	MAC            string            `ovn:"mac"`
	Enabled        bool              `ovn:"enabled"`
	IPv6RAConfigs  map[string]string `ovn:"ipv6_ra_configs"`
	Options        map[string]string `ovn:"options"`
	Peer           string            `ovn:"peer"`
	ExternalIDs    map[string]string `ovn:"external_ids"`
}

// LogicalRouterStaticRoute row
type LogicalRouterStaticRoute struct {
	UUID        string            `ovn:"uuid"`
	IPPrefix    string            `ovn:"ip_prefix"`
	Nexthop     string            `ovn:"nexthop"`
	OutputPort  string            `ovn:"output_port"`
	Policy      string            `ovn:"policy"`
	ExternalIDs map[string]string `ovn:"external_ids"`
}

// NAT row
type NAT struct {
	UUID        string            `ovn:"uuid"`
	Type        string            `ovn:"type"`
	ExternalIP  string            `ovn:"external_ip"`
	ExternalMAC string            `ovn:"external_mac"`
	LogicalIP   string            `ovn:"logical_ip"`
	LogicalPort string            `ovn:"logical_port"`
	ExternalIDs map[string]string `ovn:"external_ids"`
}

// DHCPOptions row
type DHCPOptions struct {
	UUID        string            `ovn:"uuid"`
	CIDR        string            `ovn:"cidr"`
	Options     map[string]string `ovn:"options"`
	ExternalIDs map[string]string `ovn:"external_ids"`
}

// Connection row
type Connection struct {
	UUID            string            `ovn:"uuid"`
	Target          string            `ovn:"target"`
	MaxBackoff      int               `ovn:"max_backoff"`
	InactivityProbe int               `ovn:"inactivity_probe"`
	IsConnected     bool              `ovn:"is_connected"`
	Status          map[string]string `ovn:"status"`
	OtherConfig     map[string]string `ovn:"other_config"`
	ExternalIDs     map[string]string `ovn:"external_ids"`
}

// DNS row
type DNS struct {
	UUID        string            `ovn:"uuid"`
	Records     map[string]string `ovn:"records"`
	ExternalIDs map[string]string `ovn:"external_ids"`
}

// SSL row
type SSL struct {
	UUID            string            `ovn:"uuid"`
	PrivateKey      string            `ovn:"private_key"`
	Certificate     string            `ovn:"certificate"`
	CACert          string            `ovn:"ca_cert"`
	BootstrapCACert bool              `ovn:"bootstrap_ca_cert"`
	SSLProtocols    string            `ovn:"ssl_protocols"`
	SSLCiphers      string            `ovn:"ssl_ciphers"`
	ExternalIDs     map[string]string `ovn:"external_ids"`
}

// GatewayChassis ovnnb item
type GatewayChassis struct {
	UUID        string            `ovn:"uuid"`
	Name        string            `ovn:"name"`
	ChassisName string            `ovn:"chassis_name"`
	Priority    int               `ovn:"priority"`
	Options     map[string]string `ovn:"options"`
	ExternalIDs map[string]string `ovn:"external_ids"`
}
