MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
APP_FILE      := marble_app.json
APP_NAME      := marble_cc
LOG_LEVEL     := INFO
BUILD_PATH    := /tmp/marble
CC_DEPLOY     := $(GOPATH)/src/github.com/hyperledger/fabric-samples/chaincode
SHIM_PATH     := $(dir $(SRC_PATH))shim
PATCH_PATH    := $(dir $(SRC_PATH))flogo-patch

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
	cp -Rf META-INF $(BUILD_PATH)/$(APP_NAME)/src
	rm -Rf $(CC_DEPLOY)/$(APP_NAME)
	cp -Rf $(BUILD_PATH)/$(APP_NAME)/src $(CC_DEPLOY)/$(APP_NAME)

.PHONY: package
package: 
	docker exec cli peer chaincode package -n $(APP_NAME) -v 1.0 -p github.com/chaincode/$(APP_NAME) ../../../chaincode/$(APP_NAME).cds
	mv $(CC_DEPLOY)/$(APP_NAME).cds $(SRC_PATH)
