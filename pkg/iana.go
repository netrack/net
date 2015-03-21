package pkg

const (
	PROTO_IPV4            ProtoType = 0x0800
	PROTO_ARP             ProtoType = 0x0806
	PROTO_FRAME_RELAY_ARP ProtoType = 0x0808
	PROTO_TRILL           ProtoType = 0x22F3
	PROTO_L2_IS_IS        ProtoType = 0x22F4
	PROTO_REVERSE_ARP     ProtoType = 0x8035
	PROTO_IPV6            ProtoType = 0x86DD
	PROTO_PPP             ProtoType = 0x880B
	PROTO_GSMP            ProtoType = 0x880C
	PROTO_MPLS            ProtoType = 0x8847
)

type ProtoType uint16

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
