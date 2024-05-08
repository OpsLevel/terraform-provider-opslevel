data "opslevel_webhook_actions" "all" {}

output "all_webhook_actions" {
  value = data.opslevel_webhook_actions.all
}
