data "opslevel_teams" "all" {}

data "opslevel_team" "first_team_by_alias" {
  alias = data.opslevel_teams.all.teams[0].alias
}

data "opslevel_team" "first_team_by_id" {
  id = data.opslevel_teams.all.teams[0].id
}
