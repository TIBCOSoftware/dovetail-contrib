#!/bin/bash

# install and instantiate equinix chaincode in the Fabric sample first-network
# execute this script from the scripts folder of the cli docker container

. ./utils.sh
CCNAME=${1:-"equinix_cc"}

echo "install ${CCNAME} chaincode on peer0 org1"
peer chaincode install -n ${CCNAME} -v 1.0 -p github.com/chaincode/${CCNAME}

setGlobals 0 2
echo "install ${CCNAME} chaincode on peer0 org2"
peer chaincode install -n ${CCNAME} -v 1.0 -p github.com/chaincode/${CCNAME}

echo "instantiate ${CCNAME} chaincode"
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
peer chaincode instantiate $ORDERER_ARGS -C mychannel -n ${CCNAME} -v 1.0 -c '{"Args":["init"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"
