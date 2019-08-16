#!/bin/bash
build_tag=$1
if [ -d .target/corda ]; then
  rm -rf .target/corda
fi
mkdir -p .target/corda

echo "Building corda extension"
#zip -r .target/corda/corda-extension-${build_tag////-}.zip TBD

echo "Building corda functions"
#zip -r .target/corda/corda-function-${build_tag////-}.zip TBD
