#!/bin/bash

# install and instantiate marble_private_cc in the Fabric sample first-network
# execute this script from the scripts folder of the cli docker container

. ./utils.sh
CCNAME=marble_private_cc
CC_PATH=${GOPATH}/src/github.com/chaincode
CDS_FILE=${CC_PATH}/${CCNAME}/${CCNAME}_1.0.cds

if [ ! -f "${CDS_FILE}" ]; then
  echo "cannot find cds package: ${CDS_FILE}"
  exit 1
fi

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

echo "instantiate ${CCNAME}"
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
peer chaincode instantiate $ORDERER_ARGS -C mychannel -n ${CCNAME} -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')" --collections-config ${CC_PATH}/marbles02_private/collections_config.json
