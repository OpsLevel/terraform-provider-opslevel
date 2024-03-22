mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_webhook_action_mocked_fields" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_webhook_action.mock_webhook_action.description == "mock-webhook-action-description"
    error_message = "wrong description in opslevel_webhook_action mock"
  }

  assert {
    condition     = data.opslevel_webhook_action.mock_webhook_action.identifier == "mock-webhook-action-alias"
    error_message = "wrong identifier in opslevel_webhook_action mock"
  }

  assert {
    condition     = data.opslevel_webhook_action.mock_webhook_action.headers["Accept"] == "application/json"
    error_message = "wrong 'Accept' header in opslevel_webhook_action mock"
  }

  assert {
    condition     = data.opslevel_webhook_action.mock_webhook_action.headers["Content-Type"] == "application/json"
    error_message = "wrong 'Content-Type' header in opslevel_webhook_action mock"
  }

  assert {
    condition     = data.opslevel_webhook_action.mock_webhook_action.headers["Authorization"] == "Basic cXTlcjpeYPNzd29fZA=="
    error_message = "wrong 'Authorization' header in opslevel_webhook_action mock"
  }

  assert {
    condition     = contains(["GET", "PATCH", "POST", "PUT", "DELETE"], data.opslevel_webhook_action.mock_webhook_action.method)
    error_message = "wrong method in opslevel_webhook_action mock"
  }

  assert {
    condition     = data.opslevel_webhook_action.mock_webhook_action.name == "mock-domain-name"
    error_message = "wrong name in opslevel_webhook_action mock"
  }

  assert {
    condition     = data.opslevel_webhook_action.mock_webhook_action.url == "https://www.opslevel.com/"
    error_message = "wrong url in opslevel_webhook_action mock"
  }

}
