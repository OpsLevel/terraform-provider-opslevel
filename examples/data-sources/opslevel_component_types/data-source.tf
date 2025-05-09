data "opslevel_component_types" "all" {}

output "all" {
  value = data.opslevel_component_types.all
}