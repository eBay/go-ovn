// +build travis

/**
 * Copyright (c) 2020 eBay Inc.
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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	SB_GLOBAL_OPTIONS_1_KEY = "controller-test-key"
	SB_GLOBAL_OPTIONS_1_VAL = "controller-test-val"
	SB_GLOBAL_DUMMY_OPT_KEY = "foo"
	SB_GLOBAL_DUMMY_OPT_VAL = "587c6ee2-93f9-4bd8-9794-f4a983d139a4"
)

func TestSBGlobalAPI(t *testing.T) {
	ovndbapi := getOVNClient(DBSB)
	t.Logf("Adding row to SB_Global table in OVN")
	options := make(map[string]string)
	options[SB_GLOBAL_DUMMY_OPT_KEY] = SB_GLOBAL_DUMMY_OPT_VAL

	ovn, ok := ovndbapi.(*ovndb)
	if !ok {
		t.Fatal(fmt.Errorf("Invalid type assertion"))
	}

	cmd, err := ovn.sbGlobalAdd(options)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Logf("Got an error while adding row to the sb_global table: %s", err)
		t.Fatal(err)
	}

	//Set options and verify
	options, err = ovndbapi.SBGlobalGetOptions()
	assert.Equal(t, options != nil, true)
	options[SB_GLOBAL_OPTIONS_1_KEY] = SB_GLOBAL_OPTIONS_1_VAL
	cmd, err = ovndbapi.SBGlobalSetOptions(options)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)

	if err != nil {
		t.Fatal(err)
	}
	//verify the options are set
	options, err = ovndbapi.SBGlobalGetOptions()
	if err != nil {
		t.Fatal(err)
	}
	val, ok := options[SB_GLOBAL_OPTIONS_1_KEY]
	assert.Equal(t, ok, true)
	assert.Equal(t, val, SB_GLOBAL_OPTIONS_1_VAL)

	t.Logf("Deleting row from SB_Global table in OVN")
	cmd, err = ovn.sbGlobalDel()
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	assert.Equal(t, err == nil, true)
}
