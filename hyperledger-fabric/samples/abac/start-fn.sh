#!/bin/bash
# restart fabric sample first-network, including CA servers without TLS

FABRIC_SAMPLE_PATH=${GOPATH}/src/github.com/hyperledger/fabric-samples/first-network
cd ${FABRIC_SAMPLE_PATH}

# cleanup started network if any
./byfn.sh down
docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
docker rmi $(docker images | grep dev-peer | awk '{print $3}')

# turn off TLS on CA servers
cp docker-compose-ca.yaml docker-compose-ca.yaml.orig
sed "s/FABRIC_CA_SERVER_TLS_ENABLED=true/FABRIC_CA_SERVER_TLS_ENABLED=false/g" docker-compose-ca.yaml.orig > docker-compose-ca.yaml

# start first-network
./byfn.sh up -a -s couchdb
