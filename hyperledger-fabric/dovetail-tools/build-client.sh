#!/bin/bash
# Copyright © 2018. TIBCO Software Inc.
#
# This file is subject to the license terms contained
# in the license file that is distributed with this file.

# run this script in dovetail-tools docker container to build executables for fabric client
# usage;
#   build-client.sh <model> <name> <GOOS> <GOARCH>

MODEL=${1}
NAME=${2}
TOS=${3}
TARCH=${4}
echo "build-client.sh ${MODEL} ${NAME} ${TOS} ${TARCH}"
MODEL_DIR=${WORK}/${NAME}
env

function create {
  if [ -d "/tmp/${NAME}" ]; then
    echo "cleanup old workspace /tmp/${NAME}"
    rm -rf /tmp/${NAME}
  fi
  mkdir -p /tmp/${NAME}
  cp ${MODEL_DIR}/${MODEL} /tmp/${NAME}
  cd /tmp/${NAME}
  flogo create --cv ${FLOGO_VER} -f ${MODEL} ${NAME}

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
  # must use this older version of go-kit
  go mod edit -replace=github.com/go-kit/kit=github.com/go-kit/kit@v0.8.0
  cd ..
  flogo build -e --verbose
  cd src
  go mod vendor
  GOOS=${TOS} GOARCH=${TARCH} go build -mod vendor -o ${MODEL_DIR}/${NAME}_${TOS}_${TARCH}
  echo "client executable: ${MODEL_DIR}/${NAME}_${TOS}_${TARCH}"
}

create
build
