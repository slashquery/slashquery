package slashquery

import (
	"testing"

	"github.com/slashquery/resolver"
)

func TestResolveUpstreams(t *testing.T) {
	upstreams := make(map[string]*Upstream)
	upstreams["localhost"] = &Upstream{
		Servers: []string{"127.0.0.1"},
	}
	r, err := resolver.New("8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	sq := &Slashquery{
		Upstreams: upstreams,
		Servers:   make(map[string]*Servers),
		Resolver:  r,
	}

	sq.ResolveUpstreams()

	if val, ok := sq.Servers["127.0.0.1"]; !ok {
		t.Fatalf("Expecting %q", "127.0.0.1")
	} else {
		v := val.Addresses[0]
		if v != "127.0.0.1" {
			t.Fatalf("Expecting %q", "127.0.0.1")
		}
	}
}

func TestResolveUpstream(t *testing.T) {
	config := make(map[string]string)
	config["debug"] = "yes"
	upstreams := make(map[string]*Upstream)
	upstreams["localhost"] = &Upstream{
		Servers: []string{"127.0.0.1"},
	}
	r, err := resolver.New("8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	sq := &Slashquery{
		Upstreams: upstreams,
		Servers:   make(map[string]*Servers),
		Resolver:  r,
		Config:    config,
	}

	sq.ResolveUpstream("127.0.0.1")

	if val, ok := sq.Servers["127.0.0.1"]; !ok {
		t.Fatalf("Expecting %q", "127.0.0.1")
	} else {
		v := val.Addresses[0]
		if v != "127.0.0.1" {
			t.Fatalf("Expecting %q", "127.0.0.1")
		}
	}
}
