package slashquery

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

// Balancer round-robin the upstreams
func (sq *Slashquery) Balancer(name, network, port string) (net.Conn, error) {
	upstreams := sq.Upstreams[name].Servers

	// timeout defaults to 10 seconds
	timeout := sq.Upstreams[name].Timeout
	if timeout <= 0 {
		timeout = 10
	}

	if sq.Debug() {
		log.Printf("Upstreams: %s, timeout: %v\n", upstreams, timeout)
	}

	// endpoints contain the IP's from the servers
	var endpoints []string

	// fill the endpoints
	for _, upstream := range upstreams {
		servers := sq.Servers[upstream]
		if time.Since(servers.Expire) > 0 {
			go sq.ResolveUpstream(upstream)
		}
		endpoints = append(endpoints, servers.Addresses...)
	}

	if sq.Debug() {
		log.Printf("Upstream: %q, endpoints: %s\n", name, endpoints)
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

		if sq.Debug() {
			log.Printf("Upstream: %q, using endpoint: %s\n", name, endpoint)
		}

		// Try to connect
		conn, err := net.DialTimeout(network,
			fmt.Sprintf("%s:%s", endpoint, port),
			time.Duration(timeout)*time.Second,
		)
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
