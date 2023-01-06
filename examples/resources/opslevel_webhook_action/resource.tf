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
