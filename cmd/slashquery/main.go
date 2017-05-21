package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	for r, v := range sq.Routes {
		fmt.Printf("route = %+v\n", r)
		fmt.Printf("path = %+v\n", v.Path)
		fmt.Printf("upstream = %+v\n", v.Upstream)
		fmt.Printf("plugins= %+v\n", v.Plugins)
		println()
	}

	for r, v := range sq.Upstreams {
		fmt.Printf("upstream = %+v\n", r)
		fmt.Printf("servers = %+v\n", v.Servers)
		println()
	}
}
