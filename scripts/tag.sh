#!/bin/bash
version=$(<version.txt)
if [ ! $TRAVIS_BRANCH == "master" ]
then
    version=${version}-${TRAVIS_BRANCH}.${TRAVIS_BUILD_NUMBER}
fi
echo "${version}"