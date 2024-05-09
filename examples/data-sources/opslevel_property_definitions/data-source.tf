data "opslevel_property_definitions" "all" {}

output "all" {
  value = data.opslevel_property_definitions.all.property_definitions
}

output "property_definition_names" {
  value = sort(data.opslevel_property_definitions.all.property_definitions[*].name)
}

output "property_definition_schemas" {
  value = data.opslevel_property_definitions.all.property_definitions[*].schema
}
