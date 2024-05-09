data "opslevel_teams" "all" {}

output "all" {
  value = data.opslevel_teams.all.teams
}

output "team_names" {
  value = sort(data.opslevel_teams.all.teams[*].name)
}
