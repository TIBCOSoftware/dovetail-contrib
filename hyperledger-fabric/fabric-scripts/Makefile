MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
BUILD_PATH    = $(PWD)
SHIM_PATH     := $(SRC_PATH)/../shim
PATCH_PATH    := $(SRC_PATH)/../flogo-patch

FLOGO_VER     ?= v0.9.4
DT_REPO       ?= github.com/TIBCOSoftware/dovetail-contrib

.PHONY: all
all: create build

.PHONY: clean
clean:
	rm -Rf $(BUILD_PATH)

.PHONY: create
create: 
	cp $(APP_FILE) $(BUILD_PATH)
	flogo create --cv $(FLOGO_VER) -f $(APP_FILE) $(APP_NAME)
	rm $(BUILD_PATH)/$(APP_NAME)/src/main.go
	cp $(SHIM_PATH)/chaincode_shim.go $(BUILD_PATH)/$(APP_NAME)/src/main.go
	cp -Rf META-INF $(BUILD_PATH)/$(APP_NAME)/src
	cp $(PATCH_PATH)/codegen.sh $(BUILD_PATH)/$(APP_NAME)

.PHONY: build
build: $(BUILD_PATH)/$(APP_NAME)
	cd $(BUILD_PATH)/$(APP_NAME)/src && \
	go mod edit -replace=github.com/project-flogo/core@v0.10.1=github.com/project-flogo/core@$(FLOGO_VER) && \
	go mod edit -replace=github.com/project-flogo/flow@v0.10.0=github.com/project-flogo/flow@$(FLOGO_VER)

	-cd $(BUILD_PATH)/$(APP_NAME) && flogo build -e
	
	cd $(BUILD_PATH)/$(APP_NAME)/src && \
	go mod vendor && \
	cp -R $(PATCH_PATH)/* $(BUILD_PATH)/$(APP_NAME)/src/vendor/github.com/project-flogo && \
	env GOOS=linux GOARCH=amd64 go build -mod vendor -o ../$(APP_NAME)

.PHONY: metadata
metadata:
	mkdir -p $(SRC_PATH)/contract-metadata
	fabric-tools metadata -f $(APP_FILE) -m $(SRC_PATH)/contract-metadata/metadata.json -g $(SRC_PATH)/contract-metadata/metadata.gql -o $(SRC_PATH)/contract-metadata/override.json

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
	cd $(CLIENT_PATH)/$(CLIENT_NAME)/src && ./gomodedit.sh
	cd $(CLIENT_PATH)/$(CLIENT_NAME) && flogo build -e

.PHONY: run
run:
	FLOGO_APP_PROP_RESOLVERS=env FLOGO_APP_PROPS_ENV=auto PORT=8989 FLOGO_LOG_LEVEL=DEBUG FLOGO_SCHEMA_SUPPORT=true FLOGO_SCHEMA_VALIDATION=false $(CLIENT_PATH)/$(CLIENT_NAME)/bin/$(CLIENT_NAME)
