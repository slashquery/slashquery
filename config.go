package slashquery

import (
	"io/ioutil"
	"log"

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
	s.Servers = make(map[string]*Servers)

	// Upstreams
	if s.Upstreams == nil {
		s.Upstreams = make(map[string]*Upstream)
	}

	// start resolver
	nameserver := s.Config["resolver"]
	r, err := resolver.New(nameserver)
	if err != nil {
		log.Fatalf("Error starting resolver: %s", err)
	}
	s.Resolver = r
	return &s, nil
}

// Debug enable/disable based on config["debug"] value
func (sq *Slashquery) Debug() bool {
	yes := []string{"y", "Y", "yes", "Yes", "YES", "true", "True", "TRUE", "on", "On", "ON"}
	debug := sq.Config["debug"]
	for _, item := range yes {
		if item == debug {
			return true
		}
	}
	return false
}
