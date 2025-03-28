#!/bin/bash
# setup vendor script

set -e

if [ -f "$(dirname "$0")/ETH.cfg" ]; then
 MODULE_NAME="bls-eth-go-binary"
 echo "github.com/herumi/${MODULE_NAME}"
else
 MODULE_NAME="bls-go-binary"
 echo "ETH.cfg not found, using module: github.com/herumi/${MODULE_NAME}"
fi

MODULE="github.com/herumi/${MODULE_NAME}"
echo "module : ${MODULE_NAME}"

if [ $# -ge 1 ]; then
 VERSION="$1"
 echo "Using specified version: ${VERSION}"
else
 VERSION=$(go list -m -f '{{.Version}}' ${MODULE})
fi

if [ -z "$VERSION" ]; then
  echo "Error: Could not determine module version. Please add the module to go.mod or specify version as argument."
  echo "Usage: $0 [version]"
  echo "Example: $0 v1.36.4"
  exit 1
fi

GOPATH=$(go env GOPATH)
MODULE_PATH="${GOPATH}/pkg/mod/${MODULE}@${VERSION}"

go mod vendor

mkdir -p vendor/${MODULE}/bls/include
mkdir -p vendor/${MODULE}/bls/lib/

cp -r "${MODULE_PATH}/bls/include/" "vendor/${MODULE}/bls/"
cp -r "${MODULE_PATH}/bls/lib/" "vendor/${MODULE}/bls/"

