#!/bin/bash

VERSION=$1
FILE_PATH=$2
OS=$3
ARCH=$4

mkdir -p ~/terraform/providers/test.local/opsgenie/opsgenie/$VERSION/"$OS"_"$ARCH"
cp $FILE_PATH ~/terraform/providers/test.local/opsgenie/opsgenie/$VERSION/"$OS"_"$ARCH"