export GO111MODULE=on

BUILD_TAG = $(shell git describe --abbrev=0)
BUILD_HASH = $(shell git rev-parse HEAD)
BUILD_HASH_SHORT = $(shell git rev-parse --short HEAD)
LDFLAGS += -X "github.com/mattermost/mattermost-marketplace/internal/api.buildTag=$(BUILD_TAG)"
LDFLAGS += -X "github.com/mattermost/mattermost-marketplace/internal/api.buildHash=$(BUILD_HASH)"
LDFLAGS += -X "github.com/mattermost/mattermost-marketplace/internal/api.buildHashShort=$(BUILD_HASH_SHORT)"

## Checks the code style, tests, builds and bundles.
all: check-style test build

## Generate uses statikfs to bundle the plugin.json for use with the lambda function.
.PHONY: generate
generate:
	go get github.com/rakyll/statik
	mkdir -p data/static/
	cp plugins.json data/static/
	go generate ./...

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
build-server: generate
	go build -ldflags="$(LDFLAGS)" -o dist/marketplace ./cmd/marketplace/

## Run the mattermost-marketplace
.PHONY: run
run: run-server

## Run the mattermost-marketplace
.PHONY: run-server
run-server:
	go run -ldflags="$(LDFLAGS)" ./cmd/marketplace server

## Compile the server as a lambda function
.PHONY: build-lambda
build-lambda: generate
	GOOS=linux go build -ldflags="-s -w $(LDFLAGS)" -o dist/marketplace-lambda ./cmd/lambda/

## Deploy the lambda stack
.PHONY: deploy-lambda
deploy-lambda: clean build-lambda
	sls deploy --verbose

## Deploy the lambda function only to an existing stack
.PHONY: deploy-lambda-fast
deploy-lambda-fast: clean build-lambda
	sls deploy function -f server

## Clean all generated files
.PHONY: clean
clean:
	rm -rf ./dist
