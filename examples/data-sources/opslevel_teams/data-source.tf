data "opslevel_teams" "all" {}

output "found" {
  value = data.opslevel_teams.all.ids[3]
}
