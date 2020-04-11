export GO111MODULE=on

GOTEST ?= go test

all: build ## Runs a build task.

build: ## Build a procon binary.
	go build -trimpath -ldflags "-s -w" -o bin/procon ./src

manifest: build ## Update a manifest.
	./bin/procon -manifest procon -location ./plugin/procon.vim

test: ## Runs tests.
	${GOTEST} -v -race ./src/...

