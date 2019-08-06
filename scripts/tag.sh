#!/bin/bash
version=$(<version.txt)
branch=$1
num=$2
if [ ! $branch == "master" ]
then
    version=${version}-${branch}.${num}
fi
echo "${version}"