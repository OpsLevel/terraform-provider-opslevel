resource "opslevel_integration_kubernetes" "this" {
  extract_definition   = var.extract_definition
  name                 = var.name
  transform_definition = var.transform_definition
}
