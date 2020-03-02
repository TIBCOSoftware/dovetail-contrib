MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
DEPLOY_PATH   := $(SRC_PATH)/../netop1.com/gateway
FAB_PATH      ?= $(GOPATH)/src/github.com/hyperledger/fabric-samples

.PHONY: all
all: clean build

.PHONY: build
build: codegen
	cd src && go build -o ../gateway
	cd src && GOOS=linux GOARCH=amd64 go build -o ../gateway-linux

.PHONY: dist
dist: codegen
	mkdir -p $(DEPLOY_PATH)/bin
	cp gateway-linux $(DEPLOY_PATH)/bin
	cp src/fabric.proto $(DEPLOY_PATH)
	cp -Rf swagger-ui $(DEPLOY_PATH)

.PHONY: run
run:
	CRYPTO_PATH=$(FAB_PATH)/first-network/crypto-config ./gateway -logtostderr -v 2

.PHONY: codegen
codegen:
	protoc -I=./src \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:./src/fabric \
		./src/fabric.proto
	protoc -I=./src \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:./src/fabric \
		./src/fabric.proto
	protoc -I=./src \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:./swagger-ui \
		./src/fabric.proto

.PHONY: clean
clean:
	rm -f gateway

.PHONY: depend
depend:
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
