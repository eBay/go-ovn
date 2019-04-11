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

func (ocmd *OvnCommand) Execute() error {
	return ocmd.Exe.Execute()
}

type Execution interface {
	//Excute multi-commands
	Execute(cmds ...*OvnCommand) error
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
