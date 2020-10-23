# equipment

This example demonstrates the use of [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) for tracking equipment purchasing and installations. This sample demonstrates client services using REST or GraphQL APIs. It uses the [TIBCO Flogo® Enterprise](https://www.tibco.com/products/tibco-flogo) to implement and deploy following 2 components:

- Chaincode for Hyperledger Fabric that implements the business logic for tracking equipment assets on blockchain;
- Client services that end-users can call to submit transactions, i.e., chaincode invocations for equipment tracking. Two equivalent service implementations are provided for demonstration. One service implements REST APIs, and the other implements equivalent GraphQL APIs.

Both components are implemented using [Flogo®](https://www.flogo.io/) models by visual programming with zero-code. The Flogo® models can be created, imported, edited, and/or exported by using [TIBCO Flogo® Enterprise 2.10](https://docs.tibco.com/products/tibco-flogo-enterprise-2-10-0).

This sample also demonstrates the use of Hyperledger Fabric event creation and event listners.

## Prerequisite

Follow the instructions [here](../../development.md) to setup the Dovetail development environment on Mac or Linux.

## Edit smart contract (optional)

Skip to the next section if you do not plan to modify the included chaincode model.

- Start TIBCO Flogo® Enterprise.
- Open http://localhost:8090 in Chrome web browser.
- Create new Flogo App of name `equipment` and choose `Import app` to import the model [`equipment.json`](equipment.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`equipment.json`](equipment.json) to this `equipment` sample folder.

## Build and deploy chaincode to Hyperledger Fabric

- In this `equipment` sample folder, execute `make build` to generate the chaincode source code from the flogo model [`equipment.json`](equipment.json).
- Execute `make deploy` to build and deploy the chaincode to the `fabric-samples` chaincode folder. Note that you may need to set env `FAB_PATH` or edit the [`Makefile`](Makefile) and match the installation folder of `fabric-samples` if it is not downloaded to the default location under `${WORK_PATH}`.

The detailed commands of the above steps are as follows:

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/samples/equipment
make build
make deploy
```

The `build` script uses a `dovetail-tools` docker container to build the chaincode model into a CDS package that can be installed on any fabric network. Refer [README](../../dovetail-tools/README.md) for more details about the docker container and the build script.

## Install and test chaincode using fabric sample first-network

Start Hyperledger Fabric first-network with CouchDB:

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/samples/marble
make start
```

Use `cli` docker container to install and instantiate the `equipment_cc` chaincode.

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/samples/equipment
make cli-init
```

This step installs the chaincode package, `equipment_cc_1.0.cds`, from `${CC_DEPLOY}/equipment_cc` folder, which should be copied by the previous step `make deploy`. The same `CDS` file can be used to install chaincode on any other fabric networks.

Optionally, test the chaincode from `cli` docker container, i.e.,

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/samples/equipment
make cli-test
```

You may skip this test, and follow the steps in the next section to build client apps, and then use the client app to execute the tests. If you run the `cli` tests, however, it should print out 5 successful tests with status code `200` if the `equipment_cc` chaincode is installed and instantiated successfully on the Fabric network.

Note that developers can also use Fabric dev-mode to test the chaincode (refer [dev](../marble/dev.md) for more details). For issues regarding how to work with the Fabric network, please refer the [Hyperledger Fabric docs](https://hyperledger-fabric.readthedocs.io/en/latest/build_network.html).

## Edit equipment REST service (optional)

The sample Flogo model, [`equipment_client.json`](equipment_client.json) is a REST service that invokes the `equipment` chaincode. It also includes an event listener for Fabric block events. Skip to the next section if you do not plan to modify the sample model.

The client app requires the metadata of the `equipment` chaincode. You can generate the contract metadata [`metadata.json`](contract-metadata/metadata.json) by

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/samples/equipment
make metadata
```

Following are steps to edit or view the REST service models.

- Start TIBCO Flogo® Enterprise.
- Open http://localhost:8090 in Chrome web browser.
- Create new Flogo App of name `equipment_client` and choose `Import app` to import the model [`equipment_client.json`](equipment_client.json)
- You can then add or update the service implementation using the graphical modeler of the TIBCO Flogo® Enterprise.
- Open `Connections` tab, find and edit the `equipment client` connector. Set the `Smart contract metadata file` to the [`metadata.json`](contract-metadata/metadata.json), which is generated in the previous step. Set the `Network configuration file` and `entity matcher file` to the corresponding files in [`testdata`](../../testdata).
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`equipment_client.json`](equipment_client.json) to this `equipment` sample folder.

Note: after you import the REST model, check the configuration of the REST trigger. The port should be mapped to `=$property["PORT"]`. Correcct the mapping if it is not imported correctly.

## Build and start the equipment REST service

Build and start the client app as follows:

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/samples/equipment
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

You can use the test messages in [rest.postman_collection.json](rest.postman_collection.json) for end-to-end tests. The test file can be imported and executed in [postman](https://www.getpostman.com/downloads/).

## Implement and test equipment GraphQL service

Simillar to the REST service app, you can import and view the GraphQL service models defined in [equipment_gql.json](equipment_gql.json). Build and start the service by

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/samples/equipment
make build-gql
make run
```

Test the GraphQL service by using the postman test file [graphql.postman_collection.json](graphql.postman_collection.json).

With a few clicks, you can also easily create the GraphQL service from scratch. In `TIBCO Flogo® Enterprise`, create a new app, e.g., `my_equipment_gql`, choose creating `From GraphQL Schema`, and `browse and upload` the file [`metadata.gql`](contract-metadata/metadata.gql), which is generated previously by `make metadata`.

This should create 6 Flogo flows based on the chaincode transactions defined in the `metadata`. You can then edit each flow by adding an activity `fabclient/Fabric Request`, and configure it to call the corresponding `equipment` transactions, and map the chaincode response to the `Return` activity.

Once you complete the same model as that in the sample `equipment_gql.json`, you can export, build and test it as described in the previous section. Note that the default service port is `7879`, although you can make it configurable by defining an `app property` for it.

## Cleanup the sample fabric network

After you are done testing, you can stop and cleanup the Fabric sample `first-network` as follows:

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/samples/marble
make shutdown
```

## Deploy to IBM Cloud

The CDS package, `equipment_cc_1.0.cds`, created by the build script above can be used to deploy to IBM Blockchain Platform. Refer to [fabric-tools](../../fabric-tools) for details about installing chaincode on the IBM Blockchain Platform.

The REST and GraphQL service apps can access the same `equipment` chaincode deployed in [IBM Cloud](https://cloud.ibm.com) using the [IBM Blockchain Platform](https://cloud.ibm.com/catalog/services/blockchain-platform-20). The only required update is the network configuration file. [config_ibp.yaml](../../testdata/config_ibp.yaml) is a sample network configuration that can be used by the REST and GraphQL service apps.
