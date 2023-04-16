GO            ?= go
NAME           = ddb
CONFIG_PATH    = ${HOME}/.ddb
DATA_PATH      = ${CONFIG_PATH}/data
TAG           ?= 0.1.0

.DEFAULT_GOAL := help

all: dep dep-dev tidy init fmt lint test ## Runs all of the required cleaning and verification targets.
.PHONY: all

init: ## Creates the necessary directories.
	@echo "==> Creating ${DATA_PATH} dir"
	@mkdir -p ${DATA_PATH}
	@echo "==> Installing precommit hooks"
	@pre-commit install
.PHONY: init

build: ## Builds the binary.
	@echo "==> Building binary"
	@$(GO) build -o ${NAME} main.go
.PHONY: build

build-docker: ## Builds the docker image.
	@echo "==> Building docker image"
	@docker build -t github.com/danielfsousa/${NAME}:${TAG} .
.PHONY: build-docker

run: ## Runs the ddb cli.
	@$(GO) run main.go
.PHONY: run

clean: ## Runs go clean and deletes binary.
	@echo "==> Cleaning go files"
	@$(GO) clean
	@rm ${NAME}
.PHONY: clean

test: ## Runs the tests.
	@echo "==> Testing ${NAME}"
	@$(GO) test -v ./...
.PHONY: test

test-race: ## Runs the test suite with the -race flag to identify race conditions, if they exist.
	@echo "==> Testing ${NAME} (race)"
	@$(GO) test -timeout=30s -race ./.. ${TESTARGS}
.PHONY: test-race

test-cov: ## Runs the tests with coverage.
	@echo "==> Testing ${NAME} (coverage)"
	@$(GO) test ./... -coverprofile=coverage.out
.PHONY: test-cov

dep: ## Downloads the Go module.
	@echo "==> Downloading Go module"
	@$(GO) mod download
.PHONY: dep

dep-dev: ## Downloads the necessary dev dependencies.
	@echo "==> Downloading development dependencies"
	@$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	@$(GO) install github.com/mgechev/revive@latest
	@if [[ "$$(uname)" == 'Darwin' ]]; then brew install golangci-lint buf pre-commit; fi
.PHONY: dep-dev

tidy: ## Cleans the Go module.
	@echo "==> Tidying module"
	@$(GO) mod tidy
.PHONY: tidy

lint-fast: ## Lints go files with staticcheck and revive.
	@echo "==> Linting go files with staticcheck"
	@staticcheck ./...
	@echo "==> Linting go files with revive"
	@revive
.PHONY: lintfast

lint: ## Lints go files with golangci-lint and protobuf with buf.
	@echo "==> Linting protobuf files"
	@buf lint proto
	@echo "==> Linting go files"
	@golangci-lint run
.PHONY: lint

fmt: ## Formats go files with goimports.
	@echo "==> Fixing imports"
	@goimports -l -w ./
.PHONY: fmt

gen: ## Generates protobuf files.
	@echo "==> Generating files from protobuf"
	@buf generate proto
.PHONY: gen

breaking: ## Checks for breaking changes in the protobuf files.
	@echo "==> Checking protobuf breaking changes"
	@buf breaking proto --against "../../.git#subdir=proto"
.PHONY: gen

help: ## Prints help menu.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
.PHONY: help
