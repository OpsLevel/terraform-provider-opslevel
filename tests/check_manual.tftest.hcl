variables {
  check_manual = "opslevel_check_manual"

  # -- check_manual fields --
  # required fields
  update_requires_comment = true

  # optional fields
  update_frequency = {
    starting_date = "2020-02-12T06:36:13Z"
    time_scale    = "week"
    value         = 1
  }

  # -- check base fields --
  # required fields
  category = null # sourced from module
  level    = null # sourced from module
  name     = "TF Test Check Manual"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null # sourced from module
  notes     = "Notes on manual check"
  owner     = null # sourced from module
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_filters.all,
      data.opslevel_rubric_categories.all,
      data.opslevel_rubric_levels.all,
      data.opslevel_teams.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_check_manual_create_with_all_fields" {

  variables {
    # other fields from file scoped variables block
    category = run.from_data_module.first_rubric_category.id
    filter   = run.from_data_module.first_filter.id
    level    = run.from_data_module.max_index_rubric_level.id
    owner    = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/check/manual"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_manual.this.category),
      can(opslevel_check_manual.this.description),
      can(opslevel_check_manual.this.enable_on),
      can(opslevel_check_manual.this.enabled),
      can(opslevel_check_manual.this.filter),
      can(opslevel_check_manual.this.id),
      can(opslevel_check_manual.this.level),
      can(opslevel_check_manual.this.name),
      can(opslevel_check_manual.this.notes),
      can(opslevel_check_manual.this.owner),
      can(opslevel_check_manual.this.update_frequency),
      can(opslevel_check_manual.this.update_requires_comment),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_manual)
  }

  assert {
    condition = opslevel_check_manual.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_manual.this.category,
    )
  }

  assert {
    condition = opslevel_check_manual.this.enable_on == var.enable_on
    error_message = format(
      "expected '%v' but got '%v'",
      var.enable_on,
      opslevel_check_manual.this.enable_on,
    )
  }

  assert {
    condition = opslevel_check_manual.this.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_manual.this.enabled,
    )
  }

  assert {
    condition     = startswith(opslevel_check_manual.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_manual)
  }

  assert {
    condition = opslevel_check_manual.this.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_check_manual.this.filter,
    )
  }

  assert {
    condition = opslevel_check_manual.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_manual.this.level,
    )
  }

  assert {
    condition = opslevel_check_manual.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_manual.this.name,
    )
  }

  assert {
    condition = opslevel_check_manual.this.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_manual.this.notes,
    )
  }

  assert {
    condition = opslevel_check_manual.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_manual.this.owner,
    )
  }

  assert {
    condition = opslevel_check_manual.this.update_frequency == var.update_frequency
    error_message = format(
      "expected '%v' but got '%v'",
      var.update_frequency,
      opslevel_check_manual.this.update_frequency,
    )
  }

  assert {
    condition = opslevel_check_manual.this.update_requires_comment == var.update_requires_comment
    error_message = format(
      "expected '%v' but got '%v'",
      var.update_requires_comment,
      opslevel_check_manual.this.update_requires_comment,
    )
  }

}

run "resource_check_manual_update_unset_optional_fields" {

  variables {
    # other fields from file scoped variables block
    category         = run.from_data_module.first_rubric_category.id
    enable_on        = null
    enabled          = null
    filter           = null
    level            = run.from_data_module.max_index_rubric_level.id
    notes            = null
    owner            = null
    update_frequency = null
  }

  module {
    source = "./opslevel_modules/modules/check/manual"
  }

  assert {
    condition = opslevel_check_manual.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_manual.this.category,
    )
  }

  assert {
    condition     = opslevel_check_manual.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_manual.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_manual.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_manual.this.level,
    )
  }

  assert {
    condition = opslevel_check_manual.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_manual.this.name,
    )
  }

  assert {
    condition     = opslevel_check_manual.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.this.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.this.update_frequency == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_manual.this.update_requires_comment == var.update_requires_comment
    error_message = format(
      "expected '%v' but got '%v'",
      var.update_requires_comment,
      opslevel_check_manual.this.update_requires_comment,
    )
  }

}

run "delete_check_manual_outside_of_terraform" {

  variables {
    command = "delete check ${run.resource_check_manual_create_with_all_fields.this.id}"
  }

  module {
    source = "./cli"
  }
}

run "resource_check_manual_create_with_required_fields" {

  variables {
    # other fields from file scoped variables block
    category         = run.from_data_module.first_rubric_category.id
    enable_on        = null
    enabled          = null
    filter           = null
    level            = run.from_data_module.max_index_rubric_level.id
    notes            = null
    owner            = null
    update_frequency = null
  }

  module {
    source = "./opslevel_modules/modules/check/manual"
  }

  assert {
    condition = run.resource_check_manual_create_with_all_fields.this.id != opslevel_check_manual.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_check_manual_create_with_all_fields.this.id,
      opslevel_check_manual.this.id,
    )
  }

  assert {
    condition = opslevel_check_manual.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_manual.this.category,
    )
  }

  assert {
    condition     = opslevel_check_manual.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_manual.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_manual.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_manual.this.level,
    )
  }

  assert {
    condition = opslevel_check_manual.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_manual.this.name,
    )
  }

  assert {
    condition     = opslevel_check_manual.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.this.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.this.update_frequency == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_manual.this.update_requires_comment == var.update_requires_comment
    error_message = format(
      "expected '%v' but got '%v'",
      var.update_requires_comment,
      opslevel_check_manual.this.update_requires_comment,
    )
  }

}

run "resource_check_manual_set_all_fields" {

  variables {
    # other fields from file scoped variables block
    category = run.from_data_module.first_rubric_category.id
    filter   = run.from_data_module.first_filter.id
    level    = run.from_data_module.max_index_rubric_level.id
    owner    = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/check/manual"
  }

  assert {
    condition = opslevel_check_manual.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_manual.this.category,
    )
  }

  assert {
    condition = opslevel_check_manual.this.enable_on == var.enable_on
    error_message = format(
      "expected '%v' but got '%v'",
      var.enable_on,
      opslevel_check_manual.this.enable_on,
    )
  }

  assert {
    condition = opslevel_check_manual.this.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_manual.this.enabled,
    )
  }

  assert {
    condition     = startswith(opslevel_check_manual.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_manual)
  }

  assert {
    condition = opslevel_check_manual.this.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_check_manual.this.filter,
    )
  }

  assert {
    condition = opslevel_check_manual.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_manual.this.level,
    )
  }

  assert {
    condition = opslevel_check_manual.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_manual.this.name,
    )
  }

  assert {
    condition = opslevel_check_manual.this.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_manual.this.notes,
    )
  }

  assert {
    condition = opslevel_check_manual.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_manual.this.owner,
    )
  }

  assert {
    condition = opslevel_check_manual.this.update_frequency == var.update_frequency
    error_message = format(
      "expected '%v' but got '%v'",
      var.update_frequency,
      opslevel_check_manual.this.update_frequency,
    )
  }

  assert {
    condition = opslevel_check_manual.this.update_requires_comment == var.update_requires_comment
    error_message = format(
      "expected '%v' but got '%v'",
      var.update_requires_comment,
      opslevel_check_manual.this.update_requires_comment,
    )
  }

}
