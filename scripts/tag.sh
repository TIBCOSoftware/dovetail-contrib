#!/bin/bash
version=$(<version.txt)
if [ ! $TRAVIS_BRANCH == "master" ]
then
    version=${version}_$TRAVIS_BRANCH
fi
echo "${version}"