#
.PHONY: help run
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


%:
    @:

run:
	go run ./cmd/cvwonder $(filter-out $@,$(MAKECMDGOALS))

build:
	go build -o cvwonder ./cmd/cvwonder
	chmod +x cvwonder

test:
	go test ./...

doc-install:
	cd docs/github-pages/; \
	pnpm install;

doc-build:
	cd docs/github-pages/; \
	pnpm build;

doc-serve:
	cd docs/github-pages/; \
	pnpm start;

goreleser-check:
	goreleaser check

goreleser-release:
	goreleaser release --snapshot --clean

mockery:
	mockery

##@ Docker

DOCKER_REPO ?= germainlefebvre4/cvwonder
DOCKER_TAG ?= latest
CVWONDER_VERSION ?= 0.5.0
PLATFORMS ?= linux/amd64#,linux/arm64

docker-build: ## Build multi-arch Docker image locally (no push)
	docker buildx build \
		--platform $(PLATFORMS) \
		--build-arg CVWONDER_VERSION=$(CVWONDER_VERSION) \
		-t $(DOCKER_REPO):$(DOCKER_TAG) \
		--no-cache \
		.

docker-build-load: ## Build Docker image for current platform and load into Docker
	docker buildx build \
		--load \
		--build-arg CVWONDER_VERSION=$(CVWONDER_VERSION) \
		-t $(DOCKER_REPO):$(DOCKER_TAG) \
		.

docker-build-push: ## Build and push multi-arch Docker image
	docker buildx build \
		--platform $(PLATFORMS) \
		--build-arg CVWONDER_VERSION=$(CVWONDER_VERSION) \
		--push \
		-t $(DOCKER_REPO):$(DOCKER_TAG) \
		.

docker-run: ## Run Docker container locally
	docker run --rm -v $(PWD):/cv $(DOCKER_REPO):$(DOCKER_TAG) generate --input=cv.yml --output=generated/ --theme=default

docker-buildx-setup: ## Create and use buildx builder for multi-arch builds
	docker buildx create --name cvwonder-builder --use || docker buildx use cvwonder-builder
	docker buildx inspect --bootstrap
