mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_custom_event" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_custom_event.example.name == "foo"
    error_message = "wrong value name for opslevel_check_custom_event.example"
  }

  assert {
    condition     = opslevel_check_custom_event.example.integration == var.test_id
    error_message = "wrong value integration for opslevel_check_custom_event.example"
  }

  assert {
    condition     = opslevel_check_custom_event.example.pass_pending == true
    error_message = "wrong value pass_pending for opslevel_check_custom_event.example"
  }

  assert {
    condition     = opslevel_check_custom_event.example.service_selector == ".messages[] | .incident.service.id"
    error_message = "wrong value service_selector for opslevel_check_custom_event.example"
  }

  assert {
    condition     = opslevel_check_custom_event.example.success_condition == ".messages[] |   select(.incident.service.id == $ctx.alias) | .incident.status == \"resolved\""
    error_message = "wrong value success_condition for opslevel_check_custom_event.example"
  }
}