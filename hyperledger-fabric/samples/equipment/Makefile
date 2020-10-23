MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
APP_FILE      := equipment.json
APP_NAME      := equipment_cc
WORK_PATH     ?= $(HOME)/work/dovetail
TOOL_PATH     := $(SRC_PATH)/../../dovetail-tools
FAB_PATH      ?= $(WORK_PATH)/fabric-samples
CC_DEPLOY     := $(FAB_PATH)/chaincode

CLIENT_FILE   := equipment_client.json
CLIENT_NAME   := equipment
CLIENT_OS     ?= darwin

GQL_FILE      := equipment_gql.json

.PHONY: all
all: build deploy

.PHONY: clean
clean:
	rm -Rf $(TOOL_PATH)/work/$(APP_NAME)

.PHONY: build
build: $(APP_FILE) clean
	$(TOOL_PATH)/build.sh cds -f $(APP_FILE) -n $(APP_NAME)

.PHONY: deploy
deploy: $(TOOL_PATH)/work/$(APP_HAME)
	-rm -rf $(CC_DEPLOY)/$(APP_NAME)
	cp -Rf $(TOOL_PATH)/work/$(APP_NAME)/$(APP_NAME) $(CC_DEPLOY)
	cp $(TOOL_PATH)/work/$(APP_NAME)/$(APP_NAME)_linux_amd64 $(CC_DEPLOY)/$(APP_NAME)
	cp $(TOOL_PATH)/work/$(APP_NAME)/$(APP_NAME)_1.0.cds $(CC_DEPLOY)/$(APP_NAME)
	cp fn-cli-init.sh $(FAB_PATH)/first-network/scripts/fn-init-equip.sh
	cp fn-cli-test.sh $(FAB_PATH)/first-network/scripts/fn-test-equip.sh

.PHONY: deploy-go
deploy-go:
	-mkdir -p $(CC_DEPLOY)/$(APP_NAME)
	cp chaincode.go $(CC_DEPLOY)/$(APP_NAME)
	cp fn-cli-init.sh $(FAB_PATH)/first-network/scripts/fn-init-equip.sh
	cp fn-cli-test.sh $(FAB_PATH)/first-network/scripts/fn-test-equip.sh

.PHONY: cli-init
cli-init:
	docker exec cli bash -c 'cd scripts; ./fn-init-equip.sh $(APP_NAME)'

.PHONY: cli-test
cli-test:
	docker exec cli bash -c 'cd scripts; ./fn-test-equip.sh $(APP_NAME)'

.PHONY: metadata
metadata:
	fabric-tools metadata -f $(APP_FILE) -m $(SRC_PATH)/contract-metadata/metadata.json -g $(SRC_PATH)/contract-metadata/metadata.gql -o $(SRC_PATH)/contract-metadata/override.json

.PHONY: clean-client
clean-client:
	rm -Rf $(TOOL_PATH)/work/$(CLIENT_NAME)

.PHONY: build-client
build-client: $(CLIENT_FILE) clean-client
	$(TOOL_PATH)/build.sh client -f $(CLIENT_FILE) -n $(CLIENT_NAME) -s $(CLIENT_OS)

.PHONY: run
run:
	FLOGO_APP_PROP_RESOLVERS=env FLOGO_APP_PROPS_ENV=auto PORT=9898 CCNAME=equipment_cc FLOGO_LOG_LEVEL=DEBUG FLOGO_SCHEMA_SUPPORT=true FLOGO_SCHEMA_VALIDATION=false CRYPTO_PATH=$(FAB_PATH)/first-network/crypto-config $(TOOL_PATH)/work/$(CLIENT_NAME)/$(CLIENT_NAME)_$(CLIENT_OS)_amd64

.PHONY: build-gql
build-gql: $(GQL_FILE) clean-client
	$(TOOL_PATH)/build.sh client -f $(GQL_FILE) -n $(CLIENT_NAME) -s $(CLIENT_OS)

.PHONY: start
start:
	cd $(FAB_PATH)/first-network && ./byfn.sh up -a -n -s couchdb -i 1.4.9

.PHONY: shutdown
shutdown:
	cd $(FAB_PATH)/first-network && ./byfn.sh down
	-docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
	-docker rmi $(docker images | grep dev-peer | awk '{print $3}')
