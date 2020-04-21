#!/bin/bash
if [ -f ./functions.zip ]; then
  rm -f ./functions.zip
fi
cd function
zip -r ../functions.zip dovetail

