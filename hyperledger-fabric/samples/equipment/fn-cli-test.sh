#!/bin/bash

# equipment tests executed from cli docker container of the Fabric sample first-network

. ./utils.sh
CCNAME=${1:-"equipment_cc"}
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
ORG1_ARGS="--peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
ORG2_ARGS="--peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

echo "Run cli tests for chaincode ${CCNAME} ..."

# submitPO:
# 1. Asset Unique Identifier
# 2. Description
# 3. Purchase Price
# 4. Location
# 5. Vendor
# 6. Org_ID
peer chaincode invoke $ORDERER_ARGS -C mychannel -n ${CCNAME} $ORG1_ARGS $ORG2_ARGS -c '{"Args":["submitPO","asset1","submit1","10000","loc1","vendor1","org1"]}'
sleep 5

# receiveAsset:
# 1. Asset Unique Identifier
# 2. Description
# 3. Acquisition Date
# 4. Location
# 5. Make
# 6. Vendor
# 7. Model
# 8. Serial Number
# 9. Org ID
peer chaincode invoke $ORDERER_ARGS -C mychannel -n ${CCNAME} $ORG1_ARGS $ORG2_ARGS -c '{"Args":["receiveAsset","asset1","receiveAsset1","2019-07-01","loc1","make1","vendor1","model1","sn1","org1"]}'
sleep 5

# installAsset:
# 1. Asset Unique Identifier
# 2. Description
# 3. Acquisition Date
# 4. Location
# 5. Make
# 6. Vendor
# 7. Model
# 8. Serial Number
# 9. Org ID
# 10. Install Date
peer chaincode invoke $ORDERER_ARGS -C mychannel -n ${CCNAME} $ORG1_ARGS $ORG2_ARGS -c '{"Args":["installAsset","asset1","installAsset1","2019-07-01","loc1","make1","vendor1","model1","sn1","org1","2019-07-11"]}'
sleep 5

# receiveInvoice:
# 1. Asset Unique Identifier
# 2. Description
# 3. Location
# 4. Invoice Date
# 5. Org ID
# 6. Invoice Price
# 7. Vendor
peer chaincode invoke $ORDERER_ARGS -C mychannel -n ${CCNAME} $ORG1_ARGS $ORG2_ARGS -c '{"Args":["receiveInvoice","asset1","receiveInvoice1","loc1","2019-07-21","org1","10000","vendor1"]}'
sleep 5

# updateAsset
# 1. Asset Unique Identifier
# 2. Description
# 3. Purchase Price
# 4. Location
# 5. Model
# 6. Org_ID
# 7. Install Date
# 8. Invoice Date
# 9. Net Book Value
# 10. Serial Number
peer chaincode invoke $ORDERER_ARGS -C mychannel -n ${CCNAME} $ORG1_ARGS $ORG2_ARGS -c '{"Args":["updateAsset","asset1","updateAsset1","10000","loc1","model1","org1","2019-07-11","2019-07-21","15000","sn1"]}'
