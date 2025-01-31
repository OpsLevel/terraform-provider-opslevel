mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_property_definition_mocked_fields" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_property_definition.mock_property_definition.allowed_in_config_files == true
    error_message = "allowed_in_config_files in mock opslevel_property_definition not set to true"
  }

  assert {
    condition     = data.opslevel_property_definition.mock_property_definition.description == "mock-property-definition-description"
    error_message = "wrong description in mock opslevel_property_definition"
  }

  assert {
    condition     = data.opslevel_property_definition.mock_property_definition.id != ""
    error_message = "id in mock opslevel_property_definition was not set"
  }

  assert {
    condition     = data.opslevel_property_definition.mock_property_definition.identifier == "mock-property_definition-alias"
    error_message = "wrong identifier in mock opslevel_property_definition"
  }

  assert {
    condition     = data.opslevel_property_definition.mock_property_definition.name == "mock-property-definition-name"
    error_message = "wrong name in mock opslevel_property_definition"
  }

  assert {
    condition     = data.opslevel_property_definition.mock_property_definition.property_display_status == "visible"
    error_message = "wrong property_display_status in mock opslevel_property_definition"
  }

  assert {
    condition     = data.opslevel_property_definition.mock_property_definition.locked_status == "unlocked"
    error_message = "wrong locked_status in mock opslevel_property_definition"
  }

  assert {
    condition = data.opslevel_property_definition.mock_property_definition.schema == jsonencode(
      {
        "$ref" : "#/$defs/MyProp",
        "$defs" : {
          "MyProp" : {
            "properties" : {
              "name" : {
                "type" : "string",
                "title" : "the new name",
                "description" : "The name of a friend",
                "default" : "alex",
                "examples" : ["joe", "lucy"]
              }
            },
            "additionalProperties" : false,
            "type" : "object",
            "required" : ["name"]
          }
        }
      }
    )
    error_message = "wrong schema in mock opslevel_property_definition"
  }

}
