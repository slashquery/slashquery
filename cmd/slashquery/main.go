package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/nbari/violetear"
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

	// setup gateway
	if err := sq.Setup(); err != nil {
		log.Fatalln(err)
	}

	// Get upstream IP's
	sq.ResolveUpstreams()

	// create router
	router := violetear.New()
	router.Verbose = true
	router.LogRequests = true
	if sq.Config["request-id"] != "" {
		router.RequestID = sq.Config["request-id"]
	}

	// go:generate go run genroutes.go
	sq.AddRoutes(router)

	// listen on socket or address:port
	if sq.Config["socket"] != "" {
		l, err := net.Listen("unix", sq.Config["socket"])
		if err != nil {
			log.Fatalln(err)
		}
		log.Fatalln(http.Serve(l, router))
	} else {
		log.Fatalln(http.ListenAndServe(
			fmt.Sprintf("%s:%s", sq.Config["host"], sq.Config["port"]),
			router),
		)
	}
}
