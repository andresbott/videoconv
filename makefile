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

check_env: # check for needed envs
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is undefined, create one with repo permissions here: https://github.com/settings/tokens/new?scopes=repo,write:packages)
endif

release: check_env test ## release a new version of goback
	@git diff --quiet || ( echo 'git is in dirty state' ; exit 1 )
	@[ "${version}" ] || ( echo ">> version is not set, usage: make release version=\"v1.2.3\" "; exit 1 )
	@git tag -d $(version) || true
	@git tag -a $(version) -m "Release version: $(version)"
	@git push origin $(version)
	@goreleaser --rm-dist

help: ## Show this help
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m·%-20s\033[0m %s\n", $$1, $$2}'