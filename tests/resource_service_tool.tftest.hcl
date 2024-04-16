mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_service_with_alias" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_service_tool.with_alias.category == "logs"
    error_message = "wrong category in opslevel_service_tool.with_alias"
  }

  assert {
    condition     = opslevel_service_tool.with_alias.environment == "Production"
    error_message = "wrong environment in opslevel_service_tool.with_alias"
  }

  assert {
    condition     = opslevel_service_tool.with_alias.name == "Datadog test"
    error_message = "wrong name in opslevel_service_tool.with_alias"
  }

  assert {
    condition     = opslevel_service_tool.with_alias.service_alias == "mock service"
    error_message = "wrong service_alias in opslevel_service_tool.with_alias"
  }

  assert {
    condition     = opslevel_service_tool.with_alias.url == "https://datadoghq.com"
    error_message = "wrong url in opslevel_service_tool.with_alias"
  }

}

run "resource_service_with_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_service_tool.with_id.category == "metrics"
    error_message = "wrong category in opslevel_service_tool.with_id"
  }

  assert {
    condition     = can(opslevel_service_tool.with_id.id)
    error_message = "id attribute missing from opslevel_service_tool.with_id"
  }

}
