package slashquery

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func (sq *Slashquery) SetupProxy(route Route) (*Proxy, error) {
	p := new(Proxy)
	u, err := url.Parse(fmt.Sprintf("http://%s/%s",
		route.Host,
		route.Path))
	if err != nil {
		return nil, err
	}

	targetQuery := u.RawQuery

	p.proxy = &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = route.Host
			req.URL.Host = route.Host
			req.URL.Path = route.Path
			req.URL.Scheme = "http"
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
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
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
	return p, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
	// plugins may go here
	p.proxy.ServeHTTP(w, r)
}
