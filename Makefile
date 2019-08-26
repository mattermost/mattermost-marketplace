GO ?= $(shell command -v go 2> /dev/null)
MACHINE = $(shell uname -m)
GOFLAGS ?= $(GOFLAGS:)
BUILD_TIME := $(shell date -u +%Y%m%d.%H%M%S)
BUILD_HASH := $(shell git rev-parse HEAD)

export GO111MODULE=on

## Checks the code style, tests, builds and bundles.
all: check-style

## Runs govet and gofmt against all packages.
.PHONY: check-style
check-style: govet lint
	@echo Checking for style guide compliance

## Runs lint against all packages.
.PHONY: lint
lint:
	@echo Running lint
	env GO111MODULE=off $(GO) get -u golang.org/x/lint/golint
	golint -set_exit_status ./...
	@echo lint success

## Runs govet against all packages.
.PHONY: vet
govet:
	@echo Running govet
	$(GO) vet ./...
	@echo Govet success

## Runs test against all packages.
.PHONY: test
test:
	@echo Running test
	$(GO) test ./...

## Run the mattermost-marketplace
.PHONY: run-server
run-server:
	go run ./cmd/marketplace server

