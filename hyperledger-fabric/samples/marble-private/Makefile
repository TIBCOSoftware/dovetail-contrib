MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
APP_FILE      := marble_private.json
APP_NAME      := marble_private_cc
TOOL_PATH     := $(SRC_PATH)/../../dovetail-tools
FAB_PATH      ?= $(GOPATH)/src/github.com/hyperledger/fabric-samples
CC_DEPLOY     := $(FAB_PATH)/chaincode

CLIENT_FILE   := marble_private_client.json
CLIENT_NAME   := marble_private
CLIENT_OS     ?= darwin

.PHONY: all
all: build deploy

.PHONY: clean
clean:
	rm -Rf $(TOOL_PATH)/work/$(APP_NAME)
	rm -f $(TOOL_PATH)/work/$(APP_NAME)_1.0.cds

.PHONY: build
build: $(APP_FILE) clean
	$(TOOL_PATH)/build.sh cds -f $(APP_FILE) -n $(APP_NAME)

.PHONY: deploy
deploy:
	mkdir -p $(CC_DEPLOY)/$(APP_NAME)
	cp $(TOOL_PATH)/work/$(APP_NAME)_1.0.cds $(CC_DEPLOY)/$(APP_NAME)
	cp fn-cli-init.sh $(FAB_PATH)/first-network/scripts/fn-init-marble-private.sh
	cp fn-cli-test.sh $(FAB_PATH)/first-network/scripts/fn-test-marble-private.sh

.PHONY: cli-init
cli-init:
	docker exec cli bash -c 'cd scripts; ./fn-init-marble-private.sh'

.PHONY: cli-test
cli-test:
	docker exec cli bash -c 'cd scripts; ./fn-test-marble-private.sh'

.PHONY: metadata
metadata:
	fabric-tools metadata -f $(APP_FILE) -m $(SRC_PATH)/contract-metadata/metadata.json -g $(SRC_PATH)/contract-metadata/metadata.gql -o $(SRC_PATH)/contract-metadata/override.json

.PHONY: clean-client
clean-client:
	rm -Rf $(TOOL_PATH)/work/$(CLIENT_NAME)
	rm -f $(TOOL_PATH)/work/$(CLIENT_NAME)_$(CLIENT_OS)_amd64

.PHONY: build-client
build-client: $(CLIENT_FILE) clean-client
	$(TOOL_PATH)/build.sh client -f $(CLIENT_FILE) -n $(CLIENT_NAME) -s $(CLIENT_OS)

.PHONY: run
run:
	FLOGO_APP_PROP_RESOLVERS=env FLOGO_APP_PROPS_ENV=auto PORT=8989 FLOGO_LOG_LEVEL=DEBUG FLOGO_SCHEMA_SUPPORT=true FLOGO_SCHEMA_VALIDATION=false CRYPTO_PATH=$(FAB_PATH)/first-network/crypto-config $(TOOL_PATH)/work/$(CLIENT_NAME)_$(CLIENT_OS)_amd64

