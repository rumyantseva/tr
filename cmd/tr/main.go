package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/rumyantseva/tr/pkg/trace"
)

func init() {
	flag.Usage = func() {
		flag.PrintDefaults()
	}
}

func main() {
	ipvF := flag.String("ipVersion", "", "IP version (4 or 6)")
	flag.Parse()

	host := strings.ToLower(flag.Arg(0))

	if host == "" {
		log.Fatal("Target host is not set")
	}

	var route *trace.Route
	var err error

	if ip := net.ParseIP(host); ip != nil {
		if *ipvF != "" {
			log.Printf("The given host represents an IP address, the ipVersion flag will be ignored.")
		}
		log.Printf("Looking for the route to `%s`", host)
		route, err = trace.ByIP(ip)
	} else {
		var ipv trace.IPVersion
		switch strings.ToLower(*ipvF) {
		case "":
			log.Printf("The ipVersion flag is not set, ipv4 will be used by default.")
			ipv = trace.IPv4
		case "4", "v4", "ipv4":
			ipv = trace.IPv4
		case "6", "v6", "ipv6":
			ipv = trace.IPv6
		default:
			log.Fatalf("The ipVersion flag value (`%s`) is incorrect. Please provide the IP version.", *ipvF)
		}
		log.Printf("Looking for the route to `%s` (%s)", host, ipv)
		route, err = trace.ByHost(host, ipv)
	}

	if err != nil {
		log.Fatalf("Couldn't find the route to `%s`: %s", host, err)
	}

	fmt.Print(route)
}
