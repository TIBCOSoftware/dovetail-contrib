#!/bin/bash
# restart fabric sample first-network, including CA servers without TLS

FAB_PATH=${1:-"${GOPATH}/src/github.com/hyperledger/fabric-samples"}
cd ${FAB_PATH}/first-network

cons=$(docker ps -a -f name=example.com | wc -l)
if [ $cons -gt 1 ]; then
  echo "cleanup previous running containers ..."
  ./byfn.sh down
  docker rm $(docker ps -a | grep dev-peer | awk '{print $1}')
  docker rmi $(docker images | grep dev-peer | awk '{print $3}')
fi

# turn off TLS on CA servers
cp docker-compose-ca.yaml docker-compose-ca.yaml.orig
sed "s/FABRIC_CA_SERVER_TLS_ENABLED=true/FABRIC_CA_SERVER_TLS_ENABLED=false/g" docker-compose-ca.yaml.orig > docker-compose-ca.yaml

# start first-network
./byfn.sh up -a -n -s couchdb -i 1.4.9
