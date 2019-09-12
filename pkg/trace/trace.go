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

type connProperties struct {
	address   string
	icmpType  icmp.Type
	ianaProto int
}

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
	icmpEchoReply = "echo reply"
)

func Build(host string, ipv IPVersion, maxTTL int, timeout time.Duration, callback func(hop *Hop, err error)) error {
	for ttl := 1; ttl <= maxTTL; ttl++ {
		h, reached, err := hop(host, ipv, ttl, timeout)
		callback(h, err)

		if reached {
			return nil
		}
	}

	return fmt.Errorf("max TTL exceeded, target is not reached")
}

func hop(host string, ipv IPVersion, ttl int, timeout time.Duration) (hop *Hop, reached bool, err error) {
	var network string
	switch ipv {
	case IPv4:
		network = "ip4:icmp"
	case IPv6:
		network = "ip6:ipv6-icmp"
	}

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

	if err = netConn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return
	}

	var cp *connProperties
	cp, err = setTTL(netConn, ipv, ttl)
	if err != nil {
		return
	}

	// Setup a socket to listen to ICMP Echo Reply
	packetConn, err := icmp.ListenPacket(network, cp.address)
	if err != nil {
		return
	}

	defer func() {
		errd := packetConn.Close()
		if errd != nil {
			log.Printf("Couldn't close packet conn: %v", err)
		}
	}()

	if err = packetConn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return
	}

	// Prepare and send an ICMP Echo Request
	wmsg := icmp.Message{
		Type: cp.icmpType,
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

	r := make([]byte, 1000)
	_, peerAddr, err := packetConn.ReadFrom(r)
	if err != nil {
		return
	}
	rttime := time.Since(t)

	// Parse ICMP Echo Reply to understand if we reached the destination
	rmsg, err := icmp.ParseMessage(cp.ianaProto, r)
	if err != nil {
		return
	}

	if rmsg.Code != echoReplyCode(ipv) {
		return nil, false, fmt.Errorf("the response code (%d) is not an ICMP Echo Reply", rmsg.Code)
	}

	ty := fmt.Sprint(rmsg.Type)
	switch ty {
	case icmpEchoReply:
		reached = true
	default:
		reached = false
	}

	// Try to lookup one of the names associated with the peer address
	peer := Peer{
		Addr: peerAddr,
		Name: peerAddr.String(),
	}
	names, namerr := net.LookupAddr(peerAddr.String())
	if namerr == nil && len(names) > 0 {
		peer.Name = names[0]
	}

	hop = &Hop{
		Seq:     ttl,
		Peer:    peer,
		Latency: rttime,
		Message: fmt.Sprintf("%s", rmsg.Type),
	}

	return
}

func setTTL(netConn net.Conn, ipv IPVersion, ttl int) (*connProperties, error) {
	var cp connProperties
	switch ipv {
	case IPv4:
		cp = connProperties{
			address:   "0.0.0.0",
			icmpType:  ipv4.ICMPTypeEcho,
			ianaProto: ianaProtocolICMP,
		}

		if err := ipv4.NewConn(netConn).SetTTL(ttl); err != nil {
			return nil, err
		}

	case IPv6:
		cp = connProperties{
			address:   "::0",
			icmpType:  ipv6.ICMPTypeEchoRequest,
			ianaProto: ianaProtocolIPv6ICMP,
		}

		if err := ipv6.NewConn(netConn).SetHopLimit(ttl); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("a wrong IP version is given")
	}

	return &cp, nil
}

func echoReplyCode(ipv IPVersion) int {
	switch ipv {
	case IPv4:
		return int(ipv4.ICMPTypeEchoReply)
	case IPv6:
		// return int(ipv6.ICMPTypeEchoReply)
		// fun fact: based on doc, I expected that for ipv6 echo reply code should be 129...
		// but somehow it's 0 🤔
		return 0
	}

	return -1
}
