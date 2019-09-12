package trace

import (
	"net"
	"time"
)

// Hop represents a single hop info:
// sequence number, hostname and address, latency, received icmp message.
type Hop struct {
	Seq     int
	Peer    Peer
	Latency time.Duration
	Message string
}

// Peer is a host presented by a name and address.
type Peer struct {
	Name string
	Addr net.Addr
}
