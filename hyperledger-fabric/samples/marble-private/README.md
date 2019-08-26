# marble-private
This example uses the [project Dovetail](https://tibcosoftware.github.io/dovetail/) to implement the [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) sample chaincode [marbles02_private](https://github.com/hyperledger/fabric-samples/tree/release-1.4/chaincode/marbles02_private).  This sample demonstrates the use of Hyperledger Fabric private collections.  It is implemented using [Flogo®](https://www.flogo.io/) models by visual programming with zero-code.  The Flogo® models can be created, imported, edited, and/or exported by using [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) or [Dovetail](https://github.com/TIBCOSoftware/dovetail).

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.6](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html).  If you do not have access to `Flogo Enterprise`, you may sign up a trial on [TIBCO CLOUD Integration (TCI)](https://cloud.tibco.com/), or download Dovetail v0.2.0.  This sample uses `TIBCO Flogo® Enterprise`, but all models can be imported and edited by using Dovetail v0.2.0 and above.
- [Install Go](https://golang.org/doc/install)
- Clone [Hyperledger Fabric](https://github.com/hyperledger/fabric)
- Download Hyperledger Fabric samples and executables of latest production release as described [here](https://github.com/hyperledger/fabric-samples/tree/release-1.4)
- Download and install [flogo-cli](https://github.com/project-flogo/cli)
- Clone dovetail-contrib with Flogo extension for Hyperledger Fabric

There are different ways to clone these packages.  This document assumes that you have installed these packages under $GOPATH after installing Go, i.e.,
```
go get -u github.com/hyperledger/fabric
cd $GOPATH/src/github.com/hyperledger
curl -sSL http://bit.ly/2ysbOFE | bash -s
export PATH=$GOPATH/src/github.com/hyperledger/fabric-samples/bin:$PATH
go get -u github.com/project-flogo/cli/...
go get -u github.com/TIBCOSoftware/dovetail-contrib
```
Note that the latest version of the Flogo extension for Hyperledger Fabric can be downloaded from the [`develop` branch of the `dovetail-contrib`](https://github.com/TIBCOSoftware/dovetail-contrib/tree/develop), i.e.,
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib
git checkout develop
```

## Edit smart contract (opptional)
Skip to the next section if you do not plan to modify the included sample model.

- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabricExtension.zip`](../fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can generate this `zip` by using the script [`zip-fabric.sh`](../zip-fabric.sh).
- Create new Flogo App of name `marble_private` and choose `Import app` to import the model [`marble_private.json`](marble_private.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`marble_private.json`](marble_private.json) to this `marble-private` sample folder.

Note that when a flogo model is imported to `Flogo® Enterprise v2.6.1`, a `return` activity is automatically added to the end of all branches, which could be an issue if the `return` activity is not at the end of a flow.  Thus, you need to carefully remove the mistakenly added `return` activities after the model is imported.  This issue will be fixed in a later release of the `Flogo® Enterprise`.

## Build and deploy chaincode to Hyperledger Fabric

- In the `marble-private` folder, execute `make create` to generate source code from the flogo model [`marble_private.json`](marble_private.json).
- Execute `make deploy` to deploy the chaincode to the `fabric-samples` chaincode folder.  Note that you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.
- Execute `make package` to generate `cds` package for cloud deployment, and `metadata` for client apps.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble-private
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
Use the `cli` docker container to install and instantiate the `marble_private_cc` chaincode.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble-private
make cli-init
```
Note that this script installs chaincode on 4 peer nodes using the `cli` container.  It is very slow on Mac due to slow volume mounts in the docker desktop for Mac.  The following [solution](https://docs.docker.com/compose/compose-file/#caching-options-for-volume-mounts-docker-for-mac) will speed up the chaincode installation by more than 4 times.
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
sed -i -e "s/github.com\/chaincode.*/github.com\/chaincode:cached/" ./docker-compose-cli.yaml
```

Optionally, test the chaincode from `cli` docker container, i.e.,
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble-private
make cli-test
```
You may skip this test, and follow the steps in the next section to build client apps, and then use the client app to execute the tests. If you run the `cli` tests, however, it should print out all successful tests with status code `200` if the `marble_private_cc` chaincode is installed and instantiated successfully on the Fabric network, except for 2 expected failure tests.

You may see timeout error for the private-details query, which is expected on the first execution when it waits for a new peer to start.  The request timeout can be configured in the app model.

Another failure message `Failed to get marble private details for: marble1` is expected by design. It demonstrates that the private detail data is automatically deleted from the private collection after 3 new blocks, which is configured in the private collection definition file `$GOPATH/src/github.com/hyperledger/fabric-samples/chaincode/marbles02_private/collections_config.json`.

Note that developers can also use Fabric dev-mode to test the chaincode (refer [dev](../marble/dev.md) for more details).  For issues regarding how to work with the Fabric network, please refer the [Hyperledger Fabric docs](https://hyperledger-fabric.readthedocs.io/en/latest/build_network.html).

## Edit marble-private REST service (optional)
The sample Flogo model, [`marble_private_client.json`](marble_private_client.json) is a REST service that invokes the `marble-private` chaincode.  Skip to the next section if you do not plan to modify the sample model.

The client app requires the metadata of the `marble-private` chaincode. You can generate the contract metadata [`metadata.json`](contract-metadata/metadata.json) by
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble-private
make package
```
Following are steps to edit or view the REST service models.
- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabclientExtension.zip`](../../fabclientExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can generate this `zip` by using the script [`zip-fabclient.sh`](../../zip-fabclient.sh)
- Create new Flogo App of name `marble_private_client` and choose `Import app` to import the model [`marble_private_client.json`](marble_private_client.json)
- You can then add or update the service implementation using the graphical modeler of the TIBCO Flogo® Enterprise.
- Open `Connections` tab, find and edit the `marble private client` connector.  Set the `Smart contract metadata file` to the [`metadata.json`](contract-metadata/metadata.json), which is generated in the previous step.  Set the `Network configuration file` and `entity matcher file` to the corresponding files in [`testdata`](../../testdata).
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`marble_private_client.json`](marble_private_client.json) to this `marble-private` sample folder.

## Build and start the marble-private REST service
Build and start the client app as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble-private
make create-client
make build-client
make run
```

## Test marble-private REST service and marble-private chaincode
The REST service implements the following APIs to invoke corresponding blockchain transactions of the `marble-private` chaincode:
- **Create Marble** (PUT): it creates a new marble.
- **Transfer Marble** (PUT): it transfers a marble to a new owner.
- **Delete Marble** (DELETE): it deletes the state of a specified marble.
- **Get Marble** (GET): it retrieves a marble record by its key.
- **Get Marble Price** (GET): it retrieves a marble's private details by its key.
- **Query By Owner** (GET): it queries marble records by an owner name.
- **Query By Range** (GET): it retrieves marble records in a specified range of keys.

You can use the test messages in [marble-private.postman_collection.json](marble-private.postman_collection.json) for end-to-end tests.  The test file can be imported and executed in [postman](https://www.getpostman.com/downloads/).

If you prefer, you can also use the following `curl` commands to invoke the REST APIs.
```
# insert test data
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble11","color":"blue","size":35,"owner":"tom","price":99}' http://localhost:8989/marbleprivate/create
curl -X GET http://localhost:8989/marbleprivate/key/marble11
curl -X GET http://localhost:8989/marbleprivate/price/marble11

# more inserts and transfer owner, test purge of private marble1 after 3 blocks
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble12","color":"red","size":50,"owner":"tom","price":199}' http://localhost:8989/marbleprivate/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble13","color":"blue","size":70,"owner":"tom","price":299}' http://localhost:8989/marbleprivate/create
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble12","owner":"jerry"}' http://localhost:8989/marbleprivate/transfer
# marble1 pricing is still available
curl -X GET http://localhost:8989/marbleprivate/price/marble1
curl -H 'Content-Type: application/json' -X PUT -d '{"name":"marble13","owner":"jerry"}' http://localhost:8989/marbleprivate/transfer
# marble1 private detail is purged after 3 blocks, so this returns error
curl -X GET http://localhost:8989/marbleprivate/price/marble11

# test query and delete
curl -X GET http://localhost:8989/marbleprivate/owner/jerry
curl -X GET "http://localhost:8989/marbleprivate/range?startKey=marble11&endKey=marble14"
curl -X DELETE http://localhost:8989/marbleprivate/delete/marble12
curl -X GET http://localhost:8989/marbleprivate/owner/jerry
```

Note that the operations for `delete` and `price` are allowed by only one of the 2 blockchain member orgs (i.e., org1 only), thus these 2 operations will fail if the REST service sends the request to an org2 peer.  You may retry the request a few times until it succeeds on an org1 peer.  The next Dovetail release will support endpoint override, and so these requests can be routed to an org1 peer only.

## Notes on GraphQL service
The previous step `make package` generated a `GraphQL` schema file [`metadata.gql`](contract-metadata/metadata.gql), which can be used to implement a GraphQL service to invoke the `marble_private` chaincode.  Refer to the [`equipment sample`](../equipment) for steps of creating a GraphQL service with zero-code.

## Cleanup the sample fabric network
After you are done testing, you can stop and cleanup the Fabric sample `first-network` as follows:
```
cd $GOPATH//src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```

## Deploy to IBM Cloud
To deploy the `marblle_private` chaincode to IBM Cloud, it is required to package the chaincode in `.cds` format.  The following script creates [`marble_private_cc.cds`](marble_private_cc.cds), which you can deploy to IBM Blockchain Platform.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble-private
make package
```
Refer to [fabric-tools](../../fabric-tools) for details about installing chaincode on the IBM Blockchain Platform.

The REST or GraphQL service apps can access the same `marble_private` chaincode deployed in [IBM Cloud](https://cloud.ibm.com) using the [IBM Blockchain Platform](https://cloud.ibm.com/catalog/services/blockchain-platform-20). The only required update is the network configuration file.  [config_ibp.yaml](../../testdata/config_ibp.yaml) is a sample network configuration that can be used by the REST or GraphQL service app.
