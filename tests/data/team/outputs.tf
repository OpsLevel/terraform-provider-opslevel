output "all" {
  value = data.opslevel_teams.all
}

output "first" {
  value = data.opslevel_teams.all.teams[0]
}
