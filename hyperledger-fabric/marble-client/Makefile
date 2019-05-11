MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
APP_FILE      := marble_client_app.json
APP_NAME      := marble_client
LOG_LEVEL     := INFO

.PHONY: all
all: create build

.PHONY: clean
clean:
	rm -Rf $(APP_NAME)

.PHONY: create
create: $(APP_FILE) clean
	-flogo create -f $(APP_FILE) $(APP_NAME)
	@echo "copy dep rules, and reset dependency..."
	cp Gopkg.toml $(APP_NAME)/src/$(APP_NAME)
	sleep 10
	cd $(APP_NAME)/src/$(APP_NAME) && dep ensure -v -update || \
	dep ensure -v -update

.PHONY: build
build: $(APP_NAME)
	cd $(APP_NAME) && \
	flogo build -e
#	GO111MODULE=on go mod init github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/$(APP_NAME)
#	GO111MODULE=on go mod edit -replace github.com/golang/lint@latest=golang.org/x/lint@latest
#	GO111MODULE=on go get -u -d github.com/hyperledger/fabric-sdk-go@master
#	GO111MODULE=on go build

.PHONY: run
run:
	FLOGO_LOG_LEVEL=DEBUG $(APP_NAME)/bin/$(APP_NAME)

