mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_repository_integrated" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_repository_integrated.example.name == "foo"
    error_message = "wrong value name for opslevel_check_repository_integrated.example"
  }

  assert {
    condition     = opslevel_check_repository_integrated.example.enabled == true
    error_message = "wrong value enabled on opslevel_check_repository_integrated.example"
  }

  assert {
    condition     = can(opslevel_check_repository_integrated.example.id)
    error_message = "id attribute missing from in opslevel_check_repository_integrated.example"
  }

  assert {
    condition     = can(opslevel_check_repository_integrated.example.owner)
    error_message = "owner attribute missing from in opslevel_check_repository_integrated.example"
  }

  assert {
    condition     = can(opslevel_check_repository_integrated.example.filter)
    error_message = "filter attribute missing from in opslevel_check_repository_integrated.example"
  }

  assert {
    condition     = can(opslevel_check_repository_integrated.example.category)
    error_message = "category attribute missing from in opslevel_check_repository_integrated.example"
  }

  assert {
    condition     = can(opslevel_check_repository_integrated.example.level)
    error_message = "level attribute missing from in opslevel_check_repository_integrated.example"
  }
}