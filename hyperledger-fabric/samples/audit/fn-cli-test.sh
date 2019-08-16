#!/bin/bash

# audit_cc tests executed from cli docker container of the Fabric sample first-network
APP_NAME=audit_cc

. ./utils.sh
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
ORG1_ARGS="--peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
ORG2_ARGS="--peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

# insert test data
echo "insert 4 audit logs ..."
peer chaincode invoke $ORDERER_ARGS -C mychannel -n ${APP_NAME} $ORG1_ARGS $ORG2_ARGS -c '{"Args":["createAudit","[{\"recID\":\"audit1\",\"domain\":\"tibco.oocto.dovetail.test\",\"owner\":\"oocto@tibco.com\",\"data\":\"test1\",\"hashType\":\"\",\"hashValue\":\"\"},{\"recID\":\"audit2\",\"domain\":\"tibco.oocto.dovetail.test\",\"owner\":\"oocto@tibco.com\",\"data\":\"test2\",\"hashType\":\"\",\"hashValue\":\"\"}]"]}'
sleep 5
peer chaincode invoke $ORDERER_ARGS -C mychannel -n ${APP_NAME} $ORG1_ARGS $ORG2_ARGS -c '{"Args":["createAudit","[{\"recID\":\"audit3\",\"domain\":\"tibco.oocto.dovetail.test\",\"owner\":\"oocto@tibco.com\",\"data\":\"test3\",\"hashType\":\"\",\"hashValue\":\"\"},{\"recID\":\"audit2\",\"domain\":\"tibco.oocto.dovetail.test\",\"owner\":\"oocto@tibco.com\",\"data\":\"test2 again\",\"hashType\":\"\",\"hashValue\":\"\"}]"]}'
sleep 5

# test query
data=$(peer chaincode query -C mychannel -n ${APP_NAME} -c '{"Args":["getRecordsByID","audit1"]}')
tx1=${data#*txID\":\"}
id1=${tx1%%\"*}

echo "query with first txID ${id1}"
peer chaincode query -C mychannel -n ${APP_NAME} -c '{"Args":["getRecord","'${id1}':audit1"]}'
peer chaincode query -C mychannel -n ${APP_NAME} -c '{"Args":["getRecordsByTxID","'${id1}'"]}'
peer chaincode query -C mychannel -n ${APP_NAME} -c '{"Args":["getRecordsByID","audit2"]}'
tx1=${data#*txTime\":\"}
tm1=${tx1%%\"*}

echo "query with first txTime ${tm1}"
peer chaincode query -C mychannel -n ${APP_NAME} -c '{"Args":["getRecordsByTxTime","'${tm1}'"]}'
today=`date +%F`
tomorrow=`date --date='1 day' +%F`

echo "query audit records of today ${today}"
peer chaincode query -C mychannel -n ${APP_NAME} -c '{"Args":["queryTimeRange","tibco.oocto.dovetail.test","oocto@tibco.com","'${today}'","'${tomorrow}'"]}'

# test delete
echo "delete record ${id1}:audit1"
peer chaincode invoke $ORDERER_ARGS -C mychannel -n ${APP_NAME} $ORG1_ARGS $ORG2_ARGS -c '{"Args":["deleteRecord","'${id1}':audit1"]}'
sleep 5
peer chaincode query -C mychannel -n ${APP_NAME} -c '{"Args":["getRecordsByTxID","'${id1}'"]}'

echo "delete transaction ${id1}"
peer chaincode invoke $ORDERER_ARGS -C mychannel -n ${APP_NAME} $ORG1_ARGS $ORG2_ARGS -c '{"Args":["deleteTransaction","'${id1}'"]}'
sleep 5
peer chaincode query -C mychannel -n ${APP_NAME} -c '{"Args":["getRecordsByTxID","'${id1}'"]}'
echo "query audit records of today ${today}"
peer chaincode query -C mychannel -n ${APP_NAME} -c '{"Args":["queryTimeRange","tibco.oocto.dovetail.test","oocto@tibco.com","'${today}'","'${tomorrow}'"]}'
