output "all" {
  value = data.opslevel_users.all
}

output "first" {
  value = data.opslevel_users.all.users[0]
}
