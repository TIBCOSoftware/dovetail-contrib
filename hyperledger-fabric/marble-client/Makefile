MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
APP_FILE      := marble_client_app.json
APP_NAME      := marble_client
LOG_LEVEL     := INFO
SHIM_PATH     := $(dir $(SRC_PATH))shim
BUILD_PATH    := /tmp/marble_client

.PHONY: all
all: create build

.PHONY: clean
clean:
	rm -Rf $(BUILD_PATH)

.PHONY: create
create: $(APP_FILE) clean
	mkdir -p $(BUILD_PATH) && \
	cp $(APP_FILE) $(BUILD_PATH)
	cd $(BUILD_PATH) && \
	flogo create -f $(APP_FILE) $(APP_NAME)
	cp $(SHIM_PATH)/clientappimports.go $(BUILD_PATH)/$(APP_NAME)/src
	cd $(BUILD_PATH)/$(APP_NAME)/src && \
	echo "" >> go.mod && \
	echo "replace github.com/TIBCOSoftware/dovetail-contrib => $(GOPATH)/src/github.com/TIBCOSoftware/dovetail-contrib" >> go.mod

.PHONY: build
build: $(BUILD_PATH)/$(APP_NAME)
	cd $(BUILD_PATH)/$(APP_NAME)/src && \
	go get -u -d github.com/hyperledger/fabric-sdk-go@master && \
	go get -u -d github.com/go-kit/kit@v0.8.0 && \
	go get -u -d github.com/cloudflare/cfssl@1.3.2
	cd $(BUILD_PATH)/$(APP_NAME) && \
	flogo install github.com/project-flogo/contrib/function/string && \
	flogo build -e

.PHONY: run
run:
	FLOGO_LOG_LEVEL=DEBUG FLOGO_SCHEMA_SUPPORT=true FLOGO_SCHEMA_VALIDATION=false CRYPTO_PATH=${GOPATH}/src/github.com/hyperledger/fabric-samples/first-network/crypto-config $(BUILD_PATH)/$(APP_NAME)/bin/$(APP_NAME)
