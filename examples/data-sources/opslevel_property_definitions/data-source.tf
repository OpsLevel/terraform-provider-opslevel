data "opslevel_property_definitions" "all" {}

output "all_schemas" {
  value = data.opslevel_property_definitions.all.schemas
}
