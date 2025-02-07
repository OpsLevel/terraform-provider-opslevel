# PropertyDefinition data sources

data "opslevel_property_definitions" "all" {}

data "opslevel_property_definition" "first_property_definition_by_id" {
  identifier = data.opslevel_property_definitions.all.property_definitions[0].id
}

resource "opslevel_property_definition" "test" {
  allowed_in_config_files = var.allowed_in_config_files
  description             = var.description
  name                    = var.name
  property_display_status = var.property_display_status
  locked_status           = var.locked_status
  schema                  = var.schema
}
