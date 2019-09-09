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

const IPv4 = "ip4"
const IPv6 = "ip6"

const ProtocolICMP = 1
const ProtocolIPv6ICMP = 58

const EchoTimeExceeded = "time exceeded"
const EchoReply = "echo reply"

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
		err = netConn.Close()
		if err != nil {
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
		ianaProto = ProtocolICMP

		if err = ipv4.NewConn(netConn).SetTTL(ttl); err != nil {
			return
		}

	case IPv6:
		address = "::0"
		icmpType = ipv6.ICMPTypeEchoRequest
		ianaProto = ProtocolIPv6ICMP

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
		err = packetConn.Close()
		if err != nil {
			log.Printf("Couldn't close packet conn: %v", err)
		}
	}()

	wmsg := icmp.Message{
		Type: icmpType,
		Body: &icmp.Echo{
			ID:   1,
			Seq:  1,
			Data: []byte{},
		},
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
	_, peer, err := packetConn.ReadFrom(r)
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
	case EchoTimeExceeded:
		reached = false
	case EchoReply:
		reached = true
	default:
		return nil, false, fmt.Errorf("error: %v", rmsg.Type)
	}

	hop = &Hop{
		TTL:     ttl,
		Peer:    peer,
		Latency: rttime,
	}

	return
}
