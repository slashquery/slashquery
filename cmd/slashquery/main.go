package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nbari/violetear"
	"github.com/nbari/violetear/middleware"
	"github.com/slashquery-plugins/waf"
	"github.com/slashquery/slashquery"
)

var version string

func main() {
	var (
		v = flag.Bool("v", false, fmt.Sprintf("Print version: %s", version))
		f = flag.String("f", "slashquery.yaml", "Configuration `slashquery.yaml`")
	)

	flag.Parse()

	if *v {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if _, err := os.Stat(*f); os.IsNotExist(err) {
		fmt.Printf("Cannot read configuration file: %s, use -h for more info.\n", *f)
		os.Exit(1)
	}

	sq, err := slashquery.New(*f)
	if err != nil {
		log.Fatalln(err)
	}

	// Get upstream IP's
	sq.ResolveUpstreams()

	// create router
	router := violetear.New()
	router.LogRequests = true
	for name, route := range sq.Routes {
		// TODO
		// prototyping plugin implementation
		if len(route.Plugins) > 0 {
			chain := middleware.New(waf.WAF)
			router.Handle(fmt.Sprintf("%s/*", name), chain.Then(sq.Proxy(route)))
		} else {
			// path, proxyHandler, methods
			if len(route.Methods) > 0 {
				router.Handle(fmt.Sprintf("%s/*", name), sq.Proxy(route), strings.Join(route.Methods, ", "))
			} else {
				router.Handle(fmt.Sprintf("%s/*", name), sq.Proxy(route))
			}
		}
	}
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%s", sq.Config["host"], sq.Config["port"]),
		router),
	)
}
