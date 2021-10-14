COMMIT_SHA_SHORT ?= $(shell git rev-parse --short=12 HEAD)
PWD_DIR := ${CURDIR}

help: ## help command
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

status: ## get info about the project
	@echo "current commit: ${COMMIT_SHA_SHORT}"

test: ## run go tests
	@go test ./... -v -cover

performance: ## run performance tests
	@go test ./... -v -performance


fmt: ## format go code and run mod tidy
	@go fmt ./...
	@go mod tidy