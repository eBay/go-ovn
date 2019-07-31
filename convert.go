package goovn

import (
	"fmt"
	"reflect"

	"github.com/ebay/libovsdb"
)

func cmpRows(crow map[string]interface{}, lrow map[string]interface{}) bool {
	var found bool
	for lk, lv := range lrow {
		switch lk {
		case "external_ids":
			cv, ok := crow[lk]
			if !ok {
				return false
			}

			cmap, ok := cv.(map[string]string)
			if !ok {
				return false
			}

			lmap, ok := lv.(map[string]string)
			if !ok {
				return false
			}
			if len(cmap) == 0 {
				return false
			}

		mapLoop:
			for lmk, lmv := range lmap {
				if cmv, ok := cmap[lmk]; !ok {
					return false
				} else if cmv == lmv {
					found = true
					break mapLoop
				}
			}
		default:
			if cv, ok := crow[lk]; ok {
				if !reflect.DeepEqual(cv, lv) {
					return false
				}
				found = true
			}
		}
	}

	return found
}

func rowUpdateToStruct(table string, uuid string, raw interface{}) (interface{}, error) {
	var row interface{}
	switch table {
	case "NB_Global":
		row = &NBGlobal{UUID: uuid}
	case "ACL":
		row = &ACL{UUID: uuid}
	case "Logical_Switch":
		row = &LogicalSwitch{UUID: uuid}
	case "Address_Set":
		row = &AddressSet{UUID: uuid}
	case "Port_Group":
		row = &PortGroup{UUID: uuid}
	case "Load_Balancer":
		row = &LoadBalancer{UUID: uuid}
	case "Logical_Router":
		row = &LogicalRouter{UUID: uuid}
	case "QoS":
		row = &QoS{UUID: uuid}
	case "Meter":
		row = &Meter{UUID: uuid}
	case "Meter_Band":
		row = &MeterBand{UUID: uuid}
	case "Logical_Router_Port":
		row = &LogicalRouterPort{UUID: uuid}
	case "Logical_Router_Static_Router":
		row = &LogicalRouterStaticRoute{UUID: uuid}
	case "NAT":
		row = &NAT{UUID: uuid}
	case "DHCP_Options":
		row = &DHCPOptions{UUID: uuid}
	case "Connection":
		row = &Connection{UUID: uuid}
	case "DNS":
		row = &DNS{UUID: uuid}
	case "SSL":
		row = &SSL{UUID: uuid}
	case "Gateway_Chassis":
		row = &GatewayChassis{UUID: uuid}
	default:
		return nil, fmt.Errorf("unsupported table %v update: %v", table, raw)
	}

	mp, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unsupported data %v", raw)
	}
	if err := libovsdb.MapToStruct(mp, "ovn", row); err != nil {
		return nil, err
	}

	return row, nil
}
