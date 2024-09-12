variables {
  resource_name = "opslevel_integration_google_cloud"

  # required fields
  client_email = "hello-world-tf@powerful-surf-427415-v1.iam.gserviceaccount.com"
  name         = "abc123"
  private_key  = "Xx_super_secret_xX"

  # optional fields
  ownership_tag_keys      = null
  ownership_tag_overrides = null

  # default values - computed from API
  default_ownership_tag_keys      = tolist(["owner"])
  default_ownership_tag_overrides = true
}

run "resource_integration_google_cloud_create_with_required_fields" {

  module {
    source = "./integration_google_cloud"
  }

  assert {
    condition = opslevel_integration_google_cloud.this.client_email == var.client_email
    error_message = format(
      "expected '%v' but got '%v'",
      var.client_email,
      opslevel_integration_google_cloud.this.client_email,
    )
  }

  assert {
    condition     = startswith(opslevel_integration_google_cloud.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_integration_google_cloud.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_integration_google_cloud.this.name,
    )
  }

  assert {
    condition = opslevel_integration_google_cloud.this.ownership_tag_keys == var.default_ownership_tag_keys
    error_message = format(
      "expected default '%v' but got '%v'",
      var.default_ownership_tag_keys,
      opslevel_integration_google_cloud.this.ownership_tag_keys,
    )
  }

  assert {
    condition = opslevel_integration_google_cloud.this.ownership_tag_overrides == var.default_ownership_tag_overrides
    error_message = format(
      "expected default '%v' but got '%v'",
      var.default_ownership_tag_overrides,
      opslevel_integration_google_cloud.this.ownership_tag_overrides,
    )
  }

  assert {
    condition     = opslevel_integration_google_cloud.this.private_key == var.private_key
    error_message = "expected different private_key value, not printing sensitive value"
  }

}

run "resource_integration_google_cloud_set_ownership_tag_keys" {

  variables {
    ownership_tag_keys = ["owner", "team"]
  }

  module {
    source = "./integration_google_cloud"
  }

  assert {
    condition = opslevel_integration_google_cloud.this.ownership_tag_keys == var.ownership_tag_keys
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_keys,
      opslevel_integration_google_cloud.this.ownership_tag_keys,
    )
  }

}

run "resource_integration_google_cloud_set_ownership_tag_overrides" {

  variables {
    ownership_tag_overrides = false
  }

  module {
    source = "./integration_google_cloud"
  }

  assert {
    condition = opslevel_integration_google_cloud.this.ownership_tag_overrides == var.ownership_tag_overrides
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_overrides,
      opslevel_integration_google_cloud.this.ownership_tag_overrides,
    )
  }

}
