mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_property_definition" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_property_definition.color_picker.id)
    error_message = "expected ID to be set"
  }

  assert {
    condition     = can(opslevel_property_definition.color_picker.last_updated)
    error_message = "expected last updated to be set"
  }

  assert {
    condition     = opslevel_property_definition.color_picker.allowed_in_config_files == false
    error_message = "unexpected value for allowed_in_config_files"
  }

  assert {
    condition     = opslevel_property_definition.color_picker.name == "Color Picker"
    error_message = "unexpected value for name"
  }

  assert {
    condition     = opslevel_property_definition.color_picker.property_display_status == "visible"
    error_message = "unexpected value for property_display_status"
  }

  assert {
    condition = opslevel_property_definition.color_picker.schema == jsonencode({
      "type" : "string",
      "enum" : [
        "red",
        "green",
        "blue",
      ]
    })
    error_message = "unexpected value for schema"
  }
}
