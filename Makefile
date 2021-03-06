.PHONY: all get test clean build cover generate

GO?=go
VERSION=$(shell git describe --tags --always)
DESTDIR?=/usr/local

all: clean build

# make CONFIG=testdata/default.yml
generate:
	@if test -n "${CONFIG}"; then \
	${GO} run genroutes.go -f ${CONFIG}; \
	else \
	${GO} run genroutes.go -f testdata/default.yml; \
	fi;
	${GOPATH}/bin/goimports -w routes.go

get:
	${GO} get -u github.com/go-yaml/yaml;
	${GO} get -u github.com/miekg/dns;
	${GO} get -u github.com/nbari/violetear;
	${GO} get -u github.com/slashquery/resolver;
	${GO} get -u golang.org/x/tools/cmd/goimports;

build: get generate
	${GO} get -u
	${GO} build -ldflags "-s -w -X main.version=${VERSION}" -o slashquery cmd/slashquery/main.go;

clean:
	${GO} clean -i
	@rm -rf slashquery *.debug *.out routes.go

test: get
	${GO} test -race -v

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
