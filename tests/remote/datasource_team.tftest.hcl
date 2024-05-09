run "datasource_teams_all" {

  variables {
    datasource_type = "opslevel_teams"
  }

  assert {
    condition     = can(data.opslevel_teams.all.teams)
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_teams.all.teams) > 0
    error_message = replace(var.empty_datasource_error, "TYPE", var.datasource_type)
  }

}

run "datasource_team_first" {

  variables {
    datasource_type = "opslevel_team"
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
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_team.first_team_by_alias.alias == data.opslevel_teams.all.teams[0].alias
    error_message = replace(var.wrong_alias_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_team.first_team_by_id.id == data.opslevel_teams.all.teams[0].id
    error_message = replace(var.wrong_id_error, "TYPE", var.datasource_type)
  }

}
