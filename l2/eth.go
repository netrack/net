package l2

import (
	"io"
	"net"

	"github.com/netrack/net/encoding/binary"
	"github.com/netrack/net/iana"
)

var (
	HWBcast  = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	HWUnspec = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

type EthernetII struct {
	HWDst   net.HardwareAddr
	HWSrc   net.HardwareAddr
	EthType iana.EthType
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
		&eth.HWDst, &eth.HWSrc, &eth.EthType,
	})
}

func (eth *EthernetII) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		eth.HWDst, eth.HWSrc, eth.EthType,
	})
}

type Ethernet8021q struct {
	HWDst   net.HardwareAddr
	HWSrc   net.HardwareAddr
	VLAN    VLAN
	EthType iana.EthType
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
		&eth.HWDst, &eth.HWSrc, &eth.VLAN, &eth.EthType,
	})
}

func (eth *Ethernet8021q) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		eth.HWDst, eth.HWSrc, eth.VLAN, eth.EthType,
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
