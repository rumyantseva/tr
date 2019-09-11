package trace

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type IPVersion string

const (
	IPv4 IPVersion = "ip4"
	IPv6 IPVersion = "ip6"
)

const (
	ianaProtocolICMP     = 1
	ianaProtocolIPv6ICMP = 58
)

const (
	icmpTimeExceeded = "time exceeded"
	icmpEchoReply    = "echo reply"
)

func Build(host string, ipv IPVersion, maxTTL int, callback func(hop *Hop, err error)) {
	for ttl := 1; ttl <= maxTTL; ttl++ {
		hop, reached, err := hop(host, ipv, ttl)
		callback(hop, err)

		if reached {
			break
		}
	}
}

func hop(host string, ipv IPVersion, ttl int) (hop *Hop, reached bool, err error) {
	network := string(ipv) + ":icmp"
	netConn, err := net.Dial(network, host)
	if err != nil {
		return
	}

	defer func() {
		errd := netConn.Close()
		if errd != nil {
			log.Printf("Couldn't close net conn: %v", err)
		}
	}()

	var address string
	var icmpType icmp.Type
	var ianaProto int
	switch ipv {
	case IPv4:
		address = "0.0.0.0"
		icmpType = ipv4.ICMPTypeEcho
		ianaProto = ianaProtocolICMP

		if err = ipv4.NewConn(netConn).SetTTL(ttl); err != nil {
			return
		}

	case IPv6:
		address = "::0"
		icmpType = ipv6.ICMPTypeEchoRequest
		ianaProto = ianaProtocolIPv6ICMP

		if err = ipv6.NewConn(netConn).SetHopLimit(ttl); err != nil {
			return
		}

	default:
		return nil, false, fmt.Errorf("a wrong IP version is given")
	}

	packetConn, err := icmp.ListenPacket(network, address)
	if err != nil {
		return
	}

	defer func() {
		errd := packetConn.Close()
		if errd != nil {
			log.Printf("Couldn't close packet conn: %v", err)
		}
	}()

	// Pepare and send an ICMP echo request
	wmsg := icmp.Message{
		Type: icmpType,
		Body: &icmp.Echo{},
	}

	w, err := wmsg.Marshal(nil)
	if err != nil {
		return
	}

	t := time.Now()
	_, err = netConn.Write(w)
	if err != nil {
		return
	}

	r := make([]byte, 1500)
	_, peerAddr, err := packetConn.ReadFrom(r)
	if err != nil {
		return
	}
	rttime := time.Since(t)

	rmsg, err := icmp.ParseMessage(ianaProto, r)
	if err != nil {
		return
	}

	ty := fmt.Sprint(rmsg.Type)
	switch ty {
	case icmpTimeExceeded:
		reached = false
	case icmpEchoReply:
		reached = true
	default:
		return nil, false, fmt.Errorf("received the `%s` ICMP message", rmsg.Type)
	}

	// Try to lookup one of the names associated with the peer address
	peer := Peer{
		Addr: peerAddr,
	}
	names, namerr := net.LookupAddr(peerAddr.String())
	if namerr == nil && len(names) > 0 {
		peer.Name = names[0]
	}

	hop = &Hop{
		TTL:     ttl,
		Peer:    peer,
		Latency: rttime,
	}

	return
}
