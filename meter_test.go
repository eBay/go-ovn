package goovn

import (
	"testing"
)

func TestMeter(t *testing.T) {
	var cmds []*OvnCommand
	cmd, err := ovndbapi.MeterAdd("meter1", "drop", 101, "kbps", nil, 300)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)
	cmd, err = ovndbapi.MeterAdd("meter2", "drop", 101, "kbps", nil, 300)
	if err != nil {
		t.Fatal(err)
	}
	cmds = append(cmds, cmd)
	cmd, err = ovndbapi.MeterAdd("meter3", "drop", 101, "kbps", nil, 300)
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
	if len(meter) != 3 {
		t.Fatal("Meter add Fail")
	}

	meterBands, err := ovndbapi.MeterBandsList()
	if err != nil {
		t.Fatal(err)
	}
	if len(meterBands) != 3 {
		t.Fatal("Meter bands shows Fail")
	}

	cmd, err = ovndbapi.MeterDel("meter1")
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
	if len(meter) != 2 {
		t.Fatal("Delete single Meter Error")
	}

	cmd, err = ovndbapi.MeterDel()
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}
	meter, err = ovndbapi.MeterList()
	if len(meter) != 0 {
		t.Fatal("Delete All Meter Fail")
	}

}
