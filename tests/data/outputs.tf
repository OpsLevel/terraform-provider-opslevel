output "all_domains" {
  value = data.opslevel_domains.all
}

output "all_filters" {
  value = data.opslevel_filters.all
}

output "all_repositories" {
  value = data.opslevel_repositories.all
}

output "all_services" {
  value = data.opslevel_services.all
}

output "all_teams" {
  value = data.opslevel_teams.all
}

output "all_users" {
  value = data.opslevel_users.all
}

output "first_domain" {
  value = data.opslevel_domains.all.domains[0]
}

output "first_filter" {
  value = data.opslevel_filters.all.filters[0]
}

output "first_repository" {
  value = data.opslevel_repositories.all.repositories[0]
}

output "first_service" {
  value = data.opslevel_services.all.services[0]
}

output "first_team" {
  value = data.opslevel_teams.all.teams[0]
}

output "first_user" {
  value = data.opslevel_users.all.users[0]
}
