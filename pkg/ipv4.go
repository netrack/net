package pkg

import (
	"bytes"
	"io"
	"net"

	"github.com/netrack/net/encoding/binary"
)

const (
	IPv4Version   uint8 = 0x4
	IPv4HeaderLen uint8 = 0x5
)

const (
	IPV4_PROTO_HOPOPT IPv4Proto = iota
	IPV4_PROTO_ICMP
	IPV4_PROTO_IGMP
	IPV4_PROTO_GGP
	IPV4_PROTO_IP_IN_IP
	IPV4_PROTO_ST
	IPV4_PROTO_TCP
	IPV4_PROTO_CBT
	IPV4_PROTO_EGP
	IPV4_PROTO_IGP
	IPV4_PROTO_BBN_RCC_MON
	IPV4_PROTO_NVP_II
	IPV4_PROTO_PUP
	IPV4_PROTO_ARGUS
	IPV4_PROTO_EMCON
	IPV4_PROTO_XNET
	IPV4_PROTO_CHAOS
	IPV4_PROTO_UDP
)

type IPv4Proto uint8

const (
	IPV4F_DO_NOT_FRAGMENT IPv4Flag = 1 + iota
	IPV4F_MORE_FRAGMENT
)

type IPv4Flag uint16

type IPv4 struct {
	DiffServ   uint8
	Len        uint16
	ID         uint16
	Flags      IPv4Flag
	FragOffset uint16
	TTL        uint8
	Proto      IPv4Proto
	Checksum   uint16
	Src        net.IP
	Dst        net.IP
}

func (h *IPv4) WriteTo(w io.Writer) (int64, error) {
	var flagOffset uint16 = (uint16(h.Flags) << 0xd) | (h.FragOffset & 0x1ffd)
	var buf bytes.Buffer

	n, err := binary.WriteSlice(&buf, binary.BigEndian, []interface{}{
		(IPv4Version<<4 | IPv4HeaderLen),
		h.DiffServ,
		h.Len,
		h.ID,
		flagOffset,
		h.TTL,
		h.Proto,
		uint16(0),
		h.Src,
		h.Dst,
	})

	if err != nil {
		return n, err
	}

	b := buf.Bytes()
	copy(b[10:12], Checksum(b))
	return binary.Write(w, binary.BigEndian, b)
}

func (h *IPv4) ReadFrom(r io.Reader) (int64, error) {
	var flagOffset uint16
	var versionHeaderLen uint8

	h.Src = make(net.IP, 4)
	h.Dst = make(net.IP, 4)

	n, err := binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&versionHeaderLen,
		&h.DiffServ,
		&h.Len,
		&h.ID,
		&flagOffset,
		&h.TTL,
		&h.Proto,
		&h.Checksum,
		&h.Src,
		&h.Dst,
	})

	h.Flags = IPv4Flag(flagOffset) >> 0xd
	h.FragOffset = flagOffset & 0x1ffd

	return n, err
}
