mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_has_recent_deploy" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.example.name == "foo"
    error_message = "wrong value name for opslevel_check_has_recent_deploy.example"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.example.days == 14
    error_message = "wrong value for days in opslevel_check_has_recent_deploy.example"
  }
}