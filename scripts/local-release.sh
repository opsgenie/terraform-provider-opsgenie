#!/bin/bash

VERSION=$1
FILE_PATH=$2
PROJECT_NAME=$3
FOLDER_NAME=$(echo "$FILE_PATH" | sed -E 's/(.*)\/dist\/terraform-provider-opsgenie_(.*)\/(.*)/\2/g')

mkdir -p ~/terraform/providers/test.local/opsgenie/opsgenie/$VERSION/$FOLDER_NAME
cp $FILE_PATH ~/terraform/providers/test.local/opsgenie/opsgenie/$VERSION/$FOLDER_NAME