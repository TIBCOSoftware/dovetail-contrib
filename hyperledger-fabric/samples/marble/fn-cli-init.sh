#!/bin/bash

# install and instantiate marble_cc in the Fabric sample first-network
# execute this script from the scripts folder of the cli docker container

. ./utils.sh
CCNAME=marble_cc

echo "install ${CCNAME} on peer0 org1"
peer chaincode install -n ${CCNAME} -v 1.0 -p github.com/chaincode/${CCNAME}

echo "install ${CCNAME} on peer1 org1"
setGlobals 1 1
peer chaincode install -n ${CCNAME} -v 1.0 -p github.com/chaincode/${CCNAME}

echo "install ${CCNAME} on peer0 org2"
setGlobals 0 2
peer chaincode install -n ${CCNAME} -v 1.0 -p github.com/chaincode/${CCNAME}

echo "install ${CCNAME} on peer1 org2"
setGlobals 1 2
peer chaincode install -n ${CCNAME} -v 1.0 -p github.com/chaincode/${CCNAME}

echo "instantiate ${CCNAME}"
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
peer chaincode instantiate $ORDERER_ARGS -C mychannel -n ${CCNAME} -v 1.0 -c '{"Args":["init"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"
