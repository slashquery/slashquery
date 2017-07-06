package slashquery

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/nbari/violetear"
	"github.com/slashquery/resolver"
)

func TestProxy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("slashquery /?"))
	}))
	defer ts.Close()

	r, err := resolver.New("4.2.2.2")
	if err != nil {
		t.Fatal(err)
	}

	routes := make(map[string]*Route)
	routes["test"] = &Route{
		URL: ts.URL,
	}
	sq := &Slashquery{
		Config:    make(map[string]string),
		Resolver:  r,
		Routes:    routes,
		Servers:   make(map[string]*Servers),
		Upstreams: make(map[string]*Upstream),
	}
	if err := sq.Setup(); err != nil {
		t.Error(err)
	}
	sq.ResolveUpstreams()
	testRoute := sq.Routes["test"]
	expect(t, testRoute.URL, ts.URL)
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}
	testUpstream, ok := sq.Upstreams[u.Host]
	if !ok {
		t.Errorf("Expecting upstream: %v", testUpstream)
	}
	testServers, ok := sq.Servers[u.Host]
	if !ok {
		t.Errorf("Expecting servers: %s", testServers)
	}
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		t.Error(err)
	}
	expect(t, testServers.Addresses[0], host)

	// test router
	router := violetear.New()
	router.Verbose = false
	router.LogRequests = false
	router.Handle("/*", sq.Proxy("test"), "GET,HEAD")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	expect(t, w.Body.String(), "slashquery /?")
}

func TestProxyQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s", r.URL)))
	}))
	defer ts.Close()

	r, err := resolver.New("4.2.2.2")
	if err != nil {
		t.Fatal(err)
	}

	routes := make(map[string]*Route)
	routes["test"] = &Route{
		URL: fmt.Sprintf("%s/test?slash=query", ts.URL),
	}
	sq := &Slashquery{
		Config:    make(map[string]string),
		Resolver:  r,
		Routes:    routes,
		Servers:   make(map[string]*Servers),
		Upstreams: make(map[string]*Upstream),
	}
	if err := sq.Setup(); err != nil {
		t.Error(err)
	}
	sq.ResolveUpstreams()
	testRoute := sq.Routes["test"]
	expect(t, testRoute.rawQuery, "slash=query")
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}
	expect(t, fmt.Sprintf("http://%s", u.Host), ts.URL)
	testUpstream, ok := sq.Upstreams[u.Host]
	if !ok {
		t.Errorf("Expecting upstream: %v", testUpstream)
	}
	testServers, ok := sq.Servers[u.Host]
	if !ok {
		t.Errorf("Expecting servers: %s", testServers)
	}
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		t.Error(err)
	}
	expect(t, testServers.Addresses[0], host)

	// test router
	router := violetear.New()
	router.Verbose = false
	router.LogRequests = false
	router.Handle("/*", sq.Proxy("test"), "GET,HEAD")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?foo=bar", nil)
	router.ServeHTTP(w, req)

	expect(t, w.Body.String(), "/test?slash=query&foo=bar")
}
