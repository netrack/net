package pkg

import (
	"io"
	"net"

	"github.com/netrack/net/encoding/binary"
)

type EthernetII struct {
	HWDst     net.HardwareAddr
	HWSrc     net.HardwareAddr
	EtherType ProtoType
}

func (eth *EthernetII) init() {
	if eth.HWDst == nil {
		eth.HWDst = make(net.HardwareAddr, 6)
	}

	if eth.HWSrc == nil {
		eth.HWSrc = make(net.HardwareAddr, 6)
	}
}

func (eth *EthernetII) ReadFrom(r io.Reader) (int64, error) {
	eth.init()
	return binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&eth.HWDst, &eth.HWSrc, &eth.EtherType,
	})
}

func (eth *EthernetII) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		eth.HWDst, eth.HWSrc, eth.EtherType,
	})
}

type Ethernet8021q struct {
	HWDst     net.HardwareAddr
	HWSrc     net.HardwareAddr
	VLAN      VLAN
	EtherType ProtoType
}

func (eth *Ethernet8021q) init() {
	if eth.HWDst == nil {
		eth.HWDst = make(net.HardwareAddr, 6)
	}

	if eth.HWSrc == nil {
		eth.HWSrc = make(net.HardwareAddr, 6)
	}
}

func (eth *Ethernet8021q) ReadFrom(r io.Reader) (int64, error) {
	eth.init()
	return binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&eth.HWDst, &eth.HWSrc, &eth.VLAN, &eth.EtherType,
	})
}

func (eth *Ethernet8021q) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		eth.HWDst, eth.HWSrc, eth.VLAN, eth.EtherType,
	})
}

type VLAN struct {
	TPID uint16
	TCI  uint16
}

func (v *VLAN) PCP() int {
	return int((v.TCI & 0xe000) >> 13)
}

func (v *VLAN) DEI() int {
	return int((v.TCI & 0x1000) >> 12)
}

func (v *VLAN) VID() int {
	return int((v.TCI & 0x0fff))
}
