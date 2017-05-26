package slashquery

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func (sq *Slashquery) Proxy(route Route) *httputil.ReverseProxy {
	// scheme defaults to http
	scheme := route.Scheme
	if scheme == "" {
		scheme = "http"
	}

	u, _ := url.Parse(fmt.Sprintf("%s://%s%s", scheme, route.Host, route.Path))
	targetQuery := u.RawQuery

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = route.Host
			req.URL.Host = u.Host
			req.URL.Path = u.Path
			req.URL.Scheme = u.Scheme
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
