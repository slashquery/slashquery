package slashquery

import (
	"log"
	"time"
)

// ResolveUpstreams get IP's of servers
func (sq *Slashquery) ResolveUpstreams() {
	for _, servers := range sq.Upstreams {
		for _, server := range servers.Servers {
			ans, err := sq.Resolver.Resolve(server)
			if err != nil {
				log.Printf("Could not resolve server: %q, %s", server, err)
			} else {
				sq.Servers[server] = &Servers{
					Addresses: ans.Addresses,
					Expire:    time.Now().Add(time.Duration(ans.TTL) * time.Second),
				}
			}
		}
	}
}

// ResolveUpstream get IP's for specified upstream
func (sq *Slashquery) ResolveUpstream(upstream string) {
	if sq.Debug() {
		log.Printf("Updating IP's for upstream: %q\n", upstream)
	}
	ans, err := sq.Resolver.Resolve(upstream)
	if err != nil {
		log.Printf("Could not resolve server: %q, %s", upstream, err)
	} else {
		sq.Servers[upstream] = &Servers{
			Addresses: ans.Addresses,
			Expire:    time.Now().Add(time.Duration(ans.TTL) * time.Second),
		}
	}
}
