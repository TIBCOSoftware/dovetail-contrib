#!/bin/bash
echo "Building fabric extension"
build_tag=$1
if [ -d .target/hyperledger-fabric ]; then
  rm -rf .target/hyperledger-fabric
fi
mkdir .target/hyperledger-fabric
zip -r .target/hyperledger-fabric/fabric-extension-${build_tag}.zip hyperledger-fabric/fabric

echo "Building fabric client extension"
zip -r .target/hyperledger-fabric/fabric-client-extension-${build_tag}.zip hyperledger-fabric/fabclient
