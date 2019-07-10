#!/bin/bash
# generate user key and cert with attributes for ABAC test

USER=${1:-"User3"}
ORG=org1.example.com
FABRIC_SAMPLE_PATH=${GOPATH}/src/github.com/hyperledger/fabric-samples/first-network
WORK=/tmp/ca
echo "generate key and cert for user ${USER}@${ORG}"

if [ -d "${WORK}" ]; then
  echo "cleanup ${WORK}"
  rm -R "${WORK}"
fi

# check CA server
docker ps | grep "hyperledger/fabric-ca" | grep "7054->7054/tcp"
if [ "$?" -ne 0 ]; then
  echo "CA server is not running.  Use 'start-ca.sh' to start it first."
  exit 1
fi

# check fabric-ca-client
which fabric-ca-client
if [ "$?" -ne 0 ]; then
  echo "fabric-ca-client not found. You can install fabric-ca by using 'go get -u github.com/hyperledger/fabric-ca/cmd/...'"
  exit 1
fi

# enroll CA admin
export FABRIC_CA_CLIENT_HOME=${WORK}/admin
fabric-ca-client getcainfo -u http://localhost:7054
openssl x509 -noout -text -in ${FABRIC_CA_CLIENT_HOME}/msp/cacerts/localhost-7054.pem
fabric-ca-client enroll -u http://admin:adminpw@localhost:7054

# register and enroll new user
fabric-ca-client register --id.name ''"${USER}@${ORG}"'' --id.secret ${USER}pw --id.type client --id.attrs 'abac.init=true:ecert,email='"${USER}@${ORG}"''
export FABRIC_CA_CLIENT_HOME=${WORK}/${USER}\@${ORG}
fabric-ca-client enroll -u http://${USER}@${ORG}:${USER}pw@localhost:7054 --enrollment.attrs "abac.init,email" -M ${FABRIC_CA_CLIENT_HOME}/msp
openssl x509 -noout -text -in ${WORK}/${USER}\@${ORG}/msp/signcerts/cert.pem

# copy key and cert to first-network sample crypto-config
cd ${FABRIC_SAMPLE_PATH}/crypto-config/peerOrganizations/${ORG}/users
cp -R User1\@${ORG} ${USER}\@${ORG}
cd ${USER}\@${ORG}
rm -R msp/keystore
cp -R ${WORK}/${USER}\@${ORG}/msp/keystore msp
rm msp/admincerts/User1\@${ORG}-cert.pem
cp ${WORK}/${USER}\@${ORG}/msp/signcerts/cert.pem msp/admincerts/${USER}\@${ORG}-cert.pem
rm msp/signcerts/User1\@${ORG}-cert.pem
cp ${WORK}/${USER}\@${ORG}/msp/signcerts/cert.pem msp/signcerts/${USER}\@${ORG}-cert.pem
openssl x509 -noout -text -in msp/signcerts/${USER}\@${ORG}-cert.pem
