export GO111MODULE=on

BUILD_TAG = $(shell git describe --abbrev=0)
BUILD_HASH = $(shell git rev-parse HEAD)
BUILD_HASH_SHORT = $(shell git rev-parse --short HEAD)
LDFLAGS += -X "github.com/mattermost/mattermost-marketplace/internal/api.buildTag=$(BUILD_TAG)"
LDFLAGS += -X "github.com/mattermost/mattermost-marketplace/internal/api.buildHash=$(BUILD_HASH)"
LDFLAGS += -X "github.com/mattermost/mattermost-marketplace/internal/api.buildHashShort=$(BUILD_HASH_SHORT)"

## Checks the code style, tests, builds and bundles.
all: check-style build

## Generate uses statikfs to bundle the plugin.json for use with the lambda function.
.PHONY: generate
generate:
	go get github.com/rakyll/statik
	mkdir -p data/static/
	cp plugins.json data/static/
	go generate ./...

## Runs govet and gofmt against all packages.
.PHONY: check-style
check-style: govet lint

## Runs govet against all packages.
.PHONY: govet
govet:
	go vet ./...

## Runs lint against all packages.
.PHONY: lint
lint:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	golint -set_exit_status ./...

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

## Update plugins.json
.PHONY: plugins.json
plugins.json:
	go run ./cmd/generator --github-token $(GITHUB_TOKEN) --existing plugins.json --debug | jq | sponge plugins.json

## Clean all generated files
.PHONY: clean
clean:
	rm -rf ./dist
