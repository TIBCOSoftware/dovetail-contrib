# abac (Attribute Based Access Control)
This is a Flogo app for testing ABAC of the Hyperledger Fabric. it is implemented by using a [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) model.  The model does not require any code, it contains only a JSON model file exported from the TIBCO Flogo® Enterprise.  You can download the prerequisites and then build and deploy the model to a Hyperledger Fabric network as described below.

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.6](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html)
- [Install Go](https://golang.org/doc/install)
- Clone [Hyperledger Fabric](https://github.com/hyperledger/fabric)
- Clone [Hyperledger Fabric Samples](https://github.com/hyperledger/fabric-samples)
- Install [Fabric CA binaries](https://hyperledger-fabric-ca.readthedocs.io/en/release-1.4/users-guide.html)
- Download and install [flogo-cli](https://github.com/TIBCOSoftware/flogo-cli)
- Clone dovetail-contrib with this Flogo extension

There are different ways to clone these packages.  I put them under $GOPATH after installing Go, i.e.,
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
- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabricExtension.zip`](../fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabric.sh`](../zip-fabric.sh)
- Create new Flogo App of name `abac_app` and choose `Import app` to import the model [`abac_app.json`](abac_app.json)
- You can then add or update the flows using the graphical modeler of the TIBCO Flogo® Enterprise.

## Build and deploy chaincode to Hyperledger Fabric
- Export the Flogo App, and copy the downloaded model file, i.e., [`abac_app.json`](abac_app.json) to the folder `abac`.  You can skip this step if you did not modify the app in Flogo® Enterprise.
- In the `abac` folder, execute `make create` to generate source code for the chaincode.
- Execute `make build` and `make deploy` to deploy the chaincode to the `fabric-samples` chaincode folder.  Note: you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/abac
make create
make build
make deploy
```

## Test chaincode using the fabric sample first-network
Start Hyperledger Fabric first-network with CouchDB:
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh up -s couchdb
```
Use the `cli` container to install the `abac_cc` chaincode on both `org1` and `org2`, and then instantiate it.
```
docker exec -it cli bash
. scripts/utils.sh
peer chaincode install -n abac_cc -v 1.0 -p github.com/chaincode/abac_cc
setGlobals 0 2
peer chaincode install -n abac_cc -v 1.0 -p github.com/chaincode/abac_cc
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
peer chaincode instantiate $ORDERER_ARGS -C mychannel -n abac_cc -v 1.0 -c '{"Args":["init"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"
```
Use `cli` container to send a test request, which will return client ID info with an un-authorized error message.
```
ORG1_ARGS="--peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
ORG2_ARGS="--peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
peer chaincode invoke $ORDERER_ARGS -C mychannel -n abac_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["check_abac","abac.init"]}'

# exit CLI when it is done
exit
```

## Setup user for ABAC

To test how ABAC works, we need to use Fabric CA server to create user key and certificate containing attributes for role verifications.

First, start the Fabric CA server docker container for `org1`:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/abac
./start-ca.sh
```

Then, generate a new user `User3` of the `org1` that contains an attribute `abac.init = true`, which is used by the `abac_app` for authorization.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/abac
./gen-user.sh User3
```

## Build and start the Flogo app to test ABAC
Build and start the fabric client app using the pre-created model file [`abac_lient.json`](abac_client.json):
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/abac
make create-client
make build-client
make run
```

## Test ABAC using REST API
The `abac_client` implements a simple REST API that receives the name of a test user, and use the user to invoke the ABAC blockchain transaction.  The following request specifies the `User3` that we have created in the previous step:
```
curl -X GET http://localhost:8989/abac/User3
```

## Cleanup
Stop and cleanup the Fabric `first-network`.
```
exit
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```
