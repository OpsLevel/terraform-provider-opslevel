mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_service_ownership" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_service_ownership.example.name == "foo"
    error_message = "wrong value name for opslevel_check_service_ownership.example"
  }

  assert {
    condition     = opslevel_check_service_ownership.example.require_contact_method == true
    error_message = "wrong value for require_contact_method in opslevel_check_service_ownership.example"
  }

  assert {
    condition     = opslevel_check_service_ownership.example.contact_method == "ANY"
    error_message = "wrong value for contact_method in opslevel_check_service_ownership.example"
  }

  assert {
    condition     = opslevel_check_service_ownership.example.tag_key == "team"
    error_message = "wrong value for tag_key in opslevel_check_service_ownership.example"
  }

  assert {
    condition = opslevel_check_service_ownership.example.tag_predicate == {
      type  = "equals"
      value = "frontend"
    }
    error_message = "wrong tag_predicate in opslevel_check_service_ownership.example"
  }
}