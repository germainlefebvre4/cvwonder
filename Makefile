#
.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


%:
    @:

run:
	go run ./cmd/cvwonder $(filter-out $@,$(MAKECMDGOALS))

build: ## Build
	go build -o cvwonder ./cmd/cvwonder

doc-install:
	cd docs && poetry lock && poetry install

doc:
	cd docs && poetry run python -m mkdocs serve
