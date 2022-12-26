
SHELL := /bin/bash

######################
# Variables
######################

# Golang environment
GOPATH       := $(shell go env GOPATH)
GOBIN        := $(GOPATH)/bin
GOINSTALL    := go install

######################
# Tools
######################

MAGE         := $(GOBIN)/mage
YARN         := $(shell command -v yarn 2> /dev/null)
YARN_EXISTS  := $(shell [[ "$(YARN)" == "" ]] && echo "false" || echo "true")

$(MAGE):
	$(GOINSTALL) github.com/magefile/mage

$(YARN):
ifeq (false, $(YARN_EXISTS))
	$(error Missing required binary: yarn)
endif

######################
# Build
######################

.PHONY: build
build: build-node build-go

.PHONY: build-node
build-node:
	yarn install && yarn build

.PHONY: build-go
build-go: $(MAGE)
	$(MAGE)

######################
# Clean
######################

.PHONY: clean
clean:
	git clean -dfx

