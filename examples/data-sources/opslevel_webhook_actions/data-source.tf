data "opslevel_webhook_actions" "all" {}

output "all_webhook_actions" {
  value = data.opslevel_webhook_actions.all.webhook_actions
}

output "webhook_action_names" {
  value = data.opslevel_webhook_actions.all.webhook_actions[*].name
}
