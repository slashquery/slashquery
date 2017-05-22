package slashquery

import (
	"time"

	"github.com/slashquery/resolver"
)

type Slashquery struct {
	Config    map[string]string
	Routes    map[string]Route
	Upstreams map[string]Upstream
	Servers   map[string]Servers
	Resolver  *resolver.Resolver
}

type Route struct {
	Path     string
	Upstream string
	Plugins  []string
}

type Upstream struct {
	Servers []string
}

// Servers keep IP's from upstreams (needs a resolver)
type Servers struct {
	Addresses []string
	Expire    time.Time
	last      string
}
