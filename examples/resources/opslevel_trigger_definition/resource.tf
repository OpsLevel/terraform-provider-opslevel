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
  name = "example"
  description = "something about what i do"
  liquid_template = "{\"event_type\": \"{{ input.event }}\"}"
  webhook_url = "https://gitlab.com/api/v4/{{ input.project }}/"
  http_method = "POST"
  headers = {
    content-type = "application/json"
    accept = "application/json"
  }
}

resource "opslevel_trigger_definition" "example" {
  name = "example"
  description = "something about what i do"
  owner = data.opslevel_team.platform.id
  action = opslevel_webhook_action.example.id
}
