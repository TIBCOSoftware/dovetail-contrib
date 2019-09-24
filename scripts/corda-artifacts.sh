#!/bin/bash
build_tag=$1
if [ -d .target/corda ]; then
  rm -rf .target/corda
fi
mkdir -p .target/corda

echo "Building corda contract extension"
zip -r .target/corda/corda-contract-extension-${build_tag////-}.zip corda/contract

echo "Building corda cordapp extension"
zip -r .target/corda/corda-cordapp-extension-${build_tag////-}.zip corda/cordapp

echo "Building corda general extension"
zip -r .target/corda/corda-general-extension-${build_tag////-}.zip corda/general

cd corda

echo "Building corda functions"
zip -r ../.target/corda/corda-function-${build_tag////-}.zip function
