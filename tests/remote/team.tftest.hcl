variables {
  team_one  = "opslevel_team"
  teams_all = "opslevel_teams"

  # required fields
  name = "TF Test Team"

  # optional fields
  aliases          = tolist(["test_team_one"])
  parent           = null
  responsibilities = "Team responsibilities"
  # member block
}


run "resource_team_create_with_all_fields" {

  variables {
    aliases          = var.aliases
    name             = var.name
    parent           = var.parent
    responsibilities = var.responsibilities
  }

  module {
    source = "./team"
  }

}

run "resource_team_update_unset_optional_fields" {

  variables {
    aliases          = null
    parent           = null
    responsibilities = null
  }

  module {
    source = "./team"
  }

}

run "resource_team_update_set_all_fields" {

  variables {
    aliases          = concat(var.aliases, ["test_alias"])
    name             = "${var.name} updated"
    parent           = var.parent
    responsibilities = "${var.responsibilities} updated"
  }

  module {
    source = "./team"
  }

}

run "datasource_teams_all" {

  module {
    source = "./team"
  }

  assert {
    condition     = can(data.opslevel_teams.all.teams)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.teams_all)
  }

  assert {
    condition     = length(data.opslevel_teams.all.teams) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.teams_all)
  }

}

run "datasource_team_first" {

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
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.team_one)
  }

  assert {
    condition     = data.opslevel_team.first_team_by_alias.alias == data.opslevel_teams.all.teams[0].alias
    error_message = replace(var.error_wrong_alias, "TYPE", var.team_one)
  }

  assert {
    condition     = data.opslevel_team.first_team_by_id.id == data.opslevel_teams.all.teams[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.team_one)
  }

}
