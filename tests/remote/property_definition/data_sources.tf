# PropertyDefinition data sources

data "opslevel_property_definitions" "all" {}

data "opslevel_property_definition" "first_property_definition_by_id" {
  identifier = data.opslevel_property_definitions.all.property_definitions[0].id
}
