# System Setup
SHELL = bash

# Go Stuff
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -v $(shell $(GOCMD) list ./... | grep -v /vendor/)
GOFMT=go fmt
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
	@echo analysis tasks have not been added yet
	curl -s https://s3.amazonaws.com/public.ionchannel.io/files/scripts/travis_analyze_project.sh | bash
	curl -s https://s3.amazonaws.com/public.ionchannel.io/files/scripts/travis_add_go_coverage.sh | bash
	curl -s https://s3.amazonaws.com/public.ionchannel.io/files/scripts/travis_compliance_check.sh | bash


.PHONY: build
build: fmt ## Build the project
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o $(APP) .

.PHONY: clean
clean:  ## Clean out all generated files
	-@$(GOCLEAN)
	-@rm -f $(APP)-linux $(APP)-darwin $(APP)-windows
	-@rm -rf coverage

.PHONY: coverage
coverage:  ## Generates the code coverage from all the tests
	@echo "Total Coverage: $$(make coverage_compfriendly)%"

.PHONY: coverage_compfriendly
coverage_compfriendly:  ## Generates the code coverage in a computer friendly manner
	-@mkdir -p $(COVERAGE_DIR)
	@for j in $$(go list ./... | grep -v '/vendor/' | grep -v '/ext/'); do go test -covermode=count -coverprofile=$(COVERAGE_DIR)/$$(basename $$j).out $$j > /dev/null 2>&1; done
	@echo 'mode: count' > $(COVERAGE_DIR)/full.out
	@tail -q -n +2 $(COVERAGE_DIR)/*.out >> $(COVERAGE_DIR)/full.out
	@$(GOCMD) tool cover -func=coverage/full.out | tail -n 1 | sed -e 's/^.*statements)[[:space:]]*//' -e 's/%//'

.PHONY: crosscompile
crosscompile:  ## Build the binaries for the primary OS'
	GOOS=linux $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o $(APP)-linux .
	GOOS=darwin $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o $(APP)-darwin .
	GOOS=windows $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o $(APP)-windows .

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

.PHONY: fmt
fmt:  ## Run go fmt
	$(GOFMT)
