package main

import (
	goovn "github.com/ebay/go-ovn"
)

//Chassis struct defines an object in Chassis table
type LogicalRouter struct {
	goovn.BaseModel
	Name         string            `ovs:"name"`
	StaticRoutes []string          `ovs:"static_routes"`
	Nat          []string          `ovs:"nat"`
	ExternalIds  map[string]string `ovs:"external_ids"`
	Ports        []string          `ovs:"ports"`
	LoadBalancer []string          `ovs:"load_balancer"`
	Options      map[string]string `ovs:"options"`
	Policies     []string          `ovs:"policies"`
	Enabled      []bool            `ovs:"enabled"`
}

func (lr *LogicalRouter) Table() goovn.TableName {
	return "Logical_Router"
}
