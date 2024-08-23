variables {
  resource_name = "opslevel_integration_azure_resources"

  # required fields
  client_id       = "XXX_CLIENT_ID_XXX"
  client_secret   = "XXX_CLIENT_SECRET_XXX"
  name            = "TF Test azure_resources Integration"
  subscription_id = "01234567-0123-0123-0123-012345678901"
  tenant_id       = "98765432-9876-9876-9876-987654321098"

  # optional fields
  ownership_tag_keys      = toset(["one", "two", "three", "four", "five"])
  ownership_tag_overrides = true
}

run "resource_integration_azure_resources_create_with_all_fields" {

  variables {
    client_id               = var.client_id
    client_secret           = var.client_secret
    ownership_tag_keys      = var.ownership_tag_keys
    ownership_tag_overrides = var.ownership_tag_overrides
    name                    = var.name
    subscription_id         = var.subscription_id
    tenant_id               = var.tenant_id
  }

  module {
    source = "./integration_azure_resources"
  }

  assert {
    condition = alltrue([
      can(opslevel_integration_azure_resources.test.aliases),
      can(opslevel_integration_azure_resources.test.client_id),
      can(opslevel_integration_azure_resources.test.client_secret),
      can(opslevel_integration_azure_resources.test.created_at),
      can(opslevel_integration_azure_resources.test.id),
      can(opslevel_integration_azure_resources.test.installed_at),
      can(opslevel_integration_azure_resources.test.name),
      can(opslevel_integration_azure_resources.test.ownership_tag_keys),
      can(opslevel_integration_azure_resources.test.ownership_tag_overrides),
      can(opslevel_integration_azure_resources.test.subscription_id),
      can(opslevel_integration_azure_resources.test.tenant_id),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition     = opslevel_integration_azure_resources.test.client_secret == var.client_secret
    error_message = "expected different client_secret, not printing sensitive value"
  }

  assert {
    condition = opslevel_integration_azure_resources.test.client_id == var.client_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.client_id,
      opslevel_integration_azure_resources.test.client_id,
    )
  }

  assert {
    condition     = startswith(opslevel_integration_azure_resources.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_integration_azure_resources.test.ownership_tag_keys == toset(var.ownership_tag_keys)
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_keys,
      opslevel_integration_azure_resources.test.ownership_tag_keys,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.test.ownership_tag_overrides == var.ownership_tag_overrides
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_overrides,
      opslevel_integration_azure_resources.test.ownership_tag_overrides,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.test.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_integration_azure_resources.test.name,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.test.subscription_id == var.subscription_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.subscription_id,
      opslevel_integration_azure_resources.test.subscription_id,
    )
  }

  assert {
    condition = opslevel_integration_azure_resources.test.tenant_id == var.tenant_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.tenant_id,
      opslevel_integration_azure_resources.test.tenant_id,
    )
  }

}

run "resource_integration_azure_resources_unset_optional_fields" {

  variables {
    ownership_tag_keys = null
  }

  module {
    source = "./integration_azure_resources"
  }

  assert {
    condition     = opslevel_integration_azure_resources.test.ownership_tag_keys == null
    error_message = var.error_expected_null_field
  }

}
