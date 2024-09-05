variables {
  team_one  = "opslevel_team"
  teams_all = "opslevel_teams"

  # required fields
  name = "TF Test Team"

  # optional fields
  aliases          = ["test_team_foo_bar_baz"]
  parent           = null
  responsibilities = "Team responsibilities"
  # member block
}

run "from_team_module" {
  command = plan

  variables {
    aliases          = null
    name             = ""
    parent           = null
    responsibilities = null
  }

  module {
    source = "./opslevel_modules/modules/team"
  }
}

run "resource_team_create_with_all_fields" {

  variables {
    aliases          = var.aliases
    name             = var.name
    parent           = run.from_team_module.all.teams[0].id
    responsibilities = var.responsibilities
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition = alltrue([
      can(opslevel_team.this.aliases),
      can(opslevel_team.this.id),
      can(opslevel_team.this.member),
      can(opslevel_team.this.name),
      can(opslevel_team.this.parent),
      can(opslevel_team.this.responsibilities),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.team_one)
  }

  assert {
    condition     = opslevel_team.this.aliases == toset(var.aliases)
    error_message = "wrong aliases for opslevel_team resource"
  }

  assert {
    condition     = startswith(opslevel_team.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.team_one)
  }

  assert {
    condition     = opslevel_team.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.team_one)
  }

  assert {
    condition     = opslevel_team.this.parent == var.parent
    error_message = "wrong parent for opslevel_team resource"
  }

  assert {
    condition     = opslevel_team.this.responsibilities == var.responsibilities
    error_message = "wrong responsibilities for opslevel_team resource"
  }

}

run "resource_team_create_with_empty_optional_fields" {

  variables {
    name             = "New ${var.name} with empty fields"
    responsibilities = ""
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition     = opslevel_team.this.responsibilities == ""
    error_message = var.error_expected_empty_string
  }

}

run "resource_team_update_unset_optional_fields" {

  variables {
    aliases          = null
    parent           = null
    responsibilities = null
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition     = opslevel_team.this.aliases == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_team.this.parent == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_team.this.responsibilities == null
    error_message = var.error_expected_null_field
  }

}

run "resource_team_update_set_all_fields" {

  variables {
    aliases          = setunion(var.aliases, ["test_alias"])
    name             = "${var.name} updated"
    parent           = run.from_team_module.all.teams[0].id
    responsibilities = "${var.responsibilities} updated"
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition     = opslevel_team.this.aliases == toset(var.aliases)
    error_message = "wrong aliases for opslevel_team resource"
  }

  assert {
    condition     = opslevel_team.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.team_one)
  }

  assert {
    condition     = opslevel_team.this.parent == var.parent
    error_message = "wrong parent for opslevel_team resource"
  }

  assert {
    condition     = opslevel_team.this.responsibilities == var.responsibilities
    error_message = "wrong responsibilities for opslevel_team resource"
  }

}

run "datasource_teams_all" {

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_teams.all.teams[0].alias),
      # can(data.opslevel_teams.all.teams[0].contacts),
      can(data.opslevel_teams.all.teams[0].id),
      can(data.opslevel_teams.all.teams[0].members),
      can(data.opslevel_teams.all.teams[0].name),
      can(data.opslevel_teams.all.teams[0].parent_alias),
      can(data.opslevel_teams.all.teams[0].parent_id),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.team_one)
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
