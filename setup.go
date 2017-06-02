package slashquery

import (
	"fmt"
	"net/url"
)

// Setup configures the upstream
func (sq *Slashquery) Setup() error {
	for _, route := range sq.Routes {
		if route.URL != "" {
			u, err := url.Parse(route.URL)
			if err != nil {
				return err
			}
			route.Scheme = u.Scheme
			route.Host = u.Host
			route.Path = u.Path
			route.rawQuery = u.RawQuery
			if route.Upstream == "" {
				route.Upstream = u.Host
				if _, ok := sq.Upstreams[u.Host]; !ok {
					sq.Upstreams[u.Host] = &Upstream{
						Servers: []string{u.Host},
					}
				}
			}
			continue
		}

		// set http or https
		if route.Scheme == "" {
			route.Scheme = "http"
		} else if route.Scheme != "https" {
			return fmt.Errorf("Unsuported scheme: %q", route.Scheme)
		}
		u, err := url.Parse(fmt.Sprintf("%s://%s%s", route.Scheme, route.Host, route.Path))
		if err != nil {
			return err
		}
		route.Scheme = u.Scheme
		route.Host = u.Host
		route.Path = u.Path
		route.rawQuery = u.RawQuery
	}
	return nil
}
