package slashquery

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

// Balancer round-robin the upstreams
func (sq *Slashquery) Balancer(name, network, port string, timeout time.Duration) (net.Conn, error) {
	upstreams := sq.Upstreams[name].Servers
	if sq.Debug() {
		log.Printf("Upstreams: %s\n", upstreams)
	}

	// endpoints contain the IP's from the servers
	var endpoints []string

	// fill the endpoints
	for _, upstream := range upstreams {
		servers := sq.Servers[upstream]
		if time.Since(servers.Expire) > 0 {
			// TODO check for race conditions
			go sq.ResolveUpstream(upstream)
		}
		endpoints = append(endpoints, servers.Addresses...)
	}

	if sq.Debug() {
		log.Printf("Endpoints: %s\n", endpoints)
	}

	// loop until can connect to one endpoint
	for {
		// No more endpoint, stop
		if len(endpoints) == 0 {
			break
		}

		// Select a random endpoint
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(len(endpoints))
		endpoint := endpoints[i]

		// Try to connect
		conn, err := net.DialTimeout(network, fmt.Sprintf("%s:%s", endpoint, port), timeout)
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
