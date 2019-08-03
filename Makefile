# Self-Documented Makefile https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.DEFAULT_GOAL := help

integration-setup: ## build docker images for integration tests [requires Docker Compose]
	docker-compose -f tests/docker-compose.yaml build

integration-teardown: ## destroy resources associated to integration tests [requires Docker Compose]
	docker-compose -f tests/docker-compose.yaml down

integration-logs: ## show logs from nfs-server [requires Docker Compose]
	docker-compose -f tests/docker-compose.yaml logs

integration: ## run all integration tests [requires Docker Compose]
	docker-compose -f tests/docker-compose.yaml run tester /usr/bin/bats /tests

unittests: ## run all unit tests
	@go test ./...

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
