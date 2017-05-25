package slashquery

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

// Balancer round-robin the upstreams
func (sq *Slashquery) Balancer(network, port, name string) (net.Conn, error) {
	endpoints := sq.Servers[name].Addresses
	// TODO
	// refresh DNS here

	// loop until can connect to one endpoint
	for {
		// No more endpoint, stop
		if len(endpoints) == 0 {
			break
		}
		// Select a random endpoint
		i := rand.Int() % len(endpoints)
		endpoint := endpoints[i]

		// Try to connect
		conn, err := net.Dial(network, fmt.Sprintf("%s:%s", endpoint, port))
		if err != nil {
			log.Printf("Error connecting to server %q: %s\n", endpoint, err)
			// Failure: remove the endpoint from the current list and try again.
			endpoints = append(endpoints[:i], endpoints[i+1:]...)
			continue
		}

		return conn, nil
	}

	// No available endpoint.
	return nil, fmt.Errorf("No endpoint available for %s", name)
}
