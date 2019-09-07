package trace

import "net"

type IPVersion string

const IPv4 = "ipv4"
const IPv6 = "ipv6"

func ByHost(host string, ipv IPVersion) (*Route, error) {
	hostIPs, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	return ByIP(hostIPs[0])
}

func ByIP(ip net.IP) (*Route, error) {
	return &Route{}, nil
}
