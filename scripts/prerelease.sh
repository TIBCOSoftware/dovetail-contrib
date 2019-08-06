#!/bin/bash
prerelease=false
branch=$1
if [ ! $branch == "master" ]
then
    prerelease=true
fi
export PRERELEASE=${prerelease}