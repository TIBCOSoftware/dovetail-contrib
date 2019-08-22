# marble

This example uses the [project Dovetail](https://tibcosoftware.github.io/dovetail/) to implement the [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) sample chaincode [marbles02](https://github.com/hyperledger/fabric-samples/tree/release-1.4/chaincode/marbles02).  This sample demonstrates basic features of the Hyperledger Fabric, including creeation and update of states and composite-keys, and different types of queries for state and history with pagination. It is implemented using [Flogo®](https://www.flogo.io/) models by visual programming with zero-code.  The Flogo® models can be created, imported, edited, and/or exported by using [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) or [Dovetail](https://github.com/TIBCOSoftware/dovetail).

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.6](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html). If you do not have access to `Flogo Enterprise`, you may sign up a trial on [TIBCO CLOUD Integration (TCI)](https://cloud.tibco.com/), or download Dovetail v0.2.0.  This sample uses `TIBCO Flogo® Enterprise`, but all models can be imported and edited by using Dovetail v0.2.0 and above.
- [Install Go](https://golang.org/doc/install)
- Clone [Hyperledger Fabric](https://github.com/hyperledger/fabric)
- Clone [Hyperledger Fabric Samples](https://github.com/hyperledger/fabric-samples)
- Download and install [flogo-cli](https://github.com/project-flogo/cli)
- Clone dovetail-contrib with Flogo extension for Hyperledger Fabric

There are different ways to clone these packages.  This document assumes that you have installed these packages under $GOPATH after installing Go, i.e.,
```
go get -u github.com/hyperledger/fabric
go get -u github.com/hyperledger/fabric-samples
go get -u github.com/project-flogo/cli/...
go get -u github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric
```
Note that the latest version of the Flogo extension for Hyperledger Fabric can be downloaded from the [`develop` branch of the `dovetail-contrib`](https://github.com/TIBCOSoftware/dovetail-contrib/tree/develop).

Bootstrap fabric-samples
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples
./scripts/bootstrap.sh
```

## Edit smart contract (optional)
Skip to the next section if you do not plan to modify the included sample model.

- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabricExtension.zip`](../../fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can generate this `zip` by using the script [`zip-fabric.sh`](../../zip-fabric.sh).
- Create new Flogo App of name `marble_app` and choose `Import app` to import the model [`marble_app.json`](marble_app.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`marble_app.json`](marble_app.json) to this `marble` sample folder.

Note that when a flogo model is imported to `Flogo® Enterprise v2.6.1`, a `return` activity is automatically added to the end of all branches, which could be an issue if the `return` activity is not at the end of a flow.  Thus, you need to carefully remove the mistakenly added `return` activities after the model is imported.  This issue will be fixed in a later release of the `Flogo® Enterprise`.

## Build and deploy chaincode to Hyperledger Fabric

- In this `marble` sample folder, execute `make create` to generate the chaincode source code from the flogo model [`marble_app.json`](marble_app.json).
- Execute `make deploy` to build and deploy the chaincode to the `fabric-samples` chaincode folder.  Note that you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.
- Execute `make package` to generate `cds` package for cloud deployment, and `metadata` for client apps.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make create
make deploy
make package
```

## Install and test chaincode using fabric sample first-network
Start Hyperledger Fabric first-network with CouchDB:
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh up -n -s couchdb
```
Use `cli` docker container to install and instantiate the `marble_cc` chaincode.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make cli-init
```

Optionally, test the chaincode from `cli` docker container, i.e.,
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make cli-test
```
You may skip this test, and follow the steps in the next section to build client apps, and then use the client app to execute the tests. If you run the `cli` tests, however, it should print out 17 successful tests with status code `200` if the `marble_app` chaincode is installed and instantiated successfully on the Fabric network.

Note that developers can also use Fabric dev-mode to test chaincode (refer [dev](./dev.md) for more details).  For issues regarding how to work with the Fabric network, please refer the [Hyperledger Fabric docs](https://hyperledger-fabric.readthedocs.io/en/latest/build_network.html).

## Edit marble REST service (optional)
The sample Flogo model, [`marble_client_app.json`](marble_client_app.json) is a REST service that invokes the `marble` chaincode.  Skip to the next section if you do not plan to modify the sample model.

The client app requires the metadata of the `marble-app` chaincode. You can generate the contract metadata [`metadata.json`](contract-metadata/metadata.json) by
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make package
```
Following are steps to edit or view the REST service models.
- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabclientExtension.zip`](../../fabclientExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can generate this `zip` by using the script [`zip-fabclient.sh`](../../zip-fabclient.sh).
- Create new Flogo App of name `marble_client_app` and choose `Import app` to import the model [`marble_client_app.json`](marble_client_app.json)
- You can then add or update the service implementation using the graphical modeler of the TIBCO Flogo® Enterprise.
- Open `Connections` tab, find and edit the `marble client` connector.  Set the `Smart contract metadata file` to the [`metadata.json`](contract-metadata/metadata.json) generated in the previous step. Set the `Network configuration file` and `entity matcher file` to the corresponding files in [`testdata`](../../testdata).
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`marble_client_app.json`](marble_client_app.json) to this `marble` sample folder.

## Build and start the marble REST service
Build and start the client app as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make create-client
make build-client
make run
```

## Test marble REST service and marble chaincode
The REST service implements the following APIs to invoke corresponding blockchain transactions of the `marble` chaincode:
- **Create Marble** (PUT): it creates a new marble.
- **Transfer Marble** (PUT): it transfers a marble to a new owner.
- **Transfer By Color** (PUT): it transfers all marbles of a specified color to a new owner.
- **Delete Marble** (DELETE): it deletes the state of a specified marble.
- **Get Marble** (GET): it retrieves a marble record by its key.
- **Query By Owner** (GET): it queries marble records by an owner name.
- **Query By Range** (GET): it retrieves marble records in a specified range of keys.
- **Marble History** (GET): it retrieves the history of a marble.
- **Query Range Page** (GET): it retrieves marble records in a range of keys, with pagination support.

You can use the test messages in [marble.postman_collection.json](marble.postman_collection.json) for end-to-end tests.  The test file can be imported and executed in [postman](https://www.getpostman.com/downloads/).

If you prefer, you can also use the following `curl` commands to invoke the REST APIs.
```
# insert test data
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble21","color":"blue","size":35,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble22","color":"red","size":50,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble23","color":"blue","size":70,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble24","color":"purple","size":80,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble25","color":"purple","size":90,"owner":"tom"}' http://localhost:8989/marble/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble26","color":"purple","size":100,"owner":"tom"}' http://localhost:8989/marble/create
curl -X GET http://localhost:8989/marble/key/marble22
curl -X GET "http://localhost:8989/marble/range?startKey=marble21&endKey=marble25"

# transfer marble ownership
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble22","newOwner":"jerry"}' http://localhost:8989/marble/transfer
curl -H 'Content-Type: application/json' -X PUT -d '{"color":"blue","newOwner":"jerry"}' http://localhost:8989/marble/transfercolor
curl -X GET http://localhost:8989/marble/owner/jerry
curl -X GET "http://localhost:8989/marble/range?startKey=marble21&endKey=marble25"

# delete marble state, not history
curl -X DELETE http://localhost:8989/marble/delete/marble21
curl -X GET http://localhost:8989/marble/history/marble21

# query pagination using page-size and starting bookmark
curl -X GET "http://localhost:8989/marble/rangepage?startKey=marble21&endKey=marble27&pageSize=3"
curl -X GET "http://localhost:8989/marble/rangepage?startKey=marble21&endKey=marble27&pageSize=3&bookmark=marble5"
```

## Notes on GraphQL service
The previous step `make package` generated a `GraphQL` schema file [`metadata.gql`](contract-metadata/metadata.gql), which can be used to implement a GraphQL service to invoke the `marble` chaincode.  Refer to the [`equipment sample`](../equipment) for steps of creating a GraphQL service with zero-code.

## Cleanup the sample fabric network
After you are done testing, you can stop and cleanup the Fabric sample `first-network` as follows:
```
cd $GOPATH//src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```

## Deploy to IBM Cloud
To deploy the `marblle` chaincode to IBM Cloud, it is required to package the chaincode in `.cds` format.  The following script creates [`marble_cc.cds`](marble_cc.cds), which you can deploy to IBM Blockchain Platform.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make package
```
Refer to [fabric-tools](../../fabric-tools) for details about installing chaincode on the IBM Blockchain Platform.

The REST or GraphQL service apps can access the same `marble` chaincode deployed in [IBM Cloud](https://cloud.ibm.com) using the [IBM Blockchain Platform](https://cloud.ibm.com/catalog/services/blockchain-platform-20). The only required update is the network configuration file.  [config_ibp.yaml](../../testdata/config_ibp.yaml) is a sample network configuration that can be used by the REST or GraphQL service app.
