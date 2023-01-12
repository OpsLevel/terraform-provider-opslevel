data "opslevel_users" "all" {
}

output "found" {
  value = data.opslevel_users.all.emails[0]
}