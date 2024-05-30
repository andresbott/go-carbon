COMMIT_SHA_SHORT ?= $(shell git rev-parse --short=12 HEAD)
PWD_DIR := ${CURDIR}

default: help

status: ## get info about the project
	@echo "current commit: ${COMMIT_SHA_SHORT}"

fmt: ## format go code and run mod tidy
	@go fmt ./...
	@go mod tidy

test: ## run go tests
	@go test ./... -v -cover

lint: ## run go linter
	@golangci-lint run

verify: fmt test benchmark lint ## run all verification and code structure tiers

benchmark: ## run go benchmarks
	@go test -run=^$$ -bench=. ./...

serve: ## run go tests
	@go run main.go start

build-ui: ## build the web ui
	@cd webui && \
	npm install && \
	npm run build

copy-ui:
	rm -rf ./app/spa/files/ui*
	mkdir -p ./app/spa/files/ui
	cp -r ./webui/dist/* ./app/spa/files/ui/

build: build-ui copy-ui ## build a static binary that includes the web ui



help: ## help command
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

