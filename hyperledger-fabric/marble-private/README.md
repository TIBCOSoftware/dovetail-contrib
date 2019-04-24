# marble-private
This is a Flogo version of the [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) sample chaincode, [marbles02_private](https://github.com/hyperledger/fabric-samples/tree/release-1.4/chaincode/marbles02_private) implemented by using a [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-4-0) model.  The model does not require any code, it contains only a JSON model file exported from the TIBCO Flogo® Enterprise.  You can download the prerequisites and then build and deploy the model to a Hyperledger Fabric network as described below.

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.4](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html)
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
- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.4.0/doc/pdf/TIB_flogo_2.4_users_guide.pdf?id=1)
- Upload [`fabricExtension.zip`](../fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabric.sh`](../zip-fabric.sh)
- Create new Flogo App of name `marble_private` and choose `Import app` to import the model [`marble_private.json`](marble_private.json)
- You can then add or update contract transactions using the graphical modeler of the TIBCO Flogo® Enterprise.

## Build and deploy chaincode to Hyperledger Fabric
- Export the Flogo App, and copy the downloaded model file, i.e., [`marble_private.json`](marble_private.json) to folder `marble-private`.  You can skip this step if you did not modify the app in Flogo® Enterprise.
- In the `marble-private` folder, execute `make create` to generate source code for the chaincode.  This step downloads all dependent packages, and thus may take a while depending on the network speed.
- Execute `make build` and `make deploy` to deploy the chaincode to the `fabric-samples` chaincode folder.  Note: you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/marble-private
make create
make deploy
```

## Test chaincode with multi-org fabric network
Start Hyperledger Fabric first-network with CouchDB:
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh up -s couchdb
```
Use the `cli` container to install the `marble_private` chaincode on all 4 peers of `org1` and `org2`, and then instantiate it.
```
docker exec -it cli bash
. scripts/utils.sh
peer chaincode install -n marble_private_cc -v 1.0 -p github.com/chaincode/marble_private_cc
setGlobals 1 1
peer chaincode install -n marble_private_cc -v 1.0 -p github.com/chaincode/marble_private_cc
setGlobals 0 2
peer chaincode install -n marble_private_cc -v 1.0 -p github.com/chaincode/marble_private_cc
setGlobals 1 2
peer chaincode install -n marble_private_cc -v 1.0 -p github.com/chaincode/marble_private_cc

ORDERER_ARGS="-o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
peer chaincode instantiate $ORDERER_ARGS -C mychannel -n marble_private_cc -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')" --collections-config /opt/gopath/src/github.com/chaincode/marbles02_private/collections_config.json
```
Use `cli` container to send marble transaction requests:
```
# test insert and read access permission
setGlobals 0 1
MARBLE=$(echo -n "{\"name\":\"marble1\",\"color\":\"blue\",\"size\":35,\"owner\":\"tom\",\"price\":99}" | base64 | tr -d \\n)
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_private_cc -c '{"Args":["initMarble"]}' --transient "{\"marble\":\"$MARBLE\"}"
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["readMarble","marble1"]}'
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["readMarblePrivateDetails","marble1"]}'
setGlobals 0 2
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["readMarble","marble1"]}'
# following should fail due to no read access permission 
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["readMarblePrivateDetails","marble1"]}'

# test more insert and transfer owner, and purge of marble1 after 3 blocks
setGlobals 0 1
# block +1
MARBLE=$(echo -n "{\"name\":\"marble2\",\"color\":\"red\",\"size\":50,\"owner\":\"tom\",\"price\":199}" | base64 | tr -d \\n)
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_private_cc -c '{"Args":["initMarble"]}' --transient "{\"marble\":\"$MARBLE\"}"
# block +2
MARBLE_OWNER=$(echo -n "{\"name\":\"marble2\",\"owner\":\"jerry\"}" | base64 | tr -d \\n)
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_private_cc -c '{"Args":["transferMarble"]}' --transient "{\"marble_owner\":\"$MARBLE_OWNER\"}"
# block +3
MARBLE=$(echo -n "{\"name\":\"marble3\",\"color\":\"blue\",\"size\":70,\"owner\":\"tom\",\"price\":299}" | base64 | tr -d \\n)
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_private_cc -c '{"Args":["initMarble"]}' --transient "{\"marble\":\"$MARBLE\"}"

# marble1 should still be available
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["readMarblePrivateDetails","marble1"]}'
# block +4
MARBLE_OWNER=$(echo -n "{\"name\":\"marble3\",\"owner\":\"jerry\"}" | base64 | tr -d \\n)
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_private_cc -c '{"Args":["transferMarble"]}' --transient "{\"marble_owner\":\"$MARBLE_OWNER\"}"
# marble1 details purged after 3 blocks, so this returns error
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["readMarblePrivateDetails","marble1"]}'

# test query
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["getMarblesByRange","marble1", "marble3"]}'
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["queryMarblesByOwner","jerry"]}'
MARBLE_DELETE=$(echo -n "{\"name\":\"marble2\"}" | base64 | tr -d \\n)
peer chaincode invoke $ORDERER_ARGS -C mychannel -n marble_private_cc -c '{"Args":["delete"]}' --transient "{\"marble_delete\":\"$MARBLE_DELETE\"}"
# verify deleted marble2
peer chaincode query -C mychannel -n marble_private_cc -c '{"Args":["queryMarblesByOwner","jerry"]}'
```

Exit the `cli` shell, and then stop and cleanup the Fabric `first-network`.
```
exit
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')
```
