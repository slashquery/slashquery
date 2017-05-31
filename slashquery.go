package slashquery

import (
	"time"

	"github.com/slashquery/resolver"
)

type Slashquery struct {
	Config    map[string]string
	Routes    map[string]*Route
	Upstreams map[string]*Upstream
	Servers   map[string]*Servers
	Resolver  *resolver.Resolver
}

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
	Plugins [][]string

	// Insecure is set to yes will skip the certificate verification
	Insecure bool

	// DisableKeepAlive if set to yes won't cache/reuse the connections
	DisableKeepAlive bool

	// rawQuery encoded query values, without '?'
	rawQuery string
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
