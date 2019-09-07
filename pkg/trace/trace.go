package trace

import (
	"log"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type IPVersion string

const IPv4 = "ip4"
const IPv6 = "ip6"

func Build(host string, ipv IPVersion, maxTTL int, callback func(hop *Hop, err error)) {
	for ttl := 1; ttl <= maxTTL; ttl++ {
		hop, err := hop(host, ipv, ttl)
		callback(hop, err)
	}
}

func hop(host string, ipv IPVersion, ttl int) (*Hop, error) {
	network := string(ipv) + ":icmp"
	netConn, err := net.Dial(network, host)
	if err != nil {
		return nil, err
	}

	var address string
	var icmpType icmp.Type
	switch ipv {
	case IPv4:
		address = "0.0.0.0"
		icmpType = ipv4.ICMPTypeEcho

		conn := ipv4.NewConn(netConn)
		if err = conn.SetTTL(ttl); err != nil {
			return nil, err
		}

	case IPv6:
		address = "::0"
		icmpType = ipv6.ICMPTypeEchoRequest

		conn := ipv6.NewConn(netConn)
		if err = conn.SetHopLimit(ttl); err != nil {
			return nil, err
		}
	}

	packetConn, err := icmp.ListenPacket(network, address)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = packetConn.Close()
		if err != nil {
			log.Printf("Couldn't close packet conn: %v", err)
		}
	}()

	msg := icmp.Message{
		Type: icmpType,
		Body: &icmp.Echo{
			ID:   1,
			Seq:  1,
			Data: []byte{},
		},
	}

	w, err := msg.Marshal(nil)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	_, err = netConn.Write(w)
	if err != nil {
		return nil, err
	}

	r := make([]byte, 1500)
	_, peer, err := packetConn.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	hop := Hop{
		TTL:     ttl,
		Peer:    peer,
		Latency: time.Since(t),
	}

	return &hop, nil
}
