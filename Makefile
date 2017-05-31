.PHONY: all get test clean build cover generate

GO ?= go
GO_XC = ${GOPATH}/bin/goxc -os="freebsd netbsd openbsd darwin linux"
GOXC_FILE = .goxc.json
GOXC_FILE_LOCAL = .goxc.local.json
VERSION=$(shell git describe --tags --always)
DESTDIR ?= /usr/local

all: clean build

generate:
	# make CONFIG=testdata/slashquery.yml
	@if test -n "${CONFIG}"; then \
	${GO} run genroutes.go -f ${CONFIG}; \
	else \
	${GO} run genroutes.go -f examples/slashquery.yml; \
	fi;
	goimports -w routes.go

get:
	${GO} get

build: get generate
	${GO} get -u github.com/go-yaml/yaml;
	${GO} get -u github.com/nbari/violetear;
	${GO} get -u github.com/miekg/dns;
	${GO} build -ldflags "-s -w -X main.version=${VERSION}" -o slashquery cmd/slashquery/main.go;

clean:
	${GO} clean -i
	@rm -rf slashquery *.debug *.out build debian routes.go

test: get
	${GO} test -race -v

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
