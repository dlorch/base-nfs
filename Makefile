# Self-Documented Makefile https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.DEFAULT_GOAL := help

integration-setup: ## build docker images for integration tests [requires Docker Compose]
	@docker-compose -f tests/docker-compose.yaml build

integration: ## run all integration tests [requires Docker Compose]
	@docker-compose -f tests/docker-compose.yaml run tester /usr/bin/bats -p /tests
	@docker-compose -f tests/docker-compose.yaml down 2> /dev/null

unittests: ## run all unit tests
	@go test ./...

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
