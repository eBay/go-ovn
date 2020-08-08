package goovn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient_InvalidNBTables(t *testing.T) {
	cfg := buildOvnDbConfig(DBNB)
	cfg.TableCols = map[string][]string{
		"Table1": {},
	}
	_, err := NewClient(cfg)
	assert.Error(t, err)
	t.Log(err.Error())
}

func TestNewClient_ValidNBTableInvalidCol(t *testing.T) {
	cfg := buildOvnDbConfig(DBNB)
	cfg.TableCols = map[string][]string{
		"Logical_Switch_Port": {"col1"},
	}
	_, err := NewClient(cfg)
	assert.Error(t, err)
	t.Log(err.Error())
}

func TestNewClient_ValidNBTableCols(t *testing.T) {
	cfg := buildOvnDbConfig(DBNB)
	cfg.TableCols = map[string][]string{
		"Logical_Switch": {},
	}
	api, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// create Switch with some values in external_ids
	t.Logf("Adding %s to OVN with external_ids set", LS3)
	cmd, err := api.LSAdd(LS3)
	if err != nil {
		t.Fatal(err)
	}
	err = api.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	// Add external_ids to LS3
	external_ids := map[string]string{NEUTRON_NETWORK: DUMMY, FOO: BAR}
	cmd, err = api.LSExtIdsAdd(LS3, external_ids)
	if err != nil {
		t.Fatal(err)
	}
	// execute the commands
	err = api.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	// Get the logical switch
	ls, err := api.LSGet(LS3)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, ls[0].Name, LS3)
	assert.Equal(t, ls[0].ExternalID[NEUTRON_NETWORK], external_ids[NEUTRON_NETWORK])
	assert.Equal(t, ls[0].ExternalID[FOO], external_ids[FOO])

	// Finally delete Switch
	t.Logf("Deleting LS3")
	cmd, err = ovndbapi.LSDel(LS3)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}
	t.Logf("Deleting LS3 Done")

	// Create a Logical Router and we should not be able to list the LR since
	// we didn't express interest in the Logical_Router table
	t.Logf("Adding LR %s", LR)
	cmd, err = api.LRAdd(LR, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = api.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Listing LR %s", LR)
	// We will not get the LR since we are not monitoring it.
	_, err = api.LRList()
	assert.Equal(t, err, ErrorNotFound)

	// We cannot delete the LR since the client doesn't have the info for it.
	_ = api.Close()
	cfg = buildOvnDbConfig(DBNB)
	cfg.TableCols = map[string][]string{
		"Logical_Router": {},
	}
	api, err = NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// finally delete the logical router
	t.Logf("Deleting LR %s", LR)
	cmd, err = ovndbapi.LRDel(LR)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Deleting LR %s Done", LR)
}

func TestNewClient_InvalidSBTables(t *testing.T) {
	cfg := buildOvnDbConfig(DBSB)
	cfg.TableCols = map[string][]string{
		"Table1": {},
	}
	_, err := NewClient(cfg)
	assert.Error(t, err)
	t.Log(err.Error())
}

func TestNewClient_ValidSBTableInvalidCol(t *testing.T) {
	cfg := buildOvnDbConfig(DBSB)
	cfg.TableCols = map[string][]string{
		"Chassis": {"col1"},
	}
	_, err := NewClient(cfg)
	assert.Error(t, err)
	t.Log(err.Error())
}
