data "opslevel_property_definition" "pd1" {
  id = "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8zNjg"
}

output "pd_schema" {
  value = data.opslevel_property_definition.pd1.schema
}
