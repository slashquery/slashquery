package main

import (
	"flag"
	"fmt"
	"log"
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

	sq.AddRoutes(router)

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%s", sq.Config["host"], sq.Config["port"]),
		router),
	)
}
