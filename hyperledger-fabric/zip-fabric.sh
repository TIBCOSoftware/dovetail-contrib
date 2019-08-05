#!/bin/bash
if [ -f ./fabricExtension.zip ]; then
  rm -f ./fabricExtension.zip
fi
zip -r ./fabricExtension.zip fabric
