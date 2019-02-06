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
VERSION := v0.1.5
COVERAGE_DIR=coverage

TRAVIS_BUILD_NUMBER ?= 1
TRAVIS_COMMIT ?= $(shell git rev-parse HEAD)

BUILD_NUMBER ?= $(TRAVIS_BUILD_NUMBER)
BUILD_VERSION := $(VERSION)-$(BUILD_NUMBER)
GIT_COMMIT_HASH ?= $(TRAVIS_COMMIT)


.PHONY: all
all: test build

.PHONY: build
build: ## Build the project
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o $(APP) .

.PHONY: clean
clean:  ## Clean out all generated files
	-@$(GOCLEAN)
	-@rm -f output.txt
	-@rm -rf deploy
	-@rm -rf coverage
	-@rm -f coverage.txt

.PHONY: coverage
coverage:  ## Generates the code coverage from all the tests
	@echo "Total Coverage: $$(make --no-print-directory coverage_compfriendly | tee coverage.txt)%"

.PHONY: coverage_compfriendly
coverage_compfriendly:  ## Generates the code coverage in a computer friendly manner
	-@rm -rf coverage
	-@mkdir -p $(COVERAGE_DIR)/tmp
	@for j in $$(go list ./... | grep -v '/vendor/' | grep -v '/ext/'); do go test -covermode=count -coverprofile=$(COVERAGE_DIR)/$$(basename $$j).out $$j > /dev/null 2>&1; done
	@echo 'mode: count' > $(COVERAGE_DIR)/tmp/full.out
	@tail -q -n +2 $(COVERAGE_DIR)/*.out >> $(COVERAGE_DIR)/tmp/full.out
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/tmp/full.out | tail -n 1 | sed -e 's/^.*statements)[[:space:]]*//' -e 's/%//'

.PHONY: crosscompile
crosscompile:  ## Build the binaries for the primary OS'
	GOOS=linux $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o deploy/linux/bin/$(APP) .
	GOOS=darwin $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o deploy/darwin/bin/$(APP) .
	GOOS=windows $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o deploy/windows/bin/$(APP).exe .

.PHONY: help
help:  ## Show This Help
	@for line in $$(cat Makefile | grep "##" | grep -v "grep" | sed  "s/:.*##/:/g" | sed "s/\ /!/g"); do verb=$$(echo $$line | cut -d ":" -f 1); desc=$$(echo $$line | cut -d ":" -f 2 | sed "s/!/\ /g"); printf "%-30s--%s\n" "$$verb" "$$desc"; done

.PHONY: install
install:  ## Installs the binary
	$(GOCMD) install

.PHONY: test
test: unit_test  ## Runs all available tests

.PHONY: unit_test
unit_test:  ## Run unit tests
	$(GOTEST)

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
