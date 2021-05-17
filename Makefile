# System Setup
SHELL = bash

# Go Stuff
GOCMD=go
GOLINTCMD=golint
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOLIST=$(GOCMD) list
GOVET=$(GOCMD) vet
GOTEST=$(GOCMD) test -v ./...
GOFMT=$(GOCMD) fmt
CGO_ENABLED ?= 0
GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')

# General Vars
APP := $(shell basename $(PWD) | tr '[:upper:]' '[:lower:]')
DATE := $(shell date -u +%Y-%m-%d%Z%H:%M:%S)
COVERAGE_DIR=coverage

TRAVIS_BUILD_NUMBER ?= 1
TRAVIS_COMMIT ?= $(shell git rev-parse HEAD)

GIT_COMMIT_HASH ?= $(TRAVIS_COMMIT)


.PHONY: all
all: test

.PHONY: travis_setup
travis_setup: ## Setup the travis environmnet
	@if [[ -n "$$BUILD_ENV" ]] && [[ "$$BUILD_ENV" == "testing" ]]; then echo -e "$(INFO_COLOR)THIS IS EXECUTING AGAINST THE TESTING ENVIRONMEMNT$(NO_COLOR)"; fi
	@echo "Downloading latest Ionize"
	@wget --quiet http://github.com/ion-channel/ionize/releases/latest/download/ionize_linux_amd64.tar.gz -O ionize.tar.gz
	@mkdir -p $$HOME/.local/bin && tar xvf ionize.tar.gz -C $$HOME/.local/bin
	@echo "Installing Go Linter"
	@go get -u golang.org/x/lint/golint
	# needed to avoid cnflicts of versions
	@go mod vendor

.PHONY: analyze
analyze:  ## Perform an analysis of the project
	@if [[ -n "$$BUILD_ENV" ]] && [[ "$$BUILD_ENV" == "testing" ]]; then \
		IONCHANNEL_SECRET_KEY=$$TESTING_APIKEY IONCHANNEL_ENDPOINT_URL=$$TESTING_ENDPOINT_URL ionize --config .ionize.test.yaml analyze; \
	else \
		ionize analyze; \
	fi

.PHONY: deploy
deploy: ## Deploy the artifacts
	@echo "Logging into Docker Hub"
	-@docker login -u "$(DOCKER_USERNAME)" -p "$(DOCKER_PASSWORD)"
	@ext/goreleaser release

.PHONY: clean
clean:  ## Clean out all generated files
	-@$(GOCLEAN)
	-@rm -f output.txt
	-@rm -rf deploy
	-@rm -rf coverage
	-@rm -f coverage.txt
	-@rm -rf dist

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

.PHONY: help
help:  ## Show This Help
	@for line in $$(cat Makefile | grep "##" | grep -v "grep" | sed  "s/:.*##/:/g" | sed "s/\ /!/g"); do verb=$$(echo $$line | cut -d ":" -f 1); desc=$$(echo $$line | cut -d ":" -f 2 | sed "s/!/\ /g"); printf "%-30s--%s\n" "$$verb" "$$desc"; done

.PHONY: install
install:  ## Installs the binary
	$(GOCMD) install -mod=vendor

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
