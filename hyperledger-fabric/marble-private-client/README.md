# marble-private-client
This is a sample client app for Hyperledger Fabric.  Implemented using the [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-4-0), this app interacts with the Hyperledger Fabric chaincode [`marble-private`](../marble-private) and exposes a set of REST APIs for managing private data collections on the marble blockchain network.

## Build and start the marble-private fabric network
First, complete the prerequisites as described in [`marble-private`](../marble-private).

Then, build and deploy the marble-private chaincode (assuming that the `fabric-samples` are installed under your `$GOPATH`):
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/marble-private
make create
make deploy
```

Start Hyperledger Fabric first-network with CouchDB:
```
cd $GOPATH//src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh up -s couchdb
```

Use the `cli` container to install the `marble_private_cc` chaincode on all 4 peers of `org1` and `org2`, and then instantiate it.
```
docker exec -it cli bash
. scripts/utils.sh
peer chaincode install -n marble_private_cc -v 1.0 -p github.com/chaincode/marble_private_cc
setGlobals 1 1
peer chaincode install -n marble_private_cc -v 1.0 -p github.com/chaincode/marble_private_cc
setGlobals 0 2
peer chaincode install -n marble_private_cc -v 1.0 -p github.com/chaincode/marble_private_cc
setGlobals 1 2
peer chaincode install -n marble_private_cc -v 1.0 -p github.com/chaincode/marble_private_cc

ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
peer chaincode instantiate $ORDERER_ARGS -C mychannel -n marble_private_cc -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')" --collections-config /opt/gopath/src/github.com/chaincode/marbles02_private/collections_config.json
```

## Build and start the marble-private-client app
Create, build and start the marble-private-client app from the model file [`marble_private_lient.json`](marble_private_client.json):
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/marble-private-client
make create
make build
make run
```

The step for `create` may take a few minutes because it uses `dep` to analyze and fetch Go dependencies, and `dep` is slow.  This issue will be resolved in a future Flogo release when `dep` is replaced by `Go modules`.  Sometimes, `dep` may fail on the first try, in which case, you may manually execute the `dep` one more time, i.e.,
```
cd marbleprivate_client/src/marbleprivate_client
dep ensure -v -update
cd ../..
```

## Test marble-private-client app
This app implements a set of REST APIs:
- **Create Marble** (PUT): it creates a new marble.
- **Transfer Marble** (PUT): it transfers a marble to a new owner.
- **Delete Marble** (DELETE): it deletes the state of a specified marble.
- **Get Marble** (GET): it retrieves a marble record by its key.
- **Get Marble Price** (GET): it retrieves a marble's private details by its key.
- **Query By Owner** (GET): it queries marble records by an owner name.
- **Query By Range** (GET): it retrieves marble records in a specified range of keys.

You may use the following commands to test the behavior of these REST APIs.  If you do not like command-line `curl`, you may download and use a REST client tool to submit these REST requests.  For Mac users, the [`Advanced Rest client`](https://install.advancedrestclient.com/install) is pretty user-friendly.

```
# insert test data
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble1","color":"blue","size":35,"owner":"tom","price":99}' http://localhost:8989/marbleprivate/create
curl -X GET http://localhost:8989/marbleprivate/key/marble1
curl -X GET http://localhost:8989/marbleprivate/price/marble1

# more inserts and transfer owner, test purge of private marble1 after 3 blocks
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble2","color":"red","size":50,"owner":"tom","price":199}' http://localhost:8989/marbleprivate/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble3","color":"blue","size":70,"owner":"tom","price":299}' http://localhost:8989/marbleprivate/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble2","owner":"jerry"}' http://localhost:8989/marbleprivate/transfer
# marble1 pricing is still available
curl -X GET http://localhost:8989/marbleprivate/price/marble1
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble3","owner":"jerry"}' http://localhost:8989/marbleprivate/transfer
# marble1 private detail is purged after 3 blocks, so this returns error
curl -X GET http://localhost:8989/marbleprivate/price/marble1

# test query and delete
curl -X GET http://localhost:8989/marbleprivate/owner/jerry
curl -X GET http://localhost:8989/marbleprivate/range?startKey=marble1&endKey=marble4
curl -X DELETE http://localhost:8989/marbleprivate/delete/marble2
curl -X GET http://localhost:8989/marbleprivate/owner/jerry
```

## Cleanup the marble-private fabric network
Stop and cleanup the Fabric `first-network`.
```
cd $GOPATH//src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```
