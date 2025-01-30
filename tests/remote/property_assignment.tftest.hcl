variables {
  resource_name = "opslevel_property_assignment"

  # required fields
  definition = null
  owner      = null

  # optional fields
  value = jsonencode(true)
}

run "from_property_definition_module" {
  command = plan

  variables {
    allowed_in_config_files = false
    name                    = ""
    schema                  = jsonencode(null)
    property_display_status = "visible"
    locked_status           = "unlocked"
  }

  module {
    source = "./property_definition"
  }
}

run "from_service_module" {
  command = plan

  variables {
    name = ""
  }

  module {
    source = "./service"
  }
}

run "resource_property_assignment_create_with_all_fields" {

  variables {
    definition = run.from_property_definition_module.first_property_definitions.id
    owner      = run.from_service_module.first_service.id
    value      = var.value
  }

  module {
    source = "./property_assignment"
  }

  assert {
    condition = alltrue([
      can(opslevel_property_assignment.test.definition),
      can(opslevel_property_assignment.test.id),
      can(opslevel_property_assignment.test.locked),
      can(opslevel_property_assignment.test.owner),
      can(opslevel_property_assignment.test.value),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_property_assignment.test.definition == var.definition
    error_message = format(
      "expected '%v' but got '%v'",
      var.definition,
      opslevel_property_assignment.test.definition,
    )
  }

  assert {
    condition     = startswith(opslevel_property_assignment.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_property_assignment.test.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_property_assignment.test.owner,
    )
  }

  assert {
    condition = opslevel_property_assignment.test.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_property_assignment.test.value,
    )
  }

}

run "resource_property_assignment_update_unset_optional_fields" {

  variables {
    definition = run.from_property_definition_module.first_property_definitions.id
    owner      = run.from_service_module.first_service.id
    value      = null
  }

  module {
    source = "./property_assignment"
  }

  assert {
    condition     = opslevel_property_assignment.test.value == null
    error_message = var.error_expected_null_field
  }

}
