package goovn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadTransact(t *testing.T) {
	ovndbapi := getOVNClient(DBSB)
	t.Logf("Adding Chassis to OVN SB DB")
	ocmd, err := ovndbapi.ChassisAdd(CHASSIS_NAME, CHASSIS_HOSTNAME, ENCAP_TYPES, IP, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Adding Chassis to OVN failed with err %v", err)
	}
	t.Logf("Adding Chassis to OVN Done")

	t.Logf("Adding second Chassis to OVN SB DB but with same ENCAP_TYPES and IP")
	ocmd, err = ovndbapi.ChassisAdd(CHASSIS2_NAME, CHASSIS2_HOSTNAME, ENCAP_TYPES, IP, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// expecting constraint violation error with following details -- "Transaction causes multiple
	// rows in \"Encap\" table to have identical values (stt and \"10.0.0.11\") for index on columns
	// \"type\" and \"ip\".  First row, with UUID 9860cf40-bd82-4c24-9514-05b225434934, existed in
	// the database before this transaction and was not modified by the transaction.  Second row,
	// with UUID 10d7d018-7444-48de-89fc-cb062f88e520, was inserted by this transaction."
	err = ovndbapi.Execute(ocmd)
	assert.Error(t, err)

	t.Logf("Deleting Chassis:%v", CHASSIS_NAME)
	ocmd, err = ovndbapi.ChassisDel(CHASSIS_NAME)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}

	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}
}

func TestConvertGoSetToStringArray(t *testing.T) {
	// 1. create a logical switch and add a port to it.
	// 2. get the newly added port's uuid
	// 3. make sure that portUUID is in logical_switch's ports field.
	ovndbapi := getOVNClient(DBNB)
	t.Logf("Adding LogicalSwitch to OVN NB DB")
	ocmd, err := ovndbapi.LSAdd(LSW)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Adding Logical Switch to OVN failed with err %v", err)
	}
	t.Logf("Adding Logical Switch to OVN Done")

	ocmd, err = ovndbapi.LSPAdd(LSW, LSP)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Adding Logical Switch Port to OVN failed with err %v", err)
	}
	t.Logf("Adding Logical Switch Port to OVN Done")

	lspInfo, err := ovndbapi.LSPGet(LSP)
	if err != nil {
		t.Fatal(err)
	}
	lsInfo, err := ovndbapi.LSGet(LSW)
	if err != nil {
		t.Fatal(err)
	}
	uuidFound := false
	for _, port := range lsInfo[0].Ports {
		if port == lspInfo.UUID {
			uuidFound = true
			break
		}
	}
	if !uuidFound {
		t.Fatalf("couldn't find port uuid %s in %s", lspInfo.UUID, LSW)
	}
	t.Logf("Found Logical Switch Port's UUID in Logical Switch")

	t.Logf("Deleting the logical switch " + LSW)
	ocmd, err = ovndbapi.LSDel(LSW)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Deleting Logical Switch from OVN failed with err %v", err)
	}
	t.Logf("Deleted the logical switch " + LSW)
}

func lsNameToUUID(lsName string, c Client) (string, error) {
	lsList, err := c.LSList()
	if err == nil {
		for _, ls := range lsList {
			if ls.Name == lsName {
				return ls.UUID, nil
			}
			return "", ErrorNotFound
		}
	}
	return "", err
}

func lspNameToUUID(lspName string, c Client) (string, error) {
	lsp, err := c.LSPGet(lspName)
	if err == nil {
		return lsp.UUID, nil
	} else {
		return "", ErrorNotFound
	}
}

func TestExecuteR(t *testing.T) {
	ovndbapi := getOVNClient(DBNB)

	t.Run("execute one command in an ExecuteR call", func(t *testing.T) {
		// Create Switch
		cmd, err := ovndbapi.LSAdd(PG_TEST_LS1)
		assert.Nil(t, err)
		result, err := ovndbapi.ExecuteR(cmd)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
		//Check UUID returned
		lsUUID, err := lsNameToUUID(PG_TEST_LS1, ovndbapi)
		assert.Nil(t, err)
		assert.Greater(t, len(lsUUID), 0)
		assert.Equal(t, lsUUID, result[0])

		// Delete Switch (LSPs will get deleted by OVSDB)
		cmd, err = ovndbapi.LSDel(PG_TEST_LS1)
		assert.Nil(t, err)
		result, err = ovndbapi.ExecuteR(cmd)
		assert.Nil(t, err)
		// LSDel should not return any UUIDs
		assert.Nil(t, result)
	})

	t.Run("execute multiple commands in one ExecuteR call", func(t *testing.T) {
		var cmds []*OvnCommand

		// Create switch and ports
		cmd, err := ovndbapi.LSAdd(PG_TEST_LS1)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		// Add ports
		cmd, err = ovndbapi.LSPAdd(PG_TEST_LS1, PG_TEST_LSP1)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPSetAddress(PG_TEST_LSP1, ADDR)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPSetPortSecurity(PG_TEST_LSP1, ADDR)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPAdd(PG_TEST_LS1, PG_TEST_LSP2)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPSetAddress(PG_TEST_LSP2, ADDR2)
		assert.Nil(t, err)
		cmds = append(cmds, cmd)
		cmd, err = ovndbapi.LSPSetPortSecurity(PG_TEST_LSP2, ADDR2)
		cmds = append(cmds, cmd)
		assert.Nil(t, err)

		result, err := ovndbapi.ExecuteR(cmds...)
		assert.Nil(t, err)
		// Only the 3 "Add" commands above should return a UUID
		assert.Equal(t, 3, len(result))

		//Check UUIDs returned
		lsUUID, err := lsNameToUUID(PG_TEST_LS1, ovndbapi)
		assert.Nil(t, err)
		assert.Greater(t, len(lsUUID), 0)
		assert.Equal(t, lsUUID, result[0])

		lsp1UUID, err := lspNameToUUID(PG_TEST_LSP1, ovndbapi)
		assert.Nil(t, err)
		assert.Greater(t, len(lsp1UUID), 0)
		assert.Equal(t, lsp1UUID, result[1])

		lsp2UUID, err := lspNameToUUID(PG_TEST_LSP2, ovndbapi)
		assert.Nil(t, err)
		assert.Greater(t, len(lsp2UUID), 0)
		assert.Equal(t, lsp2UUID, result[2])

		// Delete Switch (LSPs will get deleted by OVSDB)
		cmd, err = ovndbapi.LSDel(PG_TEST_LS1)
		assert.Nil(t, err)
		result, err = ovndbapi.ExecuteR(cmd)
		assert.Nil(t, err)
		// LSDel should not return any UUIDs
		assert.Nil(t, result)
	})
}
