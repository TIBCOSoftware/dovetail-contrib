# equipment
This example demonstrates the use of [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) to trace the equipment purchasing and installation process.  It uses the project Dovetail to implement and deploy the following 2 components:
- Chaincode for Hyperledger Fabric that implements the business logic for tracking equipment assets on blockchain;
- Client services that end-users can call to submit transactions, i.e., chaincode invocations for equipment tracking.  Two equivalent service implementations are provided for demonstration. One exposes REST APIs, and the other exposes GraphQL APIs.
Both components are implemented using Flogo® models by visual programming with zero-code.  The Flogo® models can be created, imported, edited, and/or exported by using [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) or [Dovetail](https://github.com/TIBCOSoftware/dovetail)

This sample also demonstrates the use of Hyperledger Fabric events.

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.6](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html) or Download Dovetail when v0.2.0 is released.  This sample uses `TIBCO Flogo® Enterprise`, but all models can be imported and edited by using Dovetail v0.2.0.
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
Note that the latest version of the Flogo extension for Hyperledger Fabric is required by this sample, and it is in the [`fabric-extension` branch of the `dovetail-contrib`](https://github.com/TIBCOSoftware/dovetail-contrib/tree/issue-36/fabric-extension).

Bootstrap fabric-samples
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples
./scripts/bootstrap.sh
```

## Edit smart contract (optional)
Skip to the next section if you do not plan to modify the included chaincode model.

- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabricExtension.zip`](../../fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabric.sh`](../../zip-fabric.sh)
- Create new Flogo App of name `equipment` and choose `Import app` to import the model [`equipment.json`](equipment.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`equipment.json`](equipment.json) to this `equipment` sample folder.

## Build and deploy chaincode to Hyperledger Fabric
- In this `equipment` sample folder, execute `make create` to generate the chaincode source code from the flogo model [`equipment.json`](equipment.json).
- Execute `make deploy` to build and deploy the chaincode to the `fabric-samples` chaincode folder.  Note that you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.
- Execute `make package` to generate `cds` package for cloud deployment, and `metadata` for client apps.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/equipment
make create
make deploy
make package
```

## Install and test chaincode using fabric sample first-network
Start Hyperledger Fabric first-network with CouchDB:
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh up -s couchdb
```
Use `cli` docker container to install and instantiate the `equipment_cc` chaincode.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/equipment
make cli-init
```

Optionally, test the chaincode from `cli` docker container, i.e.,
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/equipment
make cli-test
```
You may skip this test, and follow the steps in the next section to build client apps, and then use the client app to execute the tests. If you run the `cli` tests, however, it should print out 5 successful tests with status code `200` if the `equipment` chaincode is installed and instantiated successfully on the Fabric network.  

Note that developers can also use Fabric dev-mode to test the chaincode (refer [dev](../marble/dev.md) for more details).  For issues regarding how to work with the Fabric network, please refer the [Hyperledger Fabric docs](https://hyperledger-fabric.readthedocs.io/en/latest/build_network.html).

## Edit equipment REST service (optional)
The sample Flogo model, [`equipment_client.json`](equipment_client.json) is a REST service that invokes the `equipment` chaincode.  It also includes an event listener for Fabric block events.  Skip to the next section if you do not plan to modify the sample model.

The client app requires the metadata of the `equipment` chaincode. You can generate the contract metadata [`metadata.json`](contract-metadata/metadata.json) by
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/equipment
make package
```
Following are steps to edit or view the REST service models.
- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabclientExtension.zip`](../../fabclientExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabclient.sh`](../../zip-fabclient.sh)
- Create new Flogo App of name `equipment_client` and choose `Import app` to import the model [`equipment_client.json`](equipment_client.json)
- You can then add or update the service implementation using the graphical modeler of the TIBCO Flogo® Enterprise.
- Open `Connections` tab, find and edit the `equipment client` connector. Set the `Smart coontract metadata file` to the [`metadata.json`](contract-metadata/metadata.json), which is generated in the previous step.  Set the `Network configuration file` and `entity matcher file` to the corresponding files in [`testdata`](../../testdata).
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`equipment_client.json`](equipment_client.json) to this `equipment` sample folder.

## Build and start the equipment REST service
Build and start the client app as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/equipment
make create-client
make build-client
make run
```

## Test REST service and equipment chaincode
The REST service implements the following APIs to invoke corresponding blockchain transactions of the `equipment` chaincode:
- **submitPO** (PUT): It receives a PO, and creates an equipment asset on the blockchain, or emits a Fabric event if operation fails.
- **receiveAsset** (PUT): it records the reception of an equipment, and updates the corresponding equipment asset on the blockchain, or emits a Fabric event if operation fails.
- **installAsset** (PUT): it records the installation of an equipment, and updates the corresponding equipment asset on the blockchain, or emits a Fabric event if operation fails.
- **receiveInvoice** (PUT): it records the reception of an invoice, and updates the corresponding equipment asset on the blockchain, or emits a Fabric event if operation fails.
- **updateAsset** (PUT): it updates a specified equipment asset on the blockchain, or emits a Fabric event if operation fails.
- **queryAsset** (GET): it fetches current state of a specified equipment asset on the bllockchain.
- **eventListener**: it is not a REST API. It listens to Fabric block events, and prints out the content of a block when it is committed to blockchain.

You can use the test messages in [rest.postman_collection.json](rest.postman_collection.json) for end-to-end tests.  The test file can be imported and executed in [postman](https://www.getpostman.com/downloads/).

## Implement and test equipment GraphQL service
Simillar to the REST service app, you can import and view the GraphQL service models defined in [equipment_gql.json](equipment_gql.json).  Build and start the service by
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/equipment
make create-gql
make build-gql
make run-gql
```
Test the GraphQL service by using the postman test file [graphql.postman_collection.json](graphql.postman_collection.json).

With a few clicks, you can also easily create the GraphQL service from scratch. In `TIBCO Flogo® Enterprise`, create a new app, e.g., `my_equipment_gql`, choose creating `From GraphQL Schema`, and `browse and upload` the file [`metadata.gql`](contract-metadata/metadata.gql), which is generated previously by `make package`.

This should create 6 Flogo flows based on the chaincode transactions defined in the `metadata`.  You can then edit each flow by adding an activity `fabclient/Fabric Request`, and configure it to call the corresponding `equipment` transactions, and map the chaincode response to the `Return` activity.

Once you complete the same model as that in the sample `equipment_gql.json`, you can export, build and test it as described in the previous section.  Note that the default service port is `7879`, although you can make it configurable by defining an `app property` for it.

## Cleanup the sample fabric network
After you are done testing, you can stop and cleanup the Fabric sample `first-network` as follows:
```
cd $GOPATH//src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```

## Deploy to IBM Cloud
To deploy the `equipment` chaincode to IBM Cloud, it is required to package the chaincode in `.cds` format.  The following script creates [`equipment_cc.cds`](equipment_cc.cds), which you can deploy to IBM Blockchain Platform.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/equipment
make package
```
Refer to [fabric-tools](../../fabric-tools) for details about installing chaincode on the IBM Blockchain Platform.

The REST and GraphQL service apps can access the same `equipment` chaincode deployed in [IBM Cloud](https://cloud.ibm.com) using the [IBM Blockchain Platform](https://cloud.ibm.com/catalog/services/blockchain-platform-20). The only required update is the network configuration file.  [config_ibp.yaml](../../testdata/config_ibp.yaml) is a sample network configuration that can be used by the REST and GraphQL service apps.
