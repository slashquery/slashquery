package slashquery

type Slashquery struct {
	Config    map[string]string
	Routes    map[string]Route
	Upstreams map[string]Upstream
	Servers   map[string]Servers
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
	IPs  []string
	last string
}
