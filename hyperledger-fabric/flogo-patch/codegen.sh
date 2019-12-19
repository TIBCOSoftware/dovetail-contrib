#!/bin/bash
# This script runs in parent folder of chaincode src
# It reads flogo.json to generate dovetailimports.go and gomodedit.sh
# and puts both file in src folder
#
# caller should pass Flogo Enterprise home folder name, e.g.,
# ./codegen.sh /usr/local/tibco/flogo/2.8

cd "$(dirname "${BASH_SOURCE[0]}")"
FE_SRC=${1}/lib/core/src
FE_GENERAL=$(dirname "${1}")/data/localstack/wicontributions/Tibco/General

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
      echo "go mod edit -require=${f}@v0.0.0"
      if [[ $f == *contributions/General* ]]; then
        general="true"
        echo "go mod edit -replace=${f}@v0.0.0=${FE_GENERAL}/${f##*contributions/General/}"
      else
        echo "go mod edit -replace=${f}@v0.0.0=${FE_SRC}/${f}"
      fi
    fi
  done
  if [ ! -z "${general}" ]; then
    echo "go mod edit -require=git.tibco.com/git/product/ipaas/wi-contrib.git/engine@v0.0.0"
    echo "go mod edit -replace=git.tibco.com/git/product/ipaas/wi-contrib.git/engine@v0.0.0=${FE_SRC}/git.tibco.com/git/product/ipaas/wi-contrib.git/engine"
  fi
}

if [ -d "${FE_SRC}" ]; then
  echo "use Flogo Enterprise extension from ${FE_SRC}"
else
  echo "Flogo Enterprise extension does not exist in ${FE_SRC}"
  echo "must pass the name of FE home folder to codegen.sh"
  exit 1
fi

if [ -f "flogo.json" ]; then
  printImports > src/dovetailimports.go
  printGomod > src/gomodedit.sh
else
  echo "cannot find flogo.json. codegen.sh must run from a Flogo project root"
  exit 1
fi
