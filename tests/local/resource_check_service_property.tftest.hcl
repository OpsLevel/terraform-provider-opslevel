mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_service_property" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_service_property.example.property == "language"
    error_message = "wrong value for property in opslevel_check_service_property.example"
  }

  assert {
    condition = opslevel_check_service_property.example.predicate == {
      type  = "equals"
      value = "python"
    }
    error_message = "wrong value for predicate in opslevel_check_service_property.example"
  }
}

run "resource_check_service_property_with_property_definition_alias_only" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_service_property.with_property_definition_alias_only.property == "custom_property"
    error_message = "wrong value for property"
  }

  assert {
    condition     = opslevel_check_service_property.with_property_definition_alias_only.property_definition == "test_property_alias"
    error_message = "property_definition should be preserved in state even when API returns null (cross-component-type behavior)"
  }

  assert {
    condition     = opslevel_check_service_property.with_property_definition_alias_only.component_type == null
    error_message = "component_type should be null when not specified"
  }

  assert {
    condition     = opslevel_check_service_property.with_property_definition_alias_only.predicate.type == "exists"
    error_message = "predicate type should be 'exists'"
  }
}

run "resource_check_service_property_with_property_definition_and_component_type" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_service_property.with_property_definition_and_component_type.property == "custom_property"
    error_message = "wrong value for property"
  }

  assert {
    condition     = opslevel_check_service_property.with_property_definition_and_component_type.property_definition == "test_property_alias"
    error_message = "property_definition should be set correctly"
  }

  assert {
    condition     = opslevel_check_service_property.with_property_definition_and_component_type.component_type != null
    error_message = "component_type should be set when specified"
  }

  assert {
    condition     = opslevel_check_service_property.with_property_definition_and_component_type.predicate.type == "exists"
    error_message = "predicate type should be 'exists'"
  }
}