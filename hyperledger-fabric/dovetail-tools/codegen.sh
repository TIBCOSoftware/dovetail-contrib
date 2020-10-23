#!/bin/bash
# This script runs in parent folder of chaincode src
# It reads flogo.json to generate dovetailimports.go and gomodedit.sh
# and puts both file in src folder
#

cd "$(dirname "${BASH_SOURCE[0]}")"

function printImports {
  local imp=$(cat flogo.json | jq .imports[])
  echo "package main"
  echo "import ("
  for f in $imp; do
    if [[ $f == *github.com/project-flogo* ]]; then
      # ignore open source flogo that are automatically added
      :
    elif [[ $f == */connector/* ]]; then
      # ignore connector that are not packaged as go mod
      :
    elif [[ $f == */TIBCOSoftware/dovetail-contrib/hyperledger-fabric/* ]]; then
      # ignore open source dovetail components that sre automatically added
      :
    else
      echo "   _ $f"
    fi
  done
  echo "   _ \"github.com/project-flogo/core/app/propertyresolver\""
  echo ")"
}

function printGomod {
  echo "#!/bin/bash"
  local imp=$(cat flogo.json | jq -r .imports[])
  local general=""
  for f in $imp; do
    if [[ $f == *product/ipaas/wi-contrib.git* ]]; then
      if [ -z "${FE_HOME}" ]; then
        echo "Missing FE_HOME, cannot build component ${f}"
        exit 1
      fi
      if [ ! -d "${FE_SRC}" ]; then
        echo "Flogo Enterprise extension does not exist in ${FE_SRC}"
        exit 1
      fi
      echo "go mod edit -require=${f}@v0.0.0"
      if [[ $f == *contributions/General* ]]; then
        general="true"
        echo "go mod edit -replace=${f}=${FE_GENERAL}/${f##*contributions/General/}"
      else
        echo "go mod edit -replace=${f}=${FE_SRC}/${f}"
      fi
    elif [[ $f == */TIBCOSoftware/dovetail-contrib/hyperledger-fabric/* ]]; then
      # use local source for dovetail components
      echo "go mod edit -replace=${f}=${HOME}/${f##*TIBCOSoftware/}"
    fi
  done
  echo "go mod edit -replace=github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common=${HOME}/dovetail-contrib/hyperledger-fabric/fabric/common"
  echo "go mod edit -replace=github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabclient/common=${HOME}/dovetail-contrib/hyperledger-fabric/fabclient/common"
  if [ ! -z "${general}" ]; then
    echo "go mod edit -require=git.tibco.com/git/product/ipaas/wi-contrib.git/engine@v0.0.0"
    echo "go mod edit -replace=git.tibco.com/git/product/ipaas/wi-contrib.git/engine=${FE_SRC}/git.tibco.com/git/product/ipaas/wi-contrib.git/engine"
    echo "go mod edit -require=git.tibco.com/git/product/ipaas/wi-contrib.git/httpservice@v0.0.0"
    echo "go mod edit -replace=git.tibco.com/git/product/ipaas/wi-contrib.git/httpservice=${FE_SRC}/git.tibco.com/git/product/ipaas/wi-contrib.git/httpservice"
    echo "go mod edit -require=git.tibco.com/git/product/ipaas/wi-contrib.git/environment@v0.0.0"
    echo "go mod edit -replace=git.tibco.com/git/product/ipaas/wi-contrib.git/environment=${FE_SRC}/git.tibco.com/git/product/ipaas/wi-contrib.git/environment"
  fi
}

if [ -f "flogo.json" ]; then
  echo "generate dovetail imports"
  printImports > src/dovetailimports.go
fi

# if FE_HOME is defined by the tools image, prepare for inclusion of FE components
if [ -z "${FE_HOME}" ]; then
  echo "skip Flogo Enterprise config: FE_HOME is not specified"
else
  FE_SRC=${FE_HOME}/lib/core/src
  FE_GENERAL=$(dirname "${FE_HOME}")/data/localstack/wicontributions/Tibco/General
  if [ -d "${FE_SRC}" ]; then
    echo "use Flogo Enterprise extension from ${FE_SRC}"
  else
    echo "Flogo Enterprise extension does not exist in ${FE_SRC}"
    exit 1
  fi
fi

if [ -f "flogo.json" ]; then
  printGomod > src/gomodedit.sh
else
  echo "cannot find flogo.json. codegen.sh must run from a Flogo project root"
  exit 1
fi
