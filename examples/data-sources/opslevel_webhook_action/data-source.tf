data "opslevel_webhook_action" "by_alias" {
  identifier = "webhook_action_alias"
}

data "opslevel_webhook_action" "by_id" {
  identifier = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS83NzQ0"
}

output "found_by_alias" {
  value = data.opslevel_webhook_action.by_alias
}

output "found_by_id" {
  value = data.opslevel_webhook_action.by_id
}
