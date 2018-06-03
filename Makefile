# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

GCFLAGS_DEBUG=-gcflags=all="-N -l"

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

debug-server:
	@echo "Building the server to debug and running it with gdb (linux only)."

ifeq ($(PLATFORM),MSYS_NT-10.0)
	@cd "$(CURDIR)/cmd/comms" && $(GOBUILD) $(GCFLAGS_DEBUG) -o ../../bin/server.exe -v
else
	@cd "$(CURDIR)/cmd/comms" && $(GOBUILD) $(GCFLAGS_DEBUG) -o ../../bin/server -v
	@gdb ./bin/server
endif

run-server:
ifeq ($(PLATFORM),MSYS_NT-10.0)
	@bin/server.exe
else
	@./bin/server
endif

run-client:
	node client/client.js
