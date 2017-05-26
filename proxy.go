package slashquery

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
)

func (sq *Slashquery) Proxy(route Route) *httputil.ReverseProxy {
	// scheme defaults to http
	scheme := route.Scheme
	if scheme == "" {
		scheme = "http"
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = route.Host
			req.URL.Host = route.Host
			req.URL.Path = route.Path
			req.URL.Scheme = scheme
		},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				_, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, fmt.Errorf("Error spliting address host:port: %s", err)
				}
				return sq.Balancer(network, port, route.Upstream)
			},
		},
	}
	return proxy
}
