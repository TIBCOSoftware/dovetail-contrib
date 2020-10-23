#!/bin/bash
# Copyright © 2018. TIBCO Software Inc.
#
# This file is subject to the license terms contained
# in the license file that is distributed with this file.

# build docker image for dovetail-tools, usage:
#   docker-image.sh build [ -n name -v version -e flogo-zip ]
# or, 
#   docker-image.sh upload
# e.g.,
#   ./docker-image.sh build -e ~/work/dovetail/felib/flogo.zip
#   docker login
#   ./docker-image.sh upload -u dhuser -d

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; echo "$(pwd)")"

# print yaml for initial docker container
function printDovetailYaml {
  echo "
version: '2'
services:
  dovetail:
    container_name: dovetail
    image: hyperledger/fabric-tools:1.4
    tty: true
    stdin_open: true
    environment:
      - GOPATH=/opt/gopath
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - WORK=/root/work
      - DOVETAIL_REPO=github.com/TIBCOSoftware
      - FLOGO_REPO=github.com/yxuco
      - FLOGO_REPO_VER=v1.1.1
      - FLOGO_VER=v1.1.0
      - FE_HOME=/root/flogo/2.10
      - SHIM_PATH=/root/dovetail-contrib/hyperledger-fabric/shim
    working_dir: /root
    command: /bin/bash
    volumes:
      - /var/run/:/host/var/run/
      - .:/root/work/:cached"
}

function buildImage {
  mkdir -p ${SCRIPT_DIR}/work
  if [ ! -z "${FELIB}" ] && [ -f "${FELIB}" ]; then
    echo "copy Flogo Enterprise lib ${FELIB}"
    cp ${FELIB} ${SCRIPT_DIR}/work/flogo.zip
  fi
  # these scripts runs in container, so copy to shared volume
  cp ${SCRIPT_DIR}/dovetail-init.sh ${SCRIPT_DIR}/work
  cp ${SCRIPT_DIR}/build-cds.sh ${SCRIPT_DIR}/work
  cp ${SCRIPT_DIR}/build-client.sh ${SCRIPT_DIR}/work
  cp ${SCRIPT_DIR}/codegen.sh ${SCRIPT_DIR}/work

  printDovetailYaml > ${SCRIPT_DIR}/work/dovetail-build.yaml
  cd ${SCRIPT_DIR}/work
  echo "start dovetail image builder ..."
  docker-compose -f dovetail-build.yaml up -d 2>&1
  echo "setup dovetail image ..."
  docker exec -it dovetail bash -c "./work/dovetail-init.sh"
  echo "create docker image ..."
  docker commit dovetail ${NAME}:${VERSION}
  echo "stop docker container ..."
  docker stop dovetail
  docker rm dovetail
}

function uploadImage {
  if [ ! -z "${DHPASS}" ]; then
    echo "login to docker hub"
    docker login -u ${DHUSER} -p ${DHPASS}
  fi
  echo "push ${DHUSER}/${NAME}:${VERSION} to Docker Hub ..."
  docker tag ${NAME}:${VERSION} ${DHUSER}/${NAME}:${VERSION}
  docker push ${DHUSER}/${NAME}:${VERSION}
  if [ ! -z "${CLEANUP}" ]; then
    echo "cleanup local docker images"
    docker rmi ${DHUSER}/${NAME}:${VERSION}
    docker rmi ${NAME}:${VERSION}
  fi
}

# Print the usage message
function printHelp() {
  echo "Usage: "
  echo "  docker-image.sh <cmd> [ args ]"
  echo "    <cmd> - one of the following"
  echo "      - 'build' - build dovetail-tools image with optional args: [ -n name -v version -e flogo-zip ]"
  echo "      - 'upload' - upload image to docker hub with optional args: -u user -p passwd [ -n name -v version -d ]"
  echo "    -n <image name> - name of the docker image, e.g., dovetail-tools (default)"
  echo "    -v <image version> - version of the docker image, e.g., 'v1.1.0' (default)"
  echo "    -e <flogo lib> - path of the zip file for Flogo Enterprise library"
  echo "    -u <user> - user name for a docker hub account"
  echo "    -p <passwd> - password for a docker hub account"
  echo "    -d - flag to cleanup local docker images"
  echo "  docker-image.sh -h (print this message)"
}

NAME="dovetail-tools"
VERSION="v1.1.1"

CMD=${1}
shift
while getopts "h?n:v:e:u:p:d" opt; do
  case "$opt" in
  h | \?)
    printHelp
    exit 0
    ;;
  n)
    NAME=$OPTARG
    ;;
  v)
    VERSION=$OPTARG
    ;;
  e)
    FELIB=$OPTARG
    ;;
  u)
    DHUSER=$OPTARG
    ;;
  p)
    DHPASS=$OPTARG
    ;;
  d)
    CLEANUP=true
    ;;
  esac
done

case "${CMD}" in
build)
  echo "build image: ${NAME}:${VERSION} ${FELIB}"
  buildImage
  ;;
upload)
  if [ -z "${DHUSER}" ]; then
    echo "User name for Docker Hub must be specified"
    printHelp
    exit 1
  else
    echo "upload image: ${DHUSER}/${NAME}:${VERSION} ${CLEANUP}"
    uploadImage
  fi
  ;;
*)
  printHelp
  exit 1
esac
