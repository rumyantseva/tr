package trace

import (
	"net"
	"time"
)

type Hop struct {
	TTL     int
	Peer    Peer
	Latency time.Duration
}

type Peer struct {
	Name string
	Addr net.Addr
}
