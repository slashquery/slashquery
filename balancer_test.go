package slashquery

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/slashquery/resolver"
)

func TestBalancerErr(t *testing.T) {
	upstreams := make(map[string]*Upstream)
	upstreams["test"] = &Upstream{
		Timeout: 1,
		Servers: []string{"balancer"},
	}
	servers := make(map[string]*Servers)
	servers["balancer"] = &Servers{
		Addresses: []string{"127.1.2.3"},
	}
	r, err := resolver.New("4.2.2.2")
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
	_, err = sq.Balancer("test", "tcp", "80")
	if err == nil {
		t.Error("Expecting error")
	}
}

func TestBalancerConn(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}
	host, port, _ := net.SplitHostPort(u.Host)

	upstreams := make(map[string]*Upstream)
	upstreams["test"] = &Upstream{
		Servers: []string{"balancer"},
	}
	servers := make(map[string]*Servers)
	servers["balancer"] = &Servers{
		Addresses: []string{host},
	}
	r, err := resolver.New("4.2.2.2")
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
	_, err = sq.Balancer("test", "tcp", port)
	if err != nil {
		t.Error(err)
	}
}
