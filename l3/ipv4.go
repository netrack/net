package l3

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"

	"github.com/netrack/net/encoding/binary"
	"github.com/netrack/net/iana"
)

const (
	IPv4Version   uint8  = 0x4
	IPv4HeaderLen uint16 = 0x14
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
	Proto      iana.IPProto
	Checksum   uint16
	Src        net.IP
	Dst        net.IP
	Payload    io.Reader
}

func (h *IPv4) WriteTo(w io.Writer) (int64, error) {
	var flagOffset uint16 = (uint16(h.Flags) << 0xd) | (h.FragOffset & 0x1ffd)
	var buf bytes.Buffer

	payload, err := ioutil.ReadAll(h.Payload)
	if err != nil {
		return 0, err
	}

	n, err := binary.WriteSlice(&buf, binary.BigEndian, []interface{}{
		(IPv4Version<<4 | uint8(IPv4HeaderLen/4)),
		h.DiffServ,
		uint16(len(payload)) + uint16(IPv4HeaderLen),
		h.ID,
		flagOffset,
		uint8(255),
		h.Proto,
		uint16(0),
		h.Src,
		h.Dst,
		payload,
	})

	if err != nil {
		return n, err
	}

	b := buf.Bytes()
	copy(b[10:12], iana.Checksum(b[:IPv4HeaderLen]))

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
