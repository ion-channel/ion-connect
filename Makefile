# Go Commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -v $(shell $(GOCMD) list ./... | grep -v /vendor/)
GOFMT=go fmt

# Optional User Provided Parameters
CGO_ENABLED ?= 0
GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
TRAVIS_BUILD_NUMBER ?= 1
BUILD_NUMBER ?= $(TRAVIS_BUILD_NUMBER)

# Calculated values for building
DATE := $(shell date -u +%Y-%m-%d%Z%H:%M:%S)
APP := $(shell basename $(PWD) | tr '[:upper:]' '[:lower:]')
BUILD_VERSION := 0.0.1-$(BUILD_NUMBER)

.PHONY: all build clean test fmt

all: test build

build: fmt ## Build the project
	$(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o $(APP) .

clean:  ## Clean out all generated files
	-@$(GOCLEAN)
	-@rm $(APP)-linux
	-@rm $(APP)-darwin
	-@rm $(APP)-windows

coverage:  ## Generates the code coverage from all the tests
	@numbers=0; sum=0; for j in $$(go test -cover $$(go list ./... | grep -v '/vendor/') 2>&1 | sed -e 's/\[no\ test\ files\]/0\.0s\ coverage:\ 0%/g' -e 's/[[:space:]]/\ /g' | tr -d "%" | cut -d ":" -f 2 | cut -d " " -f 2); do ((numbers+=1)) && sum=$$(echo $$sum + $$j | bc); done; avg=$$(echo "$$sum / $$numbers" | bc -l); printf "Total Coverage: %.1f%%\n" $$avg

coverage_compfriendly:  ## Generates the code coverage in a computer friendly manner
	@echo `make coverage | cut -d " " -f 3 | tr -d "%"`
	
crosscompile:  ## Build the binaries for the primary OS'
	GOOS=linux $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o $(APP)-linux .
	GOOS=darwin $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o $(APP)-darwin .
	GOOS=windows $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o $(APP)-windows .

dockerize:  ## Create a docker image of the project
	GOOS=linux make build
	docker build \
		--build-arg BUILD_DATE=$(DATE) \
		--build-arg VERSION=$(BUILD_VERSION) \
		-t $(APP) .

help:  ## Show This Help
	@for line in $$(cat Makefile | grep "##" | grep -v "grep" | sed  "s/:.*##/:/g" | sed "s/\ /!/g"); do verb=$$(echo $$line | cut -d ":" -f 1); desc=$$(echo $$line | cut -d ":" -f 2 | sed "s/!/\ /g"); printf "%-30s--%s\n" "$$verb" "$$desc"; done

test: unit_test ## Run all available tests

unit_test:  ## Run unit tests
	$(GOTEST)

fmt:  ## Run go fmt
	$(GOFMT)
