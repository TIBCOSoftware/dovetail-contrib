#!/bin/bash
build_tag=$1
if [ -d .target/multitarget ]; then
  rm -rf .target/multitarget
fi
mkdir -p .target/multitarget

echo "Building multitarget extension"
zip -r .target/multitarget/multitarget-extension-${build_tag////-}.zip multitarget/multitarget

echo "Building multitarget general extension"
zip -r .target/multitarget/multitarget-general-extension-${build_tag////-}.zip multitarget/general

cd multitarget

echo "Building multitarget functions"
zip -r ../.target/multitarget/multitarget-function-${build_tag////-}.zip function
