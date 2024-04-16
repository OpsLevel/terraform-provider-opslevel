mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_infra_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_integration_aws.example.external_id == "mock-external-id"
    error_message = "wrong external_id for opslevel_integration_aws.example"
  }

  assert {
    condition     = opslevel_integration_aws.example.iam_role == "arn:aws:ecr:us-east-1:mock-iam-role"
    error_message = "wrong iam_role for opslevel_integration_aws.example"
  }

  assert {
    condition     = can(opslevel_integration_aws.example.id)
    error_message = "expected opslevel_integration_aws to have an ID"
  }

  assert {
    condition     = opslevel_integration_aws.example.name == "dev"
    error_message = "wrong name for opslevel_integration_aws.example"
  }

  assert {
    condition     = opslevel_integration_aws.example.ownership_tag_overrides == true
    error_message = "expected 'ownership_tag_overrides' to be 'true' for opslevel_integration_aws.example"
  }

  assert {
    condition     = opslevel_integration_aws.example.ownership_tag_keys == tolist(["owner", "team", "group"])
    error_message = "wrong ownership_tag_keys for opslevel_integration_aws.example"
  }

}
