mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_webhook_action" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_webhook_action.mock.description == "Pages the On Call"
    error_message = "wrong description for opslevel_webhook_action.mock"
  }

  assert {
    condition     = opslevel_webhook_action.mock.headers["accept"] == "application/vnd.pagerduty+json;version=2"
    error_message = "wrong 'accept' header in opslevel_webhook_action mock"
  }

  assert {
    condition     = opslevel_webhook_action.mock.headers["authorization"] == "Token token=XXXXXXXXXXXXXX"
    error_message = "wrong 'authorization' header in opslevel_webhook_action mock"
  }

  assert {
    condition     = opslevel_webhook_action.mock.headers["content-type"] == "application/json"
    error_message = "wrong 'content-type' header in opslevel_webhook_action mock"
  }

  assert {
    condition     = opslevel_webhook_action.mock.headers["from"] == "foo@opslevel.com"
    error_message = "wrong 'from' header in opslevel_webhook_action mock"
  }

  assert {
    condition     = opslevel_webhook_action.mock.method == "POST"
    error_message = "wrong http method for opslevel_webhook_action.mock"
  }

  assert {
    condition     = opslevel_webhook_action.mock.name == "Small Webhook Action"
    error_message = "wrong name for opslevel_webhook_action.mock"
  }

  assert {
    condition = opslevel_webhook_action.mock.payload == <<EOT
{
    "incident":
    {
        "type": "incident",
        "title": "{{manualInputs.IncidentTitle}}",
        "service": {
        "id": "{{ service | tag_value: 'pd_id' }}",
        "type": "service_reference"
        },
        "body": {
        "type": "incident_body",
        "details": "Incident triggered from OpsLevel by {{user.name}} with the email {{user.email}}. {{manualInputs.IncidentDescription}}"
        }
    }
}
  EOT

    error_message = "wrong payload for opslevel_webhook_action.mock"
  }

  assert {
    condition     = opslevel_webhook_action.mock.url == "https://api.pagerduty.com/incidents"
    error_message = "wrong url for opslevel_webhook_action.mock"
  }

}
