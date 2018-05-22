# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# os platform
PLATFORM := $(shell uname)

# code dir names
DIR_NAME_POD=pod
DIR_NAME_RELAY=relay
DIR_NAME_TEST=test

# binary names
BIN_NAME_POD=pod_server
BIN_NAME_RELAY=relay_server
BIN_NAME_TEST=test_client

all: get deps-update pod-server relay-server test-client

deps-update:
	@echo Updating Dependencies
	@govendor update +vendor
	@govendor sync +vendor
	@echo Dependencies Updated!

get:
	@echo Getting Required tools to build this project
	go get -u github.com/kardianos/govendor
	@echo Tools Downloaded and Built!!

pod-server:
	@echo Building $(BIN_NAME_POD) project:
ifeq ($(PLATFORM),Windows)
	@cd "$(CURDIR)/$(DIR_NAME_POD)/cmd/$(DIR_NAME_POD)" && $(GOBUILD) -o ../../../bin/$(BIN_NAME_POD).exe -v
else
	@cd "$(CURDIR)/$(DIR_NAME_POD)/cmd/$(DIR_NAME_POD)" && $(GOBUILD) -o ../../../bin/$(BIN_NAME_POD) -v
endif

	@echo Pod Server Project built!!

relay-server:
	@echo Building $(BIN_NAME_RELAY) project:
ifeq ($(PLATFORM),Windows)
	@cd "$(CURDIR)/$(DIR_NAME_RELAY)/cmd/$(DIR_NAME_RELAY)" && $(GOBUILD) -o ../../../bin/$(BIN_NAME_RELAY).exe -v
else
	@cd "$(CURDIR)/$(DIR_NAME_RELAY)/cmd/$(DIR_NAME_RELAY)" && $(GOBUILD) -o ../../../bin/$(BIN_NAME_RELAY) -v
endif

	@echo Relay Server Project built!!


test-client:
	@echo Building $(BIN_NAME_RELAY) project:
ifeq ($(PLATFORM),Windows)
	@cd "$(CURDIR)/$(DIR_NAME_TEST)/cmd/$(DIR_NAME_TEST)" && $(GOBUILD) -o ../../../bin/$(BIN_NAME_TEST).exe -v
else
	@cd "$(CURDIR)/$(DIR_NAME_TEST)/cmd/$(DIR_NAME_TEST)" && $(GOBUILD) -o ../../../bin/$(BIN_NAME_TEST) -v
endif

	@echo Test Client Project built!!

clean:
	@echo Cleaning build files:
ifeq ($(PLATFORM),Windows)
	@rmdir bin
else
	@rm -rf bin
endif

	@echo Finished cleaning!!
