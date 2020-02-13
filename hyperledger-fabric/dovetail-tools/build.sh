#!/bin/bash
# Copyright © 2018. TIBCO Software Inc.
#
# This file is subject to the license terms contained
# in the license file that is distributed with this file.

# build chaincode cds or client executable, usage:
#   build.sh cds|client -f json-file [ -n name -v version -s goos -a goarch ]
# or, to start or shutdown builder container
#   build.sh start|shutdown
# e.g.,
#   ./build.sh cds -f ../samples/iou/iou.json -n iou_cc -v 1.0
#   ./build.sh client -f ../samples/iou/iou_client.json -n iou -s darwin -a amd64
# Note: build result will be written in ./work folder

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; echo "$(pwd)")"

function buildCDS {
  startBuilder

  local modelFile=${MODEL##*/}
  local modelDir=${MODEL%/*}
  if [ "${modelDir}" == "${modelFile}" ]; then
    echo "set model file directory to PWD"
    modelDir="."
  fi
  local modelName=${NAME}
  if [ -z "${NAME}" ]; then
    modelName="${modelFile%.*}_cc"
    echo "set chaincode name ${modelName}"
  fi
  local targetDir=${SCRIPT_DIR}/work/${modelName}
  if [ ! -f "${targetDir}/${modelFile}" ]; then
    echo "copy model file to workspace"
    mkdir -p ${targetDir}
    cp ${MODEL} ${targetDir}
    if [ -d "${modelDir}/META-INF" ]; then
      echo "copy META-INF from model folder"
      rm -rf ${targetDir}/META-INF
      cp -rf ${modelDir}/META-INF ${targetDir}
    fi
  fi

  echo "execute build commmand ./build-cds.sh ./work/${modelName}/${modelFile} ${modelName} ${VERSION}"
  docker exec -it dovetail-tools bash -c "./build-cds.sh ./work/${modelName}/${modelFile} ${modelName} ${VERSION}"
}

function buildClient {
  startBuilder

  local modelFile=${MODEL##*/}
  local modelName=${NAME}
  if [ -z "${NAME}" ]; then
    modelName=${modelFile%.*}
    echo "set executable name ${modelName}"
  fi
  local targetDir=${SCRIPT_DIR}/work/${modelName}
  if [ ! -f "${targetDir}/${modelFile}" ]; then
    echo "copy model file to workspace"
    mkdir -p ${targetDir}
    cp ${MODEL} ${targetDir}
  fi

  echo "execute build commmand ./build-client.sh ./work/${modelName}/${modelFile} ${modelName} ${TOS} ${TARCH}"
  docker exec -it dovetail-tools bash -c "./build-client.sh ./work/${modelName}/${modelFile} ${modelName} ${TOS} ${TARCH}"
}

function startBuilder {
  docker ps -f name=dovetail-tools | grep dovetail-tools
  if [ $? -ne 0 ]; then
    echo "start dovetail-tools container ..."
    docker-compose -f ${SCRIPT_DIR}/dovetail-tools.yaml up -d 2>&1
  else
    echo "dovetail-tools container already started"
  fi
}

function shutdownBuilder {
  docker ps -a -f name=dovetail-tools | grep dovetail-tools
  if [ $? -eq 0 ]; then
    echo "stop and cleanup dovetail-tools container ..."
    docker stop dovetail-tools
    docker rm dovetail-tools
  else
    echo "dovetail-tools container is not running"
  fi
}

# Print the usage message
function printHelp() {
  echo "Usage: "
  echo "  build.sh <cmd> [ args ]"
  echo "    <cmd> - one of the following"
  echo "      - 'cds' - build chaincode cds with args: -f model-json [ -n cc-name -v cc-version ]"
  echo "      - 'client' - build client executable with args: -f model-json [ -n exe-name -s GOOS -a GOARCH ]"
  echo "      - 'start' - start docker-tools docker container"
  echo "      - 'shutdown' - shutdown docker-tools docker container"
  echo "    -f <model json> - path of the flogo model json file"
  echo "    -n <name> - name of the chaincode or client exe file, default <model>_cc_<version>.cds or <model>_<goos>_<goarch>"
  echo "    -v <cc-version> - version of the chaincode, e.g., '1.0' (default)"
  echo "    -s <GOOS> - GOOS platform for the client exe, e.g., linux (default), darwin, or windows"
  echo "    -a <GOARCH> - GOARCH for the the client exe, e.g., amd64 (default), or 386"
  echo "  build.sh -h (print this message)"
}

TOS="linux"
TARCH="amd64"
VERSION="1.0"

CMD=${1}
shift
while getopts "h?f:n:v:s:a:" opt; do
  case "$opt" in
  h | \?)
    printHelp
    exit 0
    ;;
  f)
    MODEL=$OPTARG
    ;;
  n)
    NAME=$OPTARG
    ;;
  v)
    VERSION=$OPTARG
    ;;
  s)
    TOS=$OPTARG
    ;;
  a)
    TARCH=$OPTARG
    ;;
  esac
done

case "${CMD}" in
cds)
  if [ -z "${MODEL}" ]; then
    echo "Flogo model file must be specified"
    printHelp
    exit 1
  fi
  echo "build cds: ${MODEL}"
  buildCDS
  ;;
client)
  if [ -z "${MODEL}" ]; then
    echo "Flogo model file must be specified"
    printHelp
    exit 1
  fi
  echo "build client: ${MODEL}"
  buildClient
  ;;
start)
  startBuilder
  ;;
shutdown)
  shutdownBuilder
  ;;
*)
  printHelp
  exit 1
esac
