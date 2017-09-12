package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"go/build"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/nbari/violetear"
	"github.com/slashquery/slashquery"
	"golang.org/x/crypto/acme/autocert"
)

var version string

func main() {
	var (
		b = flag.Bool("b", false, fmt.Sprintf("Build slashquery"))
		f = flag.String("f", "slashquery.yml", "Configuration `file`")
		i = flag.Bool("i", false, fmt.Sprintf("Install slashquery in /usr/local/bin, need to use option -b"))
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
		if err := Build(*f, *i); err != nil {
			fmt.Printf("Error while building: %s\n", err)
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
		os.Remove(sq.Config["socket"])
		l, err := net.Listen("unix", sq.Config["socket"])
		if err != nil {
			log.Fatalln(err)
		}
		log.Fatalln(http.Serve(l, router))
	}
	if sq.Config["hostwhitelist"] != "" {
		domains := strings.Fields(sq.Config["hostwhitelist"])
		cache := "/tmp/certs"
		if val, ok := sq.Config["certcache"]; ok {
			cache = val
		}
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domains...),
			Cache:      autocert.DirCache(cache)}
		s := &http.Server{
			Addr:      ":https",
			Handler:   router,
			TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		}
		log.Fatalln(s.ListenAndServeTLS("", ""))
	}
	log.Fatalln(http.ListenAndServe(
		fmt.Sprintf("%s:%s", sq.Config["host"], sq.Config["port"]),
		router),
	)
}

// Build create slashquery from custom plugins
func Build(config string, install bool) error {
	// Getting slashquery
	if err := exec.Command("go", "get", "github.com/slashquery/slashquery").Run(); err != nil {
		return err
	}
	// Getting all dependecies
	sqPath := path.Join(build.Default.GOPATH, "src", "github.com", "slashquery", "slashquery")
	cmd := exec.Command("go", "get", "-d", "./...")
	cmd.Dir = sqPath
	if err := cmd.Run(); err != nil {
		return err
	}
	// golang.org/x/tools/cmd/goimports
	if err := exec.Command("go", "get", "golang.org/x/tools/cmd/goimports").Run(); err != nil {
		return err
	}
	fmt.Println("Creating routes")
	cmd = exec.Command("go", "run", "genroutes.go", "-f", config)
	cmd.Dir = sqPath
	if err := cmd.Run(); err != nil {
		return err
	}
	goimportsPath := path.Join(build.Default.GOPATH, "bin", "goimports")
	cmd = exec.Command(goimportsPath, "-w", "routes.go")
	cmd.Dir = sqPath
	if err := cmd.Run(); err != nil {
		return err
	}
	//get current version
	cmd = exec.Command("git", "describe", "--tags", "--always")
	cmd.Dir = sqPath
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	version := fmt.Sprintf("%s-%s", bytes.TrimSpace(out), time.Now().Format(time.RFC3339))
	fmt.Printf("Building slashquery: %s\n", version)
	cmd = exec.Command("go", "build", "-ldflags", "-s -w -X main.version="+version, "-o", "slashquery", "cmd/slashquery/main.go")
	cmd.Dir = sqPath
	if err := cmd.Run(); err != nil {
		return err
	}
	if install {
		fmt.Println("Installing into /usr/local/bin/slashquery")
		cmd = exec.Command("install", "slashquery", "/usr/local/bin")
		cmd.Dir = sqPath
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
