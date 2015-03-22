package iana

import (
	"errors"
)

var ErrInvalidChecksum = errors.New("iana: invalid checksum")

const (
	ETHT_IPV4            EthType = 0x0800
	ETHT_ARP             EthType = 0x0806
	ETHT_FRAME_RELAY_ARP EthType = 0x0808
	ETHT_TRILL           EthType = 0x22F3
	ETHT_L2_IS_IS        EthType = 0x22F4
	ETHT_REVERSE_ARP     EthType = 0x8035
	ETHT_IPV6            EthType = 0x86DD
	ETHT_PPP             EthType = 0x880B
	ETHT_GSMP            EthType = 0x880C
	ETHT_MPLS            EthType = 0x8847
)

type EthType uint16

const (
	IP_PROTO_HOPOPT IPProto = iota
	IP_PROTO_ICMP
	IP_PROTO_IGMP
	IP_PROTO_GGP
	IP_PROTO_IP_IN_IP
	IP_PROTO_ST
	IP_PROTO_TCP
	IP_PROTO_CBT
	IP_PROTO_EGP
	IP_PROTO_IGP
	IP_PROTO_BBN_RCC_MON
	IP_PROTO_NVP_II
	IP_PROTO_PUP
	IP_PROTO_ARGUS
	IP_PROTO_EMCON
	IP_PROTO_XNET
	IP_PROTO_CHAOS
	IP_PROTO_UDP
)

type IPProto uint8

func Checksum(b []byte) []byte {
	var index int
	var sum uint32

	for index = 0; index < len(b)-1; index += 2 {
		sum += (uint32(b[index]) << 8) | uint32(b[index+1])
	}

	if index < len(b) {
		sum += uint32(b[index])
	}

	for sum>>16 != 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}

	sum ^= 0xffff
	return []byte{byte(sum >> 8), byte(sum)}
}
