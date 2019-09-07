package trace

import (
	"net"
	"time"
)

type Hop struct {
	TTL     int
	Peer    net.Addr
	Latency time.Duration
}
