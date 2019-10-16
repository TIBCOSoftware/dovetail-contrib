# abac (Attribute Based Access Control)
This example uses the [project Dovetail](https://tibcosoftware.github.io/dovetail/) to demonstrate Attribute Based Access Control (ABAC) in the [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric). It is implemented using [Flogo®](https://www.flogo.io/) models by visual programming with zero-code.  The Flogo® models can be created, imported, edited, and/or exported by using [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) or [Dovetail](https://github.com/TIBCOSoftware/dovetail).

## Prerequisite
Follow the instructions [here](../../development.md) to setup the Dovetail development environment on Mac or Linux.

## Edit smart contract (optional)
Skip to the next section if you do not plan to modify the included sample model.

- Start TIBCO Flogo® Enterprise or Dovetail.
- Open http://localhost:8090 in Chrome web browser.
- Create new Flogo App of name `abac_app` and choose `Import app` to import the model [`abac_app.json`](abac_app.json)
- You can then add or update the flows using the graphical modeler of the TIBCO Flogo® Enterprise.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`abac_app.json`](abac_app.json) to this `abac` sample folder.

Note that when a flogo model is imported to `Flogo® Enterprise v2.6.1`, a `return` activity is automatically added to the end of all branches, which could be an issue if the `return` activity is not at the end of a flow.  Thus, you need to carefully remove the mistakenly added `return` activities after the model is imported.  This issue will be fixed in a later release of the `Flogo® Enterprise`.

## Build and deploy chaincode to Hyperledger Fabric
- In this `abac` sample folder, execute `make create` to generate source code from the flogo model [`abac_app.json`](abac_app.json).
- Execute `make deploy` to build and deploy the chaincode to the `fabric-samples` chaincode folder.  Note that you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.
- Execute `make package` to generate `cds` package for cloud deployment, and `metadata` for client apps.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/abac
make create
make deploy
make package
```

## Install and test chaincode using fabric sample first-network
Start Hyperledger Fabric first-network and create users for ABAC tests:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/abac
make start-fn
```
This script will start the sample first-network with CA servers, and then use the CA servers to create 2 new users, Alice of Org1 and Bob of Org2. Both users's certificates will contain an attribute `abac.init = true`, which is used by the chaincode for user authorization.

Use `cli` docker container to install and instantiate the `abac_cc` chaincode.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/abac
make cli-init
```
Optionally, test the chaincode from `cli` docker container, i.e.,
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/abac
make cli-test
```
This test is expected to fail, because it uses the `Admin` user of Org1, whose certificate does not contain the `abac.init` attribute.

You may skip this test, and follow the steps in the next section to build the client app, and then use the client app to execute more interesting tests.

Note that developers can also use Fabric dev-mode to test chaincode (refer [dev](../marble/dev.md) for more details).  For issues regarding how to work with the Fabric network, please refer the [Hyperledger Fabric docs](https://hyperledger-fabric.readthedocs.io/en/latest/build_network.html).

## Edit abac REST service (optional)
The sample Flogo model, [`abac_client.json`](abac_client.json) is a REST service that invokes the `abac_app` chaincode.  Skip to the next section if you do not plan to modify the sample model.

The client app requires the metadata of the `abac-app` chaincode. You can generate the contract metadata [`metadata.json`](contract-metadata/metadata.json) by
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/abac
make package
```
Following are steps to edit or view the REST service models.
- Start TIBCO Flogo® Enterprise or Dovetail.
- Open http://localhost:8090 in Chrome web browser.
- Create new Flogo App of name `abac_client` and choose `Import app` to import the model [`abac_client.json`](abac_client.json)
- Edit `Settings` of the REST trigger to set `port` to `=$property["PORT"]`
- You can then add or update service implementation using the graphical modeler of the TIBCO Flogo® Enterprise.
- Open `Connections` tab, find and edit the `abac client` connector. Set the `Smart contract metadata file` to the [`metadata.json`](contract-metadata/metadata.json) generated in the previous step. Set the `Network configuration file` and `entity matcher file` to the corresponding files in [`testdata`](../../testdata).
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`abac_client.json`](abac_client.json) to this `abac` sample folder.

## Build and start the abac REST service
Build and start the client app as follows
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/abac
make create-client
make build-client
make run
```

## Test abac REST service and abac_app chaincode
The REST service implements a simple API that receives the name of a test user and its org, and uses the user to invoke the `check_abac` chaincode transaction.  The following requests should succeed for users `Alice@org1` and `Bob@org2`, but fail for user `User1@org2`.
```
curl -X GET http://localhost:8989/abac/org1/Alice
curl -X GET http://localhost:8989/abac/org2/Bob
curl -X GET http://localhost:8989/abac/org2/User1
```

## Cleanup the sample fabric network
After you are done testing, you can stop and cleanup the Fabric sample `first-network` as follows:
```
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```
