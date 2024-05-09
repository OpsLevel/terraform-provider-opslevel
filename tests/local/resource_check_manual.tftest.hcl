mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_manual" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_manual.example.name == "foo"
    error_message = "wrong value name for opslevel_check_manual.example"
  }

  assert {
    condition     = opslevel_check_manual.example.enabled == false
    error_message = "wrong value enabled on opslevel_check_manual.example"
  }

  assert {
    condition     = can(opslevel_check_manual.example.id)
    error_message = "id attribute missing from in opslevel_check_manual.example"
  }

  assert {
    condition     = can(opslevel_check_manual.example.owner)
    error_message = "owner attribute missing from in opslevel_check_manual.example"
  }

  assert {
    condition     = can(opslevel_check_manual.example.filter)
    error_message = "filter attribute missing from in opslevel_check_manual.example"
  }

  assert {
    condition     = can(opslevel_check_manual.example.category)
    error_message = "category attribute missing from in opslevel_check_manual.example"
  }

  assert {
    condition     = can(opslevel_check_manual.example.level)
    error_message = "level attribute missing from in opslevel_check_manual.example"
  }

  assert {
    condition     = opslevel_check_manual.example.notes == "Optional additional info on why this check is run or how to fix it"
    error_message = "wrong value for notes in opslevel_check_manual.example"
  }

  assert {
    condition = opslevel_check_manual.example.update_frequency == {
      starting_date = "2020-02-12T06:36:13Z"
      time_scale    = "week"
      value         = 1
    }
    error_message = "wrong update_frequency in opslevel_check_manual.example"
  }

  assert {
    condition     = opslevel_check_manual.example.update_requires_comment == false
    error_message = "wrong value for update_requires_comment in opslevel_check_manual.example"
  }
}