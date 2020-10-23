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
MODEL_DIR=${WORK}/${NAME}
env

function create {
  if [ -d "/tmp/${NAME}" ]; then
    echo "cleanup old workspace /tmp/${NAME}"
    rm -rf /tmp/${NAME}
  fi
  mkdir /tmp/${NAME}
  cp ${MODEL_DIR}/${MODEL} /tmp/${NAME}
  cd /tmp/${NAME}
  flogo create --cv ${FLOGO_VER} --verbose -f ${MODEL} ${NAME}
  rm ${NAME}/src/main.go
  cp ${SHIM_PATH}/chaincode_shim.go ${NAME}/src/main.go

  cd ${HOME}
  if [ -d "${MODEL_DIR}/META-INF" ]; then
    cp -rf ${MODEL_DIR}/META-INF /tmp/${NAME}/${NAME}/src
  fi

  cp ${HOME}/codegen.sh /tmp/${NAME}/${NAME}
  cd /tmp/${NAME}/${NAME}
  ./codegen.sh
  cd src
  chmod +x gomodedit.sh
  ./gomodedit.sh
}

function build {
  cd /tmp/${NAME}/${NAME}/src
  go mod edit -replace=github.com/project-flogo/core=${FLOGO_REPO}/core@${FLOGO_REPO_VER}
  go mod edit -replace=github.com/project-flogo/flow=${FLOGO_REPO}/flow@${FLOGO_REPO_VER}

  cd ..
  flogo build -e --verbose
  cd src
  go mod vendor
  go build -mod vendor -o ${MODEL_DIR}/${NAME}_linux_amd64
  if [ ! -f "${MODEL_DIR}/${NAME}_linux_amd64" ]; then
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
  fabric-tools package -n ${NAME} -v ${VERSION} -p /opt/gopath/src/github.com/chaincode/${NAME} -o ${MODEL_DIR}/${NAME}_${VERSION}.cds
  chmod +r ${MODEL_DIR}/${NAME}_${VERSION}.cds
  echo "chaincode cds package: ${MODEL_DIR}/${NAME}_${VERSION}.cds"

  if [ -d "${MODEL_DIR}/${NAME}" ]; then
    echo "cleanup old chaincode source ${MODEL_DIR}/${NAME}"
    rm -rf ${MODEL_DIR}/${NAME}
  fi
  cp -Rf /tmp/${NAME}/${NAME}/src ${MODEL_DIR}/${NAME}
}

create
build