package goovn

import (
	"fmt"
	"reflect"
)

func structToMap(iface interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(iface)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("structToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get("ovn"); len(tagv) > 0 {
			// set key of map to value in struct field
			out[tagv] = v.Field(i).Interface()
		}
	}

	return out, nil
}

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
