package slashquery

type Slashquery struct {
	Config    map[string]string
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
