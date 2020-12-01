SHELL = /bin/bash

VERSION := "v0.0.1"
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := $(GOBASE)/cmd/$(PROJECTNAME)/$(wildcard *.go)

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

.PHONY: build
## build: Builds a binary in /bin
build:
	@echo " > Building binary..."
	go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

## test: Runs tests
test:
	go test ./...

## lint: Runs golint
lint:
	golint ./...

## version: Prints the current version
version:
	@echo $(VERSION)

.PHONY: help
all: help
help: Makefile
	@echo
	@echo "Commands"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
