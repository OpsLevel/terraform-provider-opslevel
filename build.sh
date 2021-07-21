#! /bin/bash
VERSION=$(git describe --tags --match "v[0-9].*" --always)
HOSTNAME=registry.terraform.io
NAMESPACE=opslevel
NAME=opslevel
BINARY=terraform-provider-${NAME}_${VERSION}
OS_ARCH=darwin_amd64

go build -o ${BINARY}
chmod +x ${BINARY}
mkdir -p ${HOME}/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION:1}/${OS_ARCH}
mv ${BINARY} ${HOME}/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION:1}/${OS_ARCH}
echo "Built terraform provider to - ${HOME}/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION:1}/${OS_ARCH}/${BINARY}"
