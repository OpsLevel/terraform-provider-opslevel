resource "opslevel_integration_kubernetes" "this" {
  etl_definition = var.etl_definition
  name           = var.name
}
