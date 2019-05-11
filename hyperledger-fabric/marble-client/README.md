# marble-client
This is a sample client app for Hyperledger Fabric.  Implemented using the [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-4-0), this app interacts with the Hyperledger Fabric chaincode [`marble-app`](../marble-app) and exposes a set of REST APIs for managing the data on the marble blockchain network.

## Build and start the marble-app fabric network
First, complete the prerequisites as described in [`marble-app`](../marble-app).

Then, build and deploy the marble-app chaincode (assuming that the `fabric-samples` are installed under your `$GOPATH`):
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/marble-app
make create
make deploy
```

Start Hyperledger Fabric first-network with CouchDB:
```
cd $GOPATH//src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh up -s couchdb
```

Use the `cli` container to install the `marble_cc` chaincode on both `org1` and `org2`, and then instantiate it on the fabric network.
```
docker exec -it cli bash
. scripts/utils.sh
peer chaincode install -n marble_cc -v 1.0 -p github.com/chaincode/marble_cc
setGlobals 0 2
peer chaincode install -n marble_cc -v 1.0 -p github.com/chaincode/marble_cc
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
peer chaincode instantiate $ORDERER_ARGS -C mychannel -n marble_cc -v 1.0 -c '{"Args":["init"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"
exit
```

## Build and start the marble-client app
Create, build and start the marble-client app from the model file [`marble_client_app.json`](marble_client_app.json):
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/marble-client
make create
make build
make run
```

The step for `create` may take a few minutes because it uses `dep` to analyze and fetch Go dependencies, and `dep` is slow.  This issue will be resolved in a future Flogo release when `dep` is replaced by `Go modules`.  Sometimes, `dep` may fail on the first try, in which case, you may manually execute the `dep` one more time, i.e.,
```
cd marble_client/src/marble_client
dep ensure -v -update
cd ../..
```

## Test marble-client app
This app implements a set of REST APIs:
- **Create Marble** (PUT): it creates a new marble.
- **Transfer Marble** (PUT): it transfers a marble to a new owner.
- **Transfer By Color** (PUT): it transfers all marbles of a specified color to a new owner.
- **Delete Marble** (DELETE): it deletes the state of a specified marble.
- **Get Marble** (GET): it retrieves a marble record by its key.
- **Query By Owner** (GET): it queries marble records by an owner name.
- **Query By Range** (GET): it retrieves marble records in a specified range of keys.
- **Marble History** (GET): it retrieves the history of a marble.
- **Query Range Page** (GET): it retrieves marble records in a range of keys, with pagination support.

You may use the following commands to test the behavior of these REST APIs.  If you do not like command-line `curl`, you may download and use a REST client tool to submit these REST requests.  For Mac users, the [`Advanced Rest client`](https://install.advancedrestclient.com/install) is pretty user-friendly.

```
# insert test data
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble1","color":"blue","size":35,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble2","color":"red","size":50,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble3","color":"blue","size":70,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble4","color":"purple","size":80,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble5","color":"purple","size":90,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble6","color":"purple","size":100,"owner":"tom"}' http://localhost:8989/marble/create
curl -X GET http://localhost:8989/marble/key/marble2
curl -X GET http://localhost:8989/marble/range?startKey=marble1&endKey=marble5

# transfer marble ownership
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble2","newOwner":"jerry"}' http://localhost:8989/marble/transfer
curl -H 'Content-Type: application/json' -X PUT -d '{"color":"blue","newOwner":"jerry"}' http://localhost:8989/marble/transfercolor
curl -X GET http://localhost:8989/marble/owner/jerry
curl -X GET http://localhost:8989/marble/range?startKey=marble1&endKey=marble5

# delete marble state, not history
curl -X DELETE http://localhost:8989/marble/delete/marble1
curl -X GET http://localhost:8989/marble/history/marble1

# query pagination using page-size and starting bookmark
curl -X GET http://localhost:8989/marble/rangepage?startKey=marble1&endKey=marble7&pageSize=3
curl -X GET http://localhost:8989/marble/rangepage?startKey=marble1&endKey=marble7&pageSize=3&bookmark=marble5
```

## Cleanup the marble-app fabric network
Stop and cleanup the Fabric `first-network`.
```
cd $GOPATH//src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```

## Deploy to IBM Cloud
This client app can access the same marble-app chaincode deployed in [IBM Cloud](https://cloud.ibm.com) using the [IBM Blockchain Platform](https://cloud.ibm.com/catalog/services/blockchain-platform-20).  Refer to [fabric-tools](../fabric-tools) for details on the configuration.