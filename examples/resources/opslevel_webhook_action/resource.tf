resource "opslevel_webhook_action" "example" {
  name = "example"
  description = "something about what i do"
  payload = "{\"event_type\": \"{{ input.event }}\"}"
  url = "https://gitlab.com/api/v4/{{ input.project }}/"
  method = "POST"
  headers = {
    content-type = "application/json"
    accept = "application/json"
  }
}
