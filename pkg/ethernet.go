package pkg

import (
	"encoding/binary"
	"io"
	"net"
)

type EtherType uint16

type EthernetII struct {
	HWDst     net.HardwareAddr
	HWSrc     net.HardwareAddr
	EtherType EtherType
}

func (eth *EthernetII) init() {
	if eth.HWDst == nil {
		eth.HWDst = make(net.HardwareAddr, 6)
	}

	if eth.HWSrc == nil {
		eth.HWSrc = make(net.HardwareAddr, 6)
	}
}

func (eth *EthernetII) Read(r io.Reader) error {
	eth.init()
	err := binary.Read(r, binary.BigEndian, &eth.HWDst)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &eth.HWSrc)
	if err != nil {
		return err
	}

	return binary.Read(r, binary.BigEndian, &eth.EtherType)
}

type Ethernet8021q struct {
	HWDst     net.HardwareAddr
	HWSrc     net.HardwareAddr
	VLAN      VLAN
	EtherType EtherType
}

func (eth *Ethernet8021q) init() {
	if eth.HWDst == nil {
		eth.HWDst = make([]byte, 6)
	}

	if eth.HWSrc == nil {
		eth.HWSrc = make([]byte, 6)
	}
}

func (eth *Ethernet8021q) Read(r io.Reader) error {
	eth.init()
	err := binary.Read(r, binary.BigEndian, &eth.HWDst)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &eth.HWSrc)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &eth.VLAN)
	if err != nil {
		return err
	}

	return binary.Read(r, binary.BigEndian, &eth.EtherType)
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
