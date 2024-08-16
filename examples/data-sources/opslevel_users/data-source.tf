data "opslevel_users" "all" {}

data "opslevel_users" "only_active" {
  ignore_deactivated = true
}

output "all" {
  value = data.opslevel_users.all.users
}

output "only_active" {
  value = data.opslevel_users.only_active.users
}

output "user_emails" {
  value = sort(data.opslevel_users.all.users[*].email)
}

output "user_names" {
  value = sort(data.opslevel_users.all.users[*].name)
}
