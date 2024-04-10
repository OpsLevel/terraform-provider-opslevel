mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_has_documentation" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_has_documentation.example.name == "foo"
    error_message = "wrong value name for opslevel_check_has_documentation.example"
  }

  assert {
    condition     = opslevel_check_has_documentation.example.enabled == true
    error_message = "wrong value enabled on opslevel_check_has_documentation.example"
  }

  assert {
    condition     = can(opslevel_check_has_documentation.example.id)
    error_message = "id attribute missing from in opslevel_check_has_documentation.example"
  }

  assert {
    condition     = can(opslevel_check_has_documentation.example.owner)
    error_message = "owner attribute missing from in opslevel_check_has_documentation.example"
  }

  assert {
    condition     = can(opslevel_check_has_documentation.example.filter)
    error_message = "filter attribute missing from in opslevel_check_has_documentation.example"
  }

  assert {
    condition     = can(opslevel_check_has_documentation.example.category)
    error_message = "category attribute missing from in opslevel_check_has_documentation.example"
  }

  assert {
    condition     = can(opslevel_check_has_documentation.example.level)
    error_message = "level attribute missing from in opslevel_check_has_documentation.example"
  }

  assert {
    condition     = opslevel_check_has_documentation.example.notes == null
    error_message = "wrong value for notes in opslevel_check_has_documentation.example"
  }

  assert {
    condition     = opslevel_check_has_documentation.example.document_type == "api"
    error_message = "wrong document_type in opslevel_check_has_documentation.example"
  }

  assert {
    condition     = opslevel_check_has_documentation.example.document_subtype == "openapi"
    error_message = "wrong value for document_subtype in opslevel_check_has_documentation.example"
  }
}