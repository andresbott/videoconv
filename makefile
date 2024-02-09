COMMIT_SHA_SHORT ?= $(shell git rev-parse --short=12 HEAD)
PWD_DIR := ${CURDIR}

default: help;

fmt: ## Run go fmt on the project
	@go fmt ./...

test: fmt ## Run tests
	@golangci-lint run -v
	@go test -v ./...

build: ## Build the binary
	@mkdir -p target
	@go build -o target/videconv main.go

package: ## build installable packages
	@goreleaser release --rm-dist --skip-publish --skip-validate

help: ## Show this help
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m·%-20s\033[0m %s\n", $$1, $$2}'