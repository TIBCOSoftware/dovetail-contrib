#!/bin/bash
SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
BUILD_PATH=/tmp/marble
GOPATH=~/go
APP=marble_cc
APP_CONFIG=marble_app.json
CC_DEPLOY=${GOPATH}/src/github.com/hyperledger/fabric-samples/chaincode

rm -Rf ${BUILD_PATH}
mkdir -p ${BUILD_PATH}
cp ${APP_CONFIG} ${BUILD_PATH}
cd ${BUILD_PATH}
flogo create -f ${APP_CONFIG} ${APP}
rm ${BUILD_PATH}/${APP}/src/main.go
cp ${GOPATH}/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/shim/chaincode_shim.go ${BUILD_PATH}/${APP}/src/main.go
cp ${GOPATH}/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/shim/dovetailimports.go ${BUILD_PATH}/${APP}/src

cd ${BUILD_PATH}/${APP}/src
# go get -u github.com/TIBCOSoftware/dovetail-contrib@issue-36/fabric-extension
echo "" >> go.mod
echo "replace github.com/TIBCOSoftware/dovetail-contrib => ${GOPATH}/src/github.com/TIBCOSoftware/dovetail-contrib" >> go.mod

cd ${BUILD_PATH}/${APP}
flogo install github.com/project-flogo/contrib/function/string
msg=$(flogo build -e 2>&1 | egrep "undefined: NewFactory|undefined: NewActivity" | awk '{print $1}')
arr=( $msg )
for item in "${arr[@]}"; do
  errfile=${item%%:*}
  if [ -f "${errfile}" ]; then
    echo "delete file ${errfile}"
    sudo rm "${errfile}"
  fi
done

echo "build chaincode ${APP}"
cd ${BUILD_PATH}/${APP}/src
env GOOS=linux GOARCH=amd64 go build -o ${APP}
go mod vendor

cd ${BUILD_PATH}/${APP}
rm -Rf ${CC_DEPLOY}/${APP}
cp -Rf src ${CC_DEPLOY}/${APP}
