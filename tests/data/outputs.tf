output "all_domains" {
  value = data.opslevel_domains.all
}

output "all_filters" {
  value = data.opslevel_filters.all
}

output "all_lifecycles" {
  value = data.opslevel_lifecycles.all
}

output "all_repositories" {
  value = data.opslevel_repositories.all
}

output "all_services" {
  value = data.opslevel_services.all
}

output "all_systems" {
  value = data.opslevel_systems.all
}

output "all_teams" {
  value = data.opslevel_teams.all
}

output "all_tiers" {
  value = data.opslevel_tiers.all
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

output "first_lifecycle" {
  value = data.opslevel_lifecycles.all.lifecycles[0]
}

output "first_repository" {
  value = data.opslevel_repositories.all.repositories[0]
}

output "first_service" {
  value = data.opslevel_services.all.services[0]
}

output "first_system" {
  value = data.opslevel_systems.all.systems[0]
}

output "first_team" {
  value = data.opslevel_teams.all.teams[0]
}

output "first_tier" {
  value = data.opslevel_tiers.all.tiers[0]
}

output "first_user" {
  value = data.opslevel_users.all.users[0]
}
