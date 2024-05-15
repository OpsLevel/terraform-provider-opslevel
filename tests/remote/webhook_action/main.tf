data "opslevel_webhook_actions" "all" {}

data "opslevel_webhook_action" "first_webhook_action_by_id" {
  identifier = data.opslevel_webhook_actions.all.webhook_actions[0].id
}

resource "opslevel_webhook_action" "test" {
  description = var.description
  headers     = var.headers
  method      = var.method
  name        = var.name
  payload     = var.payload
  url         = var.url
}
