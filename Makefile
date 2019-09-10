export GO111MODULE=on

## Checks the code style, tests, builds and bundles.
all: check-style build

## Generate uses statikfs to bundle the plugin.json for use with the lambda function.
.PHONY: generate
generate:
	go get github.com/rakyll/statik
	mkdir data/static/
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
	go test ./...

## Build builds the various commands
.PHONY: build
build: build-server build-lambda

## Compile the server for the current platform.
.PHONY: build-server
build-server: generate
	go build -o dist/marketplace ./cmd/marketplace/

## Run the mattermost-marketplace
.PHONY: run-server
run-server:
	go run ./cmd/marketplace server

## Compile the server as a lambda function
.PHONY: build-lambda
build-lambda: generate
	GOOS=linux go build -ldflags="-s -w" -o dist/marketplace-lambda ./cmd/lambda/

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
