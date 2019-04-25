MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
APP_FILE      := marble_app.json
APP_NAME      := marble_cc
LOG_LEVEL     := INFO
CC_DEPLOY     := $(GOPATH)/src/github.com/hyperledger/fabric-samples/chaincode
SHIM_PATH     := $(dir $(SRC_PATH))shim

.PHONY: all
all: create build

.PHONY: clean
clean:
	rm -Rf $(APP_NAME)

.PHONY: create
create: $(APP_FILE) clean
	flogo create -f $(APP_FILE) $(APP_NAME) || echo "continue on error"
	cp $(SHIM_PATH)/chaincode_shim.go $(APP_NAME)/src/$(APP_NAME) && \
	cp Gopkg.toml $(APP_NAME)/src/$(APP_NAME) && \
	rm $(APP_NAME)/src/$(APP_NAME)/main.go && \
	cd $(APP_NAME)/src/$(APP_NAME) && \
	dep ensure -v -update

.PHONY: build
build: $(APP_NAME)
	cp $(SHIM_PATH)/chaincode_shim.go $(APP_NAME)/src/$(APP_NAME) && \
	cd $(APP_NAME) && \
	env GOOS=linux GOARCH=amd64 flogo build -e

.PHONY: deploy
deploy: build
ifneq (,$(wildcard ./$(APP_NAME)/bin/linux_amd64/$(APP_NAME)))
	cp $(APP_NAME)/bin/linux_amd64/$(APP_NAME) $(APP_NAME)/src/$(APP_NAME)
endif
ifneq (,$(wildcard ./$(APP_NAME)/bin/$(APP_NAME)))
	cp $(APP_NAME)/bin/$(APP_NAME) $(APP_NAME)/src/$(APP_NAME)
endif
	rm -Rf $(CC_DEPLOY)/$(APP_NAME) && \
	cp -Rf $(APP_NAME)/src/$(APP_NAME) $(CC_DEPLOY) && \
	cp -Rf META-INF $(CC_DEPLOY)/$(APP_NAME)
