package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/nbari/violetear"
	"github.com/slashquery/slashquery"
)

var version string

func main() {
	var (
		b = flag.Bool("b", false, fmt.Sprintf("Build slashquery"))
		f = flag.String("f", "slashquery.yaml", "Configuration `slashquery.yaml`")
		v = flag.Bool("v", false, fmt.Sprintf("Print version: %s", version))
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

	if *b {
		// Getting slashquery
		if err := exec.Command("go", "get", "github.com/slashquery/slashquery").Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Getting all dependecies
		sqPath := path.Join(build.Default.GOPATH, "src", "github.com", "slashquery", "slashquery")
		cmd := exec.Command("/bin/sh", "-c", "go get -d ./...")
		cmd.Dir = sqPath
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// golang.org/x/tools/cmd/goimports
		if err := exec.Command("go", "get", "golang.org/x/tools/cmd/goimports").Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Creating routes")
		cmd = exec.Command("go", "run", "genroutes.go", "-f", *f)
		cmd.Dir = sqPath
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		goimportsPath := path.Join(build.Default.GOPATH, "bin", "goimports")
		cmd = exec.Command(goimportsPath, "-w", "routes.go")
		cmd.Dir = sqPath
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//get current version
		cmd = exec.Command("git", "describe", "--tags", "--always")
		cmd.Dir = sqPath
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		version := fmt.Sprintf("%s-%s", bytes.TrimSpace(out), time.Now().Format(time.RFC3339))
		fmt.Printf("Building slashquery: %s\n", version)
		cmd = exec.Command("go", "build", "-ldflags", fmt.Sprintf("-s -w -X main.version=%s", version), "-o", "slashquery", "cmd/slashquery/main.go")
		cmd.Dir = sqPath
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
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
