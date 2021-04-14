package goovn

import (
	"testing"
)

func TestMeter(t *testing.T) {
	ovndbapi := getOVNClient(DBNB)
	var cmds []*OvnCommand
	cmd, err := ovndbapi.MeterAdd(METER1, "drop", 101, "kbps", nil, 300)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)
	cmd, err = ovndbapi.MeterAdd(METER2, "drop", 101, "kbps", nil, 300)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)
	cmd, err = ovndbapi.MeterAdd(METER3, "drop", 101, "kbps", nil, 300)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)
	err = ovndbapi.Execute(cmds...)
	if err != nil {
		t.Fatal(err)
	}

	meter, err := ovndbapi.MeterList()
	if err != nil {
		t.Fatal(err)
	}
	if len(meter) < 3 {
		t.Fatal("Meter add Fail")
	}

	meterBands, err := ovndbapi.MeterBandsList()
	if err != nil {
		t.Fatal(err)
	}
	if len(meterBands) < 3 {
		t.Fatal("Meter bands shows Fail")
	}

	defer func() {
		cmd, err = ovndbapi.MeterDel()
		if err != nil {
			t.Fatal(err)
		}
		err = ovndbapi.Execute(cmd)
		if err != nil {
			t.Fatal(err)
		}
		meter, err = ovndbapi.MeterList()
		if err != nil {
			t.Fatal(err)
		}
		if len(meter) != 0 {
			t.Fatal("Delete All Meter Fail")
		}
	}()

	defer func() {
		cmd, err = ovndbapi.MeterDel(METER1)
		if err != nil {
			t.Fatal(err)
		}
		err = ovndbapi.Execute(cmd)
		if err != nil {
			t.Fatal(err)
		}
		meter, err = ovndbapi.MeterList()
		if err != nil {
			t.Fatal(err)
		}
		if len(meter) < 2 {
			t.Fatal("Delete single Meter Error")
		}
	}()
}
