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

serve: ## start the GO service
	@go run main.go start

serve-ui: package-ui serve## build the UI and start the GO service

package-ui: build-ui ## build the web and copy into Go pacakge
	rm -rf ./app/spa/files/ui*
	mkdir -p ./app/spa/files/ui
	cp -r ./webui/dist/* ./app/spa/files/ui/
build-ui:
	@cd webui && \
	npm install && \
	export VITE_BASE="/ui" && \
	npm run build

build: build-ui package-ui ## build a static binary that includes the web ui

swagger: swagger-build ## build and serve the swagger spec
	@cd zarf/swagger && go run main.go

# this uses https://goswagger.io/go-swagger/
swagger-build: ## build the swagger spec
	rm  zarf/swagger/swagger.json
	swagger generate spec main.go > zarf/swagger/swagger.json

help: ## help command
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

