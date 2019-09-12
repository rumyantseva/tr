package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/rumyantseva/tr/pkg/trace"
)

var (
	BuildTime = "unset"
	Commit    = "unset"
	Release   = "unset"
)

func init() {
	flag.Usage = func() {
		flag.PrintDefaults()
	}
}

func main() {
	log.Printf("Build time: %s, Commit: %s, Release: %s", BuildTime, Commit, Release)

	ipvF := flag.String("ipVersion", "", "IP version (4 or 6)")
	maxTTLF := flag.Int("maxTTL", 64, "Max number of hops used in outgoing probes")
	timeoutF := flag.Duration("timeout", 3*time.Second, "Max time of probe")

	flag.Parse()

	host := strings.ToLower(flag.Arg(0))

	if host == "" {
		log.Fatal("Target host is not set")
	}

	var ipv trace.IPVersion
	switch strings.ToLower(*ipvF) {
	case "":
		log.Printf("The ipVersion flag is not set, ipv4 will be used by default.")
		ipv = trace.IPv4
	case "4", "v4", "ip4", "ipv4":
		ipv = trace.IPv4
	case "6", "v6", "ip6", "ipv6":
		ipv = trace.IPv6
	default:
		log.Fatalf("The ipVersion flag value (`%s`) is incorrect. Please provide the IP version.", *ipvF)
	}
	log.Printf("Looking for the route to `%s` (%s)", host, ipv)

	var prevLatency, maxDifference int64
	printHop := func(hop *trace.Hop, err error) {
		if err != nil {
			log.Printf("  *\t*\t*\t%+v\n", err)
			prevLatency = 0
		} else {
			log.Printf("%3d %s (%s) %v\n", hop.TTL, hop.Peer.Name, hop.Peer.Addr, hop.Latency)

			if prevLatency != 0 {
				diff := hop.Latency.Nanoseconds() - prevLatency
				if diff < 0 {
					diff = -diff
				}
				if diff > maxDifference {
					maxDifference = diff
				}
			}
			prevLatency = hop.Latency.Nanoseconds()

		}
	}

	trace.Build(host, ipv, *maxTTLF, *timeoutF, printHop)

	if maxDifference == 0 {
		log.Printf("Couldn't calculate the largest difference in response time.")
	} else {
		log.Printf("The largest difference in response time is %s.", time.Duration(maxDifference))
	}
}
