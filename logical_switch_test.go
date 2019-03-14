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

func TestLogicalSwitch(t *testing.T) {
	var cmds []*OvnCommand
	var cmd *OvnCommand
	var err error

	cmds = make([]*OvnCommand, 0)

	cmd, err = ovndbapi.LSWDel(LSW)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	_, err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err = ovndbapi.LSWAdd(LSW)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)

	_, err = ovndbapi.Execute(cmds...)
	if err != nil {
		t.Fatal(err)
	}

	lsws, err := ovndbapi.GetLogicalSwitches()
	if err != nil || len(lsws) == 0 {
		t.Fatal("ls-add failed")
	}

	assert.Equal(t, true, len(lsws) == 1 && lsws[0].Name == LSW, "test[%s]: %v", "lsw added", lsws)
	cmds = make([]*OvnCommand, 0)

	cmd, err = ovndbapi.LSWDel(LSW)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)
	_, err = ovndbapi.Execute(cmds...)
	if err != nil {
		t.Fatal(err)
	}
	lsws, err = ovndbapi.GetLogicalSwitches()
	if err != nil || len(lsws) != 0 {
		t.Fatal("ls-del failed")
	}

}
