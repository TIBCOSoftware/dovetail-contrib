#!/bin/bash
# start CA server for Fabric first-network sample

FABRIC_SAMPLE_PATH=${GOPATH}/src/github.com/hyperledger/fabric-samples/first-network

CA_KEY_PATH=$(ls ${FABRIC_SAMPLE_PATH}/crypto-config/peerOrganizations/org1.example.com/ca/*_sk)
if [ -z "${CA_KEY_PATH}" ]; then
  echo "Fabric sample network CA key does not exist.  Start the first-network sample before this script."
  exit 1
fi

CA_KEY="${CA_KEY_PATH##*/}"
sed "s/{{CA_PRIVATE_KEY}}/${CA_KEY}/g" docker-compose-ca.yaml > ${FABRIC_SAMPLE_PATH}/docker-compose-ca.yaml

cd ${FABRIC_SAMPLE_PATH}
docker-compose -f docker-compose-ca.yaml up -d
