mock_data "opslevel_webhook_action" {
  defaults = {
    description = "mock-webhook-action-description"
    headers = {
      "Accept"        = "application/json"
      "Content-Type"  = "application/json"
      "Authorization" = "Basic cXTlcjpeYPNzd29fZA=="
    }
    method  = "GET"
    name    = "mock-domain-name"
    payload = "{\"token\": \"XXX\", \"ref\":\"main\", \"action\": \"rollback\"}"
    url     = "https://www.opslevel.com/"
  }
}

