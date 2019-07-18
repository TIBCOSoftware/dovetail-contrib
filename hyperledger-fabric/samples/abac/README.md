# abac (Attribute Based Access Control)

This is a Hyperledger Fabric chaincode for demonstrating ABAC. It is implemented by using a [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) model.  The model does not require any code, it contains only a JSON model file exported from the TIBCO Flogo® Enterprise.  You can download the prerequisites and then build and deploy the model to a Hyperledger Fabric network as described below.

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.6](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html)
- [Install Go](https://golang.org/doc/install)
- Clone [Hyperledger Fabric](https://github.com/hyperledger/fabric)
- Clone [Hyperledger Fabric Samples](https://github.com/hyperledger/fabric-samples)
- Install [Fabric CA binaries](https://hyperledger-fabric-ca.readthedocs.io/en/release-1.4/users-guide.html)
- Download and install [flogo-cli](https://github.com/TIBCOSoftware/flogo-cli)
- Clone dovetail-contrib with this Flogo extension

There are different ways to clone these packages.  This document assumes that you have installed these packages under $GOPATH after installing Go, i.e.,
```
go get -u github.com/hyperledger/fabric
go get -u github.com/hyperledger/fabric-samples
go get -u github.com/hyperledger/fabric-ca/cmd/...
go get -u github.com/TIBCOSoftware/flogo-cli/...
go get -u github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric
```
Bootstrap fabric-samples
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples
./scripts/bootstrap.sh
```

## Edit smart contract (optional)
Skip to the next section if you do not plan to modify the included sample model.

- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabricExtension.zip`](../../fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabric.sh`](../../zip-fabric.sh)
- Create new Flogo App of name `abac_app` and choose `Import app` to import the model [`abac_app.json`](abac_app.json)
- You can then add or update the flows using the graphical modeler of the TIBCO Flogo® Enterprise.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`abac_app.json`](abac_app.json) to this `abac` folder.

## Build and deploy chaincode to Hyperledger Fabric

- In this `abac` folder, execute `make create` to generate source code from the flogo model [`abac_app.json`](abac_app.json).
- Execute `make deploy` to build and deploy the chaincode to the `fabric-samples` chaincode folder.  Note: you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/abac
make create
make deploy
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

You can skip this test, and follow the steps in the next section to build the client app, and then use the client app to execute more interesting tests.

Note that developers can also use Fabric dev-mode to test chaincode (refer [dev](../marble/dev.md) for more details).

## Edit abac-client app (optional)
The abac-client is a REST service that invokes the `abac_app` chaincode.  It is implemented as a Flogo model, [`abac_client.json`](abac_client.json).  Skip to the next section if you do not plan to modify the included sample model.

The client app requires the metadata of the chaincode implemented by the `abac-app`. You can generate the contract metadata [`metadata.json`](contract-metadata/metadata.json) by
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/abac
make package
```

- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabclientExtension.zip`](../../fabclientExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabclient.sh`](../../zip-fabclient.sh)
- Create new Flogo App of name `abac_client` and choose `Import app` to import the model [`abac_client.json`](abac_client.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.
- Open `Connections` tab, find and edit the `abac client` connector, and set the `Smart coontract metadata file` to the [`metadata.json`](contract-metadata/metadata.json) generated in the above step.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`abac_client.json`](abac_client.json) to this `abac` folder.

## Build and start the abac-client
Build and start the client app as follows
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/abac
make create-client
make build-client
make run
```

## Test abac-client and abac-app chaincode
The `abac_client` implements a simple REST API that receives the name of a test user and its org, and use the user to invoke the `check-abac` chaincode transaction.  The following requests should succeed for users `Alice@org1` and `Bob@org2`, but fail for user `User1@org2`.
```
curl -X GET http://localhost:8989/abac/org1/Alice
curl -X GET http://localhost:8989/abac/org2/Bob
curl -X GET http://localhost:8989/abac/org2/User1
```

## Cleanup the fabric network
After you are done testing, you can stop and cleanup the Fabric sample `first-network` as follows:
```
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```
