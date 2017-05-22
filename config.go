package slashquery

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/slashquery/resolver"
)

// New return a new slashquery instance
func New(file string) (*Slashquery, error) {
	ymlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var s Slashquery
	if err := yaml.Unmarshal(ymlFile, &s); err != nil {
		return nil, err
	}

	// to store upstream IP's
	s.Servers = make(map[string]Servers)

	// start resolver
	nameserver := s.Config["resolver"]
	r, err := resolver.New(nameserver)
	if err != nil {
		log.Fatalf("Error starting resolver: %s", err)
	}
	s.Resolver = r
	return &s, nil
}

// Resolve get upstream IP's
func (sq *Slashquery) ResolveUpstreams() {
	for upstream, servers := range sq.Upstreams {
		for _, server := range servers.Servers {
			ans, err := sq.Resolver.Resolve(server)
			if err != nil {
				log.Printf("Could not resolve server: %q, %s", err)
			} else {
				sq.Servers[upstream] = Servers{
					Addresses: ans.Addresses,
					Expire:    time.Now().Add(time.Duration(ans.TTL) * time.Second),
				}
			}
		}
	}
}
