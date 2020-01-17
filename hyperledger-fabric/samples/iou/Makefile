MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
APP_FILE      := iou.json
APP_NAME      := iou_cc
BUILD_PATH    ?= /tmp/iou
FAB_PATH      ?= $(GOPATH)/src/github.com/hyperledger/fabric-samples
PEER_CRYPTO   := $(FAB_PATH)/first-network/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com
CC_DEPLOY     := $(FAB_PATH)/chaincode
SHIM_PATH     := $(SRC_PATH)/../../shim
PATCH_PATH    := $(SRC_PATH)/../../flogo-patch

CLIENT_FILE   := iou_client.json
CLIENT_PATH   ?= /tmp/iou_client
CLIENT_NAME   := iou_client

FLOGO_VER     ?= v0.9.4
DT_REPO       ?= github.com/TIBCOSoftware/dovetail-contrib
FE_HOME       ?= $(HOME)/tibco/flogo28/2.8

.PHONY: all
all: create build

.PHONY: clean
clean:
	rm -Rf $(BUILD_PATH)

.PHONY: create
create: $(APP_FILE) clean
	mkdir -p $(BUILD_PATH)
	cp $(APP_FILE) $(BUILD_PATH)
	cd $(BUILD_PATH) && \
	flogo create --cv $(FLOGO_VER) -f $(APP_FILE) $(APP_NAME)
	rm $(BUILD_PATH)/$(APP_NAME)/src/main.go
	cp $(SHIM_PATH)/chaincode_shim.go $(BUILD_PATH)/$(APP_NAME)/src/main.go
	cp -Rf META-INF $(BUILD_PATH)/$(APP_NAME)/src
	cp collections_config.json $(BUILD_PATH)/$(APP_NAME)/src
	cp $(PATCH_PATH)/codegen.sh $(BUILD_PATH)/$(APP_NAME)
	cd $(BUILD_PATH)/$(APP_NAME) && ./codegen.sh $(FE_HOME) && chmod +x src/gomodedit.sh

.PHONY: build
build: $(BUILD_PATH)/$(APP_NAME)
	cd $(BUILD_PATH)/$(APP_NAME)/src && \
	go get -u -d github.com/project-flogo/flow/activity/subflow@master && \
	./gomodedit.sh

	-cd $(BUILD_PATH)/$(APP_NAME) && flogo build -e

	cd $(BUILD_PATH)/$(APP_NAME)/src && \
	go mod vendor && \
	cp -R $(PATCH_PATH)/* $(BUILD_PATH)/$(APP_NAME)/src/vendor/github.com/project-flogo && \
	env GOOS=linux GOARCH=amd64 go build -mod vendor -o ../$(APP_NAME)

.PHONY: deploy
deploy:
	rm -f $(CC_DEPLOY)/$(APP_NAME)_1.0.cds
	rm -Rf $(CC_DEPLOY)/$(APP_NAME)
	cp -Rf $(BUILD_PATH)/$(APP_NAME)/src $(CC_DEPLOY)/$(APP_NAME)
	cp fn-cli-init.sh $(FAB_PATH)/first-network/scripts/fn-init-iou.sh

.PHONY: start-fn
start-fn:
	./start-fn.sh $(FAB_PATH)
	./gen-user.sh $(FAB_PATH) EURBankAdmin 1
	./gen-user.sh $(FAB_PATH) Alice 1
	./gen-user.sh $(FAB_PATH) Bob 1
	./gen-user.sh $(FAB_PATH) USDBankAdmin 2
	./gen-user.sh $(FAB_PATH) Carol 2
	./gen-user.sh $(FAB_PATH) David 2

.PHONY: cli-init
cli-init:
	docker exec cli bash -c 'cd scripts; ./fn-init-iou.sh'

# Note: this does not work cross-platform, e.g., cds built on Mac cannot be used in linux container
.PHONY: package
package: $(BUILD_PATH)/$(APP_NAME)/src $(PEER_CRYPTO)
	mkdir -p $(GOPATH)/src/github.com/chaincode
	-rm -R $(GOPATH)/src/github.com/chaincode/$(APP_NAME)
	cp -R $(BUILD_PATH)/$(APP_NAME)/src $(GOPATH)/src/github.com/chaincode/$(APP_NAME)
	cp $(SHIM_PATH)/core.yaml $(PEER_CRYPTO)
	FABRIC_CFG_PATH=$(PEER_CRYPTO) $(FAB_PATH)/bin/peer chaincode package -n $(APP_NAME) -v 1.0 -l golang -p github.com/chaincode/$(APP_NAME) $(FAB_PATH)/chaincode/$(APP_NAME)_1.0.cds

.PHONY: metadata
metadata:
	fabric-tools metadata -f $(APP_FILE) -m $(SRC_PATH)/contract-metadata/metadata.json -g $(SRC_PATH)/contract-metadata/metadata.gql -o $(SRC_PATH)/contract-metadata/override.json
	cd $(SRC_PATH)/contract-metadata && \
	sed -i -e "s/redeem(/redeem(owner: String, /" ./metadata.gql && \
	sed -i -e "s/exchange(/exchange(owner: String, /" ./metadata.gql && \
	sed -i -e "s/transfer(/transfer(owner: String, /" ./metadata.gql && \
	sed -i -e "s/send(/send(receiver: String, /" ./metadata.gql && \
	rm metadata.gql-e

.PHONY: clean-client
clean-client:
	rm -Rf $(CLIENT_PATH)

.PHONY: create-client
create-client: $(CLIENT_FILE) clean-client
	mkdir -p $(CLIENT_PATH)
	cp $(CLIENT_FILE) $(CLIENT_PATH)
	cd $(CLIENT_PATH) && \
	flogo create --cv $(FLOGO_VER) -f $(CLIENT_FILE) $(CLIENT_NAME)
	cp $(PATCH_PATH)/codegen.sh $(CLIENT_PATH)/$(CLIENT_NAME)
	cd $(CLIENT_PATH)/$(CLIENT_NAME) && ./codegen.sh $(FE_HOME) && chmod +x src/gomodedit.sh

.PHONY: build-client
build-client: $(CLIENT_PATH)/$(CLIENT_NAME)
	cd $(CLIENT_PATH)/$(CLIENT_NAME)/src && \
	go get -u -d github.com/project-flogo/flow/activity/subflow@master
	cd $(CLIENT_PATH)/$(CLIENT_NAME) && flogo build -e

.PHONY: run
run:
	FLOGO_APP_PROP_RESOLVERS=env FLOGO_APP_PROPS_ENV=auto PORT=8989 FLOGO_LOG_LEVEL=DEBUG FLOGO_SCHEMA_SUPPORT=true FLOGO_SCHEMA_VALIDATION=false CRYPTO_PATH=$(FAB_PATH)/first-network/crypto-config $(CLIENT_PATH)/$(CLIENT_NAME)/bin/$(CLIENT_NAME)
