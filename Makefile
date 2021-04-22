export GO111MODULE=on

BUILD_TAG = $(shell git describe --abbrev=0)
BUILD_HASH = $(shell git rev-parse HEAD)
BUILD_HASH_SHORT = $(shell git rev-parse --short HEAD)
LDFLAGS += -X "github.com/mattermost/mattermost-marketplace/internal/api.buildTag=$(BUILD_TAG)"
LDFLAGS += -X "github.com/mattermost/mattermost-marketplace/internal/api.buildHash=$(BUILD_HASH)"
LDFLAGS += -X "github.com/mattermost/mattermost-marketplace/internal/api.buildHashShort=$(BUILD_HASH_SHORT)"
LDFLAGS += -X "main.upstreamURL=$(BUILD_UPSTREAM_URL)"
SLS_STAGE ?= "dev"

$(shell cp plugins.json ./cmd/lambda/)

## Checks the code style, tests, builds and bundles.
all: check-style test build

## Runs go vet and golangci-lint against all packages.
.PHONY: check-style
check-style:
	go vet ./...

# https://stackoverflow.com/a/677212/1027058 (check if a command exists or not)
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint is not installed. Please see https://github.com/golangci/golangci-lint#install for installation instructions."; \
		exit 1; \
	fi; \

	golangci-lint run ./...

## Runs test against all packages.
.PHONY: test
test:
	go test -ldflags="$(LDFLAGS)" ./...

## Build builds the various commands
.PHONY: build
build: build-server build-lambda

## Compile the server for the current platform.
.PHONY: build-server
build-server:
	go build -ldflags="$(LDFLAGS)" -o dist/marketplace ./cmd/marketplace/

## Run the Plugin Marketplace
.PHONY: run
run: run-server

## Run the Plugin Marketplace
.PHONY: run-server
run-server:
	go run -ldflags="$(LDFLAGS)" ./cmd/marketplace server

## Compile the server as a lambda function
.PHONY: build-lambda
build-lambda:
	GOOS=linux go build -ldflags="-s -w $(LDFLAGS)" -o dist/marketplace-lambda ./cmd/lambda/

## Deploy the lambda stack
.PHONY: deploy-lambda
deploy-lambda: clean build-lambda
	serverless deploy --verbose --stage $(SLS_STAGE)

## Deploy the lambda function only to an existing stack
.PHONY: deploy-lambda-fast
deploy-lambda-fast: clean build-lambda
	serverless deploy function -f server --stage $(SLS_STAGE)

## Update plugins.json
.PHONY: plugins.json
plugins.json:
	@echo "This command is deprecated. Use go run ./cmd/generator/ add instead."
	go run ./cmd/generator --database plugins.json --debug

## Clean all generated files
.PHONY: clean
clean:
	rm -rf ./dist
