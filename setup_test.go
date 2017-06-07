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
	sq := &Slashquery{
		Routes: routes,
	}
}
