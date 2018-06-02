# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# os platform
PLATFORM := $(shell uname)

setup-client:
	@cd "client" && npm install --save

build-server:
	@echo Building the server

	$(GOGET) -u github.com/gorilla/websocket
	$(GOGET) -u github.com/waterloop/wcomms/wjson
	$(GOGET) -u github.com/waterloop/wcomms/wbinary

ifeq ($(PLATFORM),MSYS_NT-10.0)
	@cd "$(CURDIR)/cmd/comms" && $(GOBUILD) -o ../../bin/server.exe -v
else
	@cd "$(CURDIR)/cmd/comms" && $(GOBUILD) -o ../../bin/server -v
endif

run-server:
ifeq ($(PLATFORM),MSYS_NT-10.0)
	@bin/server.exe
else
	@./bin/server
endif

run-client:
	node client/client.js
