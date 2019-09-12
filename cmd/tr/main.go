package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/rumyantseva/tr/pkg/trace"
)

// Version parameters: build time, commit hash and semantic release number.
// These values might be set via linker during compilation.
var (
	buildTime = "unset"
	commit    = "unset"
	release   = "unset"
)

func main() {
	log.Printf("Build time: %s, Commit: %s, Release: %s", buildTime, commit, release)

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
			log.Printf("%3d %s (%s) %v (%s)\n", hop.Seq, hop.Peer.Name, hop.Peer.Addr, hop.Latency, hop.Message)

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

	if err := trace.Build(host, ipv, *maxTTLF, *timeoutF, printHop); err != nil {
		log.Fatalf("The route is not found: %v", err)
	}

	if maxDifference == 0 {
		log.Printf("Couldn't calculate the largest difference in response time.")
	} else {
		log.Printf("The largest difference in response time is %s.", time.Duration(maxDifference))
	}
}
