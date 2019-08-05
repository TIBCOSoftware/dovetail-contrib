#!/bin/bash

cacert=${1-"crypto-ibp/certificateAuthorities/173.193.100.113:32727/tlsca.pem"}
user=${2:-"user2"}
userpw=${3:-"${user}pw"}

export FABRIC_CA_CLIENT_HOME=${PWD}
capath=${cacert%/*}
caurl="https://${user}:${userpw}@${capath##*/}"

userroot=$(find crypto-ibp -name users)
usermsp=${userroot}/${user}/msp

echo "ca url: ${caurl}"
echo "user msp path: ${usermsp}"

# enroll user
mkdir -p ${usermsp}
fabric-ca-client enroll -u ${caurl} --tls.certfiles ${cacert} --enrollment.attrs 'hf.EnrollmentID,hf.Type,hf.Affiliation' -M ${usermsp}

# rename user cert as user@org-cert.pem
orgpath=${userroot%/*}
pemname="${user}@${orgpath##*/}-cert.pem"
echo "sign cert name: ${pemname}"
mv ${usermsp}/signcerts/cert.pem ${usermsp}/signcerts/${pemname}