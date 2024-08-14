mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_integration_google_cloud" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition = alltrue([
      can(opslevel_integration_google_cloud.example.aliases),
      can(opslevel_integration_google_cloud.example.client_email),
      can(opslevel_integration_google_cloud.example.created_at),
      can(opslevel_integration_google_cloud.example.id),
      can(opslevel_integration_google_cloud.example.installed_at),
      can(opslevel_integration_google_cloud.example.name),
      can(opslevel_integration_google_cloud.example.ownership_tag_keys),
      can(opslevel_integration_google_cloud.example.private_key),
      can(opslevel_integration_google_cloud.example.projects),
    ])
    error_message = "expected opslevel_integration_google_cloud to have all attributes"
  }

  assert {
    condition     = opslevel_integration_google_cloud.example.client_email == "hello-world-tf@powerful-surf-427415-v1.iam.gserviceaccount.com"
    error_message = "expected opslevel_integration_google_cloud to have the correct client_email"
  }

  assert {
    condition     = opslevel_integration_google_cloud.example.name == "GCP Dev Integration"
    error_message = "expected opslevel_integration_google_cloud to have the correct name"
  }

  assert {
    condition     = opslevel_integration_google_cloud.example.ownership_tag_keys == toset(["opslevel_team", "team", "owner"])
    error_message = "expected opslevel_integration_google_cloud to have the correct ownership_tag_keys"
  }

  assert {
    condition = opslevel_integration_google_cloud.example.projects == tolist([
      {
        id   = "pauls-test-project-431117"
        name = "Paul's Test Project"
        url  = "https://console.cloud.google.com/home/dashboard?project=pauls-test-project-431117"
      },
      {
        id   = "powerful-surf-427415"
        name = "Powerful Surf"
        url  = "https://console.cloud.google.com/home/dashboard?project=powerful-surf-427415"
      }
    ])
    error_message = "expected opslevel_integration_google_cloud to have the correct projects"
  }

  assert {
    condition     = opslevel_integration_google_cloud.example.ownership_tag_overrides == true
    error_message = "expected opslevel_integration_google_cloud to have the correct ownership_tag_overrides"
  }
}

run "resource_integration_google_cloud_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition = alltrue([
      can(opslevel_integration_google_cloud.small.aliases),
      can(opslevel_integration_google_cloud.small.client_email),
      can(opslevel_integration_google_cloud.small.created_at),
      can(opslevel_integration_google_cloud.small.id),
      can(opslevel_integration_google_cloud.small.installed_at),
      can(opslevel_integration_google_cloud.small.name),
      can(opslevel_integration_google_cloud.small.ownership_tag_keys),
      can(opslevel_integration_google_cloud.small.private_key),
      can(opslevel_integration_google_cloud.small.projects),
    ])
    error_message = "expected opslevel_integration_google_cloud to have all attributes"
  }

  assert {
    condition     = opslevel_integration_google_cloud.small.client_email == "hello-world-tf@powerful-surf-427415-v1.iam.gserviceaccount.com"
    error_message = "expected opslevel_integration_google_cloud to have the correct client_email"
  }

  assert {
    condition     = opslevel_integration_google_cloud.small.name == "GCP Dev Integration"
    error_message = "expected opslevel_integration_google_cloud to have the correct name"
  }

  assert {
    condition     = length(opslevel_integration_google_cloud.small.ownership_tag_keys) == 0
    error_message = "expected opslevel_integration_google_cloud to have empty ownership_tag_keys"
  }

  assert {
    condition = opslevel_integration_google_cloud.small.projects == tolist([
      {
        id   = "pauls-test-project-431117"
        name = "Paul's Test Project"
        url  = "https://console.cloud.google.com/home/dashboard?project=pauls-test-project-431117"
      },
      {
        id   = "powerful-surf-427415"
        name = "Powerful Surf"
        url  = "https://console.cloud.google.com/home/dashboard?project=powerful-surf-427415"
      }
    ])
    error_message = "expected opslevel_integration_google_cloud to have the correct projects"
  }

  assert {
    condition     = opslevel_integration_google_cloud.small.ownership_tag_overrides == false
    error_message = "expected opslevel_integration_google_cloud to have the correct ownership_tag_overrides"
  }
}
