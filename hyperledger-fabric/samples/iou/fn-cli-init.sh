#!/bin/bash

# install and instantiate iou_cc in the Fabric sample first-network
# execute this script from the scripts folder of the cli docker container

. ./utils.sh
CCNAME=iou_cc
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

echo "instantiate ${CCNAME}"
# Note: some transactions (createAccount) can be done by only one of the orgs, and so the endorsement policy uses 'OR' for simplicity here;
# real production should use `AND` for most transactions; createAccount may use a different endorsement policy by packaging it as a second chaincode.
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
peer chaincode instantiate $ORDERER_ARGS -C mychannel -n ${CCNAME} -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')" --collections-config /opt/gopath/src/github.com/chaincode/${CCNAME}/collections_config.json
