package l3

import (
	"bytes"
	"io"

	"github.com/netrack/net/encoding/binary"
	"github.com/netrack/net/iana"
)

const ICMPHeaderLen = 0x8

const (
	// The data received in the echo message must be returned in
	// the echo reply message.
	ICMPT_ECHO_REPLY ICMPType = 0

	// If, according to the information in the gateway's routing tables,
	// the network specified in the internet destination field of a
	// datagram is unreachable, e.g., the distance to the network is
	// infinity, the gateway may send a destination unreachable message
	// to the internet source host of the datagram.
	ICMPT_DESTINATION_UNREACHABLE ICMPType = 3

	ICMPT_REDIRECT             ICMPType = 5
	ICMPT_ECHO_REQUEST         ICMPType = 8
	ICMPT_ROUTER_ADVERTISEMENT ICMPType = 9
	ICMPT_ROUTER_SOLICITATION  ICMPType = 10

	// If the gateway processing a datagram finds the time to live field
	// is zero it must discard the datagram. The gateway may also notify
	// the source host via the time exceeded message.
	ICMPT_TIME_EXCEEDED ICMPType = 11

	// If the gateway or host processing a datagram finds a problem with
	// the header parameters such that it cannot complete processing the
	// datagram it must discard the datagram.
	ICMPT_PARAMETER_PROBLEM ICMPType = 12

	ICMPT_TIMESTAMP       ICMPType = 13
	ICMPT_TIMESTAMP_REPLY ICMPType = 14
)

type ICMPType uint8

const (
	// Generated by a router if a forwarding path
	// (route) to the destination network is not available
	ICMPC_NETWORK_UNREACHABLE ICMPCode = iota

	// Generated by a router if a forwarding path (route) to the destination
	// host on a directly connected network is not available (does not
	// respond to ARP)
	ICMPC_HOST_UNREACHABLE

	// Generated if the transport protocol designated in a datagram
	// is not supported in the transport layer of the final destination
	ICMPC_PROTOCOL_UNREACHABLE

	// Generated if the designated transport protocol (e.g., UDP) is
	// unable to demultiplex the datagram in the transport layer of the
	// final destination but has no protocol mechanism to inform the sender
	ICMPC_PORT_UNREACHABLE

	// Generated if a router needs to fragment a datagram but cannot
	// since the DF flag is set
	ICMPC_FRAGMENTATION_REQUIRED

	// Generated if a router cannot forward a packet to the next hop
	// in a source route option
	ICMPC_SOURCE_ROUTE_FAILED

	// This code SHOULD NOT be generated since it would imply on the part
	// of the router that the destination network does not exist (net
	// unreachable code 0 SHOULD be used in place of code 6)
	ICMPC_NETWORK_UNKNOWN

	// Generated only when a router can determine (from link layer advice)
	// that the destination host does not exist
	ICMPC_HOST_UNKNOWN

	// Source host isolated
	ICMPC_SOURCE_HOST_ISOLATED

	// Communication with destination network administratively prohibited
	ICMPC_NETWORK_ADMINISTRATIVELY_PROHIBITED

	// Communication with destination host administratively prohibited
	ICMPC_HOST_ADMINISTRATIVELY_PROHIBITED

	// Generated by a router if a forwarding path (route) to the
	// destination network with the requested or default TOS is not available
	ICMPC_NETWORK_UNREACHABLE_FOR_TOS

	// Generated if a router cannot forward a packet because
	// its route(s) to the destination do not match either
	// the TOS requested in the datagram or the default TOS (0)
	ICMPC_HOST_UNREACHABLE_FOR_TOS

	// Generated if a router cannot forward a packet due to
	// administrative filtering
	ICMPC_COMMUNICATION_ADMINISTRATIVELY_PROHIBITED

	// Sent by the first hop router to a host to indicate that a
	// requested precedence is not permitted for the particular
	// combination of source/destination host or network, upper
	// layer protocol, and source/destination port
	ICMPC_HOST_PRECEDENCE_VIOLATION

	// The network operators have imposed a minimum level
	// of precedence required for operation, the datagram
	// was sent with a precedence below this level
	ICMPC_PRECEDECE_CUTOFF_IN_EFFECT
)

type ICMPCode uint8

type ICMP struct {
	Type     ICMPType
	Code     ICMPCode
	Checksum uint16
}

func (p *ICMP) WriteTo(w io.Writer) (int64, error) {
	return binary.Write(w, binary.BigEndian, p)
}

func (p *ICMP) ReadFrom(r io.Reader) (int64, error) {
	return binary.Read(r, binary.BigEndian, p)
}

type ICMPEcho struct {
	ICMP
	ID   uint16
	Seq  uint16
	Data []byte
}

func (p *ICMPEcho) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	n, err := binary.WriteSlice(&buf, binary.BigEndian, []interface{}{
		p.ICMP, p.ID, p.Seq, p.Data,
	})

	if err != nil {
		return n, err
	}

	b := buf.Bytes()
	copy(b[2:4], iana.Checksum(b))
	return binary.Write(w, binary.BigEndian, b)
}

func (p *ICMPEcho) ReadFrom(r io.Reader) (int64, error) {
	return binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&p.ICMP, &p.ID, &p.Seq, &p.Data,
	})
}