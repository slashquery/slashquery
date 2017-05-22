package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/slashquery/resolver"
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

	r, _ := resolver.New("8.8.8.8")
	for k, v := range sq.Upstreams {
		fmt.Printf("upstream = %+v\n", k)
		fmt.Printf("servers = %+v\n", v.Servers)
		for _, v := range v.Servers {
			ans, err := r.Resolve(v)
			if err != nil {
				fmt.Printf("err = %+v\n", err)
			}
			sq.Servers[k] = slashquery.Servers{
				Addresses: ans.Addresses,
				Expire:    time.Now().Add(time.Duration(ans.TTL) * time.Second),
			}
		}
		println()
	}

	for k, v := range sq.Servers {
		fmt.Printf("k = %+v\n", k)
		fmt.Printf("v = %+v\n", v)
	}
}
