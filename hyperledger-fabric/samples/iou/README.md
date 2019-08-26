# iou
This sample implements a simplified cross-border payment system similar to [Ripple](https://www.ripple.com/files/ripple_product_overview.pdf). Although it is a simplified network, it implements the core blockchain operations for secure cross-border fund transfer with zero-code, thanks to the visual modeling environment of the [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) and the blockchain platform of [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric).  The Flogo® models in this sample can be created, imported, edited, and/or exported by using [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) or [Dovetail](https://github.com/TIBCOSoftware/dovetail).

## Use case
Alice, a customer of a bank in Europe wants to send a payment to Bob, a customer of a bank in USA.  Although the parties involved may not trust each other, we can find a chain of intermediaries with 1-to-1 trust relationships, and so that [IOU](http://www.businessdictionary.com/definition/IOU.html)s can be exchanged along the trusted path, resulting in the specified payment amount withdrawn from Alice's account at the Euro Bank, and deposited to Bob's account at the US Bank. 

Blockchain distributed ledger is a perfect technology to record such end-to-end payment flow immutably and with cryptographic security. For simplicity of the sample, we assume that there are only 2 banks, one in Europe and the other in US, and they trust each other.

## Design
The system has 3 types of actors:
 - Bank service provider, who issues and redeems IOUs for its customers, and exchanges IOUs of different currencies according to an exchange rate.  The provider may be a bank, or a service unit within a bank, or a trusted service provider.  For simplicity, this sample assumes that each provider uses only one currency, and the exchange rate is a fixed constant.  We name 2 providers as `EURBank` and `USDBank`.  The exchange rate is configured as 1 EUR = 1.1 USD.
 - User account, that is associated with a bank service provider, and holds a balance of fund used to purchase IOUs from the associated bank service provider.  Each account is identified by a pair of crypto key and certificate.  The account name is an alias that does not necessarily match the true identity of the account owner, and thus the transactions on the blockchain are pseudonymous.
 - Network operator, who operates peer nodes of the Hyperledger Fabric network, and creates crypto keys for bank administrator and user accounts.  Multiple operators may be associated with different business entities.  An operator may be part of the same bank as a bank service provider, or an independent blockchain infrastructure provider. 

Hyperledger Fabric network config:
 - The first network operator, Org1, runs 2 peer nodes for `EURBank`, and `CA` server used to generate 3 crypto key pairs for `EURBankAdmin`, `Alice` and `Bob`, respectively.  The generated certificates contain attributes used for chaincode authorization.
 - The second network operator, Org2, runs 2 peer nodes for `USDBank`, and `CA` server used to generate 3 crypto key pairs for `USDBankAdmin`, `Carol` and `David`, respectively.  The generated certificates also contain attributes used for chaincode authorization.
 - Private collections `EURBankTransactions` and `USDBankTransactions` are configured to store balance updates for acounts of `EURBank` and `USDBank` respectively.  Both operator orgs can access these 2 private collections.
 - Private collections `EURBankAccounts` and `USDBankAccounts` are configured to store current state of user accounts of `EURBank` and `USBBank` respectively.  Each of these private collections is exclusively accessible by only its operator org.
 - CouchDB is configured to store the current state and full history of IOU's, and indexes are defined to support rich queries on IOU's.

 Data object definitions:

 | IOU               | Account          | Transaction     |
 | ----------------- | ---------------- | --------------- |
 | id:        string | name:     string | txID:    string | 
 | issuer:    string | bank:     string | txTime:  string |
 | issueDate: string | balance:  float  | account: string |
 | owner:     string | currency: string | amount:  float  |
 | amount:    float  |                  | iouRef:  string |
 | currency:  string |                  |                 |

## Basic IOU operations and rules
1. issue(bank, owner, amount)
   - Actions:
     - Reduce owner's account-balance by specified amount in the currency of the specified bank;
     - Create IOU issued by the bank to the owner with specified ammount in the bank's currency;
     - Record transaction for negative balance change of the owner's account;
     - Record transaction for positive debt increase of the bank.
   - Rules:
     - Reject the request if requestor's certificate does not match the owner;
     - Reject the request if the owner does not have an account with the bank;
     - Reject the request if the owner's account does not have enough balance.
2. buy(bank, owner, iou)
   - Actions:
     - Reduce owner's account-balance by IOU's amount converted to the bank's currency according to the currency exchange rate;
     - Change the IOU's owner to the specified new owner;
     - Record transaction for negative balance change of the owner's account;
   - Rules:
     - Reject the request if the specified IOU does not exist;
     - Reject the request if the specified IOU is not owned by the specified bank;
     - Reject the request if the requestor's certificate does not match the owner;
     - Reject the request if owner does not have an account with the bank;
     - Reject the request if the owner's account does not have enough balance.
3. transfer(iou, newOwner)
   - Actions:
     - Change the IOU's owner to the specified new owner.
   - Rules:
     - Reject the request if the specified IOU does not exist;
     - Reject the request if the requestor's certificate does not match the original owner of the IOU.
4. exchange(iou, bank)
   - Actions:
     - Change the IOU's owner to the specified bank;
     - Create IOU issued by the bank to the IOU's owner with amount converted to the bank's currency according to the exchange rate;
     - Record transaction for positive debt increase of the bank.
   - Rules:
     - Reject the request if the specified IOU does not exist;
     - Reject the request if the requestor's certificate does not match the IOU's original owner;
     - Reject the request if the IOU's currency is the same the bank's currency (i.e., no need to exchange).
5. redeem(iou, bank)
   - Actions:
     - Delete the specified IOU;
     - Increase owner's account-balance by amount of the IOU;
     - Record transaction for positive balance increase of the owner's account;
     - Record transaction for negative debt change of the bank.
   - Rules:
     - Reject the request if the specified IOU does not exist;
     - Reject the request if IOU is not issued by the bank;
     - Reject the request if the requestor's certificate does not match the IOU's owner;
     - Reject the request if IOU's owner does not have an account with the bank.

## Other chaincode operations:
Composite operation for finding or creating an equivalent IOU with the specified amount in the currency of a receiver's bank:

6. send(sender, senderBank, receiverBank, amount)
   - Actions:
     - If senderBank is the same as receiverBank, call `issue` to create IOU issued by the senderBank to the sender with specified amount;
     - Otherwise, search for IOU issued by receiverBank, owned by senderBank, with the specified amount
       - If found, call `buy` to get the IOU transferred to the sender;
       - If not found, call `issue` to create IOU issued by the senderBank with amount converted to senderBank's currency according to exchange-rate.

Account management operations are also required, and they are better packaged as a separate chaincode because they require different endorsement policies. However, for simplicity of the sample, we implement only a single operation for creating accounts, and package it in the same chaincode with IOU operations.

7. createAccount(name, bank, balance)
   - Actions:
     - create an account for an specified name at a bank with an initial balance in the bank's currency.
   - Rules:
     - Reject the request if the requestor is not the bank's admin;
     - Reject the request if an account with the same name already exists in the bank.

## Client operations:
A client app is implemented to send requests to the blockchain and verify the results.  It implements a GraphQL service interface.  Although this client app implements more test operations, only the following operations are needed to support the cross-border payment process:

1. Mutation createAccount(name, bank, balance): It initializes user accounts;
2. Mutation send(senderBank, sender, receiverBank, receiver, amount): It processes the sender's request to pay the receiver the specified amount in receiverBank's currency.  It orchestrates the process by making the following calls to the chaincode:
   - Use sender credential to call the composite chaincode operation: `send(sender, senderBank, receiverBank, amount)`;
   - Use sender credential to call chaincode operation with the returned IOU: `transfer(iou, receiver)`;
   - If IOU is issued by the senderBank:
     - Use receiver credential to call chaincode operation: `exchange(iou, receiverBank)`;
     - Use receiver credential to call chaincode operation with the new IOU: `redeem(newIOU, receiverBank)`;
   - Otherwise, if IOU is issued by the receiverBank:
     - Use receiver credential to call chaincode operation: `redeem(iou, receiverBank)`.
3. Query getBankAccounts(bank): It returns the balances of all user accounts of the specified bank;
4. Query getAccountTransactions(name|bank, bank): It returns all transactions of a user account or a bank;
5. Query getIOUHistory(iou): It returns the history of a specified IOU.

The file [iou.postman_collection.json](iou.postman_collection.json) contains sample GraphQL test messages that can be viewed and executed in [postman](https://www.getpostman.com/downloads/).

## Modeling with TIBCO Cloud Integration (TCI)
If you are already a subscriber of [TIBCO Cloud Integration (TCI)](https://cloud.tibco.com/), or you plan to sign-up for a TCI trial, you can view or edit this app by using a Chrome browser.  Refer to [Modeling with TCI](../../tci) for more detailed instructions.

## Development Prerequisite
If you want to set up development environment and execute tests, you can install these prerequisites and follow the next sections.
- Download [TIBCO Flogo® Enterprise 2.6](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html).  If you do not have access to `Flogo Enterprise`, you may sign up a trial on [TIBCO CLOUD Integration (TCI)](https://cloud.tibco.com/), or download Dovetail v0.2.0.  This sample is edited using `TIBCO Flogo® Enterprise`, but all models can be imported and edited by using Dovetail v0.2.0.
- [Install Go](https://golang.org/doc/install)
- Clone [Hyperledger Fabric](https://github.com/hyperledger/fabric)
- Install [Fabric CA binaries](https://hyperledger-fabric-ca.readthedocs.io/en/release-1.4/users-guide.html)
- Download Hyperledger Fabric samples and executables of latest production release as described [here](https://github.com/hyperledger/fabric-samples/tree/release-1.4)
- Download and install [flogo-cli](https://github.com/project-flogo/cli)
- Clone dovetail-contrib with Flogo extension for Hyperledger Fabric

There are different ways to clone these packages.  This document assumes that you have installed these packages under $GOPATH after installing Go, i.e.,
```
go get -u github.com/hyperledger/fabric
go get -u github.com/hyperledger/fabric-ca/cmd/...
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

## Edit smart contract (optional)
Skip to the next section if you do not plan to modify the included chaincode model.

- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabricExtension.zip`](../../fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can generate this `zip` by using the script [`zip-fabric.sh`](../../zip-fabric.sh).
- Create new Flogo App of name `iou` and choose `Import app` to import the model [`iou.json`](iou.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`iou.json`](iou.json) to this `iou` sample folder.

Note that when a flogo model is imported to `Flogo® Enterprise v2.6.1`, a `return` activity is automatically added to the end of all branches, which could be an issue if the `return` activity is not at the end of a flow.  Thus, you need to carefully remove the mistakenly added `return` activities after the model is imported.  This issue will be fixed in a later release of the `Flogo® Enterprise`.

## Build and deploy chaincode to Hyperledger Fabric
- In this `iou` sample folder, execute `make create` to generate the chaincode source code from the flogo model [`iou.json`](iou.json).
- Execute `make deploy` to build and deploy the chaincode to the `fabric-samples` chaincode folder.  Note that you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.
- Execute `make package` to generate `cds` package for cloud deployment, and `metadata` for client apps.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/iou
make create
make deploy
make package
```

## Install and test chaincode using fabric sample first-network
Start Hyperledger Fabric first-network with CouchDB, and create crypto key-pairs for bank-admin and test user-accounts
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/iou
make start-fn
```
Use `cli` docker container to install and instantiate the `iou_cc` chaincode.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/iou
make cli-init
```
Note that this script installs chaincode on 4 peer nodes using the `cli` container.  It is very slow on Mac due to slow volume mounts in the docker desktop for Mac.  The following [solution](https://docs.docker.com/compose/compose-file/#caching-options-for-volume-mounts-docker-for-mac) will speed up the chaincode installation by more than 4 times.
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
sed -i -e "s/github.com\/chaincode.*/github.com\/chaincode:cached/" ./docker-compose-cli.yaml
```
By configuring the `chaincode` volume in `cli` container as `cached`, the chaincode installation time can be reduced from 157 seconds to 37 seconds.

## Edit iou GraphQL service (optional)
The sample Flogo model, [`iou_client.json`](iou_client.json) is a GraphQL service that invokes the `iou_cc` chaincode.  Skip to the next section if you do not plan to modify the sample model.

The client app requires the metadata of the `iou` chaincode. You can generate the contract metadata [`metadata.json`](contract-metadata/metadata.json) by
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/iou
make package
```
Following are steps to edit or view the GraphQL service models.
- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabclientExtension.zip`](../../fabclientExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can generate this `zip` by using the script [`zip-fabclient.sh`](../../zip-fabclient.sh).
- Create new Flogo App of name `iou_client` and choose `Import app` to import the model [`iou_client.json`](iou_client.json)
- You can then add or update the service implementation using the graphical modeler of the TIBCO Flogo® Enterprise.
- Open `Connections` tab, find and edit the `iou client` connector. Set the `Smart contract metadata file` to the [`metadata.json`](contract-metadata/metadata.json), which is generated in the previous step.  Set the `Network configuration file` and `entity matcher file` to the corresponding files in [`testdata`](../../testdata).
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`iou_client.json`](iou_client.json) to this `iou` sample folder.

## Build and start the iou GraphQL service
Build and start the client app as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/iou
make create-client
make build-client
make run
```

## Test GraphQL service and iou chaincode
You can use the test messages in [iou.postman_collection.json](iou.postman_collection.json) for end-to-end tests.  The test file can be imported and executed in [postman](https://www.getpostman.com/downloads/).

With a few clicks, you can also easily re-create the GraphQL service from scratch. In `TIBCO Flogo® Enterprise`, create a new app, e.g., `my_iou_gql`, choose creating `From GraphQL Schema`, and `browse and upload` the file [`metadata.gql`](contract-metadata/metadata.gql), which is generated previously by `make package`.

This should create 11 Flogo flows based on the chaincode transactions defined in the `metadata`.  You can then edit each flow by adding an activity `fabclient/Fabric Request`, and configure it to call the corresponding `iou` transactions, and map the chaincode response to the `Return` activity. Note that the `send` operation is a little more complex because it is an orchestration process that makes multiple calls to the chaincode.

Once you complete the model similar to the sample file `iou_client.json`, you can export, build and test it as described above.  Note that the default service port is `7879`, although you can make it configurable by defining an `app property` for it.

## Cleanup the sample fabric network
After you are done testing, you can stop and cleanup the Fabric sample `first-network` as follows:
```
cd $GOPATH//src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```

## Deploy to IBM Cloud
To deploy the `iou` chaincode to IBM Cloud, it is required to package the chaincode in `.cds` format.  The script `make package` has already created [`iou_cc.cds`](iou_cc.cds), which you can deploy to IBM Blockchain Platform.  Refer to [fabric-tools](../../fabric-tools) for details about installing chaincode on the IBM Blockchain Platform.

The GraphQL service app can access the same `iou` chaincode deployed in [IBM Cloud](https://cloud.ibm.com) using the [IBM Blockchain Platform](https://cloud.ibm.com/catalog/services/blockchain-platform-20). The only required update is the network configuration file.  [config_ibp.yaml](../../testdata/config_ibp.yaml) is a sample network configuration that can be used by the GraphQL service.
