# iou-app
This is a sample chaincode, implemented by using a [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) model.  The model does not require any code, it contains only a JSON model file exported from the TIBCO Flogo® Enterprise.  You can download the prerequisites and then build and deploy the model to a Hyperledger Fabric network as described below.

This sample uses only basic functionalities of chaincodes, refer to [marble-app](../marble-app) and [marble-private](../marble-private) for more features of Hyperledber Fabric, including history, range query, rich query, pagination, and private collections.

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.6](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html)
- [Install Go](https://golang.org/doc/install)
- Clone [Hyperledger Fabric](https://github.com/hyperledger/fabric)
- Clone [Hyperledger Fabric Samples](https://github.com/hyperledger/fabric-samples)
- Download and install [flogo-cli](https://github.com/TIBCOSoftware/flogo-cli)
- Clone dovetail-contrib with this Flogo extension

There are different ways to clone these packages.  I put them under $GOPATH after installing Go, i.e.,
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

## Edit smart contract
- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabricExtension.zip`](../fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabric.sh`](../zip-fabric.sh)
- Create new Flogo App of name `iou` and choose `Import app` to import the model [`iou.json`](iou.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.

## Build and deploy chaincode to Hyperledger Fabric
- Export the Flogo App, and copy the downloaded model file, i.e., [`iou.json`](iou.json) to the folder `iou-app`.  You can skip this step if you did not modify the app in Flogo® Enterprise.
- In the `iou-app` folder, execute `make create` to generate source code for the chaincode.  This step downloads all dependent packages, and thus may take a while depending on the network speed.
- Execute `make build` and `make deploy` to deploy the chaincode to the `fabric-samples` chaincode folder.  Note: you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/iou-app
make create
make build
make deploy
```

## Test chaincode in fabric devmode
Start Hyperledger Fabric test network in dev mode:
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/chaincode-docker-devmode
rm -R chaincode
docker-compose -f docker-compose-simple.yaml up
```
In another terminal, start the chaincode:
```
docker exec -it chaincode bash
cd iou_cc
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=iou_cc:0 CORE_CHAINCODE_LOGGING_LEVEL=DEBUG ./iou_cc
```
In a third terminal, install chaincode and send test requests:
```
docker exec -it cli bash
peer chaincode install -p chaincodedev/chaincode/iou_cc -n iou_cc -v 0
peer chaincode instantiate -n iou_cc -v 0 -c '{"Args":["init"]}' -C myc

# test transactions using the following commands:
peer chaincode invoke -C myc -n iou_cc -c '{"Args":["issueIOU","iou1","State of California","1000","USD", "tom"]}'
peer chaincode invoke -C myc -n iou_cc -c '{"Args":["issueIOU","iou2","US Treasury","2000","USD", "tom"]}'
peer chaincode invoke -C myc -n iou_cc -c '{"Args":["issueIOU","iou3","State of California","3000","USD", "tom"]}'
peer chaincode invoke -C myc -n iou_cc -c '{"Args":["transferIOU","iou3","jerry"]}'
peer chaincode query -C myc -n iou_cc -c '{"Args":["readIOU","iou3"]}'
peer chaincode invoke -C myc -n iou_cc -c '{"Args":["transferIOUByOwner","tom", "don"]}'
peer chaincode query -C myc -n iou_cc -c '{"Args":["getIOUByIssuer","State of California"]}'
```

`Ctrl+C` and `exit` the docker containers, and then clean up the docker processes,
```
docker rm $(docker ps -a | grep hyperledger | awk '{print $1}')
```

## Test chaincode with multi-org fabric network
Start Hyperledger Fabric first-network with CouchDB:
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh up -s couchdb
```
Use the `cli` container to install the `iou_cc` chaincode on both `org1` and `org2`, and then instantiate it.
```
docker exec -it cli bash
. scripts/utils.sh
peer chaincode install -n iou_cc -v 1.0 -p github.com/chaincode/iou_cc
setGlobals 0 2
peer chaincode install -n iou_cc -v 1.0 -p github.com/chaincode/iou_cc
ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
peer chaincode instantiate $ORDERER_ARGS -C mychannel -n iou_cc -v 1.0 -c '{"Args":["init"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"
```
Use `cli` container to send IOU transaction requests:
```
ORG1_ARGS="--peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
ORG2_ARGS="--peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

# insert test data
peer chaincode invoke $ORDERER_ARGS -C mychannel -n iou_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["issueIOU","iou1","State of California","1000","USD", "tom"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n iou_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["issueIOU","iou2","State of California","2000","USD", "tom"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n iou_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["issueIOU","iou3","State of California","3000","USD", "tom"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n iou_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["issueIOU","iou4","US Treasury","4000","USD", "tom"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n iou_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["issueIOU","iou5","US Treasury","5000","USD", "tom"]}'

# transfer IOU ownership
setGlobals 0 1
peer chaincode query -C mychannel -n iou_cc -c '{"Args":["readIOU","iou2"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n iou_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["transferIOU","iou2", "jerry"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n iou_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["transferIOU","iou4", "jerry"]}'
peer chaincode query -C mychannel -n iou_cc -c '{"Args":["readIOU","iou2"]}'
peer chaincode invoke $ORDERER_ARGS -C mychannel -n iou_cc $ORG1_ARGS $ORG2_ARGS -c '{"Args":["transferIOUByOwner","jerry", "don"]}'
peer chaincode query -C mychannel -n iou_cc -c '{"Args":["getIOUByIssuer","US Treasury"]}'
```

Exit the `cli` shell, and then stop and cleanup the Fabric `first-network`.
```
exit
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```
