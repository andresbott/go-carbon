COMMIT_SHA_SHORT ?= $(shell git rev-parse --short=12 HEAD)
PWD_DIR := ${CURDIR}

default: help

#==========================================================================================
#  Testing
#==========================================================================================
test: ## run go tests
	@go test ./... -cover

lint: ## run go linter
	# depends on https://github.com/golangci/golangci-lint
	@golangci-lint run

benchmark: ## run go benchmarks
	@go test -run=^$$ -bench=. ./...

license-check: ## check for invalid licenses
	# depends on : https://github.com/elastic/go-licence-detector
	@go list -m -mod=readonly  -json all  | go-licence-detector -includeIndirect -validate -rules zarf/allowedLicenses.json

.PHONY: verify
verify: package-ui test lint benchmark license-check ## run all tests


#==========================================================================================
#  Running
#==========================================================================================
run: serve ## alias to make serve
serve: ## start the GO service
	@go run main.go start -c zarf/appData/config.yaml

serve-ui: package-ui serve## build the UI and start the GO service

#==========================================================================================
#  Building
#==========================================================================================
package-ui: build-ui ## build the web and copy into Go pacakge
	rm -rf ./app/spa/files/ui*
	mkdir -p ./app/spa/files/ui
	cp -r ./webui/dist/* ./app/spa/files/ui/
build-ui:
	@cd webui && \
	npm install && \
	export VITE_BASE="/ui" && \
	npm run build

build: package-ui ## use goreleaser to build
	@goreleaser build --auto-snapshot --clean

#==========================================================================================
#  Swagger
#==========================================================================================
swagger: swagger-build ## build and serve the swagger spec
	@cd zarf/swagger && go run main.go

# this uses https://github.com/swaggo/swag
swagger-build: ## build the swagger spec
	swag fmt
	swag init -g app/router/api_v0.go -ot "json" -o zarf/swagger/

#==========================================================================================
#  Docker
#==========================================================================================
docker-builder-image: ## build the base docker image used to build the project
	@docker build ./ -t carbon-builder -f zarf/Docker/base.Dockerfile

docker-test: docker-builder-image ## run tests in docker
	@docker build ./ -f zarf/Docker/test.Dockerfile

docker-build: docker-builder-image ## build a snapshot release within docker
	@docker build ./ -t carbon-build:${COMMIT_SHA_SHORT} -f zarf/Docker/build.Dockerfile
	@./zarf/Docker/dockerCP.sh carbon-build:${COMMIT_SHA_SHORT} /project/dist/ ${PWD_DIR}

.PHONY: check-git-clean
check-git-clean: # check if git repo is clen
	@git diff --quiet

check_env: # check for needed envs
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is undefined, create one with repo permissions here: https://github.com/settings/tokens/new?scopes=repo,write:packages)
endif
	@[ "${version}" ] || ( echo ">> version is not set, usage: make release version=\"v1.2.3\" "; exit 1 )

#release: check_env check-git-clean docker-test ## release a new version
release: check_env check-git-clean ## release a new version
	@git diff --quiet || ( echo 'git is in dirty state' ; exit 1 )
	@[ "${version}" ] || ( echo ">> version is not set, usage: make release version=\"v1.2.3\" "; exit 1 )
	@git tag -d $(version) || true
	@git tag -a $(version) -m "Release version: $(version)"
	@git push --delete origin $(version) || true
	@git push origin $(version) || true
	@#goreleaser --rm-dist

clean: ## clean build env
	@rm -rf dist

#==========================================================================================
#  Help
#==========================================================================================
help: ## help command
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

