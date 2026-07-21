variables {
  resource_name = "opslevel_integration_kubernetes"

  # required fields
  name = "TF Test Kubernetes Integration"

  # optional fields
  extract_definition = <<-EOT
    extractors:
      - external_kind: "v1/Namespace"
        external_id: ".metadata.uid"
  EOT

  transform_definition = <<-EOT
    transforms:
      - external_kind: v1/Namespace
        opslevel_kind: kubernetes_namespace
        opslevel_identifier: ".metadata.name"
        on_component_not_found: create
  EOT
}

run "resource_integration_kubernetes_create_with_all_fields" {

  module {
    source = "./integration_kubernetes"
  }

  assert {
    condition = alltrue([
      can(opslevel_integration_kubernetes.this.extract_definition),
      can(opslevel_integration_kubernetes.this.id),
      can(opslevel_integration_kubernetes.this.name),
      can(opslevel_integration_kubernetes.this.transform_definition),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition     = startswith(opslevel_integration_kubernetes.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_integration_kubernetes.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_integration_kubernetes.this.name,
    )
  }

  assert {
    condition = opslevel_integration_kubernetes.this.extract_definition == var.extract_definition
    error_message = format(
      "expected '%v' but got '%v'",
      var.extract_definition,
      opslevel_integration_kubernetes.this.extract_definition,
    )
  }

  assert {
    condition = opslevel_integration_kubernetes.this.transform_definition == var.transform_definition
    error_message = format(
      "expected '%v' but got '%v'",
      var.transform_definition,
      opslevel_integration_kubernetes.this.transform_definition,
    )
  }

}

run "resource_integration_kubernetes_update_definitions" {

  variables {
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
        - external_kind: apps/v1/Deployment
          opslevel_kind: kubernetes_deployment
          opslevel_identifier: ".metadata.uid"
          on_component_not_found: create
    EOT
  }

  module {
    source = "./integration_kubernetes"
  }

  assert {
    condition = run.resource_integration_kubernetes_create_with_all_fields.this.id == opslevel_integration_kubernetes.this.id
    error_message = format(
      "expected id '%v' to be unchanged after update but got '%v'",
      run.resource_integration_kubernetes_create_with_all_fields.this.id,
      opslevel_integration_kubernetes.this.id,
    )
  }

  assert {
    condition = opslevel_integration_kubernetes.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_integration_kubernetes.this.name,
    )
  }

  assert {
    condition = opslevel_integration_kubernetes.this.extract_definition == var.extract_definition
    error_message = format(
      "expected '%v' but got '%v'",
      var.extract_definition,
      opslevel_integration_kubernetes.this.extract_definition,
    )
  }

  assert {
    condition = opslevel_integration_kubernetes.this.transform_definition == var.transform_definition
    error_message = format(
      "expected '%v' but got '%v'",
      var.transform_definition,
      opslevel_integration_kubernetes.this.transform_definition,
    )
  }

}
