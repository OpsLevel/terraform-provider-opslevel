variables {
  resource_name = "opslevel_team_tag"

  # required fields
  key   = "test-tag-key"
  value = "test-tag-value"

  # optional fields (only one of 'team' or 'team_alias' may be set)
  team       = null
  team_alias = null
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

run "resource_team_tag_create_with_all_fields" {

  variables {
    key        = var.key
    team       = run.from_team_get_owner_id.first_team.id
    team_alias = null
    value      = var.value
  }

  module {
    source = "./team_tag"
  }

  assert {
    condition = alltrue([
      can(opslevel_team_tag.test.id),
      can(opslevel_team_tag.test.key),
      can(opslevel_team_tag.test.team),
      can(opslevel_team_tag.test.team_alias),
      can(opslevel_team_tag.test.value),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition     = startswith(opslevel_team_tag.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_team_tag.test.key == var.key
    error_message = format(
      "expected '%v' but got '%v'",
      var.key,
      opslevel_team_tag.test.key,
    )
  }

  assert {
    condition = opslevel_team_tag.test.team == run.from_team_get_owner_id.first_team.id
    error_message = format(
      "expected '%v' but got '%v'",
      var.team,
      opslevel_team_tag.test.team,
    )
  }

  assert {
    condition     = opslevel_team_tag.test.team_alias == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_team_tag.test.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_tag.test.value,
    )
  }

}

run "resource_team_tag_update_set_all_fields" {

  variables {
    key        = "${var.key}-updated"
    team       = null
    team_alias = run.from_team_get_owner_id.first_team.alias
    value      = "${var.value}-updated"
  }

  module {
    source = "./team_tag"
  }

  assert {
    condition = opslevel_team_tag.test.key == var.key
    error_message = format(
      "expected '%v' but got '%v'",
      var.key,
      opslevel_team_tag.test.key,
    )
  }

  assert {
    condition     = opslevel_team_tag.test.team == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_team_tag.test.team_alias == var.team_alias
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_team_tag.test.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_tag.test.value,
    )
  }

}
