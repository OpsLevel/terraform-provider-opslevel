variables {
  component_type_one = "opslevel_component_type"
  identifier         = "service"
}

run "datasource_component_type_all_fields_accessible" {
  command = plan

  variables {
    identifier = var.identifier
  }

  module {
    source = "./data/component_type"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_component_type.test.id),
      can(data.opslevel_component_type.test.name),
      can(data.opslevel_component_type.test.alias),
      can(data.opslevel_component_type.test.description),
      can(data.opslevel_component_type.test.icon),
      can(data.opslevel_component_type.test.properties),
    ])
    error_message = format("'%s' data source missing expected fields", var.component_type_one)
  }

  assert {
    condition     = data.opslevel_component_type.test.id != null && data.opslevel_component_type.test.id != ""
    error_message = format("'%s' data source should return a valid id", var.component_type_one)
  }

  assert {
    condition     = data.opslevel_component_type.test.name != null && data.opslevel_component_type.test.name != ""
    error_message = format("'%s' data source should return a valid name", var.component_type_one)
  }

  assert {
    condition     = data.opslevel_component_type.test.alias == var.identifier
    error_message = format("'%s' data source should return the correct alias", var.component_type_one)
  }
}

run "datasource_component_type_read_service" {
  command = plan

  variables {
    identifier = "service"
  }

  module {
    source = "./data/component_type"
  }

  assert {
    condition     = data.opslevel_component_type.test.name == "Service"
    error_message = format("'%s' data source should return correct name", var.component_type_one)
  }
}

