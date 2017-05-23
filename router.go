package slashquery

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nbari/violetear"
)

func handleTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.Path[1:])
}

func (sq *Slashquery) Router() {
	router := violetear.New()
	router.LogRequests = true
	for name := range sq.Routes {
		router.HandleFunc(fmt.Sprintf("%s/*", name), handleTest)
	}
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%s", sq.Config["host"], sq.Config["port"]),
		router),
	)
}
