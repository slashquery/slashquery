package slashquery

import "testing"

func TestSetup(t *testing.T) {
	routes := make(map[string]*Route)
	routes[""] = &Route{
		URL: "http://www.slashquery.org",
	}
	routes["1"] = &Route{
		URL: "https://www.slashquery.org",
	}
	routes["2"] = &Route{
		URL: "https://www.slashquery.org/path",
	}
	routes["3"] = &Route{
		URL: "https://www.slashquery.org/path?slash=query",
	}
	routes["4"] = &Route{
		Host:     "example.com",
		Upstream: "upstream",
	}
	routes["5"] = &Route{
		Host:     "example.com",
		Path:     "/get",
		Scheme:   "https",
		Upstream: "upstream",
	}
	routes["6"] = &Route{
		Host:     "example.com",
		Path:     "/path?slash=query",
		Upstream: "upstream",
	}
	sq := &Slashquery{
		Routes:    routes,
		Upstreams: make(map[string]*Upstream),
	}
	sq.Setup()
	for name, route := range sq.Routes {
		switch name {
		case "":
			expect(t, route.Scheme, "http")
			expect(t, route.Host, "www.slashquery.org")
			expect(t, route.Path, "")
			expect(t, route.rawQuery, "")
			expect(t, route.Upstream, "www.slashquery.org")
		case "1":
			expect(t, route.Scheme, "https")
			expect(t, route.Host, "www.slashquery.org")
			expect(t, route.Path, "")
			expect(t, route.rawQuery, "")
			expect(t, route.Upstream, "www.slashquery.org")
		case "2":
			expect(t, route.Scheme, "https")
			expect(t, route.Host, "www.slashquery.org")
			expect(t, route.Path, "/path")
			expect(t, route.rawQuery, "")
			expect(t, route.Upstream, "www.slashquery.org")
		case "3":
			expect(t, route.Scheme, "https")
			expect(t, route.Host, "www.slashquery.org")
			expect(t, route.Path, "/path")
			expect(t, route.rawQuery, "slash=query")
			expect(t, route.Upstream, "www.slashquery.org")
		case "4":
			expect(t, route.Host, "example.com")
			expect(t, route.Upstream, "upstream")
		case "5":
			expect(t, route.Scheme, "https")
			expect(t, route.Host, "example.com")
			expect(t, route.Path, "/get")
			expect(t, route.Upstream, "upstream")
		case "6":
			expect(t, route.Scheme, "http")
			expect(t, route.Host, "example.com")
			expect(t, route.Path, "/path")
			expect(t, route.rawQuery, "slash=query")
			expect(t, route.Upstream, "upstream")
		default:
			t.Errorf("route: %s not found", name)
		}
	}
}
