resource "opslevel_integration_custom" "dev" {
  name = "Custom Integration"

  # Optional - but omitting it creates an integration with no mapping, which
  # ingests nothing until definitions are set.
  # The two definitions are managed as a unit, so set both together.
  etl_definition = {
    extract_definition = <<-EOT
      extractors:
        - external_kind: widget
          external_id: ".id"
    EOT

    transform_definition = <<-EOT
      transforms:
        - external_kind: widget
          opslevel_kind: component
          opslevel_identifier: ".id"
          on_component_not_found: create
          properties:
            name: ".name"
    EOT
  }
}

# The endpoint to POST payloads to, for push based extractors.
output "custom_integration_webhook_url" {
  value = opslevel_integration_custom.dev.webhook_url
}
