#!/bin/bash

# marble_private_cc tests executed from cli docker container of the Fabric sample first-network

. ./utils.sh
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

# test insert and read access permission
echo "create marble1 and query test ..."
setGlobals 0 1
MARBLE=$(echo -n "{\"name\":\"marble1\",\"color\":\"blue\",\"size\":35,\"owner\":\"tom\",\"price\":99}" | base64 | tr -d \\n)
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_private_cc -c '{"Args":["initMarble"]}' --transient "{\"marble\":\"$MARBLE\"}"
sleep 5
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["readMarble","marble1"]}'
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["readMarblePrivateDetails","marble1"]}'
setGlobals 0 2
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["readMarble","marble1"]}'
# following should fail due to no read access permission
echo "org2 query should fail with no read access permission on privatedata ..."
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["readMarblePrivateDetails","marble1"]}'

# test more insert and transfer owner, and purge of marble1 after 3 blocks
setGlobals 0 1
# block +1
echo "create marble2 and transfer test ..."
MARBLE=$(echo -n "{\"name\":\"marble2\",\"color\":\"red\",\"size\":50,\"owner\":\"tom\",\"price\":199}" | base64 | tr -d \\n)
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_private_cc -c '{"Args":["initMarble"]}' --transient "{\"marble\":\"$MARBLE\"}"
sleep 5
# block +2
MARBLE_OWNER=$(echo -n "{\"name\":\"marble2\",\"owner\":\"jerry\"}" | base64 | tr -d \\n)
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_private_cc -c '{"Args":["transferMarble"]}' --transient "{\"marble_owner\":\"$MARBLE_OWNER\"}"
sleep 5
# block +3
MARBLE=$(echo -n "{\"name\":\"marble3\",\"color\":\"blue\",\"size\":70,\"owner\":\"tom\",\"price\":299}" | base64 | tr -d \\n)
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_private_cc -c '{"Args":["initMarble"]}' --transient "{\"marble\":\"$MARBLE\"}"
sleep 5

# marble1 should still be available
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["readMarblePrivateDetails","marble1"]}'
# block +4
MARBLE_OWNER=$(echo -n "{\"name\":\"marble3\",\"owner\":\"jerry\"}" | base64 | tr -d \\n)
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_private_cc -c '{"Args":["transferMarble"]}' --transient "{\"marble_owner\":\"$MARBLE_OWNER\"}"
sleep 5

# marble1 details purged after 3 blocks, so this returns error
echo "private detail should have been deleted, so the query fails to get marble private details ..."
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["readMarblePrivateDetails","marble1"]}'

# test query
echo "test query ..."
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["getMarblesByRange","marble1", "marble3"]}'
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["queryMarblesByOwner","jerry"]}'
echo "delete marble2 ..."
MARBLE_DELETE=$(echo -n "{\"name\":\"marble2\"}" | base64 | tr -d \\n)
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_private_cc -c '{"Args":["delete"]}' --transient "{\"marble_delete\":\"$MARBLE_DELETE\"}"
# verify deleted marble2
sleep 5
echo "marble2 should not be in the query result ..."
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["queryMarblesByOwner","jerry"]}'