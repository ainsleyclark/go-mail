#!/bin/bash
#
# tag.sh
#

# Set variables
version=$1
message=$2

# Check version is not empty
if [[ $version == "" ]]
  then
    echo "Add Version number"
    exit
fi

# Check commit message is not empty
if [[ $message == "" ]]
  then
    echo "Add commit message"
    exit
fi

echo "Releasing version: " $version

git tag -a "$version" -m "$message"
git push origin $version
