MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
FAB_PATH      := $(GOPATH)/src/github.com/hyperledger/fabric-samples
SHIM_PATH     := $(SRC_PATH)/../../shim

CLIENT_FILE   := event_listener.json
CLIENT_PATH   := /tmp/event_listener
CLIENT_NAME   := event_listener

.PHONY: all
all: create-client build-client

.PHONY: clean-client
clean-client:
	rm -Rf $(CLIENT_PATH)

.PHONY: create-client
create-client: $(CLIENT_FILE) clean-client
	mkdir -p $(CLIENT_PATH) && \
	cp $(CLIENT_FILE) $(CLIENT_PATH)
	cd $(CLIENT_PATH) && \
	flogo create -f $(CLIENT_FILE) $(CLIENT_NAME)
	cp $(SHIM_PATH)/clientappimports.go $(CLIENT_PATH)/$(CLIENT_NAME)/src
	cd $(CLIENT_PATH)/$(CLIENT_NAME)/src && \
	echo "" >> go.mod && \
	echo "replace github.com/TIBCOSoftware/dovetail-contrib => $(GOPATH)/src/github.com/TIBCOSoftware/dovetail-contrib" >> go.mod

.PHONY: build-client
build-client: $(CLIENT_PATH)/$(CLIENT_NAME)
	cd $(CLIENT_PATH)/$(CLIENT_NAME)/src && \
	go get -u -d github.com/hyperledger/fabric-sdk-go@master && \
	go get -u -d github.com/go-kit/kit@v0.8.0 && \
	go get -u -d github.com/cloudflare/cfssl@1.3.2
	cd $(CLIENT_PATH)/$(CLIENT_NAME) && \
	flogo install github.com/project-flogo/contrib/function/string && \
	flogo build -e

.PHONY: run
run:
	FLOGO_LOG_LEVEL=DEBUG FLOGO_SCHEMA_SUPPORT=true FLOGO_SCHEMA_VALIDATION=false CRYPTO_PATH=$(FAB_PATH)/first-network/crypto-config $(CLIENT_PATH)/$(CLIENT_NAME)/bin/$(CLIENT_NAME)
