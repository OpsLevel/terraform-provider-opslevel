data "opslevel_team" "platform" {
    alias = "platform"
}

data "opslevel_filter" "tier_1" {
  filter {
    field = "name"
    value = "Tier 1 Services"
  }
}

resource "opslevel_webhook_action" "example" {
  name = "Page The On Call"
  description = "Pages the On Call"
  url = "https://api.pagerduty.com/incidents"
  method = "POST"
  headers = {
    content-type = "application/json"
    accept = "application/vnd.pagerduty+json;version=2"
    authorization = "Token token=XXXXXXXXXXXXX"
    from = "john@opslevel.com"
  }
  payload = <<EOT
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
}

resource "opslevel_trigger_definition" "example" {
  name = "Page The On Call"
  description = "Pages the On Call"
  owner = data.opslevel_team.platform.id
  filter = data.opslevel_filter.tier_1.id
  action = opslevel_webhook_action.example.id
  access_control = "everyone"
  extended_team_access = ["team_1", "team_2"]
  manual_inputs_definition = <<EOT
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
  response_template = <<EOT
{% if response.status >= 200 and response.status < 300 %}
## Congratulations!
Your request for {{ service.name }} has succeeded. See the incident here: {{response.body.incident.html_url}}
{% else %}
## Oops something went wrong!
Please contact [{{ action_owner.name }}]({{ action_owner.href }}) for more help.
{% endif %}
  EOT
}
