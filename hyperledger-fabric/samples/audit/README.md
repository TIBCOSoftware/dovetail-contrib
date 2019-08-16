# audit
This example demonstrates the use of [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) for building audit trails across multiple organizations and multiple data domains.  The [TIBCO Cloud AuditSafe](https://www.tibco.com/products/tibco-cloud-auditsafe) uses a similar blockchain service. It uses the [project Dovetail](https://tibcosoftware.github.io/dovetail/) to implement and deploy following 2 components:
- Chaincode for Hyperledger Fabric that implements the business logic for creating and querying audit trails for multiple data domains and owners;
- GraphQL service that end-users can call to submit audit transactions.

Both components are implemented using [Flogo®](https://www.flogo.io/) models by visual programming with zero-code.  The Flogo® models can be created, imported, edited, and/or exported by using [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) or [Dovetail](https://github.com/TIBCOSoftware/dovetail).

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.6](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html).  If you do not have access to `Flogo Enterprise`, you may sign up a trial on [TIBCO CLOUD Integration (TCI)](https://cloud.tibco.com/), or download Dovetail v0.2.0 when it is released.  This sample uses `TIBCO Flogo® Enterprise`, but all models can be imported and edited by using Dovetail v0.2.0 and above.
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
- Create new Flogo App of name `audit` and choose `Import app` to import the model [`audit.json`](audit.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`audit.json`](audit.json) to this `audit` sample folder.

## Build and deploy chaincode to Hyperledger Fabric
- In this `audit` sample folder, execute `make create` to generate the chaincode source code from the flogo model [`audit.json`](audit.json).
- Execute `make deploy` to build and deploy the chaincode to the `fabric-samples` chaincode folder.  Note that you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.
- Execute `make package` to generate `cds` package for cloud deployment, and `metadata` for client apps.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/audit
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
Use `cli` docker container to install and instantiate the `audit_cc` chaincode.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/audit
make cli-init
```

Optionally, test the chaincode from `cli` docker container, i.e.,
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/audit
make cli-test
```
You may skip this test, and follow the steps in the next section to build the GraphQL service, and then use the client service to execute the tests. If you run the `cli` tests, however, it should print out all successful tests with status code `200` and query results if the `audit_cc` chaincode is installed and instantiated successfully on the Fabric network.  

Note that developers can also use Fabric dev-mode to test the chaincode (refer [dev](../marble/dev.md) for more details).  For issues regarding how to work with the Fabric network, please refer the [Hyperledger Fabric docs](https://hyperledger-fabric.readthedocs.io/en/latest/build_network.html).

## Edit audit GraphQL service (optional)
The sample Flogo model, [`audit_client.json`](audit_client.json) is a GraphQL service that invokes the `audit_cc` chaincode.  Skip to the next section if you do not plan to modify the sample model.

The client app requires the metadata of the `audit` chaincode. You can generate the contract metadata [`metadata.json`](contract-metadata/metadata.json) by
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/audit
make package
```
Following are steps to edit or view the GraphQL service models.
- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabclientExtension.zip`](../../fabclientExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabclient.sh`](../../zip-fabclient.sh)
- Create new Flogo App of name `audit_client` and choose `Import app` to import the model [`audit_client.json`](audit_client.json)
- You can then add or update the service implementation using the graphical modeler of the TIBCO Flogo® Enterprise.
- Open `Connections` tab, find and edit the `audit client` connector. Set the `Smart contract metadata file` to the [`metadata.json`](contract-metadata/metadata.json), which is generated in the previous step.  Set the `Network configuration file` and `entity matcher file` to the corresponding files in [`testdata`](../../testdata).
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`audit_client.json`](audit_client.json) to this `audit` sample folder.

## Build and start the audit GraphQL service
Build and start the client app as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/audit
make create-client
make build-client
make run
```

## Test GraphQL service and audit chaincode
The GraphQL service implements the following APIs to invoke corresponding blockchain transactions of the `audit` chaincode:
- **createAudit** (Mutation): It receives a list of audit logs, and creates an audit record on the blockchain for each audit log. Each record state is unquely identified by the combination of blockchain transaction ID the input log record ID.
- **getRecord** (Query): It retrieves an audit record from the blockchain for a specified state key.
- **getRecordsByID** (Query): It retrieves all audit records from the blockchain for a specified log record ID. Note that if a record ID is entered by multiple transactions, multiple blockchain states will be created, and thus this query may return multiple record states for the same log record ID.
- **getRecordsByTxID** (Query): It retrieves audit records from the blockchain for a specified transaction ID of audit creation. If the original `createAudit` transaction contains multiple log records, this query will return all the records if their states have not been deleted.
- **getRecordsByTxTime** (Query): It retrieves audit records from the blockchain for a specified creation timestamp.
- **queryTimeRange** (Query): It retrieves audit records that match a specified data domain and owner, and within a specified range of creation time.
- **deleteRecord** (Mutation): It deletes the current state of a audit record for a specified state key.
- **deleteTransaction** (Mutation): It deletes all current states for a specified creation transaction ID.

You can use the test messages in [audit.postman_collection.json](audit.postman_collection.json) for end-to-end tests.  The test file can be imported and executed in [postman](https://www.getpostman.com/downloads/).

With a few clicks, you can also easily re-create the GraphQL service from scratch. In `TIBCO Flogo® Enterprise`, create a new app, e.g., `my_audit_gql`, choose creating `From GraphQL Schema`, and `browse and upload` the file [`metadata.gql`](contract-metadata/metadata.gql), which is generated previously by `make package`.

This should create 8 Flogo flows based on the chaincode transactions defined in the `metadata`.  You can then edit each flow by adding an activity `fabclient/Fabric Request`, and configure it to call the corresponding `audit` transactions, and map the chaincode response to the `Return` activity.

Once you complete the same model as that in the sample `audit_client.json`, you can export, build and test it as described above.  Note that the default service port is `7879`, although you can make it configurable by defining an `app property` for it.

## Cleanup the sample fabric network
After you are done testing, you can stop and cleanup the Fabric sample `first-network` as follows:
```
cd $GOPATH//src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```

## Deploy to IBM Cloud
To deploy the `audit` chaincode to IBM Cloud, it is required to package the chaincode in `.cds` format.  The script `make package` has already created [`audit_cc.cds`](audit_cc.cds), which you can deploy to IBM Blockchain Platform.  Refer to [fabric-tools](../../fabric-tools) for details about installing chaincode on the IBM Blockchain Platform.

The GraphQL service app can access the same `audit` chaincode deployed in [IBM Cloud](https://cloud.ibm.com) using the [IBM Blockchain Platform](https://cloud.ibm.com/catalog/services/blockchain-platform-20). The only required update is the network configuration file.  [config_ibp.yaml](../../testdata/config_ibp.yaml) is a sample network configuration that can be used by the GraphQL service.
