#!/bin/bash
# Copyright © 2018. TIBCO Software Inc.
#
# This file is subject to the license terms contained
# in the license file that is distributed with this file.

# run this script in dovetail-tools docker container to build chaincode cds package
# usage;
#   build-cds.sh <model> <cc-name> <cc-version>

MODEL=${1}
NAME=${2}
VERSION=${3}
echo "build-cds.sh ${MODEL} ${NAME} ${VERSION}"
env

function create {
  local modelFile=${MODEL##*/}
  local modelDir=${MODEL%/*}

  if [ -d "/tmp/${NAME}" ]; then
    echo "cleanup old workspace /tmp/${NAME}"
    rm -rf /tmp/${NAME}
  fi
  mkdir /tmp/${NAME}
  cp ${MODEL} /tmp/${NAME}
  cd /tmp/${NAME}
  flogo create --cv ${FLOGO_VER} --verbose -f ${modelFile} ${NAME}
  rm ${NAME}/src/main.go
  cp ${SHIM_PATH}/chaincode_shim.go ${NAME}/src/main.go

  cd ${HOME}
  if [ -d "${modelDir}/META-INF" ]; then
    cp -rf ${modelDir}/META-INF /tmp/${NAME}/${NAME}/src
  fi

  if [ -d "${FE_HOME}" ]; then
    cp ${PATCH_PATH}/codegen.sh /tmp/${NAME}/${NAME}
    cd /tmp/${NAME}/${NAME}
    ./codegen.sh ${FE_HOME}
    cd src
    chmod +x gomodedit.sh
    ./gomodedit.sh
  fi
}

function build {
  cd /tmp/${NAME}/${NAME}/src
  go mod edit -replace=github.com/project-flogo/core@v0.10.1=github.com/project-flogo/core@${FLOGO_VER}
  go mod edit -replace=github.com/project-flogo/flow@v0.10.0=github.com/project-flogo/flow@${FLOGO_VER}
  go mod edit -replace=github.com/project-flogo/flow/activity/subflow@v0.9.0=github.com/project-flogo/flow/activity/subflow@master
  cd ..
  flogo build -e --verbose
  cd src
  go mod vendor
  cp -R ${PATCH_PATH}/* vendor/github.com/project-flogo
  go build -mod vendor -o ../${NAME}_linux_amd64
  if [ ! -f "../${NAME}_linux_amd64" ]; then
    echo "failed to build chaincode"
    exit 1
  fi

  echo "build chaincode cds package ..."
  cd ${HOME}
  if [ -d "/opt/gopath/src/github.com/chaincode/${NAME}" ]; then
    echo "cleanup old chaincode ${NAME}"
    rm -rf /opt/gopath/src/github.com/chaincode/${NAME}
  fi
  mkdir -p /opt/gopath/src/github.com/chaincode
  cp -Rf /tmp/${NAME}/${NAME}/src /opt/gopath/src/github.com/chaincode/${NAME}
  fabric-tools package -n ${NAME} -v ${VERSION} -p /opt/gopath/src/github.com/chaincode/${NAME} -o ${HOME}/work/${NAME}_${VERSION}.cds
  chmod +r ${HOME}/work/${NAME}_${VERSION}.cds
  echo "chaincode cds package: ./work/${NAME}_${VERSION}.cds"
}

create
build