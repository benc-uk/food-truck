# Common variables
VERSION := 0.0.7
BUILD_INFO := Manual build from makefile
SRC_DIR := ./cmd

# Most likely want to override these when calling `make image`
IMAGE_REG ?= ghcr.io
IMAGE_REPO ?= benc-uk/food-truck
IMAGE_TAG ?= latest
IMAGE_PREFIX := $(IMAGE_REG)/$(IMAGE_REPO)

# Don't change these
SWAGGER_PATH := ./bin/swagger
AIR_PATH := ./bin/air
GOLINT_PATH := ./bin/golangci-lint

# Used for deployment only
AZURE_DEPLOY_NAME ?= food-truck
AZURE_REGION ?= westeurope
AZURE_IMAGE ?= ghcr.io/benc-uk/food-truck:latest

.PHONY: help image push build run run-frontend generate test lint lint-fix install-tools deploy
.DEFAULT_GOAL := help
.EXPORT_ALL_VARIABLES:

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

install-tools: ## 🔮 Install dev tools into project bin directory
	@figlet $@ || true
	@$(GOLINT_PATH) > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin/
	@$(AIR_PATH) -v > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh
	@$(SWAGGER_PATH) version > /dev/null 2>&1 || ./scripts/download-goswagger.sh

generate: ## 🔬 Generate Swagger / OpenAPI spec
	@figlet $@ || true
	go generate ./cmd 
	cp $(SRC_DIR)/swagger.yaml api/spec.yaml

test: ## 🧪 Run unit and integration tests
	@figlet $@ || true
	go test ./... -v -count=1

test-perf: ## 📈 Run performance tests
	@figlet $@ || true
	@tests/run.sh

deploy: ## 🚀 Deploy to Azure using Bicep & Azure CLI
	@figlet $@ || true
	@./deploy/deploy.sh