MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
APP_FILE      := abac_app.json
APP_NAME      := abac_cc
LOG_LEVEL     := INFO
BUILD_PATH    := /tmp/abac
CC_DEPLOY     := $(GOPATH)/src/github.com/hyperledger/fabric-samples/chaincode
SHIM_PATH     := $(dir $(SRC_PATH))shim
PATCH_PATH    := $(dir $(SRC_PATH))flogo-patch

CLIENT_FILE   := abac_client.json
CLIENT_PATH   := /tmp/abac_client
CLIENT_NAME   := abac_client

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
	rm $(BUILD_PATH)/$(APP_NAME)/src/main.go
	cp $(SHIM_PATH)/chaincode_shim.go $(BUILD_PATH)/$(APP_NAME)/src/main.go
	cp $(SHIM_PATH)/dovetailimports.go $(BUILD_PATH)/$(APP_NAME)/src
	cd $(BUILD_PATH)/$(APP_NAME)/src && \
	echo "" >> go.mod && \
	echo "replace github.com/TIBCOSoftware/dovetail-contrib => $(GOPATH)/src/github.com/TIBCOSoftware/dovetail-contrib" >> go.mod

.PHONY: build
build: $(BUILD_PATH)/$(APP_NAME)
	cd $(BUILD_PATH)/$(APP_NAME) && \
	flogo install github.com/project-flogo/contrib/function/string && \
	flogo build -e
	cd $(BUILD_PATH)/$(APP_NAME)/src && \
	env GOOS=linux GOARCH=amd64 go build -o $(APP_NAME) && \
	go mod vendor
	cp $(PATCH_PATH)/ref.go $(BUILD_PATH)/$(APP_NAME)/src/vendor/github.com/project-flogo/core/support

.PHONY: deploy
deploy: build
	rm -Rf $(CC_DEPLOY)/$(APP_NAME)
	cp -Rf $(BUILD_PATH)/$(APP_NAME)/src $(CC_DEPLOY)/$(APP_NAME)

.PHONY: package
package: 
	docker exec cli peer chaincode package -n $(APP_NAME) -v 1.0 -p github.com/chaincode/$(APP_NAME) ../../../chaincode/$(APP_NAME).cds
	mv $(CC_DEPLOY)/$(APP_NAME).cds $(SRC_PATH)

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
	FLOGO_LOG_LEVEL=DEBUG FLOGO_SCHEMA_SUPPORT=true FLOGO_SCHEMA_VALIDATION=false CRYPTO_PATH=${GOPATH}/src/github.com/hyperledger/fabric-samples/first-network/crypto-config $(CLIENT_PATH)/$(CLIENT_NAME)/bin/$(CLIENT_NAME)
