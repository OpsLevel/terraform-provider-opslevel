mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_tool_usage" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_tool_usage.example.tool_category == "metrics"
    error_message = "wrong value for tool_category in opslevel_check_tool_usage.example"
  }

  assert {
    condition = opslevel_check_tool_usage.example.tool_name_predicate == {
      type  = "equals"
      value = "datadog"
    }
    error_message = "wrong tool_name_predicate in opslevel_check_tool_usage.example"
  }

  assert {
    condition = opslevel_check_tool_usage.example.environment_predicate == {
      type  = "equals"
      value = "production"
    }
    error_message = "wrong environment_predicate in opslevel_check_tool_usage.example"
  }
}