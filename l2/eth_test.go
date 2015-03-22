package l2

import (
	"bytes"
	"testing"
)

func TestEthernetIIRead(t *testing.T) {
	header := []byte{0x33, 0x33, 0x00, 0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x00, 0x0, 0x002, 0x86, 0xdd}
	reader := bytes.NewBuffer(header)

	var eth EthernetII
	err := eth.Read(reader)

	if err != nil {
		t.Fatalf("Failed to parse ethernet header: '%s'", err.Error())
	}

	if eth.HWDst.String() != "33:33:00:00:00:16" {
		t.Fatalf("Wrong destination MAC address")
	}

	if eth.HWSrc.String() != "00:00:00:00:00:02" {
		t.Fatalf("Wrong source MAC address")
	}

	if eth.EtherType != 0x86dd {
		t.Fatalf("Wrong ether type value")
	}
}
