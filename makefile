# Common variables
VERSION := 0.0.2
BUILD_INFO := Manual build from makefile
#SRC_DIR := cmd
SRC_DIR := ./cmd

# Most likely want to override these when calling `make image`
IMAGE_REG ?= ghcr.io
IMAGE_REPO ?= benc-uk/food-truck
IMAGE_TAG ?= latest
IMAGE_PREFIX := $(IMAGE_REG)/$(IMAGE_REPO)

.PHONY: help image push build run lint lint-fix install-tools
.DEFAULT_GOAL := help

SWAGGER_PATH := ./bin/swagger
AIR_PATH := ./bin/air
GOLINT_PATH := ./bin/golangci-lint

help: ## 💬 This help message :)
	@figlet $@ || true
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint: ## 🌟 Lint & format, will not fix but sets exit code on error
	@figlet $@ || true
	cd $(SRC_DIR); .$(GOLINT_PATH) run --modules-download-mode=mod *.go

lint-fix: ## 🔍 Lint & format, will try to fix errors and modify code
	@figlet $@ || true
	cd $(SRC_DIR); .$(GOLINT_PATH) run --modules-download-mode=mod *.go --fix

image: ## 📦 Build container image from Dockerfile
	@figlet $@ || true
	docker build --no-cache --file ./build/Dockerfile \
	--build-arg BUILD_INFO="$(BUILD_INFO)" \
	--build-arg VERSION="$(VERSION)" \
	--tag $(IMAGE_PREFIX):$(IMAGE_TAG) . 

push: ## 📤 Push container image to registry
	@figlet $@ || true
	docker push $(IMAGE_PREFIX):$(IMAGE_TAG)

build: ## 🔨 Run a local build without a container
	@figlet $@ || true
	go mod tidy
	go build -o bin/server $(SRC_DIR)/...

run: ## 🏃 Run backend server, with hot reload, for local development
	@figlet $@ || true
	$(AIR_PATH) -c .air.toml

run-frontend: ## 💻 Run frontend, with hot reload, for local development
	@figlet $@ || true
	browser-sync start --server ./web/client --no-ui --no-open --no-notify --watch

install-tools: ## 🔮 Install dev tools
	@$(GOLINT_PATH) > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin/
	@$(AIR_PATH) -v > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh
	@$(SWAGGER_PATH) version > /dev/null 2>&1 || ./scripts/download-goswagger.sh
	npm install -g browser-sync

generate: ## 🔬 Generate Swagger / OpenAPI spec
	go generate ./cmd 
	cp $(SRC_DIR)/swagger.yaml api/spec.yaml

test: generate ## 🥽 Run unit and integration tests
	go test ./... -v -count=1
