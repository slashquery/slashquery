package slashquery

//go:generate go run genroutes.go -f examples/slashquery.yml

import (
	"time"

	"github.com/slashquery/resolver"
)

// Slashquery structure
type Slashquery struct {
	Config    map[string]string
	Plugins   map[string][]string
	Resolver  *resolver.Resolver
	Routes    map[string]*Route
	Servers   map[string]*Servers
	Upstreams map[string]*Upstream
}

// Route define an upstream
type Route struct {
	// URL of the upstream, if set, Scheme, Host, Path and Upstream
	// will be set from it,  the upstream will named has the host
	// after parsing the URL
	URL string

	// Scheme http or https
	Scheme string

	// Host hostname to use when doing the request
	Host string

	// Path to use when doing the request
	Path string

	// Upstream identifier to use
	Upstream string

	// Methods list of allowed methods, example: GET, POST, HEAD
	Methods []string

	// Plugins list of plugins to use (middleware)
	Plugins []string

	// Insecure is set to yes will skip the certificate verification
	Insecure bool

	// DisableKeepAlive if set to yes won't cache/reuse the connections
	DisableKeepAlive bool

	// rawQuery encoded query values, without '?'
	rawQuery string
}

// Upstream structure
type Upstream struct {
	Servers []string
	Timeout int
}

// Servers keep IP's from upstreams (needs a resolver)
type Servers struct {
	Addresses []string
	Expire    time.Time
	last      string
}
