package l3

import (
	"io"
	"net"

	"github.com/netrack/net/encoding/binary"
	"github.com/netrack/net/iana"
)

const (
	ARPT_ETHERNET ARPType = 1 + iota
	ARPT_EXPERIMENTAL_ETHERNET
	ARPT_AR_X25
	ARPT_TOKEN_RING
	ARPT_CHAOS
	ARPT_IEEE_802
	ARPT_ARCNET
	ARPT_HYPER_CHANNEL
	ARPT_LANSTAR
	ARPT_AUTONET
	ARPT_LOCAL_TALK
	ARPT_LOCAL_NET
	ARPT_ULTRA_LINK
	ARPT_SMDS
	ARPT_FRAME_RELAY
)

type ARPType uint16

const (
	ARPOT_REQUEST ARPOperation = 1 + iota
	ARPOT_REPLY
)

type ARPOperation uint16

type ARP struct {
	HWType    ARPType
	ProtoType iana.EthType
	Operation ARPOperation
	HWSrc     net.HardwareAddr
	ProtoSrc  net.IP
	HWDst     net.HardwareAddr
	ProtoDst  net.IP
}

func (a *ARP) init() {
	if a.HWDst == nil {
		a.HWDst = make(net.HardwareAddr, 6)
	}

	if a.HWSrc == nil {
		a.HWSrc = make(net.HardwareAddr, 6)
	}

	if a.ProtoDst == nil {
		a.ProtoDst = make(net.IP, 4)
	}

	if a.ProtoSrc == nil {
		a.ProtoSrc = make(net.IP, 4)
	}
}

func (a *ARP) ReadFrom(r io.Reader) (n int64, err error) {
	var hwlen, plen uint8
	a.init()

	return binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&a.HWType,
		&a.ProtoType,
		&hwlen,
		&plen,
		&a.Operation,
		&a.HWSrc,
		&a.ProtoSrc,
		&a.HWDst,
		&a.ProtoDst,
	})
}

func (a *ARP) WriteTo(w io.Writer) (n int64, err error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		a.HWType,
		a.ProtoType,
		uint8(len(a.HWDst)),
		uint8(len(a.ProtoDst)),
		a.Operation,
		a.HWSrc,
		a.ProtoSrc,
		a.HWDst,
		a.ProtoDst,
	})
}
