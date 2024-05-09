data "opslevel_users" "all" {}

output "all" {
  value = data.opslevel_users.all.users
}

output "user_emails" {
  value = sort(data.opslevel_users.all.users[*].email)
}

output "user_names" {
  value = sort(data.opslevel_users.all.users[*].name)
}
