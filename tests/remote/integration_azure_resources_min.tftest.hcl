variables {
  resource_name = "opslevel_integration_azure_resources"

  # required fields
  client_id       = "XXX_CLIENT_ID_XXX"
  client_secret   = "XXX_CLIENT_SECRET_XXX"
  name            = "TF Test azure_resources Integration"
  subscription_id = "01234567-0123-0123-0123-012345678901"
  tenant_id       = "98765432-9876-9876-9876-987654321098"

  # optional fields
  ownership_tag_keys      = null
  ownership_tag_overrides = null

  # default values - computed from API
  default_ownership_tag_keys      = tolist(["owner"])
  default_ownership_tag_overrides = true
}

run "resource_integration_azure_create_with_required_fields" {

  module {
    source = "./integration_azure_resources"
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

run "resource_integration_azure_set_ownership_tag_keys" {

  variables {
    ownership_tag_keys = ["one", "two", "three", "four", "five"]
  }

  module {
    source = "./integration_azure_resources"
  }

  assert {
    condition = opslevel_integration_azure_resources.this.ownership_tag_keys == var.default_ownership_tag_keys
    error_message = format(
      "expected default '%v' but got '%v'",
      var.default_ownership_tag_keys,
      opslevel_integration_azure_resources.this.ownership_tag_keys,
    )
  }

}

run "resource_integration_azure_set_ownership_tag_overrides" {

  variables {
    ownership_tag_overrides = false
  }

  module {
    source = "./integration_azure_resources"
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

