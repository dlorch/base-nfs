# Self-Documented Makefile https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.DEFAULT_GOAL := help
SOURCES := $(shell find . -name '*.go')

$(TMPDIR).integration-setup: $(SOURCES)
	docker-compose -f tests/docker-compose.yaml down 2> /dev/null
	docker-compose -f tests/docker-compose.yaml build
	touch $@

integration-setup: $(TMPDIR).integration-setup ## build docker images for integration tests [requires Docker Compose]

integration-teardown: ## destroy resources associated to integration tests [requires Docker Compose]
	docker-compose -f tests/docker-compose.yaml down
	rm $(TMPDIR).integration-setup

integration-shell: integration-setup ## enter shell on tester [requires Docker Compose]
	docker-compose -f tests/docker-compose.yaml run tester /bin/sh

integration-logs-f: integration-setup
	docker-compose -f tests/docker-compose.yaml logs -f

integration-logs: integration-setup ## show logs from nfs-server [requires Docker Compose]
	docker-compose -f tests/docker-compose.yaml logs

integration: integration-setup ## run all integration tests [requires Docker Compose]
	docker-compose -f tests/docker-compose.yaml run tester /usr/bin/bats /tests

unittests: ## run all unit tests
	@go test ./...

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
