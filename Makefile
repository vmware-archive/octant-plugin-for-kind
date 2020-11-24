SHELL = /bin/bash

.PHONY: build
build:
	go build -o bin/octant-plugin-for-kind cmd/octant-plugin-for-kind/main.go
