mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_property_assignment_using_aliases" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_property_assignment.color_picker_using_aliases.id == "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MzMyOQ:Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8zNzU"
    error_message = "expected ID (legacy field) to have 2 aliases separated by a ':'"
  }

  assert {
    condition     = can(opslevel_property_assignment.color_picker_using_aliases.last_updated)
    error_message = "expected last updated to exist"
  }

  assert {
    condition     = opslevel_property_assignment.color_picker_using_aliases.owner == "some_service"
    error_message = "unexpected value for owner"
  }

  assert {
    condition     = opslevel_property_assignment.color_picker_using_aliases.definition == "some_definition"
    error_message = "unexpected value for definition"
  }

  assert {
    condition     = opslevel_property_assignment.color_picker_using_aliases.value == "\"green\""
    error_message = "unexpected value field"
  }

  assert {
    condition     = opslevel_property_assignment.color_picker_using_aliases.locked == false
    error_message = "unexpected value for locked"
  }
}

run "resource_property_assignment_using_ids" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_property_assignment.color_picker_using_id.id == "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MzMyOQ:Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8zNzU"
    error_message = "expected ID (legacy field) to have 2 aliases separated by a ':'"
  }

  assert {
    condition     = can(opslevel_property_assignment.color_picker_using_id.last_updated)
    error_message = "expected last updated to exist"
  }

  assert {
    condition     = opslevel_property_assignment.color_picker_using_id.owner == "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MzMyOQ"
    error_message = "unexpected value for owner"
  }

  assert {
    condition     = opslevel_property_assignment.color_picker_using_id.definition == "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8zNzU"
    error_message = "unexpected value for definition"
  }

  assert {
    condition     = opslevel_property_assignment.color_picker_using_id.value == "{\"hello\":\"world\",\"key\":null}"
    error_message = "unexpected value field"
  }

  assert {
    condition     = opslevel_property_assignment.color_picker_using_id.locked == false
    error_message = "unexpected value for locked"
  }
}
