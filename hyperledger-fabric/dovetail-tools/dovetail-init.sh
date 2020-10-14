#!/bin/bash
# Copyright © 2018. TIBCO Software Inc.
#
# This file is subject to the license terms contained
# in the license file that is distributed with this file.

# run this script in docker container to setup for dovetail
# assume $WORK is mounted and contain Flogo Enterprise lib named flogo.zip

if [ -f "${WORK}/flogo.zip" ]; then
  unzip ${WORK}/flogo.zip
  rm ${WORK}/flogo.zip
fi
mv ${WORK}/build-cds.sh ${HOME}
mv ${WORK}/build-client.sh ${HOME}
mv ${WORK}/codegen.sh ${HOME}

git clone https://${DOVETAIL_REPO}/dovetail-contrib.git
go get -u github.com/project-flogo/cli/...
cd dovetail-contrib/hyperledger-fabric/fabric-tools
go install
