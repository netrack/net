package pkg

import (
	"net"
)

type IPv4HeaderFlags int

const (
	DontFragment IPv4HeaderFlags = 1 + iota
	MoreFragment
)

type IPv4 struct {
	Version  int
	Len      int
	TOS      int
	TotalLen int
	ID       int
	Flags    IPv4HeaderFlags
	FragOff  int
	TTL      int
	Protocol int
	Checksum int
	Src      net.IP
	Dst      net.IP
	Options  []byte
}
