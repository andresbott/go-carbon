default: help;

run: ## run hlx up
	@npm run start

lint: ## run code linter
	@npm run lint

lint-fix: ## let the linter try to fix issue
	@npx eslint ./src --ext .js --fix

#	@npx stylelint blocks/**/*.css styles/*.css --fix

.PHONY: test
test: ## run tests
	@npm test


help: ## Show this help
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m·%-20s\033[0m %s\n", $$1, $$2}'