#!/bin/bash
ibpjson=${1:-"ibpConnection.json"}
appuser=${2:-"user2"}
appuserpw=${3:-"${appuser}pw"}

cacert=$(fabric-tools ibp config -i ${ibpjson})
echo "ca cert: ${cacert}"

./enroll.sh "${cacert}" "${appuser}" "${appuserpw}"

echo "To run a fabric client app, do the following:"
echo "  1. use config-ibp.yaml as the network config yaml"
echo "  2. copy folder crypto-ibp to client app host"
echo "  3. when running client app, set env \$CRYPTO_PATH to the location of crypto-ibp"