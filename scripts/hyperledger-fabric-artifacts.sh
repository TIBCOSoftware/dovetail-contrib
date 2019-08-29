#!/bin/bash
build_tag=$1
if [ -d .target/hyperledger-fabric ]; then
  rm -rf .target/hyperledger-fabric
fi
mkdir -p .target/hyperledger-fabric

echo "Building fabric extension"
zip -r .target/hyperledger-fabric/fabric-extension-${build_tag////-}.zip hyperledger-fabric/fabric

echo "Building fabric client extension"
zip -r .target/hyperledger-fabric/fabric-client-extension-${build_tag////-}.zip hyperledger-fabric/fabclient

cd hyperledger-fabric

echo "Building fabric functions"
zip -r ../.target/hyperledger-fabric/fabric-function-${build_tag////-}.zip function