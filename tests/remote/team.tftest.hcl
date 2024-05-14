run "datasource_teams_all" {

  variables {
    datasource_type = "opslevel_teams"
  }

  module {
    source = "./team"
  }

  assert {
    condition     = can(data.opslevel_teams.all.teams)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_teams.all.teams) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.datasource_type)
  }

}

run "datasource_team_first" {

  variables {
    datasource_type = "opslevel_team"
  }

  module {
    source = "./team"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_team.first_team_by_id.alias),
      can(data.opslevel_team.first_team_by_id.id),
      can(data.opslevel_team.first_team_by_id.members),
      can(data.opslevel_team.first_team_by_id.name),
      can(data.opslevel_team.first_team_by_id.parent_alias),
      can(data.opslevel_team.first_team_by_id.parent_id),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_team.first_team_by_alias.alias == data.opslevel_teams.all.teams[0].alias
    error_message = replace(var.error_wrong_alias, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_team.first_team_by_id.id == data.opslevel_teams.all.teams[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.datasource_type)
  }

}
