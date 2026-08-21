resource "opslevel_integration_custom" "dev" {
  name = "Custom Integration"

  # Optional - OpsLevel's default definitions are used when omitted on create.
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
