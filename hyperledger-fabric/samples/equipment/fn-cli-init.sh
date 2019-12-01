#!/bin/bash

# install and instantiate equipment chaincode in the Fabric sample first-network
# execute this script from the scripts folder of the cli docker container

. ./utils.sh
CCNAME=${1:-"equipment_cc"}
CDS_FILE=${CCNAME}_1.0.cds

echo "package chaincode ${CCNAME}:1.0"
peer chaincode package -n ${CCNAME} -v 1.0 -p github.com/chaincode/${CCNAME} ${CDS_FILE}

echo "install ${CCNAME} on peer0 org1"
peer chaincode install ${CDS_FILE}

echo "install ${CCNAME} on peer1 org1"
setGlobals 1 1
peer chaincode install ${CDS_FILE}

echo "install ${CCNAME} on peer0 org2"
setGlobals 0 2
peer chaincode install ${CDS_FILE}

echo "install ${CCNAME} on peer1 org2"
setGlobals 1 2
peer chaincode install ${CDS_FILE}

echo "instantiate ${CCNAME} chaincode"
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
peer chaincode instantiate $ORDERER_ARGS -C mychannel -n ${CCNAME} -v 1.0 -c '{"Args":["init"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"
