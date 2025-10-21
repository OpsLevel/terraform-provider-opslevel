mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_component_type_api" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_component_type.api.id)
    error_message = "id attribute missing from opslevel_component_type.api"
  }

  assert {
    condition     = opslevel_component_type.api.name == "API"
    error_message = "wrong name for opslevel_component_type.api"
  }

  assert {
    condition     = opslevel_component_type.api.alias == "api"
    error_message = "wrong alias for opslevel_component_type.api"
  }

  assert {
    condition     = opslevel_component_type.api.description == "An API component type"
    error_message = "wrong description for opslevel_component_type.api"
  }
}

run "resource_component_type_icon_structure" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_component_type.api.icon)
    error_message = "icon attribute missing from opslevel_component_type.api"
  }

  assert {
    condition     = opslevel_component_type.api.icon.color == "#F59E0B"
    error_message = "wrong icon color for opslevel_component_type.api"
  }

  assert {
    condition     = opslevel_component_type.api.icon.name == "PhCloud"
    error_message = "wrong icon name for opslevel_component_type.api"
  }
}

run "resource_component_type_properties" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_component_type.api.properties)
    error_message = "properties attribute missing from opslevel_component_type.api"
  }

  assert {
    condition     = can(opslevel_component_type.api.properties["api_version"])
    error_message = "api_version property missing from opslevel_component_type.api"
  }

  assert {
    condition     = opslevel_component_type.api.properties["api_version"].name == "API Version"
    error_message = "wrong property name for api_version in opslevel_component_type.api"
  }

  assert {
    condition     = opslevel_component_type.api.properties["api_version"].description == "The version of the API"
    error_message = "wrong property description for api_version in opslevel_component_type.api"
  }

  assert {
    condition     = opslevel_component_type.api.properties["api_version"].allowed_in_config_files == true
    error_message = "wrong allowed_in_config_files for api_version in opslevel_component_type.api"
  }

  assert {
    condition     = opslevel_component_type.api.properties["api_version"].display_status == "visible"
    error_message = "wrong display_status for api_version in opslevel_component_type.api"
  }
}

run "resource_component_type_relationships" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_component_type.api.relationships)
    error_message = "relationships attribute missing from opslevel_component_type.api"
  }
}

