# System Setup
SHELL = bash

# Go Stuff
GOCMD=go
GOLINTCMD=golint
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOLIST=$(GOCMD) list
GOVET=$(GOCMD) vet
GOTEST=$(GOCMD) test -v $(shell $(GOCMD) list ./... | grep -v /vendor/)
GOFMT=$(GOCMD) fmt
CGO_ENABLED ?= 0
GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')

# General Vars
APP := $(shell basename $(PWD) | tr '[:upper:]' '[:lower:]')
DATE := $(shell date -u +%Y-%m-%d%Z%H:%M:%S)
VERSION := 0.0.1
COVERAGE_DIR=coverage

# Build System Vars
TRAVIS_BUILD_NUMBER ?= 1
BUILD_NUMBER ?= $(TRAVIS_BUILD_NUMBER)
BUILD_VERSION := $(VERSION)-$(BUILD_NUMBER)

.PHONY: all
all: test build

.PHONY: analyze
analyze:  ## Perform an analysis of the project
	echo "Downloading latest Ionize"
	-@wget --quiet https://s3.amazonaws.com/public.ionchannel.io/files/ionize/linux/bin/ionize && chmod +x ionize
	./ionize analyze

.PHONY: build
build: ## Build the project
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o $(APP) .

.PHONY: clean
clean:  ## Clean out all generated files
	-@$(GOCLEAN)
	-@rm -f $(APP)-linux $(APP)-darwin $(APP)-windows
	-@rm -rf coverage
	-@rm coverage.txt

.PHONY: coverage
coverage:  ## Generates the code coverage from all the tests
	@echo "Total Coverage: $$(make --no-print-directory coverage_compfriendly | tee coverage.txt)%"

.PHONY: coverage_compfriendly
coverage_compfriendly:  ## Generates the code coverage in a computer friendly manner
	-@rm -rf coverage
	-@mkdir -p $(COVERAGE_DIR)/tmp
	@for j in $$(go list ./... | grep -v '/vendor/' | grep -v '/ext/'); do IONCHANNEL_SECRET_KEY="" go test -covermode=count -coverprofile=$(COVERAGE_DIR)/$$(basename $$j).out $$j > /dev/null 2>&1; done
	@echo 'mode: count' > $(COVERAGE_DIR)/tmp/full.out
	@tail -q -n +2 $(COVERAGE_DIR)/*.out >> $(COVERAGE_DIR)/tmp/full.out
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/tmp/full.out | tail -n 1 | sed -e 's/^.*statements)[[:space:]]*//' -e 's/%//'

.PHONY: crosscompile
crosscompile:  ## Build the binaries for the primary OS'
	GOOS=linux $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o compiled/ion-connect/linux/$(APP) .
	GOOS=darwin $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o compiled/ion-connect/darwin/$(APP) .
	GOOS=windows $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o compiled/ion-connect/windows/$(APP).exe .

.PHONY: dockerize
dockerize: clean  ## Create a docker image of the project
	CGO_ENABLED=0 GOOS=linux make build
	$(GOPATH)/bin/rice append --exec ion-connect -i ./lib
	docker build \
		--build-arg BUILD_DATE=$(DATE) \
		--build-arg VERSION=$(BUILD_VERSION) \
		-t ionchannel/$(APP):latest .

.PHONY: help
help:  ## Show This Help
	@for line in $$(cat Makefile | grep "##" | grep -v "grep" | sed  "s/:.*##/:/g" | sed "s/\ /!/g"); do verb=$$(echo $$line | cut -d ":" -f 1); desc=$$(echo $$line | cut -d ":" -f 2 | sed "s/!/\ /g"); printf "%-30s--%s\n" "$$verb" "$$desc"; done

.PHONY: test
test: unit_test ## Run all available tests

.PHONY: unit_test
unit_test:  ## Run unit tests
	$(GOTEST)

.PHONY: integration_test
integration_test:  ## Run integration tests
	cucumber -t ~@expected_failure

.PHONY: fmt
fmt: ## Run gofmt
	@echo "checking formatting..."
	@$(GOFMT) $(shell $(GOLIST) ./... | grep -v '/vendor/')

.PHONY: vet
vet: ## Run go vet
	@echo "vetting..."
	@$(GOVET) $(shell $(GOLIST) ./... | grep -v '/vendor/')

.PHONY: lint
lint: ## Run golint
	@echo "linting..."
	@$(GOLINTCMD) -set_exit_status $(shell $(GOLIST) ./... | grep -v '/vendor/')

.PHONY: travis_setup
travis_setup:  ## Setup the travis environment
	@sudo apt-get install -y rpm alien dpkg-dev debhelper build-essential jq
