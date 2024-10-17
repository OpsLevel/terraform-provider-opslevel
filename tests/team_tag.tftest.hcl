variables {
  resource_name = "opslevel_team_tag"

  # required fields
  key   = "test-tag-key"
  value = "test-tag-value"

  # optional fields (only one of 'team' or 'team_alias' may be set)
  team       = null
  team_alias = null
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_teams.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_team_tag_create_with_all_fields_using_id" {

  variables {
    # other fields from file scoped variables block
    team = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/team/tag"
  }

  assert {
    condition = alltrue([
      can(opslevel_team_tag.this.id),
      can(opslevel_team_tag.this.key),
      can(opslevel_team_tag.this.team),
      can(opslevel_team_tag.this.team_alias),
      can(opslevel_team_tag.this.value),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition     = startswith(opslevel_team_tag.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_team_tag.this.key == var.key
    error_message = format(
      "expected '%v' but got '%v'",
      var.key,
      opslevel_team_tag.this.key,
    )
  }

  assert {
    condition = opslevel_team_tag.this.team == var.team
    error_message = format(
      "expected '%v' but got '%v'",
      var.team,
      opslevel_team_tag.this.team,
    )
  }

  assert {
    condition     = opslevel_team_tag.this.team_alias == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_team_tag.this.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_tag.this.value,
    )
  }

}

run "resource_team_tag_unset_id_set_alias" {

  variables {
    team_alias = run.from_data_module.first_team.alias
  }

  module {
    source = "./opslevel_modules/modules/team/tag"
  }

  assert {
    condition = run.resource_team_tag_create_with_all_fields_using_id.this.id != opslevel_team_tag.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v' due to forced replacement",
      run.resource_team_tag_create_with_all_fields_using_id.this.id,
      opslevel_team_tag.this.id,
    )
  }

  assert {
    condition     = opslevel_team_tag.this.team == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_team_tag.this.team_alias == var.team_alias
    error_message = var.error_expected_null_field
  }

}

run "delete_team_tag_outside_of_terraform" {

  variables {
    resource_id   = run.resource_team_tag_unset_id_set_alias.this.id
    resource_type = "tag"
  }

  module {
    source = "./provisioner"
  }
}

run "resource_team_tag_create_with_all_fields_using_alias" {

  variables {
    team_alias = run.from_data_module.first_team.alias
  }

  module {
    source = "./opslevel_modules/modules/team/tag"
  }

  assert {
    condition = run.resource_team_tag_create_with_all_fields_using_id.this.id != opslevel_team_tag.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_team_tag_create_with_all_fields_using_id.this.id,
      opslevel_team_tag.this.id,
    )
  }

  assert {
    condition     = opslevel_team_tag.this.team == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_team_tag.this.team_alias == var.team_alias
    error_message = var.error_expected_null_field
  }

}

run "resource_team_tag_unset_alias_set_id" {

  variables {
    team = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/team/tag"
  }

  assert {
    condition = run.resource_team_tag_create_with_all_fields_using_alias.this.id != opslevel_team_tag.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_team_tag_create_with_all_fields_using_id.this.id,
      opslevel_team_tag.this.id,
    )
  }

  assert {
    condition = opslevel_team_tag.this.team == var.team
    error_message = format(
      "expected '%v' but got '%v'",
      var.team,
      opslevel_team_tag.this.team,
    )
  }

  assert {
    condition     = opslevel_team_tag.this.team_alias == null
    error_message = var.error_expected_null_field
  }

}
