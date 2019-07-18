# marble-private
This is a Flogo version of the [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) sample chaincode, [marbles02_private](https://github.com/hyperledger/fabric-samples/tree/release-1.4/chaincode/marbles02_private) implemented by using a [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-5-0) model.  The model does not require any code, it contains only a JSON model file exported from the TIBCO Flogo® Enterprise.  You can download the prerequisites and then build and deploy the model to a Hyperledger Fabric network as described below.

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.5](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html)
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

## Edit smart contract (opptional)
Skip to the next section if you do not plan to modify the included sample model.

- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabricExtension.zip`](../fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabric.sh`](../zip-fabric.sh)
- Create new Flogo App of name `marble_private` and choose `Import app` to import the model [`marble_private.json`](marble_private.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`marble_private.json`](marble_private.json) to this `marble-private` folder.

## Build and deploy chaincode to Hyperledger Fabric

- In the `marble-private` folder, execute `make create` to generate source code from the flogo model [`marble_private.json`](marble_private.json).
- Execute `make deploy` to deploy the chaincode to the `fabric-samples` chaincode folder.  Note: you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble-private
make create
make deploy
```

## Install and test chaincode using fabric sample first-network
Start Hyperledger Fabric first-network with CouchDB:
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh up -s couchdb
```
Use the `cli` docker container to install and instantiate the `marble_private` chaincode.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble-private
make cli-init
```

Optionally, test the chaincode from `cli` docker container, i.e.,
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble-private
make cli-test
```
You can skip this test, and follow the steps in the next section to build the client app, and then use the client app to execute the same tests.

Note that developers can also use Fabric dev-mode to test chaincode (refer [dev](../marble/dev.md) for more details).

## Edit marble-private-client app (optional)
The marble-private-client is a REST service that invokes the `marble-private` chaincode.  It is implemented as a Flogo model, [`marble_private_client.json`](marble_private_client.json).  Skip to the next section if you do not plan to modify the included sample model.

The client app requires the metadata of the chaincode implemented by the `marble-private`. You can generate the contract metadata [`metadata.json`](contract-metadata/metadata.json) by
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble-private
make package
```

- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabclientExtension.zip`](../../fabclientExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabclient.sh`](../../zip-fabclient.sh)
- Create new Flogo App of name `marble_private_client` and choose `Import app` to import the model [`marble_private_client.json`](marble_private_client.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.
- Open `Connections` tab, find and edit the `marble private client` connector, and set the `Smart coontract metadata file` to the [`metadata.json`](contract-metadata/metadata.json) generated in the above step.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`marble_private_client.json`](marble_private_client.json) to this `marble-private` folder.

## Build and start the marble-private-client
Build and start the client app as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble-private
make create-client
make build-client
make run
```

## Test marble-private-client and the marble-private chaincode
The marble-private-client implements the following REST APIs to invoke corresponding blockchain transactions of the `marble-private` chaincode:
- **Create Marble** (PUT): it creates a new marble.
- **Transfer Marble** (PUT): it transfers a marble to a new owner.
- **Delete Marble** (DELETE): it deletes the state of a specified marble.
- **Get Marble** (GET): it retrieves a marble record by its key.
- **Get Marble Price** (GET): it retrieves a marble's private details by its key.
- **Query By Owner** (GET): it queries marble records by an owner name.
- **Query By Range** (GET): it retrieves marble records in a specified range of keys.

You may use the following commands to invoke the REST APIs.  If you do not like command-line `curl`, you may download and use a REST client tool to submit these REST requests.  For Mac users, the [`Advanced Rest client`](https://install.advancedrestclient.com/install) is pretty user-friendly.

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
curl -X GET "http://localhost:8989/marbleprivate/range?startKey=marble1&endKey=marble4"
curl -X DELETE http://localhost:8989/marbleprivate/delete/marble2
curl -X GET http://localhost:8989/marbleprivate/owner/jerry
```

Note that the operations for `delete` and `price` is allowed by only one of the 2 blockchain member orgs (i.e., org1 only), thus these 2 operations should use a special client configuration such that the client communicates with org1 peers only.  This sample, however, is configured to pick any peer randomly for each request, and so it will fail when a peer node of org2 is used.  You may retry the same request a few times to see the effect of security control for these 2 operations.

## Cleanup the marble-private fabric network
After you are done testing, you can stop and cleanup the Fabric sample `first-network` as follows:
```
cd $GOPATH//src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```

## Deploy to IBM Cloud
This client app can access the same `marble-private` chaincode deployed in [IBM Cloud](https://cloud.ibm.com) using the [IBM Blockchain Platform](https://cloud.ibm.com/catalog/services/blockchain-platform-20).

To deploy the `marble-private` chaincode to IBM Cloud, it is required to package the chaincode in `.cds` format.  The following script creates [`marble_private_cc.cds`](marble_private_cc.cds), which you can deploy to IBM Blockchain Platform.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble-private
make package
```
Refer to [fabric-tools](../../fabric-tools) for details about installing chaincode in the IBM Blockchain Platform.