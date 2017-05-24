package slashquery

import (
	"fmt"
	"net/http"
)

func (sq *Slashquery) Proxy(upstream string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Upstream: %s\n", upstream)))
	})
}
