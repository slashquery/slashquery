package slashquery

import (
	"testing"
	"time"

	"github.com/slashquery/resolver"
)

func TestBalancer(t *testing.T) {
	upstreams := make(map[string]*Upstream)
	upstreams["test"] = &Upstream{
		Servers: []string{"balancer"},
	}
	servers := make(map[string]*Servers)
	servers["balancer"] = &Servers{
		Addresses: []string{"127.1.2.3"},
	}
	r, err := resolver.New("8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	sq := &Slashquery{
		Upstreams: upstreams,
		Servers:   servers,
		Resolver:  r,
		Config:    make(map[string]string),
	}
	sq.Config["debug"] = "yes"
	_, err = sq.Balancer("test", "tcp", "80", time.Second)
	if err == nil {
		t.Error("Expecting error")
	}
}
