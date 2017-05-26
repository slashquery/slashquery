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
	DisableKeepAlive bool
	Host             string
	Path             string
	Plugins          []string
	Scheme           string
	Upstream         string
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
