package slashquery

type Slashquery struct {
	Routes    map[string]Route
	Upstreams map[string]Upstream
}

type Route struct {
	Path     string
	Upstream string
	Plugins  []string
	Servers  []string
}

type Upstream struct {
	Servers []string
}

type Servers struct {
	host string
	port int
}
