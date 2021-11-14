#! /bin/bash
# It would be ideal if we could forever increment this so during development testing we could just constantly do terraform init -upgrade
VERSION="v1.0.0"
HOSTNAME=registry.terraform.io
NAMESPACE=opslevel
NAME=opslevel
BINARY=terraform-provider-${NAME}_${VERSION}
OS_ARCH=darwin_amd64
PROVIDER_PATH="${HOME}/.terraform.d/plugins"

go build -o ${BINARY}
chmod +x ${BINARY}
mkdir -p ${PROVIDER_PATH}/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION:1}/${OS_ARCH}
mv ${BINARY} ${PROVIDER_PATH}/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION:1}/${OS_ARCH}
echo "Built terraform provider to - ${PROVIDER_PATH}/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION:1}/${OS_ARCH}/${BINARY}"
