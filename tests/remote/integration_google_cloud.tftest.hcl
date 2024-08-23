variables {
  client_email = "hello-world-tf@powerful-surf-427415-v1.iam.gserviceaccount.com"
  private_key  = "abc123"
}

run "resource_integration_google_cloud_with_optional_fields" {
  variables {
    name                    = "GCP Integration Has Ownership Tag Keys and True Ownership Tag Overrides"
    ownership_tag_keys      = ["opslevel_team", "team", "owner"]
    ownership_tag_overrides = true
  }

  module {
    source = "./integration_google_cloud"
  }

  assert {
    condition = alltrue([
      can(opslevel_integration_google_cloud.test.aliases),
      can(opslevel_integration_google_cloud.test.client_email),
      can(opslevel_integration_google_cloud.test.created_at),
      can(opslevel_integration_google_cloud.test.id),
      can(opslevel_integration_google_cloud.test.installed_at),
      can(opslevel_integration_google_cloud.test.name),
      can(opslevel_integration_google_cloud.test.ownership_tag_keys),
      can(opslevel_integration_google_cloud.test.private_key),
      can(opslevel_integration_google_cloud.test.projects),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", "opslevel_integration_google_cloud")
  }

  assert {
    condition     = opslevel_integration_google_cloud.test.client_email == var.client_email
    error_message = replace(var.error_wrong_value, "TYPE", "opslevel_integration_google_cloud")
  }

  assert {
    condition     = opslevel_integration_google_cloud.test.name == var.name
    error_message = replace(var.error_wrong_value, "TYPE", "opslevel_integration_google_cloud")
  }

  assert {
    condition     = opslevel_integration_google_cloud.test.ownership_tag_keys == var.ownership_tag_keys
    error_message = replace(var.error_wrong_value, "TYPE", "opslevel_integration_google_cloud")
  }

  assert {
    condition     = opslevel_integration_google_cloud.test.ownership_tag_overrides == var.ownership_tag_overrides
    error_message = replace(var.error_wrong_value, "TYPE", "opslevel_integration_google_cloud")
  }
}

run "resource_integration_google_cloud_without_optional_fields" {
  variables {
    name = "GCP Integration Does Not Have Ownership Tag Keys and Ownership Tag Overrides"
  }

  module {
    source = "./integration_google_cloud"
  }

  assert {
    condition = alltrue([
      can(opslevel_integration_google_cloud.test.aliases),
      can(opslevel_integration_google_cloud.test.client_email),
      can(opslevel_integration_google_cloud.test.created_at),
      can(opslevel_integration_google_cloud.test.id),
      can(opslevel_integration_google_cloud.test.installed_at),
      can(opslevel_integration_google_cloud.test.name),
      can(opslevel_integration_google_cloud.test.ownership_tag_keys),
      can(opslevel_integration_google_cloud.test.private_key),
      can(opslevel_integration_google_cloud.test.projects),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", "opslevel_integration_google_cloud")
  }

  assert {
    condition     = opslevel_integration_google_cloud.test.client_email == var.client_email
    error_message = replace(var.error_wrong_value, "TYPE", "opslevel_integration_google_cloud")
  }

  assert {
    condition     = opslevel_integration_google_cloud.test.name == var.name
    error_message = replace(var.error_wrong_value, "TYPE", "opslevel_integration_google_cloud")
  }

  assert {
    condition     = opslevel_integration_google_cloud.test.ownership_tag_keys == tolist(["owner"])
    error_message = replace(var.error_wrong_value, "TYPE", "opslevel_integration_google_cloud")
  }

  assert {
    condition     = opslevel_integration_google_cloud.test.ownership_tag_overrides == true
    error_message = replace(var.error_wrong_value, "TYPE", "opslevel_integration_google_cloud")
  }
}

run "resource_integration_google_cloud_with_empty_optional_fields" {
  variables {
    name                    = "GCP Integration Has Empty Ownership Tag Keys and False Ownership Tag Overrides"
    ownership_tag_keys      = []
    ownership_tag_overrides = false
  }

  module {
    source = "./integration_google_cloud"
  }

  assert {
    condition = alltrue([
      can(opslevel_integration_google_cloud.test.aliases),
      can(opslevel_integration_google_cloud.test.client_email),
      can(opslevel_integration_google_cloud.test.created_at),
      can(opslevel_integration_google_cloud.test.id),
      can(opslevel_integration_google_cloud.test.installed_at),
      can(opslevel_integration_google_cloud.test.name),
      can(opslevel_integration_google_cloud.test.ownership_tag_keys),
      can(opslevel_integration_google_cloud.test.private_key),
      can(opslevel_integration_google_cloud.test.projects),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", "opslevel_integration_google_cloud")
  }

  assert {
    condition     = opslevel_integration_google_cloud.test.client_email == var.client_email
    error_message = replace(var.error_wrong_value, "TYPE", "opslevel_integration_google_cloud")
  }

  assert {
    condition     = opslevel_integration_google_cloud.test.name == var.name
    error_message = replace(var.error_wrong_value, "TYPE", "opslevel_integration_google_cloud")
  }

  assert {
    condition     = opslevel_integration_google_cloud.test.ownership_tag_keys == var.ownership_tag_keys
    error_message = replace(var.error_wrong_value, "TYPE", "opslevel_integration_google_cloud")
  }

  assert {
    condition     = opslevel_integration_google_cloud.test.ownership_tag_overrides == var.ownership_tag_overrides
    error_message = replace(var.error_wrong_value, "TYPE", "opslevel_integration_google_cloud")
  }
}

run "resource_integration_google_cloud_ownership_tag_keys_default_value" {

  variables {
    name               = "GCP Integration Default Ownership Tag Keys"
    ownership_tag_keys = null
  }

  module {
    source = "./integration_google_cloud"
  }

  assert {
    condition = opslevel_integration_google_cloud.test.ownership_tag_keys == tolist(["owner"])
    error_message = format(
      "expected '%v' but got '%v'",
      var.ownership_tag_keys,
      opslevel_integration_google_cloud.test.ownership_tag_keys,
    )
  }

}
