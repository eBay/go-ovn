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
	"testing"

	"github.com/stretchr/testify/assert"
)

func findAS(name string) bool {
	as, err := ovndbapi.GetAddressSets()
	if err != nil {
		return false
	}
	for _, a := range as {
		if a.Name == name {
			return true
		}
	}
	return false
}

func addressSetCmp(asname string, targetvalue []string) bool {
	as, err := ovndbapi.GetAddressSets()
	if err != nil {
		return false
	}
	for _, a := range as {
		if a.Name == asname {
			if len(a.Addresses) == len(targetvalue) {
				addressSetMap := map[string]bool{}
				for _, i := range a.Addresses {
					addressSetMap[i] = true
				}
				for _, t := range targetvalue {
					if _, ok := addressSetMap[t]; !ok {
						return false
					}
				}
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func TestAddressSet(t *testing.T) {
	addressList := []string{"127.0.0.1"}
	var cmd *OvnCommand
	var err error

	/*
		// can not call like:
		// ovndbapi.ASAdd("AS1", addressList, map[string][]{})
		// it will not be successful when input empty map.
	*/
	cmd, err = ovndbapi.ASAdd("AS1", addressList, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	as, err := ovndbapi.GetAddressSets()
	if err != nil || len(as) == 0 {
		t.Fatal("address_set is nil")
	}
	assert.Equal(t, true, addressSetCmp("AS1", addressList), "test[%s] and value[%v]", "address set 1 added.", as[0].Addresses)

	cmd, err = ovndbapi.ASAdd("AS2", addressList, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	as, err = ovndbapi.GetAddressSets()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, addressSetCmp("AS2", addressList), "test[%s] and value[%v]", "address set 2 added.", as[1].Addresses)

	addressList = []string{"127.0.0.4", "127.0.0.5", "127.0.0.6"}
	cmd, err = ovndbapi.ASUpdate("AS2", addressList, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	as, err = ovndbapi.GetAddressSets()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, addressSetCmp("AS2", addressList), "test[%s] and value[%v]", "address set added with different list.", as[0].Addresses)

	addressList = []string{"127.0.0.4", "127.0.0.5"}
	cmd, err = ovndbapi.ASUpdate("AS2", addressList, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	as, err = ovndbapi.GetAddressSets()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, addressSetCmp("AS2", addressList), "test[%s] and value[%v]", "address set updated.", as[0].Addresses)

	cmd, err = ovndbapi.ASDel("AS1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, false, findAS("AS1"), "test AS remove")

	cmd, err = ovndbapi.ASDel("AS2")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, false, findAS("AS2"), "test AS remove")
}
