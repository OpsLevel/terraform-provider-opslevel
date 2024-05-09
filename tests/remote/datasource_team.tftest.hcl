run "datasource_teams_all" {

  assert {
    condition     = length(data.opslevel_teams.all.teams) > 0
    error_message = "zero teams found in data.opslevel_teams"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_teams.all.teams[0].id),
    ])
    error_message = "cannot set all expected team datasource fields"
  }

}

run "datasource_team_first" {

  assert {
    condition     = data.opslevel_team.first_team_by_alias.alias == data.opslevel_teams.all.teams[0].alias
    error_message = "wrong alias on opslevel_team"
  }

  assert {
    condition     = data.opslevel_team.first_team_by_id.id == data.opslevel_teams.all.teams[0].id
    error_message = "wrong ID on opslevel_team"
  }

}
