data "opslevel_property_definition" "pd1" {
  identifier = "id_or_alias"
}

output "pd_schema" {
  value = data.opslevel_property_definition.pd1.schema
}
