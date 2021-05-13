data "opslevel_teams" "all" {

}

output "all_teams" {
  value = data.opslevel_teams.all
}
