package trace

import (
	"log"
	"time"
)

func ExampleBuild() {
	cb := func(hop *Hop, err error) {
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", hop)
	}

	if err := Build("google.com", IPv6, 64, time.Second, cb); err != nil {
		log.Fatal(err)
	}

	if err := Build("google.com", IPv4, 64, time.Second, cb); err != nil {
		log.Fatal(err)
	}
}
