#!/bin/bash
build_tag=$1
if [ -d .target/multichain ]; then
  rm -rf .target/multichain
fi
mkdir -p .target/multichain

echo "Building multichain extension"
zip -r .target/multichain/multichain-extension-${build_tag////-}.zip SmartContract

echo "Building multichain functions"
zip -r .target/multichain/multichain-function-${build_tag////-}.zip function
