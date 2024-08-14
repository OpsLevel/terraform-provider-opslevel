variables {
  client_email = "hello-world-tf@powerful-surf-427415-v1.iam.gserviceaccount.com"
  name = "TF Remote Integration GCP"
  ownership_tag_keys = toset(["opslevel_team", "team", "owner"])
  ownership_tag_overrides = true
  private_key = "abc123"
}

run "resource_integration_google_cloud_create" {
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

run "resource_integration_google_cloud_update" {
  variables {
    name = "TF Remote Integration GCP Updated"
    ownership_tag_keys = toset([])
    ownership_tag_overrides = false
    private_key = "abc123456"
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
