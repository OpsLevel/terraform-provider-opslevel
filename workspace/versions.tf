terraform {
  required_providers {
    opslevel = {
      source  = "OpsLevel/opslevel"
      version = "> 0.0.1"
    }
  }
}

provider "opslevel" {}
