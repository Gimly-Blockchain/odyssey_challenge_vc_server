#!/usr/bin/env bash

include ./scripts/env.sh

.DEFAULT_GOAL := help
.PHONY: start stop build

start: ## Start.
	DEBUG=true go run error.go resource_manager.go main.go

stop:
	lsof -i tcp:9999 | awk 'NR!=1 {print $2}' | xargs kill -9 | true;

build:
	go build -o dsb -ldflags="-s -w"


install-linters: ## Install linters
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-linters ## Run linters. Use make install-linters first.
	golangci-lint run --deadline=3m --disable-all --tests \
		-E deadcode \
		-E errcheck \
		-E staticcheck \
		-E goconst \
		-E goimports \
		-E golint \
		-E typecheck \
		-E ineffassign \
		-E maligned \
		-E misspell \
		-E nakedret \
		-E structcheck \
		-E unconvert \
		-E varcheck \
		-E govet \
		-E gosec \
		-E interfacer \
		-E staticcheck \
		-E unparam \
		-E goimports \
		-E unconvert \
		-E stylecheck \
		-E bodyclose \
		-E gosimple \
		-E unused \
		--exclude="should not use ALL_CAPS in Go names; use CamelCase instead,don't use ALL_CAPS in Go names; use CamelCase"
