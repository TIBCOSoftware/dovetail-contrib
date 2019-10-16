MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
APP_FILE      := equipment.json
APP_NAME      := equipment_cc
BUILD_PATH    := /tmp/equipment
FAB_PATH      := $(GOPATH)/src/github.com/hyperledger/fabric-samples
CC_DEPLOY     := $(FAB_PATH)/chaincode
SHIM_PATH     := $(SRC_PATH)/../../shim
PATCH_PATH    := $(SRC_PATH)/../../flogo-patch

CLIENT_FILE   := equipment_client.json
CLIENT_PATH   := /tmp/equipment_client
CLIENT_NAME   := equipment_client

GQL_FILE   := equipment_gql.json
GQL_PATH   := /tmp/equipment_gql
GQL_NAME   := equipment_gql

FLOGO_VER     ?= v0.9.2
DT_REPO       ?= github.com/TIBCOSoftware/dovetail-contrib
DT_BRANCH     ?= develop

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
	flogo create --cv $(FLOGO_VER) -f $(APP_FILE) $(APP_NAME)
	rm $(BUILD_PATH)/$(APP_NAME)/src/main.go
	cp $(SHIM_PATH)/chaincode_shim.go $(BUILD_PATH)/$(APP_NAME)/src/main.go
	cp $(SHIM_PATH)/dovetailimports.go $(BUILD_PATH)/$(APP_NAME)/src
	cd $(BUILD_PATH)/$(APP_NAME)/src && \
	echo "" >> go.mod && \
	echo "replace github.com/TIBCOSoftware/dovetail-contrib => $(DT_REPO) $(DT_BRANCH)" >> go.mod

.PHONY: build
build: $(BUILD_PATH)/$(APP_NAME)
	-cd $(BUILD_PATH)/$(APP_NAME) && flogo build -e
	cd $(BUILD_PATH)/$(APP_NAME) && \
	flogo install github.com/project-flogo/contrib/function/string && \
	flogo install github.com/project-flogo/contrib/function/array

	cd $(BUILD_PATH)/$(APP_NAME)/src && \
	go mod vendor && \
	cp -R $(PATCH_PATH)/* $(BUILD_PATH)/$(APP_NAME)/src/vendor/github.com/project-flogo && \
	env GOOS=linux GOARCH=amd64 go build -mod vendor -o $(APP_NAME)

.PHONY: deploy
deploy: build
	rm -Rf $(CC_DEPLOY)/$(APP_NAME)
	cp -Rf $(BUILD_PATH)/$(APP_NAME)/src $(CC_DEPLOY)/$(APP_NAME)
	cp fn-cli-init.sh $(FAB_PATH)/first-network/scripts
	cp fn-cli-test.sh $(FAB_PATH)/first-network/scripts

.PHONY: deploy-go
deploy-go:
	-mkdir -p $(CC_DEPLOY)/$(APP_NAME)
	cp chaincode.go $(CC_DEPLOY)/$(APP_NAME)
	cp fn-cli-init.sh $(FAB_PATH)/first-network/scripts
	cp fn-cli-test.sh $(FAB_PATH)/first-network/scripts

.PHONY: cli-init
cli-init:
	docker exec cli bash -c 'cd scripts; ./fn-cli-init.sh $(APP_NAME)'

.PHONY: cli-test
cli-test:
	docker exec cli bash -c 'cd scripts; ./fn-cli-test.sh $(APP_NAME)'

.PHONY: package
package:
	-rm -R $(GOPATH)/src/github.com/chaincode/$(APP_NAME)
	mkdir -p $(GOPATH)/src/github.com/chaincode
	cp -R $(BUILD_PATH)/$(APP_NAME)/src $(GOPATH)/src/github.com/chaincode/$(APP_NAME)
	fabric-tools package -n $(APP_NAME) -p $(GOPATH)/src/github.com/chaincode/$(APP_NAME)
	rm -R $(GOPATH)/src/github.com/chaincode/$(APP_NAME)
	fabric-tools metadata -f $(APP_FILE) -m $(SRC_PATH)/contract-metadata/metadata.json -g $(SRC_PATH)/contract-metadata/metadata.gql -o $(SRC_PATH)/contract-metadata/override.json

.PHONY: clean-client
clean-client:
	rm -Rf $(CLIENT_PATH)

.PHONY: create-client
create-client: $(CLIENT_FILE) clean-client
	mkdir -p $(CLIENT_PATH) && \
	cp $(CLIENT_FILE) $(CLIENT_PATH)
	cd $(CLIENT_PATH) && \
	flogo create --cv $(FLOGO_VER) -f $(CLIENT_FILE) $(CLIENT_NAME)
	cp $(SHIM_PATH)/clientappimports.go $(CLIENT_PATH)/$(CLIENT_NAME)/src
	cd $(CLIENT_PATH)/$(CLIENT_NAME)/src && \
	echo "" >> go.mod && \
	echo "replace github.com/TIBCOSoftware/dovetail-contrib => $(DT_REPO) $(DT_BRANCH)" >> go.mod

.PHONY: build-client
build-client: $(CLIENT_PATH)/$(CLIENT_NAME)
	cd $(CLIENT_PATH)/$(CLIENT_NAME)/src && \
	go get -u -d github.com/hyperledger/fabric-sdk-go@master && \
	go get -u -d github.com/go-kit/kit@v0.8.0 && \
	go get -u -d github.com/cloudflare/cfssl@1.3.2
	-cd $(CLIENT_PATH)/$(CLIENT_NAME) && \
	flogo install github.com/project-flogo/contrib/function/string && \
	flogo build -e
	cd $(CLIENT_PATH)/$(CLIENT_NAME)/src && \
	go mod vendor && \
	cp -R $(PATCH_PATH)/* $(CLIENT_PATH)/$(CLIENT_NAME)/src/vendor/github.com/project-flogo && \
	go build -mod vendor -o $(CLIENT_NAME)

.PHONY: run
run:
	FLOGO_APP_PROP_RESOLVERS=env FLOGO_APP_PROPS_ENV=auto PORT=9898 CCNAME=equipment_cc FLOGO_LOG_LEVEL=DEBUG FLOGO_SCHEMA_SUPPORT=true FLOGO_SCHEMA_VALIDATION=false CRYPTO_PATH=$(FAB_PATH)/first-network/crypto-config $(CLIENT_PATH)/$(CLIENT_NAME)/src/$(CLIENT_NAME)

.PHONY: clean-gql
clean-gql:
	rm -Rf $(GQL_PATH)

.PHONY: create-gql
create-gql: $(GQL_FILE) clean-gql
	mkdir -p $(GQL_PATH) && \
	cp $(GQL_FILE) $(GQL_PATH)
	cd $(GQL_PATH) && \
	flogo create --cv $(FLOGO_VER) -f $(GQL_FILE) $(GQL_NAME)
	cp $(SHIM_PATH)/clientappimports.go $(GQL_PATH)/$(GQL_NAME)/src
	cd $(GQL_PATH)/$(GQL_NAME)/src && \
	echo "" >> go.mod && \
	echo "replace github.com/TIBCOSoftware/dovetail-contrib => $(DT_REPO) $(DT_BRANCH)" >> go.mod

.PHONY: build-gql
build-gql: $(GQL_PATH)/$(GQL_NAME)
	cd $(GQL_PATH)/$(GQL_NAME)/src && \
	go get -u -d github.com/hyperledger/fabric-sdk-go@master && \
	go get -u -d github.com/go-kit/kit@v0.8.0 && \
	go get -u -d github.com/cloudflare/cfssl@1.3.2
	-cd $(GQL_PATH)/$(GQL_NAME) && \
	flogo install github.com/project-flogo/contrib/function/string && \
	flogo build -e
	cd $(GQL_PATH)/$(GQL_NAME)/src && \
	go mod vendor && \
	cp -R $(PATCH_PATH)/* $(GQL_PATH)/$(GQL_NAME)/src/vendor/github.com/project-flogo && \
	go build -mod vendor -o $(GQL_NAME)

.PHONY: run-gql
run-gql:
	FLOGO_APP_PROP_RESOLVERS=env FLOGO_APP_PROPS_ENV=auto PORT=9797 CCNAME=equipment_cc FLOGO_LOG_LEVEL=DEBUG FLOGO_SCHEMA_SUPPORT=true FLOGO_SCHEMA_VALIDATION=false CRYPTO_PATH=$(FAB_PATH)/first-network/crypto-config $(GQL_PATH)/$(GQL_NAME)/src/$(GQL_NAME)
