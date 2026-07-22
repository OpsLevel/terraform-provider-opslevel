resource "opslevel_integration_kubernetes" "dev" {
  name = "Kubernetes Integration"

  # Optional - OpsLevel's default definitions are used when unset.
  # The two definitions are managed as a unit, so set both together.
  etl_definition = {
    extract_definition = <<-EOT
      extractors:
        - external_kind: "v1/Namespace"
          external_id: ".metadata.uid"
        - external_kind: "apps/v1/Deployment"
          external_id: ".metadata.uid"
    EOT

    transform_definition = <<-EOT
      transforms:
        - external_kind: v1/Namespace
          opslevel_kind: kubernetes_namespace
          opslevel_identifier: ".metadata.name"
          on_component_not_found: create
          properties:
            name: ".metadata.name"
    EOT
  }
}
