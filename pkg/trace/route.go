package trace

import (
	"net"
	"time"
)

type Hop struct {
	Seq     int
	Peer    Peer
	Latency time.Duration
	Message string
}

type Peer struct {
	Name string
	Addr net.Addr
}
