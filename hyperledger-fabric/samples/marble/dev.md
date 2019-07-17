# Use devmode to test `marble` chaincode

## Build and deploy chaincode to Hyperledger Fabric sample folder
- Export the Flogo App, and copy the downloaded model file, i.e., [`marble_app.json`](marble_app.json) to the `marble` folder.  You can skip this step if you did not modify the flogo model in this folder.
- In the `marble` folder, execute `make create` to generate source code for the chaincode.
- Execute `make deploy` to deploy the chaincode to the `fabric-samples` chaincode folder.  Note: you may need to edit the [`Makefile`](Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make create
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
cd marble_cc
# display Flogo debug logs for debugging
FLOGO_LOG_LEVEL=DEBUG CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=marble_cc:0 CORE_CHAINCODE_LOGGING_LEVEL=DEBUG ./marble_cc
```
Note that the above command assumes that `marble_cc` is built and then mounted in the chaincode container.  To rebuild it in the chaincode container, you must use Go Modules with packages in the vendor folder, which must be done outside the `$GOPATH`, i.e.,
```
docker exec -it chaincode bash
cp -R marble_cc /tmp
cd /tmp/marble_cc
# clean all cached packages only if necessary
# go clean -cache -modcache -i -r
GO111MODULE=on GOCACHE=cache go build -mod vendor -o marble_cc
```
In a third terminal, install chaincode and send test requests:
```
docker exec -it cli bash
peer chaincode install -p chaincodedev/chaincode/marble_cc -n marble_cc -v 0
peer chaincode instantiate -n marble_cc -v 0 -c '{"Args":["init"]}' -C myc

# test transactions using the following commands:
peer chaincode invoke -C myc -n marble_cc -c '{"Args":["initMarble","marble1","blue","35","tom"]}'
peer chaincode invoke -C myc -n marble_cc -c '{"Args":["initMarble","marble2","red","50","tom"]}'
peer chaincode invoke -C myc -n marble_cc -c '{"Args":["initMarble","marble3","blue","70","tom"]}'
peer chaincode invoke -C myc -n marble_cc -c '{"Args":["transferMarble","marble2","jerry"]}'
peer chaincode query -C myc -n marble_cc -c '{"Args":["readMarble","marble2"]}'
peer chaincode query -C myc -n marble_cc -c '{"Args":["getMarblesByRange","marble1","marble3"]}'
peer chaincode invoke -C myc -n marble_cc -c '{"Args":["transferMarblesBasedOnColor","blue","jerry"]}'
peer chaincode query -C myc -n marble_cc -c '{"Args":["getHistoryForMarble","marble1"]}'
peer chaincode invoke -C myc -n marble_cc -c '{"Args":["delete","marble1"]}'
peer chaincode query -C myc -n marble_cc -c '{"Args":["getHistoryForMarble","marble1"]}'
```

`Ctrl+C` and `exit` the docker containers, and then clean up the docker processes,
```
docker rm $(docker ps -a | grep hyperledger | awk '{print $1}')
```