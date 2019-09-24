#!/bin/bash
# generate user key and cert with attribute bank, admin, alias, and email
# usage: gen-user.sh name org
#  e.g., gen-user/sh Alice 1
#  user-name ending with "Admin" acts as bank admin 

USER=${1:-"Alice"}
ORG=org1.example.com
PORT=7054
BANK=EURBank
if [ "$2" == "2" ]; then
  ORG=org2.example.com
  PORT=8054
  BANK=USDBank
fi

ADMIN="false"
if [[ "${USER}" == *Admin ]]; then
  echo "${USER} is bank admin"
  ADMIN="true"
fi

FABRIC_SAMPLE_PATH=${GOPATH}/src/github.com/hyperledger/fabric-samples/first-network
WORK=/tmp/ca
echo "generate key and cert for user ${USER}@${ORG}"

if [ -d "${WORK}" ]; then
  echo "cleanup ${WORK}"
  rm -R "${WORK}"
fi

# check CA server
docker ps | grep "hyperledger/fabric-ca" | grep "${PORT}->${PORT}/tcp"
if [ "$?" -ne 0 ]; then
  echo "CA server is not running.  Start first network with '-a' option, e.g., './byfn.sh up -a -s couchdb'."
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
fabric-ca-client getcainfo -u http://localhost:${PORT}
openssl x509 -noout -text -in ${FABRIC_CA_CLIENT_HOME}/msp/cacerts/localhost-${PORT}.pem
fabric-ca-client enroll -u http://admin:adminpw@localhost:${PORT}

# register and enroll new user
# Note: important to make id.name as user@org for signature verification!
fabric-ca-client register --id.name ''"${USER}@${ORG}"'' --id.secret ${USER}pw --id.type client --id.attrs 'admin='"${ADMIN}"':ecert,alias='"${USER}"',bank='"${BANK}"',email='"${USER}@${ORG}"''
export FABRIC_CA_CLIENT_HOME=${WORK}/${USER}\@${ORG}
fabric-ca-client enroll -u http://${USER}@${ORG}:${USER}pw@localhost:${PORT} --enrollment.attrs "admin,alias,bank,email" -M ${FABRIC_CA_CLIENT_HOME}/msp
openssl x509 -noout -text -in ${WORK}/${USER}\@${ORG}/msp/signcerts/cert.pem

# copy key and cert to first-network sample crypto-config
cd ${FABRIC_SAMPLE_PATH}/crypto-config/peerOrganizations/${ORG}/users
if [ -d "${USER}@${ORG}" ]; then
  echo "remove old crypto ${USER}@${ORG}"
  rm -Rf ${USER}\@${ORG}
fi
cp -R User1\@${ORG} ${USER}\@${ORG}
cd ${USER}\@${ORG}
rm -R msp/keystore
cp -R ${WORK}/${USER}\@${ORG}/msp/keystore msp
rm msp/admincerts/User1\@${ORG}-cert.pem
cp ${WORK}/${USER}\@${ORG}/msp/signcerts/cert.pem msp/admincerts/${USER}\@${ORG}-cert.pem
rm msp/signcerts/User1\@${ORG}-cert.pem
cp ${WORK}/${USER}\@${ORG}/msp/signcerts/cert.pem msp/signcerts/${USER}\@${ORG}-cert.pem
openssl x509 -noout -text -in msp/signcerts/${USER}\@${ORG}-cert.pem
