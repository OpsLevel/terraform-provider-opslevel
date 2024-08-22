variables {
  resource_name = "opslevel_integration_aws"

  # required fields
  external_id             = "194c7dfc-3a3f-4b0a-b898-578ce7e8f6dc"
  iam_role                = "arn:aws:iam::994866125780:user/opslevel-test"
  ownership_tag_overrides = false
  name                    = "TF Test AWS Integration"

  # optional fields
  ownership_tag_keys = ["one", "two", "three", "four", "five"]
}

run "resource_integration_aws_create_with_all_fields" {

  variables {
    external_id             = var.external_id
    iam_role                = var.iam_role
    ownership_tag_keys      = var.ownership_tag_keys
    ownership_tag_overrides = var.ownership_tag_overrides
    name                    = var.name
  }

  module {
    source = "./integration_aws"
  }

  assert {
    condition = alltrue([
      can(opslevel_integration_aws.test.external_id),
      can(opslevel_integration_aws.test.iam_role),
      can(opslevel_integration_aws.test.id),
      can(opslevel_integration_aws.test.ownership_tag_keys),
      can(opslevel_integration_aws.test.ownership_tag_overrides),
      can(opslevel_integration_aws.test.name),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_integration_aws.test.external_id == var.external_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.external_id,
      opslevel_integration_aws.test.external_id,
    )
  }

  assert {
    condition = opslevel_integration_aws.test.iam_role == var.iam_role
    error_message = format(
      "expected '%v' but got '%v'",
      var.iam_role,
      opslevel_integration_aws.test.iam_role,
    )
  }

  assert {
    condition     = startswith(opslevel_integration_aws.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_integration_aws.test.ownership_tag_keys == var.ownership_tag_keys
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_keys,
      opslevel_integration_aws.test.ownership_tag_keys,
    )
  }

  assert {
    condition = opslevel_integration_aws.test.ownership_tag_overrides == var.ownership_tag_overrides
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_overrides,
      opslevel_integration_aws.test.ownership_tag_overrides,
    )
  }

  assert {
    condition = opslevel_integration_aws.test.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_integration_aws.test.name,
    )
  }

}

run "resource_integration_aws_ownership_tag_keys_default_value" {

  variables {
    external_id             = var.external_id
    iam_role                = var.iam_role
    ownership_tag_keys      = null
    ownership_tag_overrides = var.ownership_tag_overrides
    name                    = var.name
  }

  module {
    source = "./integration_aws"
  }

  assert {
    condition = opslevel_integration_aws.test.ownership_tag_keys == tolist(["owner"])
    error_message = format(
      "expected '%v' but got '%v'",
      tolist(["owner"]),
      opslevel_integration_aws.test.ownership_tag_keys,
    )
  }

}
