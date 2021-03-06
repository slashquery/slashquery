package slashquery

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// Proxy return instances of httputil.ReverseProxy
func (sq *Slashquery) Proxy(r string) *httputil.ReverseProxy {
	route := sq.Routes[r]
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			if strings.HasSuffix(route.Path, "*") {
				// path with out *
				p := strings.TrimSuffix(route.Path, "*")
				s := strings.TrimPrefix(req.URL.Path, p)
				req.URL.Path = p + s
			} else {
				req.URL.Path = route.Path
			}
			req.Host = route.Host
			req.URL.Host = route.Host
			req.URL.Scheme = route.Scheme
			if route.rawQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = route.rawQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = route.rawQuery + "&" + req.URL.RawQuery
			}
		},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				_, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, fmt.Errorf("Error getting port from address %q: %s", addr, err)
				}
				return sq.Balancer(route.Upstream, network, port)
			},
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   30 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     route.DisableKeepAlive,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: route.Insecure},
		},
	}
	return proxy
}
