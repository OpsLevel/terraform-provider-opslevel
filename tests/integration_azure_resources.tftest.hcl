variables {
  resource_name = "opslevel_integration_azure_resources"

  # required fields
  client_id       = "XXX_CLIENT_ID_XXX"
  client_secret   = "XXX_CLIENT_SECRET_XXX"
  name            = "TF Test azure_resources Integration"
  subscription_id = "01234567-0123-0123-0123-012345678901"
  tenant_id       = "98765432-9876-9876-9876-987654321098"

  # optional fields
  ownership_tag_keys      = ["one", "two", "three", "four", "five"]
  ownership_tag_overrides = false

  # default values - computed from API
  default_ownership_tag_keys      = tolist(["owner"])
  default_ownership_tag_overrides = true
}

run "resource_integration_azure_create_with_all_fields" {

  module {
    source = "./opslevel_modules/modules/integration/azure_resources"
  }

  assert {
    condition = alltrue([
      can(opslevel_integration_azure_resources.this.aliases),
      can(opslevel_integration_azure_resources.this.client_id),
      can(opslevel_integration_azure_resources.this.client_secret),
      can(opslevel_integration_azure_resources.this.created_at),
      can(opslevel_integration_azure_resources.this.id),
      can(opslevel_integration_azure_resources.this.installed_at),
      can(opslevel_integration_azure_resources.this.name),
      can(opslevel_integration_azure_resources.this.ownership_tag_keys),
      can(opslevel_integration_azure_resources.this.ownership_tag_overrides),
      can(opslevel_integration_azure_resources.this.subscription_id),
      can(opslevel_integration_azure_resources.this.tenant_id),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_integration_azure_resources.this.client_id == var.client_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.client_id,
      opslevel_integration_azure_resources.this.client_id,
    )
  }

  assert {
    condition     = opslevel_integration_azure_resources.this.client_secret == var.client_secret
    error_message = "expected different client_secret value, not printing sensitive value"
  }

  assert {
    condition     = startswith(opslevel_integration_azure_resources.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_integration_azure_resources.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_integration_azure_resources.this.name,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.this.ownership_tag_keys == var.ownership_tag_keys
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_keys,
      opslevel_integration_azure_resources.this.ownership_tag_keys,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.this.ownership_tag_overrides == var.ownership_tag_overrides
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_overrides,
      opslevel_integration_azure_resources.this.ownership_tag_overrides,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.this.subscription_id == var.subscription_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.subscription_id,
      opslevel_integration_azure_resources.this.subscription_id,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.this.tenant_id == var.tenant_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.tenant_id,
      opslevel_integration_azure_resources.this.tenant_id,
    )
  }

}

run "resource_integration_azure_unset_optional_fields" {

  variables {
    ownership_tag_keys = null
    ownership_tag_overrides = null
  }

  module {
    source = "./opslevel_modules/modules/integration/azure_resources"
  }

  assert {
    condition = opslevel_integration_azure_resources.this.ownership_tag_keys == var.default_ownership_tag_keys
    error_message = format(
      "expected default '%v' but got '%v'",
      var.default_ownership_tag_keys,
      opslevel_integration_azure_resources.this.ownership_tag_keys,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.this.ownership_tag_overrides == var.default_ownership_tag_overrides
    error_message = format(
      "expected default '%v' but got '%v'",
      var.default_ownership_tag_overrides,
      opslevel_integration_azure_resources.this.ownership_tag_overrides,
    )
  }

}

run "delete_azure_integration_outside_of_terraform" {

  variables {
    command = "delete integration ${run.resource_integration_azure_create_with_all_fields.this.id}"
  }

  module {
    source = "./cli"
  }
}

run "resource_integration_azure_create_with_required_fields" {

  variables {
    ownership_tag_keys = null
    ownership_tag_overrides = null
  }

  module {
    source = "./opslevel_modules/modules/integration/azure_resources"
  }

  assert {
    condition = run.resource_integration_azure_create_with_all_fields.this.id != opslevel_integration_azure_resources.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_integration_azure_create_with_all_fields.this.id,
      opslevel_integration_azure_resources.this.id,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.this.client_id == var.client_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.client_id,
      opslevel_integration_azure_resources.this.client_id,
    )
  }

  assert {
    condition     = opslevel_integration_azure_resources.this.client_secret == var.client_secret
    error_message = "expected different client_secret value, not printing sensitive value"
  }

  assert {
    condition = opslevel_integration_azure_resources.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_integration_azure_resources.this.name,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.this.ownership_tag_keys == var.default_ownership_tag_keys
    error_message = format(
      "expected '%v' but got '%v'",
      var.default_ownership_tag_keys,
      opslevel_integration_azure_resources.this.ownership_tag_keys,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.this.ownership_tag_overrides == var.default_ownership_tag_overrides
    error_message = format(
      "expected '%v' but got '%v'",
      var.default_ownership_tag_overrides,
      opslevel_integration_azure_resources.this.ownership_tag_overrides,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.this.subscription_id == var.subscription_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.subscription_id,
      opslevel_integration_azure_resources.this.subscription_id,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.this.tenant_id == var.tenant_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.tenant_id,
      opslevel_integration_azure_resources.this.tenant_id,
    )
  }

}

run "resource_integration_azure_set_all_fields" {

  variables {
    ownership_tag_keys = ["one", "two", "three", "four", "five"]
    ownership_tag_overrides = false
  }

  module {
    source = "./opslevel_modules/modules/integration/azure_resources"
  }

  assert {
    condition = opslevel_integration_azure_resources.this.ownership_tag_keys == var.ownership_tag_keys
    error_message = format(
      "expected default '%v' but got '%v'",
      var.ownership_tag_keys,
      opslevel_integration_azure_resources.this.ownership_tag_keys,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.this.ownership_tag_overrides == var.ownership_tag_overrides
    error_message = format(
      "expected default '%v' but got '%v'",
      var.ownership_tag_overrides,
      opslevel_integration_azure_resources.this.ownership_tag_overrides,
    )
  }

}
