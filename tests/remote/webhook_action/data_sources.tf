data "opslevel_webhook_actions" "all" {}

data "opslevel_webhook_action" "first_webhook_action_by_id" {
  identifier = data.opslevel_webhook_actions.all.webhook_actions[0].id
}
