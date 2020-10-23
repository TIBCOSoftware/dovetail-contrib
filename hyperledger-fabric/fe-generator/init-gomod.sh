#!/bin/bash
# This script initializes Flogo Enterprise trigger/activity compoonents so apps using these components can build with go mod
#
# caller should pass Flogo Enterprise home folder name, and optionally a component name, e.g.,
# ./init-gomod.sh "/usr/local/tibco/flogo/2.10" "trigger/rest"

cd "$(dirname "${BASH_SOURCE[0]}")"
feroot=${1:-"${FE_HOME}"}
FE_SRC=${feroot}/lib/core/src
FE_GENERAL=$(dirname "${feroot}")/data/localstack/wicontributions/Tibco/General

function createEngineMod {
  local wicontribDir=${FE_SRC}/git.tibco.com/git/product/ipaas/wi-contrib.git
  for f in "environment" "httpservice" "engine"; do
    cd ${wicontribDir}/${f}
    if [ -f "go.mod" ]; then
      echo "cleanup ${wicontribDir}/${f}/go.mod"
      rm go.mod
      rm go.sum
    fi

    echo "initilize ${wicontribDir}/${f}/go.mod"
    go mod init git.tibco.com/git/product/ipaas/wi-contrib.git/${f}
    if [ "${f}" == "engine" ]; then
      go mod edit -require=git.tibco.com/git/product/ipaas/wi-contrib.git/httpservice@v0.0.0
      go mod edit -replace=git.tibco.com/git/product/ipaas/wi-contrib.git/httpservice=${FE_SRC}/git.tibco.com/git/product/ipaas/wi-contrib.git/httpservice
    fi
    echo "update go module ${wicontribDir}/${f}"
    go mod tidy
  done
}

# createGeneralMod <component>, e.g.
# createGeneralMod "trigger/rest"
function createGeneralMod {
  local compDir=${FE_GENERAL}/${1}
  if [ -f "${compDir}/go.mod" ]; then
    echo "cleanup ${compDir}/go.mod"
    rm ${compDir}/go.mod
    rm ${compDir}/go.sum
  fi
  echo "initilizing ${compDir}/go.mod"
  cd ${compDir}
  go mod init git.tibco.com/git/product/ipaas/wi-contrib.git/contributions/General/${1}
  go mod edit -require=git.tibco.com/git/product/ipaas/wi-contrib.git/engine@v0.0.0
  go mod edit -replace=git.tibco.com/git/product/ipaas/wi-contrib.git/engine=${FE_SRC}/git.tibco.com/git/product/ipaas/wi-contrib.git/engine
  go mod edit -require=git.tibco.com/git/product/ipaas/wi-contrib.git/httpservice@v0.0.0
  go mod edit -replace=git.tibco.com/git/product/ipaas/wi-contrib.git/httpservice=${FE_SRC}/git.tibco.com/git/product/ipaas/wi-contrib.git/httpservice
  go mod edit -require=git.tibco.com/git/product/ipaas/wi-contrib.git/environment@v0.0.0
  go mod edit -replace=git.tibco.com/git/product/ipaas/wi-contrib.git/environment=${FE_SRC}/git.tibco.com/git/product/ipaas/wi-contrib.git/environment

  echo "update go module ${compDir}"
  go mod tidy
}

if [ ! -d "${FE_GENERAL}" ]; then
  echo "${feroot} is not a valid FE home directory"
  exit 1
fi

if [ ! -f "fe-generator" ]; then
  echo "build fe-generator"
  go build
fi
echo "generate legacy metadata for FE general components"
./fe-generator -dir ${FE_GENERAL}

createEngineMod

if [ ! -z "${2}" ]; then
  if [ -d "${FE_GENERAL}/${2}" ]; then
    createGeneralMod ${2}
    exit 0
  else
    echo "${2} is not a valid FE component"
    exit 1
  fi
fi

echo "generate go mod for all FE components"
cd ${FE_GENERAL}
for f in trigger/*; do 
  createGeneralMod ${f}
done

cd ${FE_GENERAL}
for f in activity/*; do 
  createGeneralMod ${f}
done
