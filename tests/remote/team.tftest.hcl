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

run "from_team_get_owner_id" {
  command = plan

  variables {
    aliases          = null
    name             = ""
    parent           = null
    responsibilities = null
  }

  module {
    source = "./team"
  }
}

run "resource_team_create_with_all_fields" {

  variables {
    aliases          = var.aliases
    name             = var.name
    parent           = run.from_team_get_owner_id.first_team.id
    responsibilities = var.responsibilities
  }

  module {
    source = "./team"
  }

  assert {
    condition = alltrue([
      can(opslevel_team.test.aliases),
      can(opslevel_team.test.id),
      can(opslevel_team.test.member),
      can(opslevel_team.test.name),
      can(opslevel_team.test.parent),
      can(opslevel_team.test.responsibilities),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.team_one)
  }

  assert {
    condition     = opslevel_team.test.aliases == toset(var.aliases)
    error_message = "wrong aliases for opslevel_team resource"
  }

  assert {
    condition     = startswith(opslevel_team.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.team_one)
  }

  assert {
    condition     = opslevel_team.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.team_one)
  }

  assert {
    condition     = opslevel_team.test.parent == var.parent
    error_message = "wrong parent for opslevel_team resource"
  }

  assert {
    condition     = opslevel_team.test.responsibilities == var.responsibilities
    error_message = "wrong responsibilities for opslevel_team resource"
  }

}

run "resource_team_create_with_empty_optional_fields" {

  variables {
    name             = "New ${var.name} with empty fields"
    responsibilities = ""
  }

  module {
    source = "./team"
  }

  assert {
    condition     = opslevel_team.test.responsibilities == ""
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
    source = "./team"
  }

  assert {
    condition     = opslevel_team.test.aliases == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_team.test.parent == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_team.test.responsibilities == null
    error_message = var.error_expected_null_field
  }

}

run "resource_team_update_set_all_fields" {

  variables {
    aliases          = setunion(var.aliases, ["test_alias"])
    name             = "${var.name} updated"
    parent           = run.from_team_get_owner_id.first_team.id
    responsibilities = "${var.responsibilities} updated"
  }

  module {
    source = "./team"
  }

  assert {
    condition     = opslevel_team.test.aliases == toset(var.aliases)
    error_message = "wrong aliases for opslevel_team resource"
  }

  assert {
    condition     = opslevel_team.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.team_one)
  }

  assert {
    condition     = opslevel_team.test.parent == var.parent
    error_message = "wrong parent for opslevel_team resource"
  }

  assert {
    condition     = opslevel_team.test.responsibilities == var.responsibilities
    error_message = "wrong responsibilities for opslevel_team resource"
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
