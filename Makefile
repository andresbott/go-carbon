COMMIT_SHA_SHORT ?= $(shell git rev-parse --short=12 HEAD)
PWD_DIR := ${CURDIR}

default: help


test: ## run go tests
	@go test ./... -v -cover

lint: ## run go linter
	@golangci-lint run

fmt:
	@go fmt ./...
	@go mod tidy

verify: fmt test benchmark lint ## run all verification and code structure tiers

benchmark: ## run go benchmarks
	@go test -run=^$$ -bench=. ./...

run: serve ## alias to make serve
serve: ## start the GO service
	@go run main.go start -c zarf/appData/config.yaml

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

docker-builder: ## build the docker image used to build the project
	@docker build ./ -t carbon-builder -f zarf/Docker/build.Dockerfile

swagger: swagger-build ## build and serve the swagger spec
	@cd zarf/swagger && go run main.go

# this uses https://github.com/swaggo/swag
swagger-build: ## build the swagger spec
	swag fmt
	swag init -g app/router/api_v0.go -ot "json" -o zarf/swagger/


help: ## help command
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

