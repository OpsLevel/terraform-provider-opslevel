resource "opslevel_integration_kubernetes" "dev" {
  name = "Kubernetes Integration"

  # Both definitions are optional - OpsLevel's defaults are used when unset.
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
