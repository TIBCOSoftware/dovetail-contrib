#!/bin/bash
build_tag=$1
if [ -d .target/libraries ]; then
  rm -rf .target/libraries
fi
mkdir -p .target/libraries

echo "Building corda java"
cd libraries/corda-java
mvn clean
mvn package
cp target/*.jar ../../.target/libraries
