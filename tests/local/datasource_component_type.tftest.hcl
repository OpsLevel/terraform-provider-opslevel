mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_component_type_all_fields_accessible" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition = alltrue([
      can(data.opslevel_component_type.mock_component_type.id),
      can(data.opslevel_component_type.mock_component_type.name),
      can(data.opslevel_component_type.mock_component_type.alias),
      can(data.opslevel_component_type.mock_component_type.description),
      can(data.opslevel_component_type.mock_component_type.icon),
      can(data.opslevel_component_type.mock_component_type.properties),
      can(data.opslevel_component_type.mock_component_type.relationships),
    ])
    error_message = "Not all expected fields are accessible from opslevel_component_type data source"
  }

  assert {
    condition     = data.opslevel_component_type.mock_component_type.name == "Service"
    error_message = "component_type data source should return correct name"
  }

  assert {
    condition     = data.opslevel_component_type.mock_component_type.alias == "service"
    error_message = "component_type data source should return correct alias"
  }
}

run "datasource_component_type_icon_structure" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(data.opslevel_component_type.mock_component_type.icon.color)
    error_message = "component_type icon should have color field"
  }

  assert {
    condition     = can(data.opslevel_component_type.mock_component_type.icon.name)
    error_message = "component_type icon should have name field"
  }
}

run "datasource_component_type_relationships_structure" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(data.opslevel_component_type.mock_component_type.relationships)
    error_message = "component_type should have relationships field"
  }

  assert {
    condition     = data.opslevel_component_type.mock_component_type.relationships != null
    error_message = "component_type relationships should not be null"
  }
}

