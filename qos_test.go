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
	"github.com/stretchr/testify/assert"
	"testing"
)

const LSW3 = "TEST_LSW3"

func TestQoS(t *testing.T) {
	ovndbapi := getOVNClient(DBNB)
	var cmd *OvnCommand
	var err error

	cmd, err = ovndbapi.LSAdd(LSW3)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.QoSAdd(LSW3, "to-lport", 1001, `inport=="lp3"`, nil, map[string]int{"rate": 1234, "burst": 12345}, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	cmd, err = ovndbapi.QoSAdd(LSW3, "from-lport", 1002, `inport=="lp3"`, nil, map[string]int{"rate": 1234, "burst": 12345}, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	cmd, err = ovndbapi.QoSAdd(LSW3, "to-lport", 1003, `inport=="lp3"`, nil, map[string]int{"rate": 1234, "burst": 12345}, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	qosrules, err := ovndbapi.QoSList(LSW3)
	if err != nil {
		t.Fatal(err)
	}

	if len(qosrules) != 3 {
		t.Fatalf("qos rules not inserted %v", qosrules)
	} else {
		for _, rule := range qosrules {
			t.Logf("qos rule inserted %v\n", rule)
		}
	}

	cmd, err = ovndbapi.QoSDel(LSW3, "to-lport", -1, "")
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	qosrules, err = ovndbapi.QoSList(LSW3)
	if err != nil {
		t.Fatal(err)
	}

	if len(qosrules) != 1 {
		for _, rule := range qosrules {
			t.Logf("qos rule not deleted %v\n", rule)
		}
		t.Fatalf("qos rules not deleted %v", qosrules)
	}
	if qosrules[0].Priority != 1002 {
		t.Fatalf("invalid qos rule deleted %#+v\n", qosrules[0])
	}

	cmd, err = ovndbapi.LSDel(LSW3)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	// verify QOS list for non-existing switch
	_, err = ovndbapi.QoSList(FAKENOSWITCH)
	if err != nil {
		assert.EqualError(t, ErrorNotFound, err.Error())
	}
}
