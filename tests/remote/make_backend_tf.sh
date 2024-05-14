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

provider "opslevel" {
  api_token = var.api_token
}

variable "api_token" {
  type      = string
  sensitive = true
}
EOF

# create symlinks for each test module. Used to match provider
for TEST_DIR in */
do
  cd "$TEST_DIR" || continue
  ln -s ../backend.tf &> /dev/null|| true
  cd ..
done
