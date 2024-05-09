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