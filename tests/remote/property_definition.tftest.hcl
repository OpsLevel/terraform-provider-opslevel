variables {
  property_definition_one  = "opslevel_property_definition"
  property_definitions_all = "opslevel_property_definitions"

  # required fields
  allowed_in_config_files = true
  name                    = "TF Test Property Definition"
  property_display_status = "visible"
  schema                  = "{\"type\":\"boolean\"}"

  # optional fields
  description = "Property Definition description"
}

run "resource_property_definition_create_with_all_fields" {

  variables {
    allowed_in_config_files = var.allowed_in_config_files
    description             = var.description
    name                    = var.name
    property_display_status = var.property_display_status
    schema                  = var.schema
  }

  module {
    source = "./property_definition"
  }

  assert {
    condition = alltrue([
      can(opslevel_property_definition.test.allowed_in_config_files),
      can(opslevel_property_definition.test.description),
      can(opslevel_property_definition.test.id),
      can(opslevel_property_definition.test.last_updated),
      can(opslevel_property_definition.test.name),
      can(opslevel_property_definition.test.property_display_status),
      can(opslevel_property_definition.test.schema),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.property_definition_one)
  }

  assert {
    condition     = opslevel_property_definition.test.allowed_in_config_files == var.allowed_in_config_files
    error_message = "wrong allowed_in_config_files for opslevel_property_definition resource"
  }

  assert {
    condition     = opslevel_property_definition.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.property_definition_one)
  }

  assert {
    condition     = startswith(opslevel_property_definition.test.id, "Z2lkOi8v")
    error_message = replace(var.error_wrong_id, "TYPE", var.property_definition_one)
  }

  assert {
    condition     = opslevel_property_definition.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.property_definition_one)
  }

  assert {
    condition     = opslevel_property_definition.test.property_display_status == var.property_display_status
    error_message = "wrong property_display_status for opslevel_property_definition resource"
  }

  assert {
    condition     = opslevel_property_definition.test.schema == var.schema
    error_message = "wrong schema for opslevel_property_definition resource"
  }

}

run "resource_property_definition_update_unset_optional_fields" {

  variables {
    description = null
  }

  module {
    source = "./property_definition"
  }

  assert {
    condition     = opslevel_property_definition.test.description == null
    error_message = var.error_expected_null_field
  }

}

run "resource_property_definition_update_all_fields" {

  variables {
    allowed_in_config_files = !var.allowed_in_config_files
    description             = "${var.description} updated"
    name                    = "${var.name} updated"
    property_display_status = "hidden"
    schema                  = "{\"type\":\"string\"}"
  }

  module {
    source = "./property_definition"
  }

  assert {
    condition     = opslevel_property_definition.test.allowed_in_config_files == var.allowed_in_config_files
    error_message = "wrong allowed_in_config_files for opslevel_property_definition resource"
  }

  assert {
    condition     = opslevel_property_definition.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.property_definition_one)
  }

  assert {
    condition     = opslevel_property_definition.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.property_definition_one)
  }

  assert {
    condition     = opslevel_property_definition.test.property_display_status == var.property_display_status
    error_message = "wrong property_display_status for opslevel_property_definition resource"
  }

  assert {
    condition     = opslevel_property_definition.test.schema == var.schema
    error_message = "wrong schema for opslevel_property_definition resource"
  }

}

run "datasource_property_definitions_all" {

  module {
    source = "./property_definition"
  }

  assert {
    condition     = can(data.opslevel_property_definitions.all.property_definitions)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.property_definitions_all)
  }

  assert {
    condition     = length(data.opslevel_property_definitions.all.property_definitions) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.property_definitions_all)
  }

}

run "datasource_property_definition_first" {

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
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.property_definition_one)
  }

  assert {
    condition     = data.opslevel_property_definition.first_property_definition_by_id.id == data.opslevel_property_definitions.all.property_definitions[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.property_definition_one)
  }

}
