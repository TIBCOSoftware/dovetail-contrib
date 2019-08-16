#!/bin/bash

# install and instantiate marble_cc in the Fabric sample first-network
# execute this script from the scripts folder of the cli docker container
APP_NAME=audit_cc

. ./utils.sh

echo "install ${APP_NAME} on peer0 org1"
peer chaincode install -n ${APP_NAME} -v 1.0 -p github.com/chaincode/${APP_NAME}

setGlobals 0 2
echo "install ${APP_NAME} on peer0 org2"
peer chaincode install -n ${APP_NAME} -v 1.0 -p github.com/chaincode/${APP_NAME}

echo "instantiate ${APP_NAME}"
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
peer chaincode instantiate $ORDERER_ARGS -C mychannel -n ${APP_NAME} -v 1.0 -c '{"Args":["init"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"
