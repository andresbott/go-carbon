COMMIT_SHA_SHORT ?= $(shell git rev-parse --short=12 HEAD)
PWD_DIR := ${CURDIR}

default: help


#==========================================================================================
##@ Testing
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
##@ Running
#==========================================================================================
run: ## start the GO service
	@CARBON_LOG_LEVEL="debug" go run main.go start -c zarf/appData/config.yaml

run-ui: package-ui serve## build the UI and start the GO service

#==========================================================================================
##@ Building
#==========================================================================================
package-ui: build-ui ## build the web and copy into Go pacakge
	rm -rf ./app/spa/files/ui*
	mkdir -p ./app/spa/files/ui
	cp -r ./webui/dist/* ./app/spa/files/ui/
	touch ./app/spa/files/ui/.gitkeep
build-ui:
	@cd webui && \
	npm install && \
	export VITE_BASE="/ui" && \
	npm run build

build: package-ui ## use goreleaser to build
	@goreleaser build --auto-snapshot --clean

build-linux: package-ui ## use goreleaser to build a single linux target
	@goreleaser release --clean --auto-snapshot --skip publish
#==========================================================================================
##@  Swagger
#==========================================================================================
swagger:  ## build and serve the swagger spec
	@cd zarf/swagger && go run main.go

swagger-editor: ## run swagger editor in docker
	docker run  -p 8087:8080 -v ${PWD_DIR}/zarf/swagger:/tmp -e SWAGGER_FILE=/tmp/swagger.yaml  swaggerapi/swagger-editor


#==========================================================================================
##@   Docker
#==========================================================================================
docker-base: ## build the base docker image used to build the project
	@docker build ./ -t carbon-builder:latest -f zarf/Docker/base.Dockerfile

docker-test: docker-base ## run tests in docker
	@docker build ./ -f zarf/Docker/test.Dockerfile

docker-build: docker-base ## build a snapshot release within docker
	@docker build ./ -t carbon-build:${COMMIT_SHA_SHORT} -f zarf/Docker/build.Dockerfile
	@./zarf/Docker/dockerCP.sh carbon-build:${COMMIT_SHA_SHORT} /project/dist/ ${PWD_DIR}

.PHONY: check-git-clean
check-git-clean: # check if git repo is clen
	@git diff --quiet

.PHONY: check-branch
check-branch:
	@current_branch=$$(git symbolic-ref --short HEAD) && \
	if [ "$$current_branch" != "main" ]; then \
		echo "Error: You are on branch '$$current_branch'. Please switch to 'main'."; \
		exit 1; \
	fi

check_env: # check for needed envs
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is undefined, create one with repo permissions here: https://github.com/settings/tokens/new?scopes=repo,write:packages)
endif
	@[ "${version}" ] || ( echo ">> version is not set, usage: make release version=\"v1.2.3\" "; exit 1 )

release: check_env check-branch check-git-clean docker-test ## release a new version
	@git diff --quiet || ( echo 'git is in dirty state' ; exit 1 )
	@[ "${version}" ] || ( echo ">> version is not set, usage: make release version=\"v1.2.3\" "; exit 1 )
	@git tag -d $(version) || true
	@git tag -a $(version) -m "Release version: $(version)"
	@git push --delete origin $(version) || true
	@git push origin $(version) || true
	@GITHUB_TOKEN=${GITHUB_TOKEN} docker build -t carbon-release:${COMMIT_SHA_SHORT} --secret id=GITHUB_TOKEN ./ -f zarf/Docker/release.Dockerfile
	@ echo "using goreleaser config in zarf/.goreleaser-all.yaml"
	@./zarf/Docker/dockerCP.sh carbon-build:${COMMIT_SHA_SHORT} /project/dist/ ${PWD_DIR}

clean: ## clean build env
	@rm -rf dist

#==========================================================================================
#  Help
#==========================================================================================
.PHONY: help
help: # Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
