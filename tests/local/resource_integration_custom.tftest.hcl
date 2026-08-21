mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_integration_custom_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_integration_custom.example.id)
    error_message = "expected opslevel_integration_custom to have an ID"
  }

  assert {
    condition     = can(opslevel_integration_custom.example.webhook_url)
    error_message = "expected opslevel_integration_custom to expose webhook_url"
  }

  assert {
    condition     = opslevel_integration_custom.example.name == "dev"
    error_message = "wrong name for opslevel_integration_custom.example"
  }

  assert {
    condition     = opslevel_integration_custom.example.etl_definition.extract_definition == <<-EOT
      extractors:
        - external_kind: widget
          external_id: ".id"
    EOT
    error_message = "wrong extract_definition for opslevel_integration_custom.example"
  }

  assert {
    condition     = opslevel_integration_custom.example.etl_definition.transform_definition == <<-EOT
      transforms:
        - external_kind: widget
          opslevel_kind: component
          opslevel_identifier: ".id"
          on_component_not_found: create
          properties:
            name: ".name"
    EOT
    error_message = "wrong transform_definition for opslevel_integration_custom.example"
  }
}
