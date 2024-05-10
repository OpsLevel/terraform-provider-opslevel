run "datasource_property_definitions_all" {

  variables {
    datasource_type = "opslevel_property_definitions"
  }

  module {
    source = "./property_definition"
  }

  assert {
    condition     = can(data.opslevel_property_definitions.all.property_definitions)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_property_definitions.all.property_definitions) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.datasource_type)
  }

}

run "datasource_property_definition_first" {

  variables {
    datasource_type = "opslevel_property_definition"
  }

  module {
    source = "./property_definition"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_property_definition.first_property_definition_by_id.allowed_in_config_files),
      can(data.opslevel_property_definition.first_property_definition_by_id.description),
      can(data.opslevel_property_definition.first_property_definition_by_id.id),
      can(data.opslevel_property_definition.first_property_definition_by_id.identifier),
      can(data.opslevel_property_definition.first_property_definition_by_id.name),
      can(data.opslevel_property_definition.first_property_definition_by_id.property_display_status),
      can(data.opslevel_property_definition.first_property_definition_by_id.schema),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_property_definition.first_property_definition_by_id.id == data.opslevel_property_definitions.all.property_definitions[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.datasource_type)
  }

}
