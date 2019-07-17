# marble
This is a Flogo version of the [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) sample chaincode, [marbles02](https://github.com/hyperledger/fabric-samples/tree/release-1.4/chaincode/marbles02) implemented by using a [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) model.  The model does not require any code, it contains only a JSON model file exported from the TIBCO Flogo® Enterprise.  You can download the prerequisites and then build and deploy the model to a Hyperledger Fabric network as described below.

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.6](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html)
- [Install Go](https://golang.org/doc/install)
- Clone [Hyperledger Fabric](https://github.com/hyperledger/fabric)
- Clone [Hyperledger Fabric Samples](https://github.com/hyperledger/fabric-samples)
- Download and install [flogo-cli](https://github.com/TIBCOSoftware/flogo-cli)
- Clone dovetail-contrib with this Flogo extension

There are different ways to clone these packages.  This document assumes that you have installed these packages under $GOPATH after installing Go, i.e.,
```
go get -u github.com/hyperledger/fabric
go get -u github.com/hyperledger/fabric-samples
go get -u github.com/TIBCOSoftware/flogo-cli/...
go get -u github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric
```
Bootstrap fabric-samples
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples
./scripts/bootstrap.sh
```

## Edit smart contract (optional)
Skip to the next section if you do not plan to modify this sample at this time.

- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabricExtension.zip`](../fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabric.sh`](../zip-fabric.sh)
- Create new Flogo App of name `marble_app` and choose `Import app` to import the model [`marble_app.json`](marble_app.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`marble_app.json`](marble_app.json) to this `marble` folder.

## Build and deploy chaincode to Hyperledger Fabric

- In this `marble` folder, execute `make create` to generate the chaincode source code from the flogo model [`marble_app.json`](marble_app.json).
- Execute `make deploy` to build and deploy the chaincode to the `fabric-samples` chaincode folder.  Note: you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make create
make deploy
```

## Install and test chaincode using fabric sample first-network
Start Hyperledger Fabric first-network with CouchDB:
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh up -s couchdb
```
Use `cli` docker container to install and instantiate the resulting `marble_cc` chaincode.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make cli-init
```

Optionally, test the chaincode from `cli` docker container, i.e.,
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make cli-test
```
You can skip this test, and follow the steps in the next section to build the client app, and then use the client app to execute the same tests.

Note that developers can also use Fabric dev-mode to test chaincode (refer [dev](./dev.md) for more details).

## Build and start the marble-client
The marble-client is a REST service that invokes the `marble` chaincode.  It is implemented as a Flogo model, [`marble_client_app.json`](marble_client_app.json).  Build and start this client app as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make create-client
make build-client
make run
```

## Test marble-client app and marble chaincode
The marble-client app implements the following REST APIs to invoke corresponding blockchain transactions of the `marble` chaincode:
- **Create Marble** (PUT): it creates a new marble.
- **Transfer Marble** (PUT): it transfers a marble to a new owner.
- **Transfer By Color** (PUT): it transfers all marbles of a specified color to a new owner.
- **Delete Marble** (DELETE): it deletes the state of a specified marble.
- **Get Marble** (GET): it retrieves a marble record by its key.
- **Query By Owner** (GET): it queries marble records by an owner name.
- **Query By Range** (GET): it retrieves marble records in a specified range of keys.
- **Marble History** (GET): it retrieves the history of a marble.
- **Query Range Page** (GET): it retrieves marble records in a range of keys, with pagination support.

You may use the following commands to invoke the REST APIs. If you do not like command-line `curl`, you may download and use a REST client tool to submit these REST requests.  For Mac users, the [`Advanced Rest client`](https://install.advancedrestclient.com/install) is pretty user-friendly.

```
# insert test data
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble1","color":"blue","size":35,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble2","color":"red","size":50,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble3","color":"blue","size":70,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble4","color":"purple","size":80,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble5","color":"purple","size":90,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble6","color":"purple","size":100,"owner":"tom"}' http://localhost:8989/marble/create
curl -X GET http://localhost:8989/marble/key/marble2
curl -X GET "http://localhost:8989/marble/range?startKey=marble1&endKey=marble5"

# transfer marble ownership
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble2","newOwner":"jerry"}' http://localhost:8989/marble/transfer
curl -H 'Content-Type: application/json' -X PUT -d '{"color":"blue","newOwner":"jerry"}' http://localhost:8989/marble/transfercolor
curl -X GET http://localhost:8989/marble/owner/jerry
curl -X GET "http://localhost:8989/marble/range?startKey=marble1&endKey=marble5"

# delete marble state, not history
curl -X DELETE http://localhost:8989/marble/delete/marble1
curl -X GET http://localhost:8989/marble/history/marble1

# query pagination using page-size and starting bookmark
curl -X GET "http://localhost:8989/marble/rangepage?startKey=marble1&endKey=marble7&pageSize=3"
curl -X GET "http://localhost:8989/marble/rangepage?startKey=marble1&endKey=marble7&pageSize=3&bookmark=marble5"
```

## Cleanup the marble-app fabric network
After you are done testing, you can stop and cleanup the Fabric sample `first-network` as follows:
```
cd $GOPATH//src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```

## Deploy to IBM Cloud
This client app can access the same `marble` chaincode deployed in [IBM Cloud](https://cloud.ibm.com) using the [IBM Blockchain Platform](https://cloud.ibm.com/catalog/services/blockchain-platform-20).  Refer to [fabric-tools](../../fabric-tools) for details about installation of the `marble` chaincode in IBM Cloud.
