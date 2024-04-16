mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_alert_source_usage" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition = opslevel_check_alert_source_usage.example.alert_name_predicate == {
      type  = "contains"
      value = "dev"
    }
    error_message = "wrong alert_name_predicate in opslevel_check_alert_source_usage.example"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.example.alert_type == "pagerduty"
    error_message = "wrong value for alert_type in opslevel_check_alert_source_usage.example"
  }
}