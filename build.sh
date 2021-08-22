#! /bin/bash
VERSION="v99.99.99"
HOSTNAME=registry.terraform.io
NAMESPACE=opslevel
NAME=opslevel
BINARY=terraform-provider-${NAME}_${VERSION}
OS_ARCH=darwin_amd64

go build -o ${BINARY}
chmod +x ${BINARY}
mkdir -p ./.terraform/providers/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION:1}/${OS_ARCH}
mv ${BINARY} ./.terraform/providers/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION:1}/${OS_ARCH}
echo "Built terraform provider to - ./.terraform/providers/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION:1}/${OS_ARCH}/${BINARY}"
