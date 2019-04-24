# fabric-contract
This is a sample smart contract for [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) implemented by a [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-4-0) model.

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.4](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html)
- [Install Go](https://golang.org/doc/install)
- Clone [Hyperledger Fabric Samples](https://github.com/hyperledger/fabric-samples)
- Clone [This Flogo extension](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric)

## Edit smart contract
- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.4.0/doc/pdf/TIB_flogo_2.4_users_guide.pdf?id=1)
- Upload [`fabticExtension.zip`](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/blob/master/fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabric.sh`](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/blob/master/zip-fabric.sh)
- Create new Flogo App of name `fabric_contract` and choose `Import app` to import the model [`fabric_contract.json`](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/blob/master/fabric-contract/fabric_contract.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.

## Build and deploy chaincode to Hyperledger Fabric
- Export the Flogo App, and copy the downloaded model file, i.e., [`fabric_contract.json`](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/blob/master/fabric-contract/fabric_contract.json) to folder `fabric-contract`.  You can skip this step if you did not modify the app in Flogo® Enterprise.
- In the `fabric-contract` folder, execute `make create` to generate source code for the chaincode.  This step downloads all dependent packages, and thus may take a while depending on the network speed.
- Execute `make deploy` to deploy the chaincode to the `fabric-samples` chaincode folder.  Note: you may need to edit the [`Makefile`](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/blob/master/fabric-contract/Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.

## Test smart contract
Start Hyperledger Fabric test network in dev mode:
```
cd $GOPATH//src/github.com/hyperledger/fabric-samples/chaincode-docker-devmode
docker-compose -f docker-compose-simple.yaml up
```
In another terminal, start the chaincode:
```
docker exec -it chaincode bash
cd flogo_cc
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=flogo_cc:0 CORE_CHAINCODE_LOGGING_LEVEL=DEBUG ./flogo_cc
```
In a third terminal, install chaincode and send tests:
```
docker exec -it cli bash
peer chaincode install -p chaincodedev/chaincode/flogo_cc -n flogo_cc -v 0
peer chaincode instantiate -n flogo_cc -v 0 -c '{"Args":["init"]}' -C myc

# test transient attributes, which must be encoded as base64
export SECRET=$(echo -n "\"MyTransientSecret\"" | base64)
export PIN=$(echo -n "1054" | base64)
peer chaincode invoke -n flogo_cc -c '{"Args":["put_record","user_txn_1","hello_1","SHA256","hash_1"]}' -C myc --transient "{\"secret\": \"$SECRET\", \"pin\": \"$PIN\"}"
peer chaincode invoke -n flogo_cc -c '{"Args":["put_records","[{\"user_txn_id\":\"trans_1\",\"data\":\"hello_1\"}]"]}' -C myc --transient "{\"secret\": \"$SECRET\"}"
```
