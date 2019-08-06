#!/bin/bash

# marble_cc tests executed from cli docker container of the Fabric sample first-network

. ./utils.sh
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
ORG1_ARGS="--peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
ORG2_ARGS="--peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

# insert test data
echo "insert 6 marbles ..."
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["initMarble","marble1","blue","35","tom"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["initMarble","marble2","red","50","tom"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["initMarble","marble3","blue","70","tom"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["initMarble","marble4","purple","80","tom"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["initMarble","marble5","purple","90","tom"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["initMarble","marble6","purple","100","tom"]}'

# transfer marble ownership
setGlobals 0 1
echo "test transfer marbles ..."
sleep 5
peer chaincode query -C mychannel -n marble_cc -c '{"Args":["readMarble","marble2"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["transferMarble","marble2","jerry"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["transferMarblesBasedOnColor","blue","jerry"]}'
sleep 5
peer chaincode query -C mychannel -n marble_cc -c '{"Args":["getMarblesByRange","marble1","marble5"]}'

# delete marble state, not history
echo "test delete and history"
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["delete","marble1"]}'
sleep 5
peer chaincode query -C mychannel -n marble_cc -c '{"Args":["getHistoryForMarble","marble1"]}'

# rich query
echo "test rich query ..."
peer chaincode query -C mychannel -n marble_cc -c '{"Args":["queryMarblesByOwner","jerry"]}'

# query pagination using page-size and starting bookmark
echo "test pagination ..."
peer chaincode query -C mychannel -n marble_cc -c '{"Args":["getMarblesByRangeWithPagination","marble1","marble9", "3", ""]}'
peer chaincode query -C mychannel -n marble_cc -c '{"Args":["getMarblesByRangeWithPagination","marble1","marble9", "3", "marble5"]}'
peer chaincode query -C mychannel -n marble_cc -c '{"Args":["queryMarbles","{\"selector\":{\"docType\":\"marble\",\"owner\":\"tom\"}}"]}'
peer chaincode query -C mychannel -n marble_cc -c '{"Args":["queryMarblesWithPagination","{\"selector\":{\"docType\":\"marble\",\"owner\":\"tom\"}}", "2", ""]}'
