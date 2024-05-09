#!/usr/bin/env sh

cat << EOF > backend.tf
terraform {
  required_providers {
    opslevel = {
      source  = "OpsLevel/opslevel"
      version = "> ${OPSLEVEL_TERRAFORM_SOURCE_VERSION:-0.0.1}"
    }
  }
}
EOF
