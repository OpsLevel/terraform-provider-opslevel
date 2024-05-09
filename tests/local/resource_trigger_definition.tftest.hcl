mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_trigger_definition_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_trigger_definition.small.access_control == "everyone"
    error_message = "wrong access_control for opslevel_trigger_definition.small"
  }

  assert {
    condition     = opslevel_trigger_definition.small.action == var.test_id
    error_message = "wrong id for 'action' in opslevel_trigger_definition.small"
  }

  assert {
    condition     = opslevel_trigger_definition.small.name == "Small Trigger Definition"
    error_message = "wrong name for opslevel_trigger_definition.small"
  }

  assert {
    condition     = opslevel_trigger_definition.small.owner == var.test_id
    error_message = "wrong id for 'owner' in opslevel_trigger_definition.small"
  }

  assert {
    condition     = opslevel_trigger_definition.small.published == true
    error_message = "published should be set to 'true' for opslevel_trigger_definition.small"
  }

}

run "resource_trigger_definition_big" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_trigger_definition.big.description == "Pages the On Call"
    error_message = "wrong description for opslevel_trigger_definition.big"
  }

  assert {
    condition     = opslevel_trigger_definition.big.entity_type == "SERVICE"
    error_message = "wrong entity_type for opslevel_trigger_definition.big"
  }

  assert {
    condition     = opslevel_trigger_definition.big.extended_team_access == tolist(["team_1", "team_2"])
    error_message = "wrong extended_team_access for opslevel_trigger_definition.big"
  }

  assert {
    condition     = opslevel_trigger_definition.big.filter == var.test_id
    error_message = "wrong id for 'filter' in opslevel_trigger_definition.big"
  }

  assert {
    condition     = opslevel_trigger_definition.big.manual_inputs_definition == <<EOT
---
version: 1
inputs:
  - identifier: IncidentTitle
    displayName: Title
    description: Title of the incident to trigger
    type: text_input
    required: true
    maxLength: 60
    defaultValue: Service Incident Manual Trigger
  - identifier: IncidentDescription
    displayName: Incident Description
    description: The description of the incident
    type: text_area
    required: true
  EOT
    error_message = "wrong manual_inputs_definition in opslevel_trigger_definition.big"
  }

  assert {
    condition     = opslevel_trigger_definition.big.response_template == <<EOT
{% if response.status >= 200 and response.status < 300 %}
## Congratulations!
Your request for {{ service.name }} has succeeded. See the incident here: {{response.body.incident.html_url}}
{% else %}
## Oops something went wrong!
Please contact [{{ action_owner.name }}]({{ action_owner.href }}) for more help.
{% endif %}
  EOT
    error_message = "wrong 'manual_inputs_definition' in opslevel_trigger_definition.big"
  }

  assert {
    condition     = opslevel_trigger_definition.big.published == false
    error_message = "published should be set to 'false' for opslevel_trigger_definition.big"
  }
}
