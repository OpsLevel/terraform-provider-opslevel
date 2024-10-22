variables {
  resource_name = "opslevel_integration_aws"

  # required fields
  external_id = "194c7dfc-3a3f-4b0a-b898-578ce7e8f6dc"
  iam_role    = "arn:aws:iam::994866125780:user/opslevel-test"
  name        = "TF Test AWS Integration"

  # optional fields
  ownership_tag_overrides = false
  ownership_tag_keys      = ["one", "two", "three", "four", "five"]
  region_override         = ["eu-west-1", "us-east-1"]

  # default values - computed from API
  default_ownership_tag_keys      = tolist(["owner"])
  default_ownership_tag_overrides = true
}

run "resource_integration_aws_create_with_all_fields" {

  module {
    source = "./opslevel_modules/modules/integration/aws"
  }

  assert {
    condition = alltrue([
      can(opslevel_integration_aws.this.external_id),
      can(opslevel_integration_aws.this.iam_role),
      can(opslevel_integration_aws.this.id),
      can(opslevel_integration_aws.this.ownership_tag_keys),
      can(opslevel_integration_aws.this.ownership_tag_overrides),
      can(opslevel_integration_aws.this.name),
      can(opslevel_integration_aws.this.region_override),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_integration_aws.this.external_id == var.external_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.external_id,
      opslevel_integration_aws.this.external_id,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.iam_role == var.iam_role
    error_message = format(
      "expected '%v' but got '%v'",
      var.iam_role,
      opslevel_integration_aws.this.iam_role,
    )
  }

  assert {
    condition     = startswith(opslevel_integration_aws.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_integration_aws.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_integration_aws.this.name,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.ownership_tag_keys == var.ownership_tag_keys
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_keys,
      opslevel_integration_aws.this.ownership_tag_keys,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.ownership_tag_overrides == var.ownership_tag_overrides
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_overrides,
      opslevel_integration_aws.this.ownership_tag_overrides,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.region_override == var.region_override
    error_message = format(
      "expected '%v' but got '%v'",
      var.region_override,
      opslevel_integration_aws.this.region_override,
    )
  }

}

run "resource_integration_aws_unset_optional_fields" {

  variables {
    ownership_tag_keys      = null
    ownership_tag_overrides = null
    region_override         = null
  }

  module {
    source = "./opslevel_modules/modules/integration/aws"
  }

  assert {
    condition = opslevel_integration_aws.this.ownership_tag_keys == var.default_ownership_tag_keys
    error_message = format(
      "expected default '%v' but got '%v'",
      var.default_ownership_tag_keys,
      opslevel_integration_aws.this.ownership_tag_keys,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.ownership_tag_overrides == var.default_ownership_tag_overrides
    error_message = format(
      "expected default '%v' but got '%v'",
      var.default_ownership_tag_overrides,
      opslevel_integration_aws.this.ownership_tag_overrides,
    )
  }

  assert {
    condition     = opslevel_integration_aws.this.region_override == null
    error_message = var.error_expected_null_field
  }

}

run "delete_aws_integration_outside_of_terraform" {

  variables {
    command = "delete integration ${run.resource_integration_aws_create_with_all_fields.this.id}"
  }

  module {
    source = "./cli"
  }
}

run "resource_integration_aws_create_with_required_fields" {

  variables {
    ownership_tag_overrides = null
    ownership_tag_keys      = null
    region_override         = null
  }

  module {
    source = "./opslevel_modules/modules/integration/aws"
  }

  assert {
    condition = run.resource_integration_aws_create_with_all_fields.this.id != opslevel_integration_aws.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_integration_aws_create_with_all_fields.this.id,
      opslevel_integration_aws.this.id,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.external_id == var.external_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.external_id,
      opslevel_integration_aws.this.external_id,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.iam_role == var.iam_role
    error_message = format(
      "expected '%v' but got '%v'",
      var.iam_role,
      opslevel_integration_aws.this.iam_role,
    )
  }

  assert {
    condition     = startswith(opslevel_integration_aws.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_integration_aws.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_integration_aws.this.name,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.ownership_tag_keys == var.default_ownership_tag_keys
    error_message = format(
      "expected '%v' but got '%v'",
      var.default_ownership_tag_keys,
      opslevel_integration_aws.this.ownership_tag_keys,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.ownership_tag_overrides == var.default_ownership_tag_overrides
    error_message = format(
      "expected '%v' but got '%v'",
      var.default_ownership_tag_overrides,
      opslevel_integration_aws.this.ownership_tag_overrides,
    )
  }

  assert {
    condition     = opslevel_integration_aws.this.region_override == null
    error_message = var.error_expected_null_field
  }

}

run "resource_integration_aws_set_all_fields" {

  module {
    source = "./opslevel_modules/modules/integration/aws"
  }

  assert {
    condition = opslevel_integration_aws.this.external_id == var.external_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.external_id,
      opslevel_integration_aws.this.external_id,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.iam_role == var.iam_role
    error_message = format(
      "expected '%v' but got '%v'",
      var.iam_role,
      opslevel_integration_aws.this.iam_role,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_integration_aws.this.name,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.ownership_tag_keys == var.ownership_tag_keys
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_keys,
      opslevel_integration_aws.this.ownership_tag_keys,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.ownership_tag_overrides == var.ownership_tag_overrides
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_overrides,
      opslevel_integration_aws.this.ownership_tag_overrides,
    )
  }

  assert {
    condition = opslevel_integration_aws.this.region_override == var.region_override
    error_message = format(
      "expected '%v' but got '%v'",
      var.region_override,
      opslevel_integration_aws.this.region_override,
    )
  }

}
