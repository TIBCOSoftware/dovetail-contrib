#!/bin/bash

# install and instantiate marble_private_cc in the Fabric sample first-network
# execute this script from the scripts folder of the cli docker container

. ./utils.sh
CCNAME=marble_private_cc

echo "install ${CCNAME} on peer0 org1"
peer chaincode install -n ${CCNAME} -v 1.0 -p github.com/chaincode/marble_private_cc

echo "install ${CCNAME} on peer1 org1"
setGlobals 1 1
peer chaincode install -n ${CCNAME} -v 1.0 -p github.com/chaincode/marble_private_cc

echo "install ${CCNAME} on peer0 org2"
setGlobals 0 2
peer chaincode install -n ${CCNAME} -v 1.0 -p github.com/chaincode/marble_private_cc

echo "install ${CCNAME} on peer1 org2"
setGlobals 1 2
peer chaincode install -n ${CCNAME} -v 1.0 -p github.com/chaincode/marble_private_cc

echo "instantiate ${CCNAME}"
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
peer chaincode instantiate $ORDERER_ARGS -C mychannel -n ${CCNAME} -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')" --collections-config /opt/gopath/src/github.com/chaincode/marbles02_private/collections_config.json
